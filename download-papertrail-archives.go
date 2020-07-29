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
	"regexp"
	"sort"
	"time"
)

var (
	token     = flag.String("token", "", "Papertrail HTTP API token from https://papertrailapp.com/account/profile")
	outputDir = flag.String("o", ".", "Output directory")
	timeout   = flag.Duration("timeout", 30*time.Second, "timeout for HTTP operations")
	quiet     = flag.Bool("q", false, "quiet operation (don't print anything)")

	before, since Date
)

const (
	archivesEndpoint = "https://papertrailapp.com/api/v1/archives.json"
)

func main() {
	flag.Var(&before, "before", "only download logs before (NOT including) this date in YYYY-MM-DD format")
	flag.Var(&since, "since", "only download logs on or after this date in YYYY-MM-DD format")
	log.SetFlags(0)
	flag.Parse()

	if *token == "" {
		log.Fatalf("** Papertrail token required (-token <TOKEN>), see https://papertrailapp.com/account/profile")
	}

	client := &http.Client{
		Timeout: *timeout,
	}
	headers := http.Header{
		"X-Papertrail-Token": []string{*token},
	}

	var archives []*Archive

	if !*quiet {
		log.Printf("Fetching a list of archives...")
	}
	err := downloadJSON(archivesEndpoint, headers, client, &archives)
	if err != nil {
		log.Fatalf("** ERROR: cannot fetch a list of archives: %v", err)
	}

	sort.Slice(archives, func(i, j int) bool {
		return archives[i].Filename > archives[j].Filename
	})

	stats := struct {
		Existing   int
		Failed     int
		Downloaded int
		Skipped    int
	}{}

	for _, a := range archives {
		fn := filepath.Join(*outputDir, a.Filename)
		url := a.Links.Download.URL
		if url == "" {
			log.Printf("WARNING: Skipping entry with no URL: %v", a.Filename)
			continue
		}

		if !since.IsZero() && a.Filename < since.String() {
			stats.Skipped++
		}
		if !before.IsZero() && a.Filename >= before.String() {
			stats.Skipped++
		}

		if _, err := os.Stat(fn); err == nil {
			stats.Existing++
		} else {
			if !*quiet {
				log.Printf("Downloading %v...", a.Filename)
			}

			raw, err := download(url, headers, client)
			if err != nil {
				log.Printf("WARNING: cannot download %v: %v", a.Filename, err)
				stats.Failed++
				continue
			}

			err = ioutil.WriteFile(fn, raw, 0644)
			if err != nil {
				log.Fatalf("** ERROR: cannot save %v: %v", a.Filename, err)
			}

			stats.Downloaded++
		}
	}

	if !*quiet {
		log.Printf("Done: %d downloaded, %d failed, %d already exist, %d skipped.", stats.Downloaded, stats.Failed, stats.Existing, stats.Skipped)
	}
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

var (
	dayRe = regexp.MustCompile(`^(\d\d\d\d)-(\d\d)-(\d\d)$`)
)

type Date string

func FromTime(tm time.Time) Date {
	return Date(tm.UTC().Format("2006-01-02"))
}
func Today() Date {
	return FromTime(time.Now())
}
func Yesterday() Date {
	return FromTime(time.Now().AddDate(0, 0, -1))
}

func (v Date) IsZero() bool {
	return len(v) == 0
}

func (v Date) String() string {
	return string(v)
}

func ParseDate(s string) (Date, error) {
	switch s {
	case "tod", "today":
		return Today(), nil
	case "yest", "yesterday":
		return Yesterday(), nil
	}
	if m := dayRe.FindStringSubmatch(s); m != nil {
		return Date(m[0]), nil
	} else {
		return Date(""), fmt.Errorf("invalid date: %q", s)
	}
}

func (v *Date) Set(s string) error {
	d, err := ParseDate(s)
	*v = d
	return err
}
