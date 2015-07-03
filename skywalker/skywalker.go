package skywalker

import (
    "github.com/michaelhakansson/skylark/db"
    "github.com/michaelhakansson/skylark/services"
    "github.com/michaelhakansson/skylark/structures"
    "github.com/michaelhakansson/skylark/utils"
    "log"
    "os"
    "time"
)

const (
    useragent string = "mozilla/5.0 (iphone; cpu iphone os 7_0_2 like mac os x) applewebkit/537.51.1 (khtml, like gecko) version/7.0 mobile/11a501 safari/9537.53"
)

// Start starts the background update schedule of Skywalker
func Start() {
    go func() {
        timer := time.Tick(24 * time.Hour)
        for _ = range timer {
            log.Println("Sync new shows")
            SyncNew()
            log.Println("Syncing of new shows was completed")
        }
    }()

    go func() {
        timer := time.Tick(10 * time.Minute)
        for _ = range timer {
            log.Println("Sync outdated shows")
            SyncOutdated()
            log.Println("Syncing of outdated shows was completed")
        }
    }()
}

// SyncNew synchronises all the new (previously not added to database) shows on
// all implemented services
func SyncNew() {
    var newShows []structures.Show
    serviceslist := services.GetServices()
    for _, service := range serviceslist {
        ids := services.GetIdsWithService(service)
        for _, id := range ids {
            added, show := db.AddShow(id, service)
            if added {
                newShows = append(newShows, show)
            }
        }
    }
    _, err := os.Stdout.Write([]byte("Updating new shows\n"))
    utils.Checkerr(err)
    processShows(newShows)
    updateChangeFrequencyForAll()
}

// SyncOutdated synchronises all outdated (as specified by the synchronisation
// policy) shows in the database
func SyncOutdated() {
    var showsToUpdate []structures.Show
    ids := db.GetAllShowIDs()
    for _, id := range ids {
        show := db.GetShowByPlayID(id)
        // Maximum allowed time since update
        maxTimeSinceUpdate := (24 / show.ChangeFrequency)
        if time.Now().Sub(show.LastUpdated).Hours() > maxTimeSinceUpdate {
            showsToUpdate = append(showsToUpdate, show)
        }
    }
    _, err := os.Stdout.Write([]byte("Updating outdated shows\n"))
    utils.Checkerr(err)
    processShows(showsToUpdate)
}
