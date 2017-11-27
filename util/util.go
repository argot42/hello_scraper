package util

import (
    "flag"
    "bytes"
    "strings"
    "path"
)

type Configuration struct {
    ARCHIVE string
    RSS string
    Files_path string
    Download_dir string
    Verbose bool
    All bool
}

/* globals */
var ARCHIVE_URL string = "https://www.hellointernet.fm/archive/"
var RSS_URL string = "http://www.hellointernet.fm/podcast?format=rss"
var DOWNLOAD_ALL bool = false
var FILES_PATH string = "~/.local/share/HI_scraper"
var DOWNLOAD_DIR string = "~/downloads/HI"
var VERBOSE bool = false
/***********/

func Get_parameters() Configuration {
    archive_ptr := flag.String("archive", ARCHIVE_URL, "Archive's URL")
    rss_ptr := flag.String("rss", RSS_URL, "RSS URL")
    all_ptr := flag.Bool("all", DOWNLOAD_ALL, "Download all available episodes")
    files_ptr := flag.String("files", FILES_PATH, "Control files directory")
    download_dir_ptr := flag.String("d", DOWNLOAD_DIR, "Directory where episodes will be stored")
    verbose_ptr := flag.Bool("v", VERBOSE, "Verbose output")

    flag.Parse()

    return Configuration {*archive_ptr, *rss_ptr, *files_ptr, *download_dir_ptr, *verbose_ptr, *all_ptr}
}

func Concat (str ...string) string {
    var buf bytes.Buffer

    for _,s := range str {
        buf.WriteString(s)
    }

    return buf.String()
}

func Get_name (url string) string {
    name := path.Base (url)

    return name[:strings.LastIndex(name, "?")]
}
