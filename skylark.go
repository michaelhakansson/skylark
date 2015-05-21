package main

import(
    "html/template"
    "log"
    "net/http"
    "strconv"
    "time"
    "github.com/gorilla/mux"
    "github.com/michaelhakansson/skylark/db"
    "github.com/michaelhakansson/skylark/skywalker"
    "github.com/michaelhakansson/skylark/structures"
)

type Page struct {
    Title string
    Services map[string]int
    Show structures.Show
    Shows []structures.Show
    Episode structures.Episode
}

func main() {
    go func() {
        timer := time.Tick(24 * time.Hour)
        for _ = range timer {
            log.Println("Sync new shows")
            skywalker.SyncNew()
            log.Println("Syncing of new show completed")
        }
    }()

    go func() {
        timer := time.Tick(10 * time.Minute)
        for _ = range timer {
            log.Println("Sync outdated shows")
            ids := db.GetAllShowIds()
            for _, id := range ids {
                show := db.GetShowByPlayId(id)
                // Max time since update allowed
                maxTimeSinceUpdate := (24 / show.ChangeFrequency)
                if time.Now().Sub(show.LastUpdated).Hours() > maxTimeSinceUpdate {
                    skywalker.SyncShow(show.PlayId, show.PlayService)
                }
            }
            log.Println("Syncing of outdated shows completed")
        }
    }()

    r := mux.NewRouter()
    r.HandleFunc("/", HomeHandler)
    r.HandleFunc("/service/{id}", ServiceHandler)
    r.HandleFunc("/show/{id}", ShowHandler)
    r.HandleFunc("/video/{showid}/{playid}", VideoHandler)
    log.Println("Port: 4321")
    panic(http.ListenAndServe(":4321", r))
}

func HomeHandler(w http.ResponseWriter, r *http.Request) {
    services := db.GetAllServices()
    p := &Page{Title: "Home", Services: services}
    t, _ := template.ParseFiles("layouts/index.html")
    t.Execute(w, p)
}

func ServiceHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    service := vars["id"]
    shows := db.GetShowsByPlayService(service)
    var watchableShows []structures.Show
    for _, show := range shows {
        if len(show.Episodes) > 0 {
            watchableShows = append(watchableShows, show)
        }
    }
    p := &Page{Title: service, Shows: watchableShows}
    t, _ := template.ParseFiles("layouts/service.html")
    t.Execute(w, p)
}

func ShowHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    showId := vars["id"]
    show := db.GetShowByPlayId(showId)
    episodes := structures.SortEpisodesByDate(show.Episodes)
    show.Episodes = episodes
    p := &Page{Title: show.Title, Show: show}
    t := template.Must(template.New("show.html").Funcs(funcMap).ParseFiles("layouts/show.html"))
    t.Execute(w, p)
}

func VideoHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    showid := vars["showid"]
    episodeId, _ := strconv.ParseInt(vars["playid"], 10, 64)
    show := db.GetShowByPlayId(showid)
    episode := db.GetEpisodeByPlayId(showid, episodeId)
    p := &Page{Title: show.Title + " - " + episode.Title, Show: show, Episode: episode}
    t := template.Must(template.New("video.html").Funcs(funcMap).ParseFiles("layouts/video.html"))
    t.Execute(w, p)
}

var funcMap = template.FuncMap{
    "timeString": timeString,
    "trimText": trimText,
    "zeroPaddingString": zeroPaddingString,
    "zeroPadding": zeroPadding,
}

func timeString(t time.Time) string {
    return t.Format("2006-01-02 15:04")
}

func trimText(s string) string {
    var result string
    for i, c := range s {
        if i >= 78 {
            break
        } else {
            result += string(c)
        }
    }
    result += "..."
    return result
}

func zeroPaddingString(i string) string {
    number, err := strconv.ParseInt(string(i), 0, 64)
    if err != nil {
        return i
    }
    return zeroPadding(number)
}

func zeroPadding(i int64) string {
    if i < 10 {
        return "0" + strconv.FormatInt(i, 10)
    } else {
        return strconv.FormatInt(i, 10)
    }
}
