package svtplay

import (
    "testing"
    "time"
)

type testshow struct {
    testurl   string
    id        string
    title     string
    thumbnail string
    number    int
}

type testepisode struct {
    testurl         string
    id              string
    broadcastedtime string
    category        string
    description     string
    episodenumber   string
    length          string
    live            bool
    playid          int64
    season          string
    thumbnail       string
    title           string
    videourl        string
}

type testlength struct {
    length int64
    text   string
}

type testdescription struct {
    testurl     string
    description string
}

type testdatetime struct {
    date     string
    time     string
    datetime time.Time
}

type testseasonepisode struct {
    seasonepisode string
    season        string
    episode       string
}

var showsXML = []testshow{
    {"http://www.svtplay.se/vetenskapens-varld/rss.xml", "vetenskapens-varld", "Vetenskapens värld", "http://www.svt.se/cachable_image/1402560543000/vetenskapens-varld/article2114658.svt/ALTERNATES/extralarge/vetenskapensvarld-victoria.jpg", 6},
    {"http://www.svtplay.se/uppdrag-granskning/rss.xml", "uppdrag-granskning", "Uppdrag granskning", "http://www.svt.se/cachable_image/1403794482000/ug/article2149590.svt/ALTERNATES/extralarge/uppdraggranskning.jpg", 23},
}

var showsPage = []testshow{
    {"http://www.svtplay.se/vetenskapens-varld", "vetenskapens-varld", "Vetenskapens värld", "http://www.svt.se/cachable_image/1402560543000/vetenskapens-varld/article2114658.svt/ALTERNATES/extralarge/vetenskapensvarld-victoria.jpg", 6},
    {"http://www.svtplay.se/uppdrag-granskning", "uppdrag-granskning", "Uppdrag granskning", "http://www.svt.se/cachable_image/1403794482000/ug/article2149590.svt/ALTERNATES/extralarge/uppdraggranskning.jpg", 15},
}

var episodes = []testepisode{
    {"http://www.svtplay.se/video/3056634?output=json", "3056634", "20150702", "kultur-och-nöje", "Del 16 av 20. Komikerna Kodjo Akolor, David Druid, Victor Linnèr, Sara Kinberg och Camilla Fågelborg snackar om aktuella och knasiga nyheter de hittar på nätet. Det är udda, konstigt och mycket mycket roligt! Kan ses till tis 29 dec (180 dagar kvar)", "16", "8m2s", false, 3056634, "0", "http://www.svt.se/cachable_image/1435307101000/svts/article3056633.svt/ALTERNATES/extralarge/default_title", "Avsnitt 16", "http://svtplay15m-f.akamaihd.net/i/world/open/20150626/1369824-016A/EPISOD-1369824-016A-3d784c0d4aacadbf_,892,144,252,360,540,1584,2700,.mp4.csmil/master.m3u8"},
    {"http://www.svtplay.se/video/3050595?output=json", "3050595", "20150702", "nyheter", "Kan ses till tor 9 jul (7 dagar kvar)", "0", "10m1s", false, 3050595, "0", "http://www.svt.se/cachable_image/1435854360000/svts/article3076351.svt/ALTERNATES/extralarge/default_title", "2/7 18.00", "http://svtplay2k-f.akamaihd.net/i/world/open/20150702/1368588-157A/EPISOD-1368588-157A-b45a0af0f4d3ae62_,892,144,252,360,540,1584,2700,.mp4.csmil/master.m3u8"},
}

var lengths = []testlength{
    {60, "1m0s"},
    {61, "1m1s"},
    {661, "11m1s"},
}

var datetimes = []testdatetime{
    {"20150503", "1200", time.Date(2015, time.May, 3, 12, 0, 0, 0, time.UTC)},
    {"19860130", "2112", time.Date(1986, time.January, 30, 21, 12, 0, 0, time.UTC)},
}

var seasonepisodes = []testseasonepisode{
    {"sasong-2-avsnitt-3-1", "2", "3"},
    {"3-5-12-00", "0", "0"},
}

var descriptions = []testdescription{
    {"http://www.svtplay.se/video/3013183", "Del 19 av 19. För ett år sedan lade TV4 ned sina lokala nyhetssändningar. Samtidigt lovar de att sända regionala program. Lever de upp till sitt löfte? Dessutom: Kommungranskarna återvänder till Örnsköldsvik. Programledare: Karin Mattisson. "},
    {"http://www.svtplay.se/video/55372", "När det är sovdags för Ernie och Bert tar deras sängar dem ut i världen på spännande äventyr. Ibland får Berts duva eller Ernies gummianka följa  med. Svenska röster: Magnus Ehrner och Steve Kratz."},
}

