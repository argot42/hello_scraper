package main

import (
    "fmt"
    "net/http"
    "io"
    "os"
    "log"
    "path/filepath"
    "strconv"
    "encoding/xml"
    "golang.org/x/net/html"
    "./util"
)

func main() {
    config := util.Get_parameters()
    download_episodes (config)
}

func download_episodes (config util.Configuration) {
    var episodes []string

    episodes = fetch_episodes (config)

    if len(episodes) == 0 {
        if config.Verbose {
            fmt.Println("No new episodes!")
        }
        return
    }

    // if verbose count fetched links
    if config.Verbose {
        fmt.Printf("%d episode(s) fetched!\n", len(episodes))
    }

    download (episodes, config)
}

func download (eps []string, config util.Configuration) {
    for i,url := range eps {
        // send http request
        res,err := http.Get(url)
        if err != nil {
            log.Fatal(err)
        }
        // get file size and convert it to mb
        file_size,_ := strconv.Atoi (res.Header.Get("Content-Length"))

        // create file
        f,err := os.Create ( filepath.Join (config.Download_dir, util.Get_name(url)) )
        if err != nil {
            log.Fatal(err)
        }

        // create a buffer and start downloading
        b := make([]byte, 4*1024)
        total := 0 // written until now
        for {
            // read from stream
            n,err := res.Body.Read(b)
            if err != nil && err != io.EOF {
                log.Fatal(err)
            }

            if err == io.EOF {
                break
            }

            // write to file
            written,err := f.Write(b[:n])
            if err != nil {
                log.Fatal(err)
            }
            total += written

            // if verbose update bar
            if config.Verbose {
                util.Update_bar (i, len(eps), total, file_size)
            }
        }

        // if verbose finish bar
        if config.Verbose {
            util.Finish_bar (i, len(eps))
        }

        // close open stuff
        res.Body.Close()
        f.Close()
    }
}

func fetch_episodes (config util.Configuration) (urls []string) {
    var watched []int

    if (!config.All) {
        fmt.Println("all")
    }

    // get rss xml
    res,err := http.Get (config.RSS)
    if err != nil {
        log.Fatal(err)
    }
    defer res.Body.Close()

    dec := xml.NewDecoder (res.Body)

    for {
        // read tokens from XML in body stream
        t,err := dec.Token()
        if err != nil {
            log.Fatal(err)
        }

        if t == nil {
            break
        }

        // check link tag
        switch s := t.(type) {
        case xml.StartElement:
            if s.Name.Local == "enclosure" {
                if watched_ep (s.Attr[0].Value, watched) {
                    urls = append (urls, s.Attr[0].Value)
                }
            }
        }
    }
}
