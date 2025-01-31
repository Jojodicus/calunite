package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func fetchFile(name string) (string, error) {
	content, err := os.ReadFile(name)
	return string(content), err
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func fetchUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func fetch(thing string) (string, error) {
	_, err := os.Stat(thing)
	if err == nil {
		// file exists
		return fetchFile(thing)
	}
	if IsUrl(thing) {
		return fetchUrl(thing)
	}

	// neither
	return "", fmt.Errorf("not a valid format: \"%s\"", thing)
}

func extractVEVENT(calendar string) string {
	extracted := ""

	insideEvent := false
	for _, line := range strings.Split(calendar, "\n") {
		// note: `line` keeps its \r

		// start of event section
		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			insideEvent = true
		}

		// append line if inside event
		if insideEvent {
			extracted += line + "\n"
		}

		// end of event section
		if strings.HasPrefix(line, "END:VEVENT") {
			insideEvent = false
		}
	}

	return extracted
}

func fetchAndMerge(entry CalEntry) (string, error) {
	// header
	merged := "BEGIN:VCALENDAR\r\n"
	merged += "VERSION:2.0\r\n"
	merged += fmt.Sprintf("X-WR-CALNAME:%s\r\n", entry.Title)
	merged += fmt.Sprintf("PRODID:-//%s//NONSGML v1.0//EN\r\n", ProdID)

	for _, thing := range entry.Urls {
		content, err := fetch(thing)
		if err != nil {
			return "", err
		}

		// assume well-formatted RFC 5545
		merged += extractVEVENT(content)
	}

	// footer
	merged += "END:VCALENDAR\r\n"
	return merged, nil
}

func unite(calmap CalMap) func() {
	// closures are awesome!
	return func() {
		for calendar, entry := range calmap {
			// get merged calendar
			merged, err := fetchAndMerge(entry)
			if err != nil {
				log.Print(err)
				continue
			}

			// create directory if it doesn't exist
			err = os.MkdirAll(filepath.Dir(calendar), os.ModePerm)
			if err != nil {
				log.Print(err)
				continue
			}
			// write merged calendar
			err = os.WriteFile(calendar, []byte(merged), 0666)
			if err != nil {
				log.Print(err)
			}
		}
	}
}
