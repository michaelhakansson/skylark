package db

import(
    "log"
    "time"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

const(
    uri string = "mongodb://localhost:27017/skylark"
    db string = "skylark"
)

type Show struct {
    Id bson.ObjectId `bson:"_id,omitempty"`
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

func Connect() *mgo.Session {
    session, err := mgo.Dial(uri)
    session.SetSafe(&mgo.Safe{})
    if err != nil {
        log.Fatal(err)
    }
    return session
}

