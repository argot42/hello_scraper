package util

import (
	"flag"
	"fmt"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
)

type Configuration struct {
	RSS          string
	Files_path   string
	Download_dir string
	Verbose      bool
	All          bool
}

/* globals */
var RSS_URL string = "http://www.hellointernet.fm/podcast?format=rss"
var DOWNLOAD_ALL bool = false
var FILES_PATH string = filepath.Join(get_home_dir(), ".local/share/HI")
var DOWNLOAD_DIR string = filepath.Join(get_home_dir(), "downloads/HI")
var VERBOSE bool = false

// bar
var BAR_WIDTH int = 40

/***********/

func Get_parameters() Configuration {
	rss_ptr := flag.String("rss", RSS_URL, "RSS URL")
	all_ptr := flag.Bool("a", DOWNLOAD_ALL, "Download all available episodes")
	files_ptr := flag.String("files", FILES_PATH, "Control files directory")
	download_dir_ptr := flag.String("d", DOWNLOAD_DIR, "Directory where episodes will be stored")
	verbose_ptr := flag.Bool("v", VERBOSE, "Verbose output")

	flag.Parse()

	return Configuration{*rss_ptr, *files_ptr, *download_dir_ptr, *verbose_ptr, *all_ptr}
}

func Get_name(url string) string {
	name := path.Base(url)

	if li := strings.LastIndex(name, "?"); li > -1 {
		return name[:li]
	}

	return name
}

func get_home_dir() string {
	usr, _ := user.Current()
	return usr.HomeDir
}

/* math? */
func Round(x, unit float64) float64 {
	if x > 0 {
		return float64(int64(x/unit+0.5)) * unit
	}
	return float64(int64(x/unit-0.5)) * unit
}

func count_digits(n int) (nd int) {
	s := strconv.Itoa(n)
	nd = len(s)

	if n < 0 {
		nd--
	}
	return
}

/* bar */
func Update_bar(file_count int, total_files int, total_written int, file_size int) {
	percent := (float64(total_written) / float64(file_size)) * 100.0
	barfiller := int((percent / 100) * float64(BAR_WIDTH))
	dynamic_space := count_digits(file_count) + count_digits(total_files)

	fmt.Printf("[%s%s] %3.f%% (%d/%d)",
		strings.Repeat("-", barfiller),
		strings.Repeat(" ", BAR_WIDTH-barfiller),
		percent,
		file_count,
		total_files,
	)

	// move cursor back
	fmt.Printf(strings.Repeat("\b", BAR_WIDTH+11+dynamic_space))
}

func Finish_bar(file_count int, total_files int) {
	fmt.Printf("[%s] 100%% (%d/%d)\n",
		strings.Repeat("-", BAR_WIDTH),
		file_count+1,
		total_files,
	)
}

func Del(index int, l *[]string) {
	aux := *l
	aux = append(aux[:index], aux[index+1:]...)
	*l = aux
}

func Ordered_insert(l *[]string, item string) {
	position := Find_position(item, *l, 0, len(*l)-1)

	if position == len(*l) {
		*l = append(*l, item)

	} else {
		var aux []string
		aux = append(aux, (*l)[:position]...)
		aux = append(aux, item)
		aux = append(aux, (*l)[position:]...)
		*l = aux
	}
}

func Find_position(letter string, l []string, start, end int) int {
	if end > start {
		mid := start + (end-start)/2

		if l[mid] == letter {
			return mid
		}

		if l[mid] > letter {
			return Find_position(letter, l, start, mid-1)
		}

		return Find_position(letter, l, mid+1, end)

	} else if end == start {
		if l[start] >= letter {
			return start
		}

		return start + 1
	}

	return -1
}
