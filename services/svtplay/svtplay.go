package svtplay

import(
    "bytes"
    "encoding/json"
    "encoding/xml"
    "log"
    "io/ioutil"
    "net/http"
    "regexp"
    "strings"
    "sync"
    "time"
    "github.com/michaelhakansson/skylark/structures"
    "github.com/PuerkitoBio/goquery"
)

const(
    useragent string = "mozilla/5.0 (iphone; cpu iphone os 7_0_2 like mac os x) applewebkit/537.51.1 (khtml, like gecko) version/7.0 mobile/11a501 safari/9537.53"
    playService string = "svtplay"
    playUrlBase string = "http://www.svtplay.se/"
    videoUrlBase string = playUrlBase + "video/"
    jsonVideoOutputString string = "?output=json&format=json"
    allProgramsPage string = playUrlBase + "program"
    rssUrl string = "/rss.xml"
)

// structs for an episode's json output
type Program struct {
    Context Context
    Statistics Statistics
    VideoId int64
    Video Video
}

type Context struct {
    Title string
    ProgramTitle string
    ThumbnailImage string
}

type Statistics struct {
    BroadcastDate string
    BroadcastTime string
    Category string
    Title string
}

type Video struct {
    AvailableOnMobile bool
    Live bool
    MaterialLength int64
    Position int64
    VideoReferences []VideoReferences
}

type VideoReferences struct {
    Bitrate int64
    PlayerType string
    Url string
}

// structs for rss feed
type Channel struct {
    XMLName xml.Name `xml:"rss"`
    Title string `xml:"channel>title"`
    Item []Item `xml:"channel>item"`
}

type Item struct {
    Title string `xml:"title"`
    Link string `xml:"link"`
    Description string `xml:"description"`
    PubDate string `xml:"pubDate"`
    Guid int64 `xml:"guid"`
}

// GetAllProgramIds fetches from the provider all of the programs id's
// By parsing the "all program page" of the provider
// Returns an array of all the id's in the form of a string array
func GetAllProgramIds() (ids []string) {
    page := getPage(allProgramsPage)
    ids = parseAllProgramsPage(page)
    return
}

// parseAllProgramsPage parses all the programs that are available on the service
func parseAllProgramsPage(page []byte) (ids []string) {
    reader := bytes.NewReader(page)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    doc.Find(".play_alphabetic-list ul li").Each(func(i int, s *goquery.Selection) {
        link, _ := s.Find("a").Attr("href")
        link = strings.Replace(link, "/", "", -1)
        ids = append(ids, link)
    })
    return
}

// GetShow fetches the information and all the episodes for a show
func GetShow(showId string) (show structures.Show, episodes []structures.Episode) {
    xmlUrl :=  playUrlBase + showId + rssUrl
    b := getPage(xmlUrl)
    var episodeIds []string
    if len(b) > 0 {
        show, episodeIds = parseShowXML(b, showId)
    } else {
        pageUrl := playUrlBase + showId
        show, episodeIds = parseShowPage(getPage(pageUrl), showId)
    }
    var wg sync.WaitGroup
    for _, id := range episodeIds {
        wg.Add(1)
        cleanId := strings.Replace(string(id), "/", "", 2)
        go func() {
            defer wg.Done()
            e := GetEpisode(cleanId)
            episodes = append(episodes, e)
        }()
    }
    wg.Wait()
    return
}

// parseShowXML parses the rss feed for a show
// Returns the show information and all the episodes
func parseShowXML(page []byte, showId string) (show structures.Show, episodeIds []string) {
    var c Channel
    err := xml.Unmarshal(page, &c)
    checkerr(err)
    show.Title = strings.Replace(c.Title, " - Senaste program", "", 1)
    show.PlayService = playService
    show.PlayId = showId
    r, err := regexp.Compile(`\/\d+\/`)
    checkerr(err)
    for _, item := range c.Item {
        shortLink := r.FindString(item.Link)
        episodeIds = append(episodeIds, shortLink)
    }
    return
}

// parseShowPage parses the website for a show
// Returns the show information and all the episodes
func parseShowPage(page []byte, showId string) (show structures.Show, episodeIds []string) {
    reader := bytes.NewReader(page)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    r, err := regexp.Compile(`\/\d+\/`)
    checkerr(err)
    show.Title = doc.Find(".play_title-page-info__header-title").Text()
    show.PlayService = playService
    show.PlayId = showId
    doc.Find(".play_vertical-list").First().Find("li").Each(func(i int, s *goquery.Selection) {
        link, _ := s.Find(".play_vertical-list__header-link").Attr("href")
        digi := r.FindString(link)
        if len(digi) > 0 {
            episodeIds = append(episodeIds, digi)
        }
    })
    return
}

