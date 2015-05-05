package tv3play

import(
    //    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "time"
    "github.com/michaelhakansson/skylark/structures"
    //    "github.com/PuerkitoBio/goquery"
)

const(
    useragent string = "mozilla/5.0 (iphone; cpu iphone os 7_0_2 like mac os x) applewebkit/537.51.1 (khtml, like gecko) version/7.0 mobile/11a501 safari/9537.53"
    playService string = "tv3play"
    playUrlBase string = "http://www.tv3play.se/"
    apiUrlBase string = "http://playapi.mtgx.tv/v3/"
    videoUrlBase string = playUrlBase + "program/"
    allProgramsPage string = playUrlBase + "program"
    // Should be followed by the episode id
    jsonVideoOutputString string = apiUrlBase + "videos/"
    // The show URL should be followed by the show id after 'format='
    jsonShowUrl string = "http://playapi.mtgx.tv/v3/formats/"
    // The format= should be followed by the show id
    jsonSeasonsInShow string = apiUrlBase + "seasons?format="
    // season= followed by season id
    jsonEpisodesInSeason string = apiUrlBase + "videos?season="
    // Stream url followed by episode id
    jsonStreamUrl string = apiUrlBase + "videos/stream/"
    xmlDateLayout string = "Mon, 2 Jan 2006 15:04:05 MST"
    thumbnailSize string = "200x200"
    allProgramsMobileApi string = "http://legacy.tv3play.se/mobileapi/format"
)

// structs for an episode's json output
type Program struct {
    Id int64 // Id of the video
    Title string // The title of the episode
    Format_position Format_position // Used for season and episode number
    Format_categories []Format_categories // Used for show category
    Embedded Embedded `json:"_embedded"` // Used to get thumbnail url
    Summary string
    Duration int64
    Broadcasts []Broadcasts // Used to get the air date
    Publish_at string // Used as backup if air date doesn't exist
    Sharing Sharing // Used as backup if HLS stream doesn't exist
}

//Contains season and episode number
type Format_position struct {
    Season int64
    Episode string
}

// Contains the category of the show
type Format_categories struct {
    Name string
}

// Used to get the thumbnail URL
type Embedded struct {
    Format Format
}

type Format struct {
    Links Links `json:"_links"`
}

type Links struct {
    Image Image
}

type Image struct {
    Href string
}

type Broadcasts struct {
    Air_at string
}

type Sharing struct {
    Url string
}

// Structs for seasons
type AllSeasons struct {
    EmbeddedSeasons EmbeddedSeasons `json:"_embedded"`
}

type EmbeddedSeasons struct {
    Seasons []Seasons
}

type Seasons struct {
    Id int64
}

// Structs for episodes
type AllEpisodes struct {
    EmbeddedEpisodes EmbeddedEpisodes `json:"_embedded"`
}

type EmbeddedEpisodes struct {
    Videos []Videos
}

type Videos struct {
    Id int64
}

// Structs for streams
type AllStreams struct {
    Streams Streams
}

type Streams struct {
    Hls string
}


type Api struct {
    Alphabet []string
    Sections []Content
}

type Content struct {
    Title string
    Formats []Formats
}

type Formats struct {
    Id string
}

// GetAllProgramIds fetches from the provider all of the programs id's
// By parsing the "all program page" of the provider
// Returns an array of all the id's in the form of a string array
func GetAllProgramIds() (programs []string, playservice string) {
    b := getPage(allProgramsMobileApi)
    var api Api
    err := json.Unmarshal(b, &api)
    checkerr(err)
    for _, content := range api.Sections {
        for _, id := range content.Formats {
            programs = append(programs, id.Id)
        }
    }

    playservice = playService

    /*// Get all program links from the program list
    b := getPage(allProgramsPage)
    reader := bytes.NewReader(b)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    var links []string
    doc.Find(".list-section").Each(func(i int, section *goquery.Selection) {
        section.Find("a").Each(func(j int, show *goquery.Selection) {
            link, _ := show.Attr("href")
            links = append(links, link)
        })
    })

    // Fetch all program ids by visiting all links and extracting
    // the id from the shows.
    for _, link := range links {
        d := getPage(link)
        reader = bytes.NewReader(d)
        doc, err = goquery.NewDocumentFromReader(reader)
        checkerr(err)
        id, _ := doc.Find("section").Attr("data-format-id")
        if len(id) > 0 {
            programs = append(programs, id)
        }
    }*/

    return
}

