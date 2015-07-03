package utils

import (
    "go/build"
    "log"
    "runtime"
)

// GetPath returns the path to the build directory
func GetPath() string {
    p, _ := build.Default.Import("github.com/michaelhakansson/skylark", "", build.FindOnly)
    return p.Dir
}

// Checkerr checks for errors
func Checkerr(err error) {
    if err != nil {
        pc, fn, line, _ := runtime.Caller(1)

        log.Printf("[error] in %s[%s:%d] %v", runtime.FuncForPC(pc).Name(), fn, line, err)
        log.Fatal(err)
    }
}
