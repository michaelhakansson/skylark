package main

import(
    "log"
    "sync"
    "time"
    "github.com/michaelhakansson/skylark/services/svtplay"
    "github.com/michaelhakansson/skylark/services/tv3play"
    "github.com/michaelhakansson/skylark/services/kanal5play"
    "github.com/michaelhakansson/skylark/db"
    "github.com/michaelhakansson/skylark/structures"
)

const freshnessLimit float64 = 0.5
var services []string = []string{"svtplay", "tv3play", "kanal5play"}

func syncAll() {
    for _, service := range services {
        ids := getIdsWithService(service)
        // Add all shows and episodes to the DB
        for _, id := range ids {
            syncShow(id, service)
        }
    }
}

func syncShow(showId string, playservice string) {
    show, episodes := getShowWithService(showId, playservice)
    _, dbShowObject := db.AddShow(show.Title, show.PlayId, show.PlayService)

    for _, episode := range episodes {
        db.AddEpisode(dbShowObject.Id, episode)
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

func getShowFreshness(episodes []structures.Episode) (freshness float64) {
    for _, episode := range episodes {
        freshness = (freshness + episode.Freshness)
    }
    freshness /= float64(len(episodes))
    return
}

func syncLowFreshnessShows() {
    var wg sync.WaitGroup
    var showsToSync []structures.Show
    showIds := db.GetAllShowIds()
    for _, showId := range showIds {
        wg.Add(1)
        show := db.GetShowByPlayId(showId)
        go func() {
            defer wg.Done()
            freshness := getShowFreshness(show.Episodes)
            if freshness < freshnessLimit {
                showsToSync = append(showsToSync, show)
            }
        }()
        wg.Wait()
        for _, show := range showsToSync {
            syncShow(show.PlayId, show.PlayService)
        }
    }
}

func main() {
    go func() {
        timer := time.Tick(15 * time.Minute)
        for now := range timer {
            log.Println(now)
            syncLowFreshnessShows()
        }
    }()
}
