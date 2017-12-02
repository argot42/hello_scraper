package main

import (
    "fmt"
    "testing"
    "os"
)

func TestExtract (t *testing.T) {
    var path string
    fmt.Scan(&path)

    f,err := os.Open(path)
    if err != nil {
        return
    }

    urls := extract_links(f)

    if len(urls) != 92 {
        t.Error("Expected 92 episodes, extracted ", len(urls))
    }
}