// GetShow fetches the information and all the episodes for a show
func GetShow(showId string) (show structures.Show, episodes []structures.Episode) {
    // 1. Build show info using API call
    url := jsonShowUrl + showId
    b := getPage(url)
    err := json.Unmarshal(b, &show)
    checkerr(err)
    show.PlayId = showId
    show.PlayService = playService

    // 2. Fetch all seasons and id's via another API call
    var s AllSeasons
    url = jsonSeasonsInShow + showId
    b = getPage(url)
    err = json.Unmarshal(b, &s)
    checkerr(err)

    var seasonIds []string
    for _, season := range s.EmbeddedSeasons.Seasons {
        seasonIds = append(seasonIds, strconv.FormatInt(season.Id, 10))
    }

    // 3. Fetch all episodes in season for all seasons via a third API call
    var allEpisodes AllEpisodes
    for _, seasonId := range seasonIds {
        url = jsonEpisodesInSeason + seasonId
        b = getPage(url)
        err = json.Unmarshal(b, &allEpisodes)
        checkerr(err)

        // Populate episodes array via GetEpisode call for each episode
        for _, episode := range allEpisodes.EmbeddedEpisodes.Videos {
            episodes = append(episodes, GetEpisode(strconv.FormatInt(episode.Id, 10)))
        }
    }
    return
}

// GetEpisode parses the information for an episode of a show
// Returns the episode information
func GetEpisode(episodeId string) (e structures.Episode) {
    url := jsonVideoOutputString + episodeId
    b := getPage(url)
    var p Program
    err := json.Unmarshal(b, &p)
    checkerr(err)

    if len(p.Broadcasts) > 0 { // If broadcast date exist, use it
        e.Broadcasted = parseDateTime(p.Broadcasts[0].Air_at)
    } else { // Else use the upload date
        e.Broadcasted = parseDateTime(p.Publish_at)
    }
    e.Category = p.Format_categories[0].Name
    e.Description = p.Summary
    epNumber := p.Format_position.Episode
    if len(epNumber) > 0 { // If episode number exists, use it
        e.EpisodeNumber = epNumber
        checkerr(err)
    } else { // Else set to 0
        e.EpisodeNumber = "0"
    }
    e.Season = strconv.FormatInt(p.Format_position.Season, 10)
    e.Length = (time.Duration(p.Duration) * time.Second).String()
    e.Live = false // Always set to false, since TV3Play has no live streams
    e.PlayId = p.Id
    e.Title = p.Title
    e.Thumbnail = fixThumbnailUrl(p.Embedded.Format.Links.Image.Href)


    // Try to get HLS stream from stream link
    url = jsonStreamUrl + episodeId
    b = getPage(url)
    var s AllStreams
    err = json.Unmarshal(b, &s)
    checkerr(err)

    if len(s.Streams.Hls) > 0 { // If HLS stream exists, use it
        e.VideoUrl = fixHlsUrl(s.Streams.Hls)
    } else { // Else use the sharing url to link to the provider page
        e.VideoUrl = p.Sharing.Url
    }
    return
}

func fixHlsUrl(url string) (fixedUrl string) {
    parts := strings.Split(url, ",")
    fixedUrl = parts[0] + "," + parts[len(parts)-2] + "," + parts[len(parts)-1]
    return
}


// Replaces the size variable in the thumbnail url with the actually wanted size
func fixThumbnailUrl(url string) string {
    return strings.Replace(url, "{size}", thumbnailSize, 1)
}

// parseDateTime parses the date and time for when an episode was broadcasted
// Returns the date and time as an time object
func parseDateTime(d string) (datetime time.Time){
    year := d[0:4]
    month := d[5:7]
    day := d[8:10]
    hour := d[11:13]
    minute := d[14:16]
    datetime, _ = time.Parse("2006 01 02 15:04", year + " " + month + " " + day + " " + hour + ":" + minute)
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
