package tv3play

import(
    "testing"
    "time"
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

func TestFixHlsUrl(t *testing.T) {
    originalUrl := "http://mtgxpitcher01-vh.akamaihd.net/i/open/201502/04/V55886_mtgx_f906f926_,48,260,460,900,1800,.mp4.csmil/master.m3u8?__b__=300&cc1=name=Svenska~default=yes~forced=no~lang=sv~uri=http://cdn.subtitles.mtgx.tv/pitcher/V5xxxx/V558xx/V55886/0000003529/V55886_sub_sv.m3u8"
    expectedUrl := "http://mtgxpitcher01-vh.akamaihd.net/i/open/201502/04/V55886_mtgx_f906f926_,1800,.mp4.csmil/master.m3u8?__b__=300&cc1=name=Svenska~default=yes~forced=no~lang=sv~uri=http://cdn.subtitles.mtgx.tv/pitcher/V5xxxx/V558xx/V55886/0000003529/V55886_sub_sv.m3u8"
    actualUrl := fixHlsUrl(originalUrl)
    if actualUrl != expectedUrl {
        t.Error(
            "For", "fixHlsUrl",
            "expected", expectedUrl,
            "got", actualUrl,
        )
    }
}

func TestFixThumbnailUrl(t *testing.T) {
    originalUrl := "http://test.com/{size}/foo"
    expectedUrl := "http://test.com/" + thumbnailSize + "/foo"
    actualUrl := fixThumbnailUrl(originalUrl)
    if actualUrl != expectedUrl {
        t.Error(
            "For", "fixThumbnailUrl",
            "expected", expectedUrl,
            "got", actualUrl,
        )
    }

}

func TestParseDateTime(t *testing.T) {
    originalDate := "2015-10-31T21:33:00+00:00"
    formattedOriginalDate := "2015 10 31 21 33"
    expectedDate, _ := time.Parse("2006 01 02 15 04", formattedOriginalDate)
    actualDate := parseDateTime(originalDate)
    if actualDate != expectedDate {
        t.Error(
            "For", "parseDateTime",
            "expected", expectedDate,
            "got", actualDate,
        )
    }
}

// func TestPrograms(t *testing.T) {
//     programs := GetAllPrograms()
//     if len(programs) != 455 {
//         t.Error(
//             "For", "GetAllPrograms",
//             "expected", 455,
//             "got", len(programs),
//         )
//     }
// }

// func TestXMLParse(t *testing.T) {
//     for _, pair := range showsXML {
//         xmlUrl :=  playUrlBase + pair.id + rssUrl
//         b := getPage(xmlUrl)
//         s, e := parseShowXML(b, pair.id)
//         if s.Title != pair.title {
//             t.Error(
//                 "For", pair.id,
//                 "expected", pair.title,
//                 "got", s.Title,
//             )
//         }
//         if len(e) != pair.number {
//             t.Error(
//                 "For", pair.id,
//                 "expected", pair.number,
//                 "got", len(e),
//             )
//         }
//     }
// }

// func TestShowPageParse(t *testing.T) {
//     for _, pair := range showsPage {
//         pageUrl := playUrlBase + pair.id
//         b := getPage(pageUrl)
//         s, e := parseShowPage(b, pair.id)
//         if s.Title != pair.title {
//             t.Error(
//                 "For", pair.id,
//                 "expected", pair.title,
//                 "got", s.Title,
//             )
//         }
//         if len(e) != pair.number {
//             t.Error(
//                 "For", pair.id,
//                 "expected", pair.number,
//                 "got", len(e),
//             )
//         }
//     }
// }
