package structures

import(
    "time"
    "gopkg.in/mgo.v2/bson"
)

type Show struct {
    Id bson.ObjectId `bson:"_id,omitempty"`
    Title string `bson:"title"`
    PlayId string `bson:"playid"`
    PlayService string `bson:"playservice"`
    Episodes []Episode `bson:"episodes"`
}

type Episode struct {
    Freshness float64 `bson:"freshness"`
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
