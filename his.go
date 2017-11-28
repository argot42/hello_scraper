package main

import (
    "fmt"
    "net/http"
    "io"
    "os"
    "log"
    "path/filepath"
    "strconv"
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

        // close open stuff
        res.Body.Close()
        f.Close()

        // if verbose finish bar
        if config.Verbose {
            util.Finish_bar (i, len(eps))
        }
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
