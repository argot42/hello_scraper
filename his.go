package main

import (
	"./util"
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
)

func main() {
	config := util.Get_parameters()
	download_episodes(config)
}

func download_episodes(config util.Configuration) {
	var episodes []string

	episodes = fetch_episodes(config)

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

	download(episodes, config)
}

func download(eps []string, config util.Configuration) {
	for i, url := range eps {
		// send http request
		res, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		// get file size and convert it to mb
		file_size, _ := strconv.Atoi(res.Header.Get("Content-Length"))

		// create file
		f, err := os.Create(filepath.Join(config.Download_dir, util.Get_name(url)))
		if err != nil {
			log.Fatal(err)
		}

		// create a buffer and start downloading
		b := make([]byte, 4*1024)
		total := 0 // written until now
		for {
			// read from stream
			n, err := res.Body.Read(b)
			if err != nil && err != io.EOF {
				log.Fatal(err)
			}

			if err == io.EOF {
				break
			}

			// write to file
			written, err := f.Write(b[:n])
			if err != nil {
				log.Fatal(err)
			}
			total += written

			// if verbose update bar
			if config.Verbose {
				util.Update_bar(i, len(eps), total, file_size)
			}
		}

		// if verbose finish bar
		if config.Verbose {
			util.Finish_bar(i, len(eps))
		}

		// close open stuff
		res.Body.Close()
		f.Close()
	}
}

func fetch_episodes(config util.Configuration) (urls []string) {
	// get rss xml
	res, err := http.Get(config.RSS)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()

	urls = extract_links(res.Body)

	if !config.All {
		//watched := get_watched (filepath.Join(config.Files_path, "watched_list"))
		urls = filter_watched(urls, config.Files_path, "watched_list")
	}

	return
}

func extract_links(r io.Reader) (urls []string) {
	dec := xml.NewDecoder(r)

	for {
		// read tokens from XML in stream
		t, err := dec.Token()
		if err != nil && err != io.EOF {
			log.Fatal(err)
		} else if err == io.EOF {
			break
		}

		// check link tag
		switch s := t.(type) {
		case xml.StartElement:
			if s.Name.Local == "enclosure" {
				urls = append(urls, s.Attr[0].Value)
			}
		}
	}

	return
}

func filter_watched(urls []string, path, filename string) (filtered_urls []string) {
	name := filepath.Join(path, filename)

	// if file doesn't exist create it
	if _, err := os.Stat(name); os.IsNotExist(err) {
		// try to create directory structure
		os.MkdirAll(path, os.ModePerm)
	}

	// open file
	f, err := os.OpenFile(name, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Split(bufio.ScanLines)

	// read file by line and delete matching lines from urls
	var new_watched []string
	for scanner.Scan() {
		line := scanner.Text()
		new_watched = append(new_watched, line)

		for i := 0; i < len(urls); i++ {
			if line == urls[i] {
				util.Del(i, &urls)
			}
		}
	}

	// add downloaded episodes to the watched file
	for _, ep := range urls {
		util.Ordered_insert(&new_watched, ep)
	}

	// write to file
	f.Truncate(0)
	f.Seek(0, 0)
	for _, ep := range new_watched {
		_, err := f.WriteString(ep + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}
	f.Sync()

	// return episodes to download
	filtered_urls = urls
	return
}
