package skylark

import (
    "github.com/gorilla/mux"
    "github.com/michaelhakansson/skylark/db"
    "github.com/michaelhakansson/skylark/structures"
    "github.com/michaelhakansson/skylark/utils"
    "html/template"
    "io"
    "net/http"
    "os"
    "strconv"
)

// Page describes the structure of a webpage
type Page struct {
    Title    string
    Services map[string]int
    Show     structures.Show
    Shows    []structures.Show
    Episode  structures.Episode
}

var funcMap = template.FuncMap{
    "timeString":          utils.TimeString,
    "prettifyServiceText": utils.PrettifyServiceText,
    "zeroPaddingString":   utils.ZeroPaddingString,
    "zeroPadding":         utils.ZeroPadding,
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    services := db.GetAllServices()
    p := &Page{Title: "Home", Services: services}
    path := utils.GetPath()
    t := template.Must(template.New("index.tmpl").Funcs(funcMap).ParseFiles(path+"/layouts/index.tmpl", path+"/layouts/header.tmpl"))
    err := t.Execute(w, p)
    utils.Checkerr(err)
}

func serviceHandler(w http.ResponseWriter, r *http.Request) {
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
    path := utils.GetPath()
    t := template.Must(template.New("service.tmpl").Funcs(funcMap).ParseFiles(path+"/layouts/service.tmpl", path+"/layouts/header.tmpl"))
    err := t.Execute(w, p)
    utils.Checkerr(err)
}

func showHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    showid := vars["id"]
    show := db.GetShowByPlayID(showid)
    episodes := structures.SortEpisodesByDate(show.Episodes)
    show.Episodes = episodes
    p := &Page{Title: show.Title, Show: show}
    path := utils.GetPath()
    t := template.Must(template.New("show.tmpl").Funcs(funcMap).ParseFiles(path+"/layouts/show.tmpl", path+"/layouts/header.tmpl"))
    err := t.Execute(w, p)
    utils.Checkerr(err)
}

func videoHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    showid := vars["showid"]
    episodeid, _ := strconv.ParseInt(vars["playid"], 10, 64)
    show := db.GetShowByPlayID(showid)
    episode := db.GetEpisodeByPlayID(showid, episodeid)
    p := &Page{Title: "Now playing: " + show.Title + " - " + episode.Title, Show: show, Episode: episode}
    path := utils.GetPath()
    t := template.Must(template.New("video.tmpl").Funcs(funcMap).ParseFiles(path+"/layouts/video.tmpl", path+"/layouts/header.tmpl"))
    err := t.Execute(w, p)
    utils.Checkerr(err)
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    image := vars["image"]
    file, _ := os.Open(utils.GetPath() + "/img/" + image)
    defer file.Close()
    io.Copy(w, file)
}
