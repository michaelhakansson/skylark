package svtplay

import (
    "bytes"
    "encoding/json"
    "encoding/xml"
    "github.com/PuerkitoBio/goquery"
    "github.com/michaelhakansson/skylark/structures"
    "io/ioutil"
    "log"
    "net/http"
    "regexp"
    "strings"
    "time"
)

// SVTPlay service struct
type SVTPlay struct {
}

const (
    useragent             string = "mozilla/5.0 (iphone; cpu iphone os 7_0_2 like mac os x) applewebkit/537.51.1 (khtml, like gecko) version/7.0 mobile/11a501 safari/9537.53"
    playService           string = "svtplay"
    playURLBase           string = "http://www.svtplay.se/"
    videoURLBase          string = playURLBase + "video/"
    jsonVideoOutputString string = "?output=json"
    allProgramsPage       string = playURLBase + "program"
    rssURL                string = "/rss.xml"
)

// Program and other structs below is for an episode's json output
type Program struct {
    Context    Context
    Statistics Statistics
    VideoID    int64
    Video      Video
}

// Context see above for explanation
type Context struct {
    Title          string
    ProgramTitle   string
    ThumbnailImage string
}

// Statistics see above for explanation
type Statistics struct {
    BroadcastDate string
    BroadcastTime string
    Category      string
    Title         string
}

// Video see above for explanation
type Video struct {
    AvailableOnMobile bool
    Live              bool
    MaterialLength    int64
    Position          int64
    VideoReferences   []VideoReferences
}

// VideoReferences see above for explanation
type VideoReferences struct {
    Bitrate    int64
    PlayerType string
    URL        string
}

// Channel and item are structs for rss feed
type Channel struct {
    XMLName xml.Name `xml:"rss"`
    Title   string   `xml:"channel>title"`
    Item    []Item   `xml:"channel>item"`
}

// Item see above for explanation
type Item struct {
    Title       string `xml:"title"`
    Link        string `xml:"link"`
    Description string `xml:"description"`
    PubDate     string `xml:"pubDate"`
    GUID        int64  `xml:"guid"`
}

// GetName returns the name of the playservice
func (s SVTPlay) GetName() string {
    return playService
}

