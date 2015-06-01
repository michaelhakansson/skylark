package skywalker

import(
    "io"
    "log"
    "math"
    "net/http"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "sync"
    "sync/atomic"
    "time"
    "github.com/michaelhakansson/skylark/services/svtplay"
    "github.com/michaelhakansson/skylark/services/tv3play"
    "github.com/michaelhakansson/skylark/services/kanal5play"
    "github.com/michaelhakansson/skylark/db"
    "github.com/michaelhakansson/skylark/structures"
)

const(
    freshnessLimit float64 = 0.5
    imgLocation string = "./img/"
    progressBarLength int = 68
    useragent string = "mozilla/5.0 (iphone; cpu iphone os 7_0_2 like mac os x) applewebkit/537.51.1 (khtml, like gecko) version/7.0 mobile/11a501 safari/9537.53"
)
var services []string = []string{"svtplay", "tv3play", "kanal5play"}
var counter = new(int64)
var total = new(int64)


func SyncNew() {
    var newShows []structures.Show
    for _, service := range services {
        ids := getIdsWithService(service)
        for _, id := range ids {
            added, show := db.AddShow(id, service)
            if added {
                newShows = append(newShows, show)
            }
        }
    }
    os.Stdout.Write([]byte("Updating new shows\n"))
    processShows(newShows)
    updateChangeFrequencyForAll()
}

func SyncOutdated() {
    var showsToUpdate []structures.Show
    ids := db.GetAllShowIds()
    for _, id := range ids {
        show := db.GetShowByPlayId(id)
        // Maximum allowed time since update
        maxTimeSinceUpdate := (24 / show.ChangeFrequency)
        if time.Now().Sub(show.LastUpdated).Hours() > maxTimeSinceUpdate {
            showsToUpdate = append(showsToUpdate, show)
        }
    }
    os.Stdout.Write([]byte("Updating outdated shows\n"))
    processShows(showsToUpdate)
}

func processShows(shows []structures.Show) {
    atomic.StoreInt64(total, 0)
    atomic.StoreInt64(counter, 0)
    atomic.StoreInt64(total, int64(len(shows)))
    in := gen(shows)
    var wg sync.WaitGroup
    for i := 0; i < 20; i++ {
        wg.Add(1)
        updateWorker(in, &wg)
    }
    wg.Wait()
    printCompletedBar()
}

func gen(shows []structures.Show) <-chan structures.Show {
    out := make(chan structures.Show)
    go func() {
        for _, show := range shows {
            out <- show
        }
        close(out)
    }()
    return out
}

func updateWorker(in <-chan structures.Show, wg *sync.WaitGroup) {
    go func() {
        defer wg.Done()
        for show := range in {
            SyncShow(show.PlayId, show.PlayService)
            atomic.AddInt64(counter, 1)
            c := atomic.LoadInt64(counter)
            t := atomic.LoadInt64(total)
            printProgressBar(c, t)
        }
    }()
}

func printProgressBar(progress int64, complete int64) {
    amount := int((int64(progressBarLength - 1) * progress) / complete)
    rest := (progressBarLength - 1) - amount
    bar := strings.Repeat("=", amount) + ">" + strings.Repeat(".", rest)
    os.Stdout.Write([]byte("Progress: [" + bar  +"]\r"))
}

func printCompletedBar() {
    os.Stdout.Write([]byte("Progress: [" + strings.Repeat("=", progressBarLength)  +"]\r"))
    os.Stdout.Write([]byte("Completed \n"))
}

func updateChangeFrequencyForAll() {
    os.Stdout.Write([]byte("Updating change frequency\n"))
    ids := db.GetAllShowIds()
    t := int64(len(ids))
    for i, id := range ids {
        show := db.GetShowByPlayId(id)
        changefrequency := calcChangeFrequency(show)
        db.UpdateShowChangeFrequency(show.Id, changefrequency)
        printProgressBar(int64(i), t)
    }
    printCompletedBar()
}

