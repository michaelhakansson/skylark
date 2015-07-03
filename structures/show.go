package structures

import (
    "gopkg.in/mgo.v2/bson"
    "time"
)

// Show describes the structure of a show
type Show struct {
    ID              bson.ObjectId `bson:"_id,omitempty"`
    ChangeFrequency float64       `bson:"changefrequency"`
    LastUpdated     time.Time     `bson:"lastupdated"`
    Title           string        `bson:"title"`
    PlayID          string        `bson:"playid"`
    PlayService     string        `bson:"playservice"`
    Thumbnail       string        `bson:"thumbnail"`
    Episodes        []Episode     `bson:"episodes"`
}
