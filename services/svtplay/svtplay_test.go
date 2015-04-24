package svtplay

import(
    "testing"
)


type testshow struct {
    id string
    title string
    number int
}

type testepisode struct {
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

var showsXML = []testshow{
    {"vetenskapens-varld", "Vetenskapens värld", 17},
    {"uppdrag-granskning", "Uppdrag granskning", 15},
}

var showsPage = []testshow{
    {"vetenskapens-varld", "Vetenskapens värld", 15},
    {"uppdrag-granskning", "Uppdrag granskning", 15},
}

var episodes = []testepisode{
    {"2843612", "", "kultur-och-nöje", "Del 4 av 10. Gruppfinal, med bland annat Sisyfos-tävlingen - vätskefyllda pilatesbollar som ska rullas uppför en backe. Några av Sveriges främsta idrottsmän och idrottskvinnor möts i fysiska och psykiska utmaningar för att kora Mästarnas mästare 2015. I Grupp 1 ingår Anette Norberg, Anna Olsson, Magnus Muhrén, Danijela Rundqvist, Glenn Hysén och Björn Lind. Programledare: Micke Leijnegard.",
    "4", "58m37s", false, 2843612, "7",
    "http://www.svt.se/cachable_image/1429226401000/svts/article2849370.svt/ALTERNATES/extralarge/default_title",
    "Avsnitt 4", "http://svtplay18p-f.akamaihd.net/i/se/open/20150417/1360782-004A/EPISOD-1360782-004A-2ae7758f8108a631_,892,144,252,360,540,1584,2700,.mp4.csmil/master.m3u8?cc1=name=Svenska~default=yes~forced=no~uri=http://media.svt.se/download/mcc/wp3/undertexter-wsrt/1360782/1360782-004A/C(sv)/index.m3u8~lang=sv"},
    {"2867878", "", "nyheter", "Kan ses till sön 3 maj (9 dagar kvar)", "11:00", "1m30s", false, 2867878, "23/4", "http://www.svt.se/cachable_image/1429781701000/svts/article2867877.svt/ALTERNATES/extralarge/default_title",
    "23/4 11.00", "http://svtplay19i-f.akamaihd.net/i/world/open/20150423/1368669-074A/EPISOD-1368669-074A-208212fd95c96099_,892,144,252,360,540,1584,2700,.mp4.csmil/master.m3u8"},
}

func TestProgramIds(t *testing.T) {
    programs := GetAllProgramIds()
    if len(programs) != 553 {
        t.Error(
            "For", "GetAllProgramIds",
            "expected", 553,
            "got", len(programs),
        )
    }
}

func TestXMLParse(t *testing.T) {
    for _, pair := range showsXML {
        xmlUrl :=  playUrlBase + pair.id + rssUrl
        b := getPage(xmlUrl)
        s, e := parseShowXML(b, pair.id)
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
        pageUrl := playUrlBase + pair.id
        b := getPage(pageUrl)
        s, e := parseShowPage(b, pair.id)
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

func TestGetEpisode(t *testing.T) {
    for _, pair := range episodes {
        e := GetEpisode(pair.id)
        if e.Description != pair.description {
            t.Error(
                "For", pair.id,
                "expected", pair.description,
                "got", e.Description,
            )
        }
        if e.EpisodeNumber != pair.episodenumber {
            t.Error(
                "For", pair.id,
                "expected", pair.episodenumber,
                "got", e.EpisodeNumber,
            )
        }
        if e.Length != pair.length {
            t.Error(
                "For", pair.id,
                "expected", pair.length,
                "got", e.Length,
            )
        }
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
        if e.Season != pair.season {
            t.Error(
                "For", pair.id,
                "expected", pair.season,
                "got", e.Season,
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
        if e.VideoUrl != pair.videourl {
            t.Error(
                "For", pair.id,
                "expected", pair.videourl,
                "got", e.VideoUrl,
            )
        }
    }
}