// GetAllProgramIDs fetches from the provider all of the programs id's
// By parsing the "all program page" of the provider
// Returns an array of all the id's in the form of a string array
func (s SVTPlay) GetAllProgramIDs() (ids []string) {
    page, status := getPage(allProgramsPage)
    if status != 200 {
        return
    }
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

// GetShow fetches the information and all the episode ids for a show
func (s SVTPlay) GetShow(showid string) (show structures.Show, episodes []string) {
    pageURL := playURLBase + showid
    xmlURL := playURLBase + showid + rssURL
    page, status := getPage(xmlURL)
    var episodeids []string
    if status == 200 {
        show, episodeids = parseShowXML(page, showid)
    } else {
        page, status = getPage(pageURL)
        if status == 200 {
            show, episodeids = parseShowPage(page, showid)
        }
    }
    page, status = getPage(pageURL)
    if status == 200 {
        show.Thumbnail = parseShowThumbnail(page)
    }

    for _, id := range episodeids {
        cleanid := strings.Replace(string(id), "/", "", 2)
        episodes = append(episodes, cleanid)
    }
    return
}

// parseShowXML parses the rss feed for a show
// Returns the show information and all the episodes
func parseShowXML(page []byte, showid string) (show structures.Show, episodeids []string) {
    var c Channel
    err := xml.Unmarshal(page, &c)
    checkerr(err)
    show.Title = strings.Replace(c.Title, " - Senaste program", "", 1)
    show.PlayService = playService
    show.PlayID = showid
    r, err := regexp.Compile(`\/\d+\/`)
    checkerr(err)
    for _, item := range c.Item {
        shortLink := r.FindString(item.Link)
        episodeids = append(episodeids, shortLink)
    }
    return
}

// parseShowPage parses the website for a show
// Returns the show information and all the episodes
func parseShowPage(page []byte, showid string) (show structures.Show, episodeids []string) {
    reader := bytes.NewReader(page)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    r, err := regexp.Compile(`\/\d+\/`)
    checkerr(err)
    show.Title = doc.Find(".play_title-page-info__header-title").Text()
    show.PlayService = playService
    show.PlayID = showid
    doc.Find(".play_vertical-list").First().Find("li").Each(func(i int, s *goquery.Selection) {
        link, _ := s.Find(".play_vertical-list__header-link").Attr("href")
        digi := r.FindString(link)
        if len(digi) > 0 {
            episodeids = append(episodeids, digi)
        }
    })
    return
}

func parseShowThumbnail(page []byte) (thumbnail string) {
    reader := bytes.NewReader(page)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    thumbnail, _ = doc.Find(".play_title-page-trailer__image").Attr("data-imagename")
    if strings.Contains(thumbnail, "public/images/default/play_default_998x561.jpg") {
        thumbnail = "http://www.svtplay.se" + thumbnail
    } else {
        thumbnail = "http:" + thumbnail
    }
    return
}

// GetEpisode parses the information for an episode of a show
// Returns the episode information
func (s SVTPlay) GetEpisode(episodeid string) (episode structures.Episode) {
    url := videoURLBase + episodeid + jsonVideoOutputString
    page, status := getPage(url)
    if status != 200 {
        return
    }
    url = videoURLBase + episodeid
    descriptionpage, statusdescription := getPage(url)
    if statusdescription != 200 {
        return
    }
    episode = parseEpisode(page, descriptionpage)
    return
}

// Parses the json episode data
func parseEpisode(page []byte, descriptionpage []byte) (episode structures.Episode) {
    program := parseJSON(page)
    episode = parseBasicEpisodeInformation(program)
    episode.Broadcasted = parseDateTime(program.Statistics.BroadcastDate, program.Statistics.BroadcastTime)
    episode.Description = parseDescription(descriptionpage)
    episode.Length = convertLengthToString(program.Video.MaterialLength)
    episode.Season, episode.EpisodeNumber = parseSeasonEpisodeNumbers(program.Statistics.Title)
    episode.VideoURL = getVideoURL(program.Video.VideoReferences)
    return
}

func parseJSON(page []byte) (program Program) {
    err := json.Unmarshal(page, &program)
    checkerr(err)
    return
}

// Parses the basic information (no conversions needs) from the Program object to Episode object
func parseBasicEpisodeInformation(program Program) (episode structures.Episode) {
    episode.Category = program.Statistics.Category
    episode.Live = program.Video.Live
    episode.PlayID = program.VideoID
    episode.Title = program.Context.Title
    episode.Thumbnail = program.Context.ThumbnailImage
    if strings.Contains(episode.Thumbnail, "public/images/default/play_default_998x561.jpg") {
        episode.Thumbnail = "http://www.svtplay.se" + episode.Thumbnail
    }
    return
}

// Gets the url to "ios-friendly" video
func getVideoURL(vrefs []VideoReferences) string {
    for _, vref := range vrefs {
        if vref.PlayerType == "ios" {
            return vref.URL
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
func parseDateTime(d string, t string) (datetime time.Time) {
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
    datetime, _ = time.Parse("2006 01 02 15:04", year+" "+month+" "+day+" "+hour+":"+minute)
    return
}

// parseSeasonEpisodeNumbers parses the season and episode number of an episode
// If the episode do not have these numbers, the season number is set to the date
// it was broadcast and the episode number is set to the time it was broadcasted
// Returns the numbers as strings
func parseSeasonEpisodeNumbers(seasonepisode string) (season string, episode string) {
    season = "0"
    episode = "0"
    seasonepisoderegex, err := regexp.Compile(`(song-([0-9]+)-)*avsnitt-([0-9]+)`)
    checkerr(err)
    found := seasonepisoderegex.FindStringSubmatch(seasonepisode)
    if len(found) > 0 {
        if found[2] != "" {
            season = found[2]
        }
        if found[3] != "" {
            episode = found[3]
        }
    }
    return
}

// getPage fetches the content from a specified url
func getPage(url string) ([]byte, int) {
    client := &http.Client{}
    req, err := http.NewRequest("GET", url, nil)
    checkerr(err)
    req.Header.Set("User-Agent", useragent)
    resp, err := client.Do(req)
    defer resp.Body.Close()
    checkerr(err)
    b, _ := ioutil.ReadAll(resp.Body)
    return b, resp.StatusCode
}

// checkerr checks if an error has occured and logs it if has.
func checkerr(err error) {
    if err != nil {
        log.Println(err)
    }
}