// GetEpisode parses the information for an episode of a show
// Returns the episode information
func GetEpisode(episodeId string) (episode structures.Episode) {
    url := videoUrlBase + episodeId + jsonVideoOutputString
    page := getPage(url)
    url = videoUrlBase + episodeId
    descriptionpage := getPage(url)
    episode = parseEpisode(page, descriptionpage, episodeId)
    return
}

// Parses the json episode data
func parseEpisode(page []byte, descriptionpage []byte, episodeId string) (episode structures.Episode) {
    program := parseJSON(page)
    episode = parseBasicEpisodeInformation(program, episodeId)
    episode.Broadcasted = parseDateTime(program.Statistics.BroadcastDate, program.Statistics.BroadcastTime)
    episode.Description = parseDescription(descriptionpage)
    episode.Length = convertLengthToString(program.Video.MaterialLength)
    episode.Season, episode.EpisodeNumber = parseSeasonEpisodeNumbers(program.Statistics.Title)
    episode.VideoUrl = getVideoUrl(program.Video.VideoReferences)
    return
}

func parseJSON(page []byte) (program Program) {
    err := json.Unmarshal(page, &program)
    checkerr(err)
    return
}

// Parses the basic information (no conversions needs) from the Program object to Episode object
func parseBasicEpisodeInformation(program Program, episodeId string) (episode structures.Episode) {
    episode.Category = program.Statistics.Category
    episode.Live = program.Video.Live
    episode.PlayId = program.VideoId
    episode.Title = program.Context.Title
    episode.Thumbnail = program.Context.ThumbnailImage
    return
}

// Gets the url to "ios-friendly" video
func getVideoUrl(vrefs []VideoReferences) string {
    for _, vref := range vrefs {
        if vref.PlayerType == "ios" {
            return vref.Url
        }
    }
    return ""
}

// Converts the length value to "human-readable" string
func convertLengthToString(length int64) string {
    return (time.Duration(length) * time.Second).String()
}

// parseDescription parses the description for an episode
// Returns the description as a string
func parseDescription(page []byte) (description string) {
    reader := bytes.NewReader(page)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    description = doc.Find(".play_video-area-aside__info-text").First().Text()
    return
}

// parseDateTime parses the date and time for when an episode was broadcasted
// Returns the date and time as an time object
func parseDateTime(d string, t string) (datetime time.Time){
    if len(d) == 0 {
        d = "19840124"
    }
    if len(t) == 0 {
        t = "0800"
    }
    year := d[0:4]
    month := d[4:6]
    day := d[6:8]
    hour := t[0:2]
    minute := t[2:4]
    datetime, _ = time.Parse("2006 01 02 15:04", year + " " + month + " " + day + " " + hour + ":" + minute)
    return
}

// parseSeasonEpisodeNumbers parses the season and episode number of an episode
// If the episode do not have these numbers, the season number is set to the date
// it was broadcast and the episode number is set to the time it was broadcasted
// Returns the numbers as strings
func parseSeasonEpisodeNumbers(seasonepisode string) (season string, episode string) {
    numbers, _ := regexp.Compile(`([0-9])`)
    foundnumbers := numbers.MatchString(seasonepisode)
    if !foundnumbers {
        season = "0"
        episode = seasonepisode
        return
    }
    letters, _ := regexp.Compile(`^([a-z])`)
    foundLetters := letters.MatchString(seasonepisode)
    s := strings.Split(seasonepisode, "-")
    if foundLetters {
        if len(s) >= 4 {
            season = s[1]
            episode = s[3]
        } else {
            season = "0"
            episode = s[1]
        }
    } else {
        if len(s) >= 4 {
            season = s[0] + "/" + s[1]
            episode = s[2] + ":" + s[3]
        } else {
            season = "0"
            episode = s[0] + "/" + s[1]
        }
    }
    return
}

// getPage fetches the content from a specified url
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

// checkerr checks if an error has occured and logs it if has.
func checkerr(err error) {
    if err != nil {
        log.Println(err)
    }
}

