package utils

import(
    "strconv"
    "time"
)

// TimeString converts a time to a readable string
func TimeString(t time.Time) string {
    return t.Format("2006-01-02 15:04")
}

// PrettifyServiceText makes the (play)service id string to a nice looking
// string, for use in GUIs
func PrettifyServiceText(service string) (prettyText string) {
    prettyText = service
    switch service {
    case "svtplay":
        prettyText = "SVT Play"
    case "tv3play":
        prettyText = "TV3 Play"
    case "kanal5play":
        prettyText = "Kanal 5 Play"
    }
    return
}

// ZeroPaddingString adds an initial zero before string values between 1 and 9
func ZeroPaddingString(i string) string {
    number, err := strconv.ParseInt(string(i), 0, 64)
    if err != nil {
        return i
    }
    return ZeroPadding(number)
}

// ZeroPadding adds an initial zero to value between 1 and 9 and return it as a
// string
func ZeroPadding(i int64) string {
    if i < 10 {
        return "0" + strconv.FormatInt(i, 10)
    }
    return strconv.FormatInt(i, 10)
}
