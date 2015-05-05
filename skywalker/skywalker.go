package main

import(
    // "log"
    "time"
    "github.com/michaelhakansson/skylark/services/svtplay"
    "github.com/michaelhakansson/skylark/services/tv3play"
    "github.com/michaelhakansson/skylark/services/kanal5play"
    "github.com/michaelhakansson/skylark/db"
)

const freshnessLimit float64 = 0.5

func syncAll() {
    for _, playService := range services {
        showsIds := playService.GetAllProgramIds()

        // Add all shows and episodes to the DB
        for _, id := range showIds {
            syncShow(id)
        }
    }
}

func syncShow(showId string) {
    show, episodes := playService.GetShow(showId)
    _, dbShowObject := db.AddShow(show.Title, show.PlayId, show.PlayService)

    for _, episode := range episodes {
        db.AddEpisode(dbShowObject.Id, episode.Broadcasted, episode.Category,
            episode.Description, episode.Episodenumber, episode.Length, episode.Live, 
            episode.Playid, episode.Season, episode.Thumbnail, episode.Title, episode.Videourl)
    }
}

func getShowFreshness(show db.Show) (freshness float64) {
    for _, episode := range show.Episodes {
        freshness = (freshness + episode.Freshness)
    }
    freshness /= len(show)
    return
}

func main() {
    go func() {
        timer := time.Tick(15 * time.Minute)
        for now := range timer {
            showIds := db.GetAllShowIds()
            for _, showId := range showIds {
                go func() {
                    show := db.GetShowByPlayId(showId)
                    freshness := getShowFreshness(show)
                    if freshness < freshnessLimit {
                        syncShow(showId)
                    }
                }()
            }
        }
    }()
}
