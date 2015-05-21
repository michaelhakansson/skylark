package structures

import(
    "sort"
    "time"
    "gopkg.in/mgo.v2/bson"
)

type Show struct {
    Id bson.ObjectId `bson:"_id,omitempty"`
    ChangeFrequency float64 `bson:"changefrequency"`
    LastUpdated time.Time `bson:"lastupdated"`
    Title string `bson:"title"`
    PlayId string `bson:"playid"`
    PlayService string `bson:"playservice"`
    Episodes []Episode `bson:"episodes"`
}

type Episode struct {
    Broadcasted time.Time `bson:"broadcasted"`
    Category string `bson:"category"`
    Description string `bson:"description"`
    EpisodeNumber string `bson:"episodenumber"`
    Length string `bson:"length"`
    Live bool `bson:"live"`
    PlayId int64 `bson:"playid"`
    Season string `bson:"season"`
    Thumbnail string `bson:"thumbnail"`
    Title string `bson:"title"`
    VideoUrl string `bson:"videourl"`
}

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

func SortEpisodesByDate(episodes []Episode) (Episodes) {
    episodes_sorted := make(Episodes, 0, len(episodes))
    for _, episode := range episodes {
        episodes_sorted = append(episodes_sorted, episode)
    }
    sort.Sort(episodes_sorted)
    return episodes_sorted
}
