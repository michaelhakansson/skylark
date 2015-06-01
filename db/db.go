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

func GetAllServices() (map[string]int) {
    services := make(map[string]int)
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    var result []string
    err := c.Find(nil).Distinct("playservice", &result)
    if err != nil {
        log.Fatal(err)
    }
    for _, service := range result {
        count, err := c.Find(bson.M{"playservice": service}).Count()
        if err != nil {
            log.Fatal(err)
        }
        services[service] = count
    }
    return services
}

func AddShow(playid string, playservice string) (result bool, show structures.Show) {
    result = false
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    count, err := c.Find(bson.M{"playid": playid}).Count()
    if err != nil || count > 0 {
        err = c.Find(bson.M{"playid": playid}).One(&show)
        //log.Printf("Show %s (%s) already exists", show.Title, show.PlayId)
        return
    }
    show = structures.Show{Id: bson.NewObjectId(), ChangeFrequency: 1, LastUpdated: time.Now(), PlayId: playid, PlayService: playservice}

    err = c.Insert(show)
    if err != nil {
        log.Fatal(err)
        return
    }
    //log.Printf("Added %s", show.PlayId)
    result = true
    return
}

func AddShowInfo(title string, thumbnail string, playid string, playservice string) (result bool, show structures.Show) {
    result = false;
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    err := c.Find(bson.M{"playid": playid}).One(&show)
    if err != nil {
        log.Fatal(err)
        return
    }

    show.Title = title
    show.Thumbnail = thumbnail
    show.PlayId = playid
    show.PlayService = playservice
    show.LastUpdated = time.Now()
    var episodes []structures.Episode
    show.Episodes = episodes

    _, err = c.UpsertId(show.Id, show)
    if err != nil {
        log.Fatal(err)
        return
    }
    result = true
    //log.Printf("Added/updated show information to %s", show.Title)
    return
}

func UpdateShow(showid bson.ObjectId) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show := GetShowById(showid)
    show.LastUpdated = time.Now()
    _, err := c.UpsertId(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    //log.Printf("Updated show %s", show.Title)
    return true
}

func UpdateShowChangeFrequency(showid bson.ObjectId, changefrequency float64) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show := GetShowById(showid)
    show.ChangeFrequency = changefrequency
    _, err := c.UpsertId(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    //log.Printf("Updated show %s change frequency", show.Title)
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
    show.Episodes = append(show.Episodes, episode)
    _, err := c.UpsertId(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    //log.Printf("Added episode %d to show %s", episode.PlayId, show.Title)
    return true
}

func UpdateEpisode(showid bson.ObjectId, episode structures.Episode) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show := GetShowById(showid)
    for i, e := range show.Episodes {
        if string(e.PlayId) == string(episode.PlayId) {
            show.Episodes[i] = episode
        }
    }
    _, err := c.UpsertId(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    //log.Printf("Updated episode %d", episode.PlayId)
    return true
}

func GetEpisodeByPlayId(showid string, playid int64) structures.Episode {
    show := GetShowByPlayId(showid)
    var episode structures.Episode
    for _, e := range show.Episodes {
        if string(e.PlayId) == string(playid) {
            episode = e
        }
    }
    return episode
}

func ContainsEpisodeWithPlayId(showid string, playid int64) bool {
    show := GetShowByPlayId(showid)
    for _, e := range show.Episodes {
        if string(e.PlayId) == string(playid) {
            return true
        }
    }
    return false
}
