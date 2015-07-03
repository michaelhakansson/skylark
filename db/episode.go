package db

import(
    "github.com/michaelhakansson/skylark/structures"
    "gopkg.in/mgo.v2/bson"
    "log"
)

// AddEpisode adds an episode to a specific show
func AddEpisode(showid bson.ObjectId, episode structures.Episode) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show := GetShowByID(showid)
    show.Episodes = append(show.Episodes, episode)
    _, err := c.UpsertId(showid, show)
    if err != nil {
        log.Fatal(err)
        return false
    }
    //log.Printf("Added episode %d to show %s", episode.PlayId, show.Title)
    return true
}

// UpdateEpisode updates the episode information for an episode of a show
func UpdateEpisode(showid bson.ObjectId, episode structures.Episode) bool {
    session := Connect()
    defer session.Close()
    c := session.DB(db).C("shows")
    show := GetShowByID(showid)
    for i, e := range show.Episodes {
        if string(e.PlayID) == string(episode.PlayID) {
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

// GetEpisodeByPlayID returns an episode with specified PlayID
func GetEpisodeByPlayID(showid string, playid int64) structures.Episode {
    show := GetShowByPlayID(showid)
    var episode structures.Episode
    for _, e := range show.Episodes {
        if string(e.PlayID) == string(playid) {
            episode = e
        }
    }
    return episode
}

// ContainsEpisodeWithPlayID returns whether a show has an episode with the
// specified PlayID
func ContainsEpisodeWithPlayID(showid string, playid int64) bool {
    show := GetShowByPlayID(showid)
    for _, e := range show.Episodes {
        if string(e.PlayID) == string(playid) {
            return true
        }
    }
    return false
}
