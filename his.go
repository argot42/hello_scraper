package main

import (
    "fmt"
    "flag"
    "net/http"
    "io"
    "os"
    "log"
    "./bar"
)

type configuration struct {
    ARCHIVE string
    RSS string
    files_path string
    download_dir string
    verbose bool
    all bool
}

/* globals */
var ARCHIVE_URL string = "https://www.hellointernet.fm/archive/"
var RSS_URL string = "http://www.hellointernet.fm/podcast?format=rss"
var DOWNLOAD_ALL bool = false
var FILES_PATH string = "~/.local/share/HI_scraper"
var DOWNLOAD_DIR string = "~/downloads/HI"
var VERBOSE bool = false
/***********/

func main() {
    config := get_parameters()
    download_episodes (config)
}

func get_parameters() configuration {
    archive_ptr := flag.String("archive", ARCHIVE_URL, "Archive's URL")
    rss_ptr := flag.String("rss", RSS_URL, "RSS URL")
    all_ptr := flag.Bool("all", DOWNLOAD_ALL, "Download all available episodes")
    files_ptr := flag.String("files", FILES_PATH, "Control files directory")
    download_dir_ptr := flag.String("d", DOWNLOAD_DIR, "Directory where episodes will be stored")
    verbose_ptr := flag.Bool("v", VERBOSE, "Verbose output")

    flag.Parse()

    return configuration {*archive_ptr, *rss_ptr, *files_ptr, *download_dir_ptr, *verbose_ptr, *all_ptr}
}

func download_episodes (config configuration) {
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

func download (eps []string, config configuration) {
    for _,url := range eps {
        res,err := http.Get(url)
        if err != nil {
            log.Fatal(err)
        }

        // create file
        f,err := os.Create ( concat (config.download_dir, get_name(url)) )
        if err != nil {
            log.Fatal(err)
        }

        // if verbose, set up progress bar
        if verbose {
            bar.Setup(100)
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
        }

        // close things
        res.Body.Close()
        f.Close()
    }
}
