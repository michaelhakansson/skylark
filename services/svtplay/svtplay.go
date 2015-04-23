package svtplay

import(
    "bytes"
    "encoding/json"
    "encoding/xml"
    "fmt"
    "io/ioutil"
    "net/http"
    "regexp"
    "strings"
    "time"
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
    xmlDateLayout string = "Mon, 2 Jan 2006 15:04:05 MST"
)

type Show struct {
    PlayId string
    PlayService string
    Title string
}

type Episode struct {
    Broadcasted time.Time
    Category string
    Description string
    EpisodeNumber int64
    Length string
    Live bool
    PlayId int64
    Season int64
    Thumbnail string
    Title string
    VideoUrl string
}

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

/*func main() {
    /*url := "http://www.svtplay.se/video/2843361?output=json&format=json"
    p := GetEpisode(url)
    fmt.Println(p)
    showId := "uppdrag-granskning"
    s, e := GetShow(showId)
    fmt.Println(s)
    fmt.Println(e)
    fmt.Println(len(e))

    /*programs := GetAllPrograms()
    for _, pro := range programs {
        fmt.Println(pro)
    }
}*/

func GetAllPrograms() (programs []string) {
    b := getPage(allProgramsPage)
    reader := bytes.NewReader(b)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    doc.Find(".play_alphabetic-list ul li").Each(func(i int, s *goquery.Selection) {
        link, _ := s.Find("a").Attr("href")
        link = strings.Replace(link, "/", "", -1)
        programs = append(programs, link)
    })
    return
}

func GetShow(showId string) (Show, []Episode) {
    xmlUrl :=  playUrlBase + showId + rssUrl
    b := getPage(xmlUrl)
    if len(b) > 0 {
        return parseShowXML(b, showId)
    } else {
        pageUrl := playUrlBase + showId
        return parseShowPage(getPage(pageUrl), showId)
    }
}
func parseShowXML(page []byte, showId string) (show Show, episodes []Episode) {
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
        episodeId := strings.Replace(string(shortLink), "/", "", 2)
        e := GetEpisode(episodeId)
        episodes = append(episodes, e)
    }
    return
}

func parseShowPage(page []byte, showId string) (show Show, episodes []Episode) {
    var ids []string
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
        cleanId := strings.Replace(string(digi), "/", "", 2)
        ids = append(ids, cleanId)
    })
    for _, id := range ids {
        e := GetEpisode(id)
        episodes = append(episodes, e)
    }
    return
}

func GetEpisode(episodeId string) (e Episode) {
    url := videoUrlBase + episodeId + jsonVideoOutputString
    b := getPage(url)
    var p Program
    err := json.Unmarshal(b, &p)
    checkerr(err)
    e.Broadcasted = time.Now()
    e.Category = p.Statistics.Category
    e.Description = ""
    e.EpisodeNumber = 0
    e.Length = (time.Duration(p.Video.MaterialLength) * time.Second).String()
    e.Live = p.Video.Live
    e.PlayId = p.VideoId
    e.Season = 0
    e.Title = p.Context.Title
    e.Thumbnail = p.Context.ThumbnailImage
    for _, vref := range p.Video.VideoReferences {
        if vref.PlayerType == "ios" {
            e.VideoUrl = vref.Url
        }
    }
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

