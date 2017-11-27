package main

import (
    "fmt"
    "net/http"
    "io"
    "os"
    "log"
    "path/filepath"
    "./util"
)

func main() {
    config := util.Get_parameters()
    download_episodes (config)
}

func download_episodes (config util.Configuration) {
    var episodes []string

    if (config.All) {
        episodes = fetch_all_episodes(config)

    }else {
        episodes = fetch_new_episodes(config)
    }

    if len(episodes) == 0 {
        if config.Verbose {
            fmt.Println("No new episodes!")
        }
        return
    }

    // if verbose count fetched links
    if config.Verbose {
        fmt.Printf("%d episodes fetched!\n", len(episodes))
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
        f,err := os.Create ( filepath.Join (config.Download_dir, util.Get_name(url)) )
        if err != nil {
            log.Fatal(err)
        }

        // if verbose, set up progress bar
        if config.Verbose {
            //bar.Setup(100)
            fmt.Println("setup bar")
        }

        // create a buffer and start downloading
        b := make([]byte, 4*1024)
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

            // if verbose update bar
            if config.Verbose {
                fmt.Printf("update bar: %d\n", written)
            }
        }

        // close open stuff
        res.Body.Close()
        f.Close()

    }
}

func fetch_all_episodes (config util.Configuration) (urls []string) {
    urls = []string{"http://hwcdn.libsyn.com/p/a/1/9/a195a79840d036db/HI92.mp3?c_id=17903889&expiration=1511807679&hwt=f34d3136a5b0534f9727ea29636cf9ec",}
    return
}

func fetch_new_episodes (config util.Configuration) (urls []string) {
    var s []string
    return s
}
