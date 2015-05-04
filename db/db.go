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

func AddShow(title string, playid string, playservice string) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    count, err := c.Find(bson.M{"playid": playid}).Count()
    if err != nil || count > 0 {
        return false
    }
    var episodes []Episode
    err = c.Insert(&Show{Id: bson.NewObjectId(), Title: title, PlayId: playid,
    PlayService: playservice, Episodes: episodes})
    if err != nil {
        log.Fatal(err)
        return false
    }
    log.Printf("Added %s", title)
    return true
}

func GetShowByPlayId(playid string) Show {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    var s Show
    err := c.Find(bson.M{"playid": playid}).One(&s)
    if err != nil {
        log.Fatal(err)
    }
    return s
}

func GetShowById(showid bson.ObjectId) Show {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    var s Show
    err := c.Find(bson.M{"_id": showid}).One(&s)
    if err != nil {
        log.Fatal(err)
    }
    return s
}

func AddEpisode(showid bson.ObjectId, broadcasted time.Time, category string,
description string, episodenumber string, length string, live bool, playid int64,
season string, thumbnail string, title string, videourl string) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    episode := &Episode{Broadcasted: broadcasted,
    Category: category, Description: description, EpisodeNumber: episodenumber,
    Length: length, Live: live, PlayId: playid, Season: season, Thumbnail:
    thumbnail, Title: title, VideoUrl: videourl}
    show := GetShowById(showid)
    for _, e := range show.Episodes {
        if e.PlayId == playid {
            UpdateEpisode(showid, *episode)
            return true
        }
    }
    show.Episodes = append(show.Episodes, *episode)
    _, err := c.UpsertId(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    log.Printf("Added episode %s to show %s", episode.PlayId, show.Title)
    return true
}

func UpdateEpisode(showid bson.ObjectId, episode Episode) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show := GetShowById(showid)
    for _, e := range show.Episodes {
        if e.PlayId == episode.PlayId {
            e = episode
        }
    }
    _, err := c.UpsertId(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    log.Printf("Updated episode %s", episode.Title)
    return true
}
