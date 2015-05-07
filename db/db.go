package db

import(
    "log"
    "time"
    "github.com/michaelhakansson/skylark/structures"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

const(
    uri string = "mongodb://localhost:27017/skylark"
    db string = "skylark"
)

func Connect() *mgo.Session {
    session, err := mgo.Dial(uri)
    session.SetSafe(&mgo.Safe{})
    if err != nil {
        log.Fatal(err)
    }
    return session
}

func AddShow(title string, playid string, playservice string) (result bool, show structures.Show) {
    result = false
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    count, err := c.Find(bson.M{"playid": playid}).Count()
    if err != nil || count > 0 {
        err = c.Find(bson.M{"playid": playid}).One(&show)
        log.Printf("Show %s (%s) already exists", title, playid)
        return
    }
    var episodes []structures.Episode
    show = structures.Show{Id: bson.NewObjectId(), ChangeFrequence: 1, LastUpdated: time.Now(), Title: title, PlayId: playid,
    PlayService: playservice, Episodes: episodes}

    err = c.Insert(show)
    if err != nil {
        log.Fatal(err)
        return
    }
    log.Printf("Added %s", title)
    result = true
    return
}

func UpdateShow(showid bson.ObjectId) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show := GetShowById(showid)
    show.LastUpdated = time.Now()
    _, err := c.Upsert(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    log.Printf("Updated show %s", show.Title)
    return true
}

func UpdateShowWithData(show structures.Show) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show.LastUpdated = time.Now()
    _, err := c.UpsertId(show.Id, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    log.Printf("Updated show %s with new data", show.Title)
    return true
}

func GetShowByPlayId(playid string) structures.Show {
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

func GetShowById(showid bson.ObjectId) structures.Show {
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

func GetAllShowIds() (showIds []string) {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    var shows []structures.Show
    err := c.Find(bson.M{}).All(&shows)
    if err != nil {
        log.Fatal(err)
    }
    for _, show := range shows {
        showIds = append(showIds, show.PlayId)
    }
    return
}

func AddEpisode(showid bson.ObjectId, episode structures.Episode) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show := GetShowById(showid)
    for _, e := range show.Episodes {
        if e.PlayId == episode.PlayId {
            UpdateEpisode(showid, episode)
            return true
        }
    }
    show.Episodes = append(show.Episodes, episode)
    _, err := c.UpsertId(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    log.Printf("Added episode %d to show %s", episode.PlayId, show.Title)
    return true
}

func UpdateEpisode(showid bson.ObjectId, episode structures.Episode) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show := GetShowById(showid)
    for i, e := range show.Episodes {
        if e.PlayId == episode.PlayId {
            show.Episodes[i] = episode
        }
    }
    _, err := c.UpsertId(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    log.Printf("Updated episode %d", episode.PlayId)
    return true
}
