package db

import(
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "log"
)

// Connect establishes an connection to the database
func Connect() *mgo.Session {
    session, err := mgo.Dial(uri)
    session.SetSafe(&mgo.Safe{})
    if err != nil {
        log.Fatal(err)
    }
    return session
}

// GetAllServices returns all distinct values for the field 'playservice' among
// all shows in the database
func GetAllServices() map[string]int {
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
