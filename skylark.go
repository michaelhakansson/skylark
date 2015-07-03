package skylark

import (
    "github.com/gorilla/mux"
    "github.com/michaelhakansson/skylark/skywalker"
    "log"
    "net/http"
)

// Start starts the skylark web server and skywalker scraper
func Start() {
    skywalker.Start()
    r := mux.NewRouter()
    r.HandleFunc("/", homeHandler)
    r.HandleFunc("/service/{id}", serviceHandler)
    r.HandleFunc("/show/{id}", showHandler)
    r.HandleFunc("/video/{showid}/{playid}", videoHandler)
    r.HandleFunc("/img/{image}", imageHandler)
    log.Println("Port: 4321")
    panic(http.ListenAndServe(":4321", r))
}
