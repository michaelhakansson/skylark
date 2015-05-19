package main

import(
    "log"
    "math"
    "sort"
    //    "sync"
    "time"
    "github.com/michaelhakansson/skylark/services/svtplay"
    "github.com/michaelhakansson/skylark/services/tv3play"
    "github.com/michaelhakansson/skylark/services/kanal5play"
    "github.com/michaelhakansson/skylark/db"
    "github.com/michaelhakansson/skylark/structures"
)

const freshnessLimit float64 = 0.5
var services []string = []string{"svtplay", "tv3play", "kanal5play"}

func syncNew() {
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
        syncShow(show.PlayId, show.PlayService)
    }
    updateChangeFrequencyForAll()
}

func updateChangeFrequencyForAll() {
    ids := db.GetAllShowIds()
    for _, id := range ids {
        show := db.GetShowByPlayId(id)
        show.ChangeFrequency = calcChangeFrequency(show)
        db.UpdateShowWithData(show)
    }
}

func syncShow(showId string, playservice string) {
    show, episodes := getShowWithService(showId, playservice)
    _, dbShowObject := db.AddShowInfo(show.Title, show.PlayId, show.PlayService)

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

func calcChangeFrequency(show structures.Show) float64 {
    var tot float64
    episodes := sortEpisodesByDate(show.Episodes)
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
    return cf
}

func sortEpisodesByDate(episodes []structures.Episode) (structures.Episodes) {
    episodes_sorted := make(structures.Episodes, 0, len(episodes))
    for _, episode := range episodes {
        episodes_sorted = append(episodes_sorted, episode)
    }
    sort.Sort(episodes_sorted)
    return episodes_sorted
}

func main() {
    syncNew()
    go func() {
        timer := time.Tick(10 * time.Minute)
        for now := range timer {
            log.Println(now)
            // Sync shows according to lastupdated and change frequence
        }
    }()
}
