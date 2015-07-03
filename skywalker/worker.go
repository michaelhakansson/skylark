package skywalker

import (
    "github.com/michaelhakansson/skylark/db"
    "github.com/michaelhakansson/skylark/services"
    "github.com/michaelhakansson/skylark/structures"
    "github.com/michaelhakansson/skylark/utils"
    "strconv"
    "sync"
    "sync/atomic"
)

var counter = new(int64)
var total = new(int64)

func processShows(shows []structures.Show) {
    atomic.StoreInt64(total, 0)
    atomic.StoreInt64(counter, 0)
    atomic.StoreInt64(total, int64(len(shows)))
    in := genShowChan(shows)
    var wg sync.WaitGroup
    for i := 0; i < 20; i++ {
        wg.Add(1)
        updateWorker(in, &wg)
    }
    wg.Wait()
    utils.PrintCompletedBar()
}

func genShowChan(shows []structures.Show) <-chan structures.Show {
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
            syncShow(show.PlayID, show.PlayService)
            atomic.AddInt64(counter, 1)
            c := atomic.LoadInt64(counter)
            t := atomic.LoadInt64(total)
            utils.PrintProgressBar(c, t)
        }
    }()
}

func syncShow(showid string, playservice string) {
    show, episodeids := services.GetShowWithService(showid, playservice)
    if show.PlayID == "" {
        return
    }
    _, dbShowObject := db.AddShowInfo(show.Title, show.Thumbnail, show.PlayID, show.PlayService)

    for _, episodeid := range episodeids {
        eid, err := strconv.ParseInt(episodeid, 10, 64)
        utils.Checkerr(err)
        episode := services.GetEpisodeWithService(episodeid, playservice)
        if db.ContainsEpisodeWithPlayID(dbShowObject.PlayID, eid) {
            db.UpdateEpisode(dbShowObject.ID, episode)
        } else {
            db.AddEpisode(dbShowObject.ID, episode)
        }
    }
    SyncThumbnail(show)
}
