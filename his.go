package main

import (
    "fmt"
    "net/http"
    "io"
    "os"
    "log"
    "./util"
)

func main() {
    config := get_parameters()
    download_episodes (config)
}

func download_episodes (config util.Configuration) {
    var episodes []string

    if (config.all) {
        episodes = fetch_all_episodes()

    }else {
        episodes = fetch_new_episodes()
    }

    // if verbose count fetched links
    if verbose {
        fmt.Prinf("%d episodes fetched!\n", len(episodes))
    }

    download (episodes, config)
}

func download (eps []string, config util.Configuration) {
    for _,url := range eps {
        res,err := http.Get(url)
        if err != nil {
            log.Fatal(err)
        }

        // create file
        f,err := os.Create ( util.Concat (config.Download_dir, util.Get_name(url)) )
        if err != nil {
            log.Fatal(err)
        }

        // if verbose, set up progress bar
        if verbose {
            //bar.Setup(100)
            fmt.Pritln("setup bar")
        }

        // create a buffer and start downloading
        b := make([]byte, 4*1024)
        for {
            // read from stream
            n,err := res.Body.Read(b)
            if err != nil {
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

            // if verbose update bar
            if verbose {
                fmt.Printf("update bar: %s\n", written)
            }
        }

        // close open stuff
        res.Body.Close()
        f.Close()
    }
}
