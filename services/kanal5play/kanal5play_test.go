package kanal5play

import(
    "io/ioutil"
    "log"
    "testing"
)

type testshow struct {
    testfile string
    title string
    linktoseasonspage string
}

type testseason struct {
    testfile string
    lastseasonelement string
    length int
}

type testepisodelinks struct {
    testfile string
    lastepisodelink string
    length int
}

type testepisode struct {
    testfile string
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

var shows = []testshow {
    {"show1.html", "Arga snickaren", "/program/226052/sasong/12"},
    {"show2.html", "Berg & Meltzer i Europa", "/program/3154653787/sasong/2"},
}

var seasons = []testseason {
    {"show1seasonpage.html", "/program/226052/sasong/1", 12},
    {"show2seasonpage.html", "/program/3154653787/sasong/2", 1},
}

var episodelinks = []testepisodelinks {
    {"show1seasonpage.html", "/program/226052/video/2453043232", 8},
    {"show2seasonpage.html", "/program/3154653787/video/3393402030", 7},
}

var episodes = []testepisode {
    {"show1episode.json", "", "", "Lars bor med sina tre barn i ett trev\u00e5ningshus p\u00e5 190 kvm i Olofstorp. Lars b\u00f6rjade bygga med sin f\u00f6re detta fru men n\u00e4r paret gick skilda v\u00e4gar avstannade bygget. Nu st\u00e5r Lars ensam med ett gigantiskt skal till hus utan fungerande k\u00f6k eller badrum.",
    "8", "43m37s", false, 2589533249, "12", "http://lh3.googleusercontent.com/iLDVLeYu3Jx05gcp3TWCAnbGTRKN14T10X6P7AjhKThsTUlnZOpHIr7zEQEoWGqg9VPmvi8_yfqLzHRyvbjZvw",
    "Familjen bor i ett skal", "http://hls0.00607-od0.dna.qbrick.com/00607-od0/_definst_/smil:20141202/20141202091544808-17bvu2qz6l78v9ud5o6sc6emq-967_ipad/playlist.m3u8"},
    {"show2episode.json", "", "", "Tjejerna forts\u00e4tter till Cypern. H\u00e4r ska de medverka i Cyperns st\u00f6rsta tv-s\u00e5pa. Vad de inte vet \u00e4r att det inneb\u00e4r k\u00e4rleksscener. Dessutom blir det lektioner i att krossa vattenmeloner med huvudet - l\u00e4rare \u00e4r ingen mindre \u00e4n Guinness rekordm\u00e4stare.",
    "7", "43m53s", false, 3398002034, "2", "http://lh3.googleusercontent.com/xCZzz8eHqDgk_evbAA2DR7_0ndOZxDNldia7sbkAny92dTu6ykY4mcEUH1QVQCQjExjtgp46bAJD99kbw9E",
    "Tjejerna forts\u00e4tter till Cypern", "http://hls0.00607-od0.dna.qbrick.com/00607-od0/_definst_/smil:20150423/20150423091014000-3guobpbya7hwvm831ohm0xxti-897_ipad/playlist.m3u8"},
}

func TestProgramIds(t *testing.T) {
    page, err := ioutil.ReadFile("testFiles/allprograms.html")
    if err != nil {
        log.Fatal(err)
    }
    ids := parseAllProgramsPage(page)
    if len(ids) != 173 {
        t.Error(
            "For", "parseAllProgramsPage, length",
            "expected", 173,
            "got", len(ids),
        )
    }
    if ids[5] != "226052" {
        t.Error(
            "For", "parseAllProgramsPage, id",
            "expected", "226052",
            "got", ids[5],
        )
    }
}

func TestShowParser(t *testing.T) {
    for _, pair := range shows {
        page, err := ioutil.ReadFile("testFiles/" + pair.testfile)
        if err != nil {
            log.Fatal(err)
        }
        show, linkToSeasonsPage := parseShowInfo(page, "")
        if show.Title != pair.title {
            t.Error(
                "For", pair.testfile,
                "expected", pair.title,
                "got", show.Title,
            )
        }
        if linkToSeasonsPage != pair.linktoseasonspage {
            t.Error(
                "For", pair.testfile,
                "expected", pair.linktoseasonspage,
                "got", linkToSeasonsPage,
            )
        }
    }
}

func TestSeasonParser(t *testing.T) {
    for _, pair := range seasons {
        page, err := ioutil.ReadFile("testFiles/" + pair.testfile)
        if err != nil {
            log.Fatal(err)
        }
        seasonLinks := parseSeasonLinks(page)
        if seasonLinks[len(seasonLinks) - 1] != pair.lastseasonelement {
            t.Error(
                "For", pair.testfile,
                "expected", pair.lastseasonelement,
                "got", seasonLinks[len(seasonLinks) - 1],
            )
        }
        if len(seasonLinks) != pair.length {
            t.Error(
                "For", pair.testfile,
                "expected", pair.length,
                "got", len(seasonLinks),
            )
        }
    }
}

func TestEpisodeLinksParser(t *testing.T) {
    for _, pair := range episodelinks {
        page, err := ioutil.ReadFile("testFiles/" + pair.testfile)
        if err != nil {
            log.Fatal(err)
        }
        episodeLinks := parseEpisodeLinksOnSeasonPage(page)
        if episodeLinks[len(episodeLinks) - 1] != pair.lastepisodelink {
            t.Error(
                "For", pair.testfile,
                "expected", pair.lastepisodelink,
                "got", episodeLinks[len(episodeLinks) - 1],
            )
        }
        if len(episodeLinks) != pair.length {
            t.Error(
                "For", pair.testfile,
                "expected", pair.length,
                "got", len(episodeLinks),
            )
        }
    }
}
func TestGetEpisode(t *testing.T) {
    for _, pair := range episodes {
        page, err := ioutil.ReadFile("testFiles/" + pair.testfile)
        if err != nil {
            log.Fatal(err)
        }
        episode := parseEpisode(page)
        if episode.Description != pair.description {
            t.Error(
                "For", pair.testfile,
                "expected", pair.description,
                "got", episode.Description,
            )
        }
        if episode.EpisodeNumber != pair.episodenumber {
            t.Error(
                "For", pair.testfile,
                "expected", pair.episodenumber,
                "got", episode.EpisodeNumber,
            )
        }
        if episode.Length != pair.length {
            t.Error(
                "For", pair.testfile,
                "expected", pair.length,
                "got", episode.Length,
            )
        }
        if episode.Live != pair.live {
            t.Error(
                "For", pair.testfile,
                "expected", pair.live,
                "got", episode.Live,
            )
        }
        if episode.PlayId != pair.playid {
            t.Error(
                "For", pair.testfile,
                "expected", pair.playid,
                "got", episode.PlayId,
            )
        }
        if episode.Season != pair.season {
            t.Error(
                "For", pair.testfile,
                "expected", pair.season,
                "got", episode.Season,
            )
        }
        if episode.Thumbnail != pair.thumbnail {
            t.Error(
                "For", pair.testfile,
                "expected", pair.thumbnail,
                "got", episode.Thumbnail,
            )
        }
        if episode.Title != pair.title {
            t.Error(
                "For", pair.testfile,
                "expected", pair.title,
                "got", episode.Title,
            )
        }
        if episode.VideoUrl != pair.videourl {
            t.Error(
                "For", pair.testfile,
                "expected", pair.videourl,
                "got", episode.VideoUrl,
            )
        }
    }
}
