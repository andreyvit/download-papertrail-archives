package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	urlpkg "net/url"
	"os"
	"path/filepath"
	"sort"
	"time"
)

var token = flag.String("token", "", "Papertrail HTTP API token from https://papertrailapp.com/account/profile")

var outputDir = flag.String("o", ".", "Output directory (default: .)")

const (
	archivesEndpoint = "https://papertrailapp.com/api/v1/archives.json"
)

func main() {
	log.SetFlags(0)
	flag.Parse()

	if *token == "" {
		log.Fatalf("** Papertrail token required (-token <TOKEN>)")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	headers := http.Header{
		"X-Papertrail-Token": []string{*token},
	}

	var archives []*Archive

	err := downloadJSON(archivesEndpoint, headers, client, &archives)
	if err != nil {
		log.Fatalf("** Fetching archives.json: %v", err)
	}

	sort.Slice(archives, func(i, j int) bool {
		return archives[i].Filename > archives[j].Filename
	})

	stats := struct {
		Existing   int
		Failed     int
		Downloaded int
	}{}

	for _, a := range archives {
		fn := filepath.Join(*outputDir, a.Filename)
		url := a.Links.Download.URL
		if url == "" {
			log.Printf("No URL: %v", a.Filename)
			continue
		}

		if _, err := os.Stat(fn); err == nil {
			log.Printf("Exists: %v", a.Filename)
			stats.Existing++
		} else {
			log.Printf("Downloading: %v", a.Filename)

			raw, err := download(url, headers, client)
			if err != nil {
				log.Printf("Failed to download %v: %v", a.Filename, err)
				stats.Failed++
				continue
			}

			err = ioutil.WriteFile(fn, raw, 0644)
			if err != nil {
				log.Fatalf("Failed to save %v: %v", a.Filename, err)
			}

			stats.Downloaded++
		}
	}

	log.Printf("%d downloaded, %d failed, %d existing skipped", stats.Downloaded, stats.Failed, stats.Existing)
}

type Archive struct {
	Filename string `json:"filename"`
	Links    struct {
		Download struct {
			URL string `json:"href"`
		} `json:"download"`
	} `json:"_links"`
}

func download(url string, headers http.Header, client *http.Client) ([]byte, error) {
	u, err := urlpkg.Parse(url)
	if err != nil {
		return nil, err
	}

	r := &http.Request{
		Method: http.MethodGet,
		URL:    u,
		Header: headers,
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if !(resp.StatusCode >= 200 && resp.StatusCode < 300) {
		return nil, fmt.Errorf("HTTP %s", resp.Status)
	}

	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func downloadJSON(url string, headers http.Header, client *http.Client, resp interface{}) error {
	b, err := download(url, headers, client)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, resp)
}
