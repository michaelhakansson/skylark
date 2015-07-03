package structures

import (
    "sort"
    "time"
)

// Episode describes the structure of an episode
type Episode struct {
    Broadcasted   time.Time `bson:"broadcasted"`
    Category      string    `bson:"category"`
    Description   string    `bson:"description"`
    EpisodeNumber string    `bson:"episodenumber"`
    Length        string    `bson:"length"`
    Live          bool      `bson:"live"`
    PlayID        int64     `bson:"playid"`
    Season        string    `bson:"season"`
    Thumbnail     string    `bson:"thumbnail"`
    Title         string    `bson:"title"`
    VideoURL      string    `bson:"videourl"`
}

// Episodes is a list with objects of type episode
type Episodes []Episode

func (e Episodes) Len() int {
    return len(e)
}

func (e Episodes) Less(i, j int) bool {
    return e[i].Broadcasted.After(e[j].Broadcasted)
}

func (e Episodes) Swap(i, j int) {
    e[i], e[j] = e[j], e[i]
}

// SortEpisodesByDate sorts the episodes in a list by the date of which they
// were broadcasted on
func SortEpisodesByDate(episodes []Episode) Episodes {
    episodessorted := make(Episodes, 0, len(episodes))
    for _, episode := range episodes {
        episodessorted = append(episodessorted, episode)
    }
    sort.Sort(episodessorted)
    return episodessorted
}
