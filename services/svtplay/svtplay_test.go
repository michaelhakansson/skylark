package svtplay

import(
    "io/ioutil"
    "log"
    "testing"
    "time"
)


type testshow struct {
    testfile string
    id string
    title string
    number int
}

type testepisode struct {
    testfile string
    id string
    broadcastedtime string
    category string
    description string
    episodenumber string
    length string
    live bool
    playid int64
    season string
    thumbnail string
    title string
    videourl string
}

type testlength struct {
    length int64
    text string
}

type testdescription struct {
    testfile string
    description string
}

type testdatetime struct {
    date string
    time string
    datetime time.Time
}

type testseasonepisode struct {
    seasonepisode string
    season string
    episode string
}

var showsXML = []testshow{
    {"show1rss.xml", "vetenskapens-varld", "Vetenskapens värld", 18},
    {"show2rss.xml", "uppdrag-granskning", "Uppdrag granskning", 16},
}

var showsPage = []testshow{
    {"show1page.html", "vetenskapens-varld", "Vetenskapens värld", 15},
    {"show2page.html", "uppdrag-granskning", "Uppdrag granskning", 15},
}

var episodes = []testepisode{
    {"episode1.json", "2843612", "", "kultur-och-nöje", "Del 4 av 10. Gruppfinal, med bland annat Sisyfos-tävlingen - vätskefyllda pilatesbollar som ska rullas uppför en backe. Några av Sveriges främsta idrottsmän och idrottskvinnor möts i fysiska och psykiska utmaningar för att kora Mästarnas mästare 2015. I Grupp 1 ingår Anette Norberg, Anna Olsson, Magnus Muhrén, Danijela Rundqvist, Glenn Hysén och Björn Lind. Programledare: Micke Leijnegard.",
    "4", "58m37s", false, 2843612, "7",
    "http://www.svt.se/cachable_image/1429226401000/svts/article2849370.svt/ALTERNATES/extralarge/default_title",
    "Avsnitt 4", "http://svtplay18p-f.akamaihd.net/i/se/open/20150417/1360782-004A/EPISOD-1360782-004A-2ae7758f8108a631_,892,144,252,360,540,1584,2700,.mp4.csmil/master.m3u8?cc1=name=Svenska~default=yes~forced=no~uri=http://media.svt.se/download/mcc/wp3/undertexter-wsrt/1360782/1360782-004A/C(sv)/index.m3u8~lang=sv"},
    {"episode2.json", "2867878", "", "nyheter", "Kan ses till imorgon 23.59 (1 dag kvar)", "11:00", "1m30s", false, 2867878, "23/4", "http://www.svt.se/cachable_image/1429781701000/svts/article2867877.svt/ALTERNATES/extralarge/default_title",
    "23/4 11.00", "http://svtplay19i-f.akamaihd.net/i/world/open/20150423/1368669-074A/EPISOD-1368669-074A-208212fd95c96099_,892,144,252,360,540,1584,2700,.mp4.csmil/master.m3u8"},
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
    {"3-5-12-00", "3/5", "12:00"},
}

var descriptions = []testdescription{
    {"episode1description.html", "Del 4 av 10. Gruppfinal, med bland annat Sisyfos-tävlingen - vätskefyllda pilatesbollar som ska rullas uppför en backe. Några av Sveriges främsta idrottsmän och idrottskvinnor möts i fysiska och psykiska utmaningar för att kora Mästarnas mästare 2015. I Grupp 1 ingår Anette Norberg, Anna Olsson, Magnus Muhrén, Danijela Rundqvist, Glenn Hysén och Björn Lind. Programledare: Micke Leijnegard."},
    {"episode2description.html", "Kan ses till ikväll 23.59 (12 timmar kvar)"},
}

func TestProgramIds(t *testing.T) {
    page, err := ioutil.ReadFile("testFiles/allprograms.html")
    if err != nil {
        log.Fatal(err)
    }
    ids := parseAllProgramsPage(page)
    if len(ids) != 556 {
        t.Error(
            "For", "parseAllProgramPage, length",
            "expected", 556,
            "got", len(ids),
        )
    }
    if ids[5] != "alla-ar-fotografer" {
        t.Error(
            "For", "parseAllProgramsPage, id",
            "expected", "alla-ar-fotografer",
            "got", ids[5],
        )
    }
}

func TestXMLParse(t *testing.T) {
    for _, pair := range showsXML {
        page, err := ioutil.ReadFile("testFiles/" + pair.testfile)
        if err != nil {
            log.Fatal(err)
        }
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

func TestShowPageParse(t *testing.T) {
    for _, pair := range showsPage {
        page, err := ioutil.ReadFile("testFiles/" + pair.testfile)
        if err != nil {
            log.Fatal(err)
        }
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
        page, err := ioutil.ReadFile("testFiles/" + pair.testfile)
        if err != nil {
            log.Fatal(err)
        }
        program := parseJSON(page)
        if program.VideoId != pair.playid {
            t.Error(
                "For", pair.id,
                "expected", pair.playid,
                "got", program.VideoId,
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
        page, err := ioutil.ReadFile("testFiles/" + pair.testfile)
        if err != nil {
            log.Fatal(err)
        }
        program := parseJSON(page)
        e := parseBasicEpisodeInformation(program, pair.id)
        if e.Live != pair.live {
            t.Error(
                "For", pair.id,
                "expected", pair.live,
                "got", e.Live,
            )
        }
        if e.PlayId != pair.playid {
            t.Error(
                "For", pair.id,
                "expected", pair.playid,
                "got", e.PlayId,
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

func TestGetVideoUrl(t *testing.T) {
    for _, pair := range episodes {
        page, err := ioutil.ReadFile("testFiles/" + pair.testfile)
        if err != nil {
            log.Fatal(err)
        }
        program := parseJSON(page)
        s := getVideoUrl(program.Video.VideoReferences)
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
        page, err := ioutil.ReadFile("testFiles/" + pair.testfile)
        if err != nil {
            log.Fatal(err)
        }
        desc := parseDescription(page)
        if desc != pair.description {
            t.Error(
                "For", pair.testfile,
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
