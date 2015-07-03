package kanal5play

import (
    "bytes"
    "encoding/json"
    "github.com/PuerkitoBio/goquery"
    "github.com/michaelhakansson/skylark/structures"
    "io/ioutil"
    "log"
    "net/http"
    "strconv"
    "strings"
    "time"
)

// Kanal5Play service struct
type Kanal5Play struct {
}

const (
    useragent             string = "mozilla/5.0 (iphone; cpu iphone os 7_0_2 like mac os x) applewebkit/537.51.1 (khtml, like gecko) version/7.0 mobile/11a501 safari/9537.53"
    playService           string = "kanal5play"
    playURLBase           string = "http://www.kanal5play.se/"
    jsonVideoOutputString string = playURLBase + "api/getVideo?format=IPAD&videoId="
    allProgramsPage       string = playURLBase + "program"
)

// API describes the structure of the response from the json api call
type API struct {
    Description            string
    EpisodeNumber          int64
    ID                     int64
    IsLive                 bool
    Length                 int64
    PosterURL              string
    Premium                bool
    SeasonNumber           int64
    ShownOnTvDateTimestamp int64
    Streams                []Streams
    Title                  string
}

// Streams describes the structure for streams
type Streams struct {
    Format string
    Source string
}

// GetName returns the name of the playservice
func (k Kanal5Play) GetName() string {
    return playService
}

// GetAllProgramIDs fetches from the provider all of the programs id's
// By parsing the "all program page" of the provider
// Returns an array of all the id's in the form of a string array
func (k Kanal5Play) GetAllProgramIDs() (ids []string) {
    page := getPage(allProgramsPage)
    ids = parseAllProgramsPage(page)
    return
}

// parseAllProgramsPage parses all the programs that are available on the service
func parseAllProgramsPage(page []byte) (ids []string) {
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

// GetShow fetches the information and all the episode ids for a show
func (k Kanal5Play) GetShow(showid string) (show structures.Show, episodes []string) {
    page := getPage(allProgramsPage + "/" + showid)

    show, linkToSeasonsPage := parseShowInfo(page, showid)

    page = getPage(playURLBase + linkToSeasonsPage)
    seasonLinks := parseSeasonLinks(page)
    var episodeLinks []string
    for _, sLink := range seasonLinks {
        page = getPage(playURLBase + sLink)
        eLinks := parseEpisodeLinksOnSeasonPage(page)
        episodeLinks = append(episodeLinks, eLinks...)
    }

    for _, link := range episodeLinks {
        split := strings.Split(link, "/")
        cleanid := split[len(split)-1]
        episodes = append(episodes, cleanid)
    }

    return
}

// parseShowInfo parses the information about a show on the show page
func parseShowInfo(page []byte, showid string) (show structures.Show, linkToSeasonsPage string) {
    reader := bytes.NewReader(page)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    show.Title = doc.Find(".content-header h1").Text()
    show.PlayService = playService
    show.PlayID = showid
    show.Thumbnail, _ = doc.Find(".sbs-program-info-content img").Attr("src")
    linkToSeasonsPage, _ = doc.Find(".season .season-info a").First().Attr("href")
    return
}

// parseSeasonLinks parses all links to the available season of a show
func parseSeasonLinks(page []byte) (linksToSeasons []string) {
    reader := bytes.NewReader(page)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    doc.Find(".season-intro a").Each(func(i int, s *goquery.Selection) {
        if s.HasClass("paging-button") {
            season, _ := s.Attr("href")
            linksToSeasons = append(linksToSeasons, season)
        }
    })
    return
}

// parseEpisodeLinksOnSeasonPage parses all links to episodes that are available on the season page
func parseEpisodeLinksOnSeasonPage(page []byte) (episodeLinks []string) {
    reader := bytes.NewReader(page)
    doc, err := goquery.NewDocumentFromReader(reader)
    checkerr(err)
    doc.Find(".sbs-video-season-episode-teaser .right-column .title a").Each(func(i int, s *goquery.Selection) {
        episode, _ := s.Attr("href")
        episodeLinks = append(episodeLinks, episode)
    })
    return
}

// GetEpisode fetches the information for an episode of a show
// Returns the episode information
func (k Kanal5Play) GetEpisode(episodeid string) (episode structures.Episode) {
    url := jsonVideoOutputString + episodeid
    page := getPage(url)
    episode = parseEpisode(page)
    return
}

// parseEpisode parses the episode information provided by the api
func parseEpisode(page []byte) (episode structures.Episode) {
    var a API
    err := json.Unmarshal(page, &a)
    checkerr(err)
    episode.Broadcasted = time.Unix(a.ShownOnTvDateTimestamp/1000, 0)
    episode.Category = ""
    episode.Description = a.Description
    episode.EpisodeNumber = strconv.FormatInt(a.EpisodeNumber, 10)
    episode.Length = (time.Duration(a.Length/1000) * time.Second).String()
    episode.Live = a.IsLive
    episode.PlayID = a.ID
    episode.Season = strconv.FormatInt(a.SeasonNumber, 10)
    episode.Title = a.Title
    episode.Thumbnail = a.PosterURL
    for _, vStream := range a.Streams {
        if vStream.Format == "IPAD" {
            episode.VideoURL = vStream.Source
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
