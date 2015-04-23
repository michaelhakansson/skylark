//package tv3play
package main

import(
    "bytes"
    "encoding/json"
    // "encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
    // "regexp"
    "strconv"
    "strings"
    "time"
    "github.com/PuerkitoBio/goquery"
)

const(
    useragent string = "mozilla/5.0 (iphone; cpu iphone os 7_0_2 like mac os x) applewebkit/537.51.1 (khtml, like gecko) version/7.0 mobile/11a501 safari/9537.53"
    playService string = "tv3play"
    playUrlBase string = "http://www.tv3play.se/"
    videoUrlBase string = playUrlBase + "program/"
    // Should be followed by the episode id
    jsonVideoOutputString string = "http://playapi.mtgx.tv/v3/videos/"
    allProgramsPage string = playUrlBase + "program"
    rssUrl string = "/rss.xml"
    // The show URL should be followed by the show id after 'format='
    showJsonUrl string = "http://playapi.mtgx.tv/v1/sections?sections=videos.one,seasons.videolist&format="
    xmlDateLayout string = "Mon, 2 Jan 2006 15:04:05 MST"
    thumbnailSize string = "200x200"
)

// struct for show information
type Show struct {
    PlayId string
    PlayService string
    Title string
}

// struct for episode information
type Episode struct {
    Broadcasted time.Time   // TODO: Maybe changed to "air_at"
    Category string         // DONE
    Description string      // DONE
    EpisodeNumber int64     // DONE
    Length string           // DONE
    Live bool               // Can't find
    PlayId int64            // DONE
    Season int64            // DONE
    Thumbnail string        // DONE
    Title string            // DONE
    VideoUrl string         // TODO: If not HLS found --> use sharing.url URL
}

// structs for an episode's json output
type Program struct {
    Id int64 // Id of the video
    Format_title string // The title of the show (without s. and ep. number)
    Format_position Format_position
    Embedded Embedded `json:"_embedded"` // Used to get thumbnail url
    Format_categories []Format_categories
    Description string
    Duration int64
    Publish_at string // Used as the broadcast date
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

// GetAllProgramIds fetches from the provider all of the programs id's
// By parsing the "all program page" of the provider
// Returns an array of all the id's in the form of a string array
func GetAllProgramIds() (programs []string) {
    // Get all program links from the program list
    b := getPage(allProgramsPage)
    reader := bytes.NewReader(b)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    var links []string
    doc.Find(".list-section").Each(func(i int, section *goquery.Selection) {
        section.Find("a").Each(func(j int, show *goquery.Selection) {
            link, _ := show.Attr("href")
            links = append(links, link)
            //fmt.Println(links)
        })
    })
    
    // Fetch all program ids by visiting all links and extracting
    // the id from the shows.
    for _, link := range links {
        fmt.Println(link)
        d := getPage(link)
        reader = bytes.NewReader(d)
        doc, err = goquery.NewDocumentFromReader(reader)
        checkerr(err)
        id, _ := doc.Find("section").Attr("data-format-id")
        if len(id) > 0 {
            programs = append(programs, id)
        }
    }

    return
}

// // GetShow fetches the information and all the episodes for a show
// func GetShow(showId string) (Show, []Episode) {
//     url := showJsonUrl + showId
//     b := getPage(url)
//     if len(b) > 0 {
//         return parseShowXML(b, showId)
//     } else {
//         pageUrl := playUrlBase + showId
//         return parseShowPage(getPage(pageUrl), showId)
//     }
// }

// func parseShowXML(page []byte, showId string) (show Show, episodes []Episode) {
//     var c Channel
//     err := xml.Unmarshal(page, &c)
//     checkerr(err)
//     show.Title = strings.Replace(c.Title, " - Senaste program", "", 1)
//     show.PlayService = playService
//     show.PlayId = showId
//     r, err := regexp.Compile(`\/\d+\/`)
//     checkerr(err)
//     for _, item := range c.Item {
//         shortLink := r.FindString(item.Link)
//         episodeId := strings.Replace(string(shortLink), "/", "", 2)
//         e := GetEpisode(episodeId)
//         episodes = append(episodes, e)
//     }
//     return
// }

// func parseShowPage(page []byte, showId string) (show Show, episodes []Episode) {
//     var ids []string
//     reader := bytes.NewReader(page)
//     doc, err := goquery.NewDocumentFromReader(reader)
//     checkerr(err)
//     r, err := regexp.Compile(`\/\d+\/`)
//     checkerr(err)
//     show.Title = doc.Find(".play_title-page-info__header-title").Text()
//     show.PlayService = playService
//     show.PlayId = showId
//     doc.Find(".play_vertical-list").First().Find("li").Each(func(i int, s *goquery.Selection) {
//         link, _ := s.Find(".play_vertical-list__header-link").Attr("href")
//         digi := r.FindString(link)
//         cleanId := strings.Replace(string(digi), "/", "", 2)
//         ids = append(ids, cleanId)
//     })
//     for _, id := range ids {
//         e := GetEpisode(id)
//         episodes = append(episodes, e)
//     }
//     return
// }

// GetEpisode parses the information for an episode of a show
// Returns the episode information
func GetEpisode(episodeId string) (e Episode) {
    url := jsonVideoOutputString + episodeId
    b := getPage(url)
    var p Program
    err := json.Unmarshal(b, &p)
    checkerr(err)
    e.Broadcasted = parseDateTime(p.Publish_at)
    e.Category = p.Format_categories[0].Name
    e.Description = p.Description
    e.EpisodeNumber, err = strconv.ParseInt(p.Format_position.Episode, 0, 64)
    checkerr(err)
    e.Season = p.Format_position.Season
    e.Length = (time.Duration(p.Duration) * time.Second).String()
    // e.Live = p.Video.Live
    e.PlayId = p.Id
    e.Title = p.Format_title
    e.Thumbnail = fixThumbnailAdress(p.Embedded.Format.Links.Image.Href)
    return
}


// Replaces the size variable in the thumbnail url with the actually wanted size
func fixThumbnailAdress(url string) string {
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

func main() {
    // Test getting all program ids
    // programs := GetAllProgramIds()
    // for _, prog := range programs {
    //     fmt.Println(prog)
    // }
    
    // Test getting episode info
    fmt.Println(GetEpisode("556495"))
}