package skywalker

import (
    "github.com/michaelhakansson/skylark/db"
    "github.com/michaelhakansson/skylark/structures"
    "github.com/michaelhakansson/skylark/utils"
    "math"
    "os"
)

func updateChangeFrequencyForAll() {
    _, err := os.Stdout.Write([]byte("Updating change frequency\n"))
    utils.Checkerr(err)
    ids := db.GetAllShowIDs()
    t := int64(len(ids))
    for i, id := range ids {
        show := db.GetShowByPlayID(id)
        changefrequency := calcChangeFrequency(show)
        db.UpdateShowChangeFrequency(show.ID, changefrequency)
        utils.PrintProgressBar(int64(i), t)
    }
    utils.PrintCompletedBar()
}

func calcChangeFrequency(show structures.Show) float64 {
    var tot float64
    episodes := structures.SortEpisodesByDate(show.Episodes)
    n := len(episodes)
    for i := 0; i < n-1; i++ {
        ed1 := episodes[i].Broadcasted
        ed2 := episodes[i+1].Broadcasted
        delta := ed1.Sub(ed2).Hours()
        delta = delta / 168
        if delta > 2 || delta == 0 {
            delta = 1
        }
        tot += delta
    }
    avg := tot / float64(n-1)
    if avg <= 0 {
        avg = 1
    }
    cf := (1 / avg)
    if math.IsNaN(cf) {
        cf = 1
    }
    return math.Min(cf, 100)
}
