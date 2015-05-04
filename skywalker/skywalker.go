package main

import(
    "log"
    "github.com/michaelhakansson/skylark/services/svtplay"
    "github.com/michaelhakansson/skylark/services/tv3play"
    "github.com/michaelhakansson/skylark/services/kanal5play"
    "github.com/michaelhakansson/skylark/db"
)

// Services []interface{} = [svtplay, tv3play, kanal5play]

func main() {

for i, _ := range []interface{}{svtplay,tv3play,kanal5play} {
    log.Println(i)
} 

// A. For playService : services
    // 1. Set freshness == 0 for all existing shows in the database

    // 2. Get all program ids
// showIds := GetAllProgramIds()

// 3. Add all shows to the database - set freshness == 1
    for _, e := range showIds {
        // playService.GetShow(id string)
        db.AddShow(title string, playid string, playservice string)
        // TODO: Make AddShow return the object
        db.AddEpisode(showid bson.ObjectId, broadcasted time.Time, category string,
            description string, episodenumber string, length string, live bool, playid int64,
            season string, thumbnail string, title string, videourl string)
    }

// B. The more often run loop for keeping often updated items fresher
    // 1. For show : db
    //   2. for episode : show
    //     3. db.UpdateEpisode(service.GetEpisode)

}
