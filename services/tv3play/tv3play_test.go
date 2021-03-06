package tv3play

import (
    "log"
    "testing"
    "time"
)

var tv3 = TV3Play{}

func TestFixHlsURL(t *testing.T) {
    log.Print("TestFixHlsUrl")
    originalURL := "http://mtgxpitcher01-vh.akamaihd.net/i/open/201502/04/V55886_mtgx_f906f926_,48,260,460,900,1800,.mp4.csmil/master.m3u8?__b__=300&cc1=name=Svenska~default=yes~forced=no~lang=sv~uri=http://cdn.subtitles.mtgx.tv/pitcher/V5xxxx/V558xx/V55886/0000003529/V55886_sub_sv.m3u8"
    expectedURL := "http://mtgxpitcher01-vh.akamaihd.net/i/open/201502/04/V55886_mtgx_f906f926_,1800,.mp4.csmil/master.m3u8?__b__=300&cc1=name=Svenska~default=yes~forced=no~lang=sv~uri=http://cdn.subtitles.mtgx.tv/pitcher/V5xxxx/V558xx/V55886/0000003529/V55886_sub_sv.m3u8"
    actualURL := fixHlsURL(originalURL)
    if actualURL != expectedURL {
        t.Error(
            "For", "fixHlsUrl",
            "expected", expectedURL,
            "got", actualURL,
        )
    }
}

func TestFixThumbnailURL(t *testing.T) {
    log.Print("TestFixThumbnailUrl")
    originalURL := "http://test.com/{size}/foo"
    expectedURL := "http://test.com/" + thumbnailSize + "/foo"
    actualURL := fixThumbnailURL(originalURL)
    if actualURL != expectedURL {
        t.Error(
            "For", "fixThumbnailUrl",
            "expected", expectedURL,
            "got", actualURL,
        )
    }

}

func TestParseDateTime(t *testing.T) {
    log.Print("TestParseDateTime")
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

// TESTS DEPENDENT ON LIVE DATA
type testshow struct {
    id        string
    title     string
    thumbnail string
    number    int
}

type testepisode struct {
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

var testNumberOfEpisodes = []testshow{
    {"804", "Adaktusson", "http://cdn.playapi.mtgx.tv/imagecache/1000x675/cloud/content-images/sites/viastream.viasat.tv/files/category_pictures/adaktusson_s6.jpg", 84},
    {"9471", "Mozart in the Jungle", "http://cdn.playapi.mtgx.tv/imagecache/1000x675/cloud/content-images/seasons/11519/season/mozart_in_the_jungle_sasong1-bCeSdEa.jpg", 3},
}

var episodes = []testepisode{
    {"469483", "2015 04 01 20:30:00 +0000 UTC", "Humor",
    "I första avsnittet handlar det om kroppen och hur vi använder den. Svenska folket har i en enkät fått svara på vad de tycker om naken överkropp på stan, om hur vi är nakna tillsammans, om kroppsljud och kroppslukter och hur det egentligen funkar med den berömda svenska kompiskramen.",
    "1", "22m51s", false, 469483, "1",
    "http://cdn.playapi.mtgx.tv/imagecache/" + thumbnailSize + "/cloud/content-images/seasons/9825/season/inteok.jpg",
    "Inte OK S01E01",
    "http://mtgxse02-vh.akamaihd.net/i/open/201410/24/V43645_mtgx_b8f111a1_,2800,.mp4.csmil/master.m3u8?__b__=300"},

    {"23636", "2009 09 01 21:30:00 +0000 UTC", "Samhälle och aktualitet",
    "I säsongspremiären av Adaktusson undrar vi varför svenska domare får sitta kvar trots att de dömts för allvarliga brott.",
    "1", "24m49s", false, 23636, "6",
    "http://cdn.playapi.mtgx.tv/imagecache/" + thumbnailSize + "/cloud/content-images/sites/viastream.viasat.tv/files/category_pictures/adaktusson_s6.jpg",
    "Adaktusson 1/9", "http://www.tv8play.se/program/adaktusson/23636"},
}

func TestGetShow(t *testing.T) {
    log.Print("TestGetShow")
    for _, pair := range testNumberOfEpisodes {
        s, e := tv3.GetShow(pair.id)
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
    log.Print("TestGetEpisode")
    for _, pair := range episodes {
        e := tv3.GetEpisode(pair.id)
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
        if e.PlayID != pair.playid {
            t.Error(
                "For", pair.id,
                "expected", pair.playid,
                "got", e.PlayID,
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
        if e.VideoURL != pair.videourl {
            t.Error(
                "For", pair.id,
                "expected", pair.videourl,
                "got", e.VideoURL,
            )
        }
    }
}

func TestProgramIDs(t *testing.T) {
    log.Print("TestProgramIDs")
    ids := tv3.GetAllProgramIDs()
    if len(ids) != 129 {
        t.Error(
            "For", "GetAllProgramsIDs",
            "expected", 129,
            "got", len(ids),
        )
    }
}

// END TESTS DEPENDENT ON LIVE DATA
