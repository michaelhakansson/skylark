package db

import (
    "github.com/michaelhakansson/skylark/structures"
    "gopkg.in/mgo.v2/bson"
    "log"
    "time"
)

const (
    uri string = "mongodb://localhost:27017/skylark"
    db  string = "skylark"
)

// AddShow adds a show to the database
func AddShow(playid string, playservice string) (result bool, show structures.Show) {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    count, err := c.Find(bson.M{"playid": playid}).Count()
    if err != nil || count > 0 {
        err = c.Find(bson.M{"playid": playid}).One(&show)
        //log.Printf("Show %s (%s) already exists", show.Title, show.PlayId)
        return false, show
    }
    show = structures.Show{ID: bson.NewObjectId(), ChangeFrequency: 1, LastUpdated: time.Now(), PlayID: playid, PlayService: playservice}

    err = c.Insert(show)
    if err != nil {
        log.Fatal(err)
        return false, show
    }
    //log.Printf("Added %s", show.PlayId)
    return true, show
}

// AddShowInfo adds information to show with PlayID = 'playid'
func AddShowInfo(title string, thumbnail string, playid string, playservice string) (result bool, show structures.Show) {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    err := c.Find(bson.M{"playid": playid}).One(&show)
    if err != nil {
        log.Fatal(err)
        return false, show
    }

    show.Title = title
    show.Thumbnail = thumbnail
    show.PlayID = playid
    show.PlayService = playservice
    show.LastUpdated = time.Now()
    var episodes []structures.Episode
    show.Episodes = episodes

    _, err = c.UpsertId(show.ID, show)
    if err != nil {
        log.Fatal(err)
        return false, show
    }
    //log.Printf("Added/updated show information to %s", show.Title)
    return true, show
}

// UpdateShow updates the LastUpdated field with the current time
func UpdateShow(showid bson.ObjectId) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show := GetShowByID(showid)
    show.LastUpdated = time.Now()
    _, err := c.UpsertId(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    //log.Printf("Updated show %s", show.Title)
    return true
}

// UpdateShowChangeFrequency updates the ChangeFrequence field for the
// specified show
func UpdateShowChangeFrequency(showid bson.ObjectId, changefrequency float64) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show := GetShowByID(showid)
    show.ChangeFrequency = changefrequency
    _, err := c.UpsertId(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    //log.Printf("Updated show %s change frequency", show.Title)
    return true
}

// GetShowByPlayID returns a show with specified PlayID
func GetShowByPlayID(playid string) structures.Show {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    var s structures.Show
    err := c.Find(bson.M{"playid": playid}).One(&s)
    if err != nil {
        log.Fatal(err)
    }
    return s
}

// GetShowByID returns a show with specified ObjectID
func GetShowByID(showid bson.ObjectId) structures.Show {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    var s structures.Show
    err := c.Find(bson.M{"_id": showid}).One(&s)
    if err != nil {
        log.Fatal(err)
    }
    return s
}

// GetShowsByPlayService returns a list of shows with specified PlayService
func GetShowsByPlayService(playservice string) (shows []structures.Show) {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    err := c.Find(bson.M{"playservice": playservice}).Sort("title").All(&shows)
    if err != nil {
        log.Fatal(err)
    }
    return
}

// GetAllShowIDs returns a list of PlayIDs for all the shows in the database
func GetAllShowIDs() (showIds []string) {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    var shows []structures.Show
    err := c.Find(bson.M{}).All(&shows)
    if err != nil {
        log.Fatal(err)
    }
    for _, show := range shows {
        showIds = append(showIds, show.PlayID)
    }
    return
}
