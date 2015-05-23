package skywalker

import(
    "io"
    "log"
    "math"
    "net/http"
    "os"
    "os/exec"
    "strconv"
    //"sort"
    //"sync"
    //"time"
    "github.com/michaelhakansson/skylark/services/svtplay"
    "github.com/michaelhakansson/skylark/services/tv3play"
    "github.com/michaelhakansson/skylark/services/kanal5play"
    "github.com/michaelhakansson/skylark/db"
    "github.com/michaelhakansson/skylark/structures"
)

const(
    freshnessLimit float64 = 0.5
    imgLocation string = "./img/"
    useragent string = "mozilla/5.0 (iphone; cpu iphone os 7_0_2 like mac os x) applewebkit/537.51.1 (khtml, like gecko) version/7.0 mobile/11a501 safari/9537.53"
)
var services []string = []string{"svtplay", "tv3play", "kanal5play"}

func SyncNew() {
    var showsToUpdate []structures.Show
    for _, service := range services {
        ids := getIdsWithService(service)
        // Add all shows and episodes to the DB
        for _, id := range ids {
            added, show := db.AddShow(id, service)
            if added {
                showsToUpdate = append(showsToUpdate, show)
            }
        }
    }

    for _, show := range showsToUpdate {
        SyncShow(show.PlayId, show.PlayService)
    }
    updateChangeFrequencyForAll()
}

func updateChangeFrequencyForAll() {
    ids := db.GetAllShowIds()
    for _, id := range ids {
        show := db.GetShowByPlayId(id)
        changefrequency := calcChangeFrequency(show)
        db.UpdateShowChangeFrequency(show.Id, changefrequency)
    }
}

func SyncShow(showId string, playservice string) {
    show, episodes := getShowWithService(showId, playservice)
    _, dbShowObject := db.AddShowInfo(show.Title, show.Thumbnail, show.PlayId, show.PlayService)

    for _, episode := range episodes {
        db.AddEpisode(dbShowObject.Id, episode)
    }
    SyncThumbnail(show)
}

func SyncThumbnails() {
    ids := db.GetAllShowIds()
    for _, id := range ids {
        show := db.GetShowByPlayId(id)
        SyncThumbnail(show)
    }
}

func SyncThumbnail(show structures.Show) {
    log.Printf("Started sync of thumbnails for %s", show.Title)
    sizes := []string{"96", "256"} // Small and large thumbnails
    prepareThumbnails(show.PlayId, show.Thumbnail, sizes)
    for _, episode := range show.Episodes {
        eid := strconv.FormatInt(episode.PlayId, 10)
        prepareThumbnails(eid, episode.Thumbnail, sizes)
    }
    log.Printf("Sync of thumbnails for %s is completed", show.Title)
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
    if err != nil {
        log.Fatal(err)
    }
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    if err != nil {
        log.Fatal(err)
    }
    req.Header.Set("User-Agent", useragent)
    resp, err := client.Do(req)
    defer resp.Body.Close()
    if err != nil {
        log.Fatal(err)
    }
    _, err = io.Copy(out, resp.Body)
    if err != nil {
        log.Fatal(err)
    }
}

func getIdsWithService(playservice string) (ids []string) {
    switch playservice {
        case "svtplay": ids = svtplay.GetAllProgramIds()
        case "tv3play": ids = tv3play.GetAllProgramIds()
        case "kanal5play": ids = kanal5play.GetAllProgramIds()
    }
    return
}

func getShowWithService(showId string, playservice string) (show structures.Show, episodes []structures.Episode) {
    switch playservice {
        case "svtplay": show, episodes = svtplay.GetShow(showId)
        case "tv3play": show, episodes = tv3play.GetShow(showId)
        case "kanal5play": show, episodes = kanal5play.GetShow(showId)
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

/*func main() {
    syncNew()
    go func() {
        timer := time.Tick(10 * time.Minute)
        for now := range timer {
            log.Println(now)
            ids := db.GetAllShowIds()
            for _, id := range ids {
                show := db.GetShowByPlayId(id)
                // Max time since update allowed (in seconds)
                maxTimeSinceUpdate := (24 / show.ChangeFrequency)
                if time.Now().Sub(show.LastUpdated).Hours() > maxTimeSinceUpdate {
                    syncShow(show.PlayId, show.PlayService)
                    db.UpdateShowWithData(show)
                }
            }
        }
    }()
}*/
