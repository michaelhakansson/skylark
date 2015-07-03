package skywalker

import (
    "github.com/michaelhakansson/skylark/db"
    "github.com/michaelhakansson/skylark/structures"
    "github.com/michaelhakansson/skylark/utils"
    "io"
    "net/http"
    "os"
    "os/exec"
    "strconv"
)

var imgLocation = utils.GetPath() + "/img/"

func init() {
    os.MkdirAll(imgLocation, 0755)
}

// SyncThumbnails synchronises the thumbnails of all the shows in the database
func SyncThumbnails() {
    _, err := os.Stdout.Write([]byte("Synchronising thumbnails\n"))
    utils.Checkerr(err)
    ids := db.GetAllShowIDs()
    t := int64(len(ids))
    for i, id := range ids {
        show := db.GetShowByPlayID(id)
        SyncThumbnail(show)
        utils.PrintProgressBar(int64(i), t)
    }
    utils.PrintCompletedBar()
}

// SyncThumbnail synchronises a thumbnail for the specified show
func SyncThumbnail(show structures.Show) {
    //log.Printf("Started sync of thumbnails for %s", show.Title)
    sizes := []string{"96"} // Support for multiple sizes
    createThumbnails(show.PlayID, show.Thumbnail, sizes)
    for _, episode := range show.Episodes {
        eid := strconv.FormatInt(episode.PlayID, 10)
        createThumbnails(eid, episode.Thumbnail, sizes)
    }
    //log.Printf("Sync of thumbnails for %s is completed", show.Title)
}

func createThumbnails(id string, url string, sizes []string) {
    if len(url) == 0 {
        return
    }
    for _, size := range sizes {
        // TODO replace exists check
        _, err := os.Stat(imgLocation + id + "-" + size + ".png")
        if os.IsNotExist(err) {
            _, err = os.Stat(imgLocation + id + "-org.jpg")
            if os.IsNotExist(err) {
                downloadThumbnail(id, url)
            }
            resizeThumbnail(id, size)
            err := os.Remove(imgLocation + id + "-org.jpg")
            utils.Checkerr(err)
        }
    }
    return
}

func resizeThumbnail(id string, size string) {
    filename := id + "-" + size + ".png"
    args := []string{"-s", size, "-o", imgLocation + filename, imgLocation + id + "-org.jpg"}
    cmd := exec.Command("vipsthumbnail", args...)
    err := cmd.Start()
    utils.Checkerr(err)
    err = cmd.Wait()
    utils.Checkerr(err)
}

func downloadThumbnail(id string, url string) {
    out, err := os.Create(imgLocation + id + "-org.jpg")
    defer out.Close()
    utils.Checkerr(err)
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    utils.Checkerr(err)
    req.Header.Set("User-Agent", useragent)
    resp, err := client.Do(req)
    defer resp.Body.Close()
    utils.Checkerr(err)
    _, err = io.Copy(out, resp.Body)
    utils.Checkerr(err)
}
