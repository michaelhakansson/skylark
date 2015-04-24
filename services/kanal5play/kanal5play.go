package kanal5play

import(
    "bytes"
    "encoding/json"
    "log"
    "io/ioutil"
    "net/http"
    "strings"
    "time"
    "github.com/PuerkitoBio/goquery"
)

const(
    useragent string = "mozilla/5.0 (iphone; cpu iphone os 7_0_2 like mac os x) applewebkit/537.51.1 (khtml, like gecko) version/7.0 mobile/11a501 safari/9537.53"
    playService string = "kanal5play"
    playUrlBase string = "http://www.kanal5play.se/"
    videoUrlBase string = playUrlBase + "video/"
    jsonVideoOutputString string = playUrlBase + "api/getVideo?format=IPAD&videoId="
    allProgramsPage string = playUrlBase + "program"
    rssUrl string = "/rss.xml"
)

// struct for show information
type Show struct {
    PlayId string
    PlayService string
    Title string
}

// struct for episode information
type Episode struct {
    Broadcasted time.Time
    Category string
    Description string
    EpisodeNumber string
    Length string
    Live bool
    PlayId int64
    Season string
    Thumbnail string
    Title string
    VideoUrl string
}

// struct for json api
type Api struct {
    Description string
    EpisodeNumber int64
    Id int64
    IsLive bool
    Length int64
    PosterUrl string
    Premium bool
    SeasonNumber int64
    ShownOnTvDateTimestamp int64
    Streams []Streams
    Title string
}

type Streams struct {
    Format string
    Source string
}

// GetAllProgramIds fetches from the provider all of the programs id's
// By parsing the "all program page" of the provider
// Returns an array of all the id's in the form of a string array
func GetAllProgramIds() (ids []string) {
    page := getPage(allProgramsPage)
    reader := bytes.NewReader(page)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    doc.Find(".sbs-program-teaser").Each(func(i int, s *goquery.Selection) {
        link, _ := s.Find("a").Attr("href")
        link = strings.Replace(link, "/program/", "", 1)
        ids = append(ids, link)
    })
    return
}

// GetShow fetches the information and all the episodes for a show
func GetShow(showId string) (show Show, episodes []Episode) {
    page := getPage(allProgramsPage + "/" + showId)
    reader := bytes.NewReader(page)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    show.Title = doc.Find(".content-header h1").Text()
    show.PlayService = playService
    show.PlayId = showId
    linkToSeason, _ := doc.Find(".season .season-info a").First().Attr("href")
    var seasonLinks []string
    page = getPage(playUrlBase + linkToSeason)
    reader = bytes.NewReader(page)
    doc, err = goquery.NewDocumentFromReader(reader)
    doc.Find(".season-intro a").Each(func(i int, s *goquery.Selection) {
        if (s.HasClass("paging-button")) {
            season, _ := s.Attr("href")
            seasonLinks = append(seasonLinks, season)
        }
    })
    var episodeLinks []string
    for _, sLink := range seasonLinks {
        page = getPage(playUrlBase + sLink)
        reader = bytes.NewReader(page)
        doc, err = goquery.NewDocumentFromReader(reader)
        doc.Find(".sbs-video-season-episode-teaser .right-column .title a").Each(func(i int, s *goquery.Selection) {
            episode, _ := s.Attr("href")
            episodeLinks = append(episodeLinks, episode)
        })
    }
    for _, eLink := range episodeLinks {
        split := strings.Split(eLink, "/")
        cleanId := split[len(split) - 1]
        e := GetEpisode(cleanId)
        episodes = append(episodes, e)
    }
    return
}

// GetEpisode parses the information for an episode of a show
// Returns the episode information
func GetEpisode(episodeId string) (episode Episode) {
    url := jsonVideoOutputString + episodeId
    page := getPage(url)
    var a Api
    err := json.Unmarshal(page, &a)
    checkerr(err)
    episode.Broadcasted = time.Unix(a.ShownOnTvDateTimestamp/1000, 0)
    episode.Category = ""
    episode.Description = a.Description
    episode.EpisodeNumber = strconv.FormatInt(a.EpisodeNumber, 10)
    episode.Length = (time.Duration(a.Length/1000) * time.Second).String()
    episode.Live = a.IsLive
    episode.PlayId = a.Id
    episode.Season = strconv.FormatInt(a.SeasonNumber, 10)
    episode.Title = a.Title
    episode.Thumbnail = a.PosterUrl
    for _, vStream := range a.Streams {
        if vStream.Format == "IPAD" {
            episode.VideoUrl = vStream.Source
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