package tv3play

import (
    "encoding/json"
    "fmt"
    "github.com/michaelhakansson/skylark/structures"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// TV3Play service struct
type TV3Play struct {
}

const (
    useragent   string = "mozilla/5.0 (iphone; cpu iphone os 7_0_2 like mac os x) applewebkit/537.51.1 (khtml, like gecko) version/7.0 mobile/11a501 safari/9537.53"
    playService string = "tv3play"
    apiURLBase  string = "http://playapi.mtgx.tv/v3/"
    // Should be followed by the episode id
    jsonVideoOutputString string = apiURLBase + "videos/"
    // The show URL should be followed by the show id after 'format='
    jsonShowURL string = "http://playapi.mtgx.tv/v3/formats/"
    // The format= should be followed by the show id
    jsonSeasonsInShow string = apiURLBase + "seasons?format="
    // season= followed by season id
    jsonEpisodesInSeason string = apiURLBase + "videos?season="
    // Stream url followed by episode id
    jsonStreamURL        string = apiURLBase + "videos/stream/"
    thumbnailSize        string = "200x200"
    allProgramsMobileAPI string = "http://legacy.tv3play.se/mobileapi/format"
)

// Program and structs below are for an episode's json output
type Program struct {
    ID               int64              // Id of the video
    Title            string             // The title of the episode
    FormatPosition   FormatPosition     `json:"format_position"`   // Used for season and episode number
    FormatCategories []FormatCategories `json:"format_categories"` // Used for show category
    Embedded         Embedded           `json:"_embedded"`         // Used to get thumbnail url
    Summary          string
    Duration         int64
    Broadcasts       []Broadcasts // Used to get the air date
    PublishAt        string       `json:"publish_at"` // Used as backup if air date doesn't exist
    Sharing          Sharing      // Used as backup if HLS stream doesn't exist
}

// FormatPosition contains season and episode number
type FormatPosition struct {
    Season  int64
    Episode string
}

// FormatCategories contains the category of the show
type FormatCategories struct {
    Name string
}

// Embedded is used to get the thumbnail URL
type Embedded struct {
    Format Format
}

// Format see above for explanation
type Format struct {
    Links Links `json:"_links"`
}

// Links see above for explanation
type Links struct {
    Image Image
}

// Image see above for explanation
type Image struct {
    Href string
}

// Broadcasts see above for explanation
type Broadcasts struct {
    AirAt string `json:"air_at"`
}

// Sharing see above for explanation
type Sharing struct {
    URL string
}

// AllSeasons and structs below are for seasons
type AllSeasons struct {
    EmbeddedSeasons EmbeddedSeasons `json:"_embedded"`
}

// EmbeddedSeasons see above for explanation
type EmbeddedSeasons struct {
    Seasons []Seasons
}

// Seasons see above for explanation
type Seasons struct {
    ID int64
}

// AllEpisodes and structs below are for episodes
type AllEpisodes struct {
    EmbeddedEpisodes EmbeddedEpisodes `json:"_embedded"`
}

// EmbeddedEpisodes see above for explanation
type EmbeddedEpisodes struct {
    Videos []Videos
}

// Videos see above for explanation
type Videos struct {
    ID   int64
    Type string
}

// AllStreams and structs below are for streams
type AllStreams struct {
    Streams Streams
}

// Streams see above for explanation
type Streams struct {
    Hls string
}

// API and structs below are for API json response
type API struct {
    Alphabet []string
    Sections []Content
}

// Content see above for explanation
type Content struct {
    Title   string
    Formats []Formats
}

// Formats see above for explanation
type Formats struct {
    ID string
}

// GetName returns the name of the playservice
func (tv TV3Play) GetName() string {
    return playService
}

// GetAllProgramIDs fetches from the provider all of the programs id's
// By parsing the "all program page" of the provider
// Returns an array of all the id's in the form of a string array
func (tv TV3Play) GetAllProgramIDs() (ids []string) {
    b := getPage(allProgramsMobileAPI)
    var api API
    err := json.Unmarshal(b, &api)
    checkerr(err)
    for _, content := range api.Sections {
        for _, id := range content.Formats {
            ids = append(ids, id.ID)
        }
    }
    return
}

// GetShow fetches the information and all the episode ids for a show
func (tv TV3Play) GetShow(showid string) (show structures.Show, episodes []string) {
    // 1. Build show info using API call
    url := jsonShowURL + showid
    b := getPage(url)
    var sh struct {
        Title string
        Image string
    }
    err := json.Unmarshal(b, &sh)
    checkerr(err)
    show.Title = sh.Title
    show.PlayID = showid
    show.PlayService = playService
    show.Thumbnail = sh.Image

    // 2. Fetch all seasons and id's via another API call
    var s AllSeasons
    url = jsonSeasonsInShow + showid
    b = getPage(url)
    err = json.Unmarshal(b, &s)
    checkerr(err)

    var seasonids []string
    for _, season := range s.EmbeddedSeasons.Seasons {
        seasonids = append(seasonids, strconv.FormatInt(season.ID, 10))
    }

    // 3. Fetch all episodes in season for all seasons via a third API call
    var allepisodes AllEpisodes
    for _, seasonid := range seasonids {
        url = jsonEpisodesInSeason + seasonid
        b = getPage(url)
        err = json.Unmarshal(b, &allepisodes)
        checkerr(err)

        // Get all the episode ids
        for _, episode := range allepisodes.EmbeddedEpisodes.Videos {
            cleanid := strconv.FormatInt(episode.ID, 10)
            if episode.Type != "clip" {
                episodes = append(episodes, cleanid)
            }
        }
    }
    return
}

// GetEpisode parses the information for an episode of a show
// Returns the episode information
func (tv TV3Play) GetEpisode(episodeid string) (e structures.Episode) {
    url := jsonVideoOutputString + episodeid
    b := getPage(url)
    var p Program
    err := json.Unmarshal(b, &p)
    checkerr(err)

    if len(p.Broadcasts) > 0 { // If broadcast date exist, use it
        e.Broadcasted = parseDateTime(p.Broadcasts[0].AirAt)
    } else { // Else use the upload date
        e.Broadcasted = parseDateTime(p.PublishAt)
    }
    e.Category = p.FormatCategories[0].Name
    e.Description = p.Summary
    epNumber := p.FormatPosition.Episode
    if len(epNumber) > 0 { // If episode number exists, use it
        e.EpisodeNumber = epNumber
        checkerr(err)
    } else { // Else set to 0
        e.EpisodeNumber = "0"
    }
    e.Season = strconv.FormatInt(p.FormatPosition.Season, 10)
    e.Length = (time.Duration(p.Duration) * time.Second).String()
    e.Live = false // Always set to false, since TV3Play has no live streams
    e.PlayID = p.ID
    e.Title = p.Title
    e.Thumbnail = fixThumbnailURL(p.Embedded.Format.Links.Image.Href)

    // Try to get HLS stream from stream link
    url = jsonStreamURL + episodeid
    b = getPage(url)
    var s AllStreams
    err = json.Unmarshal(b, &s)
    checkerr(err)

    if len(s.Streams.Hls) > 0 { // If HLS stream exists, use it
        e.VideoURL = fixHlsURL(s.Streams.Hls)
    } else { // Else use the sharing url to link to the provider page
        e.VideoURL = p.Sharing.URL
    }
    return
}

func fixHlsURL(url string) (fixedURL string) {
    parts := strings.Split(url, ",")
    if len(parts) > 1 {
        fixedURL = parts[0] + "," + parts[len(parts)-2] + "," + parts[len(parts)-1]
    } else {
        fixedURL = url
    }
    return
}

// Replaces the size variable in the thumbnail url with the actually wanted size
func fixThumbnailURL(url string) string {
    return strings.Replace(url, "{size}", thumbnailSize, 1)
}

// parseDateTime parses the date and time for when an episode was broadcasted
// Returns the date and time as an time object
func parseDateTime(d string) (datetime time.Time) {
    year := d[0:4]
    month := d[5:7]
    day := d[8:10]
    hour := d[11:13]
    minute := d[14:16]
    datetime, _ = time.Parse("2006 01 02 15:04", year+" "+month+" "+day+" "+hour+":"+minute)
    return
}

func getPage(url string) []byte {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    checkerr(err)
    req.Header.Set("User-Agent", useragent)
    resp, err := client.Do(req)
    defer resp.Body.Close()
    checkerr(err)
    b, _ := ioutil.ReadAll(resp.Body)
    return b
}

func checkerr(err error) {
    if err != nil {
        fmt.Println(err)
    }
}
