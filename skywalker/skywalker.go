package main

import(
    "log"
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

func getShowFreshness(show structures.Show) (freshness float64) {
    for _, episode := range show.Episodes {
        freshness = (freshness + episode.Freshness)
    }
    freshness /= float64(len(show.Episodes))
    return
}

func main() {
    syncAll()
    go func() {
        timer := time.Tick(15 * time.Minute)
        for now := range timer {
            log.Println(now)
            showIds := db.GetAllShowIds()
            for _, showId := range showIds {
                go func() {
                    show := db.GetShowByPlayId(showId)
                    freshness := getShowFreshness(show)
                    if freshness < freshnessLimit {
                        syncShow(showId, show.PlayService)
                    }
                }()
            }
        }
    }()
}