func SyncShow(showId string, playservice string) {
    show, episodeIds := getShowWithService(showId, playservice)
    _, dbShowObject := db.AddShowInfo(show.Title, show.Thumbnail, show.PlayId, show.PlayService)

    for _, episodeId := range episodeIds {
        eId, err := strconv.ParseInt(episodeId, 10, 64)
        checkerr(err)
        episode := getEpisodeWithService(episodeId, playservice)
        if db.ContainsEpisodeWithPlayId(dbShowObject.PlayId, eId) {
            db.UpdateEpisode(dbShowObject.Id, episode)
        } else {
            db.AddEpisode(dbShowObject.Id, episode)
        }
    }
    SyncThumbnail(show)
}

func SyncThumbnails() {
    os.Stdout.Write([]byte("Synchronising thumbnails\n"))
    ids := db.GetAllShowIds()
    t := int64(len(ids))
    for i, id := range ids {
        show := db.GetShowByPlayId(id)
        SyncThumbnail(show)
        printProgressBar(int64(i), t)
    }
    printCompletedBar()
}

func SyncThumbnail(show structures.Show) {
    //log.Printf("Started sync of thumbnails for %s", show.Title)
    sizes := []string{"96"} // Support for multiple sizes
    prepareThumbnails(show.PlayId, show.Thumbnail, sizes)
    for _, episode := range show.Episodes {
        eid := strconv.FormatInt(episode.PlayId, 10)
        prepareThumbnails(eid, episode.Thumbnail, sizes)
    }
    //log.Printf("Sync of thumbnails for %s is completed", show.Title)
}

func prepareThumbnails(id string, url string, sizes []string) {
    if len(url) == 0 {
        return
    }
    for _, size := range sizes {
        f1, err := os.Open(imgLocation + id + "-" + size + ".png")
        f1.Close()
        if err != nil {
            f2, err := os.Open(imgLocation + id + "-org.jpg")
            f2.Close()
            if err != nil {
                downloadThumbnail(id, url)
            }
            resizeThumbnail(id, size)
        }
    }
    os.Remove("./img/" + id + "-org.jpg")
    return
}

func resizeThumbnail(id string, size string) {
    filename := id + "-" + size + ".png"
    args := []string{"-s", size, "-o", "./" + filename, imgLocation + id + "-org.jpg"}
    cmd := exec.Command("vipsthumbnail", args...)
    cmd.Start()
    cmd.Wait()
}

func downloadThumbnail(id string, url string) {
    out, err := os.Create(imgLocation + id + "-org.jpg")
    defer out.Close()
    checkerr(err)
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    checkerr(err)
    req.Header.Set("User-Agent", useragent)
    resp, err := client.Do(req)
    defer resp.Body.Close()
    checkerr(err)
    _, err = io.Copy(out, resp.Body)
    checkerr(err)
}

func getIdsWithService(playservice string) (ids []string) {
    switch playservice {
        case "svtplay": ids = svtplay.GetAllProgramIds()
        case "tv3play": ids = tv3play.GetAllProgramIds()
        case "kanal5play": ids = kanal5play.GetAllProgramIds()
    }
    return
}

func getShowWithService(showId string, playservice string) (show structures.Show, episodeIds []string) {
    switch playservice {
        case "svtplay": show, episodeIds = svtplay.GetShow(showId)
        case "tv3play": show, episodeIds = tv3play.GetShow(showId)
        case "kanal5play": show, episodeIds = kanal5play.GetShow(showId)
    }
    return
}

func getEpisodeWithService(episodeId string, playservice string) (episode structures.Episode) {
    switch playservice {
        case "svtplay": episode = svtplay.GetEpisode(episodeId)
        case "tv3play": episode = tv3play.GetEpisode(episodeId)
        case "kanal5play": episode = kanal5play.GetEpisode(episodeId)
    }
    return
}

func calcChangeFrequency(show structures.Show) float64 {
    var tot float64
    episodes := structures.SortEpisodesByDate(show.Episodes)
    n := len(episodes)
    for i := 0; i < n - 1; i++ {
        ed1 := episodes[i].Broadcasted
        ed2 := episodes[i + 1].Broadcasted
        delta := ed1.Sub(ed2).Hours()
        delta = delta / 168
        if delta > 2 || delta == 0 {
            delta = 1
        }
        tot += delta
    }
    avg := tot / float64(n - 1)
    if avg <= 0 {
        avg = 1
    }
    cf := (1 / avg)
    if math.IsNaN(cf) {
        cf = 1
    }
    return math.Min(cf,100)
}

func checkerr(err error) {
    if err != nil {
        log.Fatal(err)
    }
}

