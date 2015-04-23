package svtplay

import(
    "testing"
)


type testshow struct {
    id string
    title string
    number int
}

var showsXML = []testshow{
    {"vetenskapens-varld", "Vetenskapens värld", 18},
    {"uppdrag-granskning", "Uppdrag granskning", 15},
}

var showsPage = []testshow{
    {"vetenskapens-varld", "Vetenskapens värld", 15},
    {"uppdrag-granskning", "Uppdrag granskning", 15},
}

func TestPrograms(t *testing.T) {
    programs := GetAllPrograms()
    if len(programs) != 553 {
        t.Error(
            "For", "GetAllPrograms",
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