func TestProgramIds(t *testing.T) {
    page, _ := getPage("http://www.svtplay.se/program")
    ids := parseAllProgramsPage(page)
    if len(ids) != 532 {
        t.Error(
            "For", "parseAllProgramPage, length",
            "expected", 532,
            "got", len(ids),
        )
    }
    if ids[5] != "allsang-pa-skansen" {
        t.Error(
            "For", "parseAllProgramsPage, id",
            "expected", "allsang-pa-skansen",
            "got", ids[5],
        )
    }
}

func TestXMLParse(t *testing.T) {
    for _, pair := range showsXML {
        page, _ := getPage(pair.testurl)
        s, e := parseShowXML(page, pair.id)
        if s.Title != pair.title {
            t.Error(
                "For", pair.id,
                "expected", pair.title,
                "got", s.Title,
            )
        }
        if len(e) != pair.number {
            t.Error(
                "For", pair.id,
                "expected", pair.number,
                "got", len(e),
            )
        }
    }
}

func TestShowThumbnail(t *testing.T) {
    for _, pair := range showsPage {
        page, _ := getPage(pair.testurl)
        thumbnail := parseShowThumbnail(page)
        if thumbnail != pair.thumbnail {
            t.Error(
                "For", pair.id,
                "expected", pair.thumbnail,
                "got", thumbnail,
            )
        }
    }
}

func TestShowPageParse(t *testing.T) {
    for _, pair := range showsPage {
        page, _ := getPage(pair.testurl)
        s, e := parseShowPage(page, pair.id)
        if s.Title != pair.title {
            t.Error(
                "For", pair.id,
                "expected", pair.title,
                "got", s.Title,
            )
        }
        if len(e) != pair.number {
            t.Error(
                "For", pair.id,
                "expected", pair.number,
                "got", len(e),
            )
        }
    }
}

func TestParseJSON(t *testing.T) {
    for _, pair := range episodes {
        page, _ := getPage(pair.testurl)
        program := parseJSON(page)
        if program.VideoID != pair.playid {
            t.Error(
                "For", pair.id,
                "expected", pair.playid,
                "got", program.VideoID,
            )
        }
        if program.Video.Live != pair.live {
            t.Error(
                "For", pair.id,
                "expected", pair.live,
                "got", program.Video.Live,
            )
        }
    }
}

func TestParseBasicEpisode(t *testing.T) {
    for _, pair := range episodes {
        page, _ := getPage(pair.testurl)
        program := parseJSON(page)
        e := parseBasicEpisodeInformation(program)
        if e.Live != pair.live {
            t.Error(
                "For", pair.id,
                "expected", pair.live,
                "got", e.Live,
            )
        }
        if e.PlayID != pair.playid {
            t.Error(
                "For", pair.id,
                "expected", pair.playid,
                "got", e.PlayID,
            )
        }
        if e.Thumbnail != pair.thumbnail {
            t.Error(
                "For", pair.id,
                "expected", pair.thumbnail,
                "got", e.Thumbnail,
            )
        }
        if e.Title != pair.title {
            t.Error(
                "For", pair.id,
                "expected", pair.title,
                "got", e.Title,
            )
        }
    }
}

func TestGetVideoURL(t *testing.T) {
    for _, pair := range episodes {
        page, _ := getPage(pair.testurl)
        program := parseJSON(page)
        s := getVideoURL(program.Video.VideoReferences)
        if s != pair.videourl {
            t.Error(
                "For", pair.id,
                "expected", pair.videourl,
                "got", s,
            )
        }
    }
}

func TestConvertLengthToString(t *testing.T) {
    for _, pair := range lengths {
        s := convertLengthToString(pair.length)
        if s != pair.text {
            t.Error(
                "For", pair.length,
                "expected", pair.text,
                "got", s,
            )
        }
    }
}

func TestParseDescription(t *testing.T) {
    for _, pair := range descriptions {
        page, _ := getPage(pair.testurl)
        desc := parseDescription(page)
        if desc != pair.description {
            t.Error(
                "For", pair.testurl,
                "expected", pair.description,
                "got", desc,
            )
        }
    }
}

func TestParseDateTime(t *testing.T) {
    for _, pair := range datetimes {
        datetime := parseDateTime(pair.date, pair.time)
        if !datetime.Equal(pair.datetime) {
            t.Error(
                "For", "Date Time",
                "expected", pair.datetime.String(),
                "got", datetime.String(),
            )
        }
    }
}

func TestParseSeasonEpisodeNumbers(t *testing.T) {
    for _, pair := range seasonepisodes {
        s, e := parseSeasonEpisodeNumbers(pair.seasonepisode)
        if s != pair.season {
            t.Error(
                "For", "Season",
                "expected", pair.season,
                "got", s,
            )
        }
        if e != pair.episode {
            t.Error(
                "For", "Episode",
                "expected", pair.episode,
                "got", e,
            )
        }
    }
}
