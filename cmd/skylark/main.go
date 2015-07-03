package main

import (
    "flag"
    "fmt"
    "github.com/michaelhakansson/skylark"
    "os"
)

var exit = os.Exit

var (
    usage      = "Usage: skylark [OPTIONS]"
    options    = "Options:\n-h, -help \t Print this help text and exit \n-v, -version \t Print program version and exit"
    version    = "2015.07.03"
    help       = fmt.Sprintf("%s\nVersion: %s\n%s", usage, version, options)
    cliVersion = flag.Bool("version", false, version)
    cliHelp    = flag.Bool("help", false, help)
)

func init() {
    flag.BoolVar(cliVersion, "v", false, version)
    flag.BoolVar(cliHelp, "h", false, help)
}

func main() {
    flag.Parse()

    if *cliVersion {
        fmt.Println(flag.Lookup("version").Usage)
        exit(0)
        return
    }
    if *cliHelp {
        fmt.Println(flag.Lookup("help").Usage)
        exit(0)
        return
    }

    skylark.Start()
    exit(0)
    return
}
