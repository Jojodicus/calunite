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

const DEFAULT_TITLE string = "Calunite Calendar"

func fetchFile(name string) (string, error) {
	content, err := os.ReadFile(name)
	return string(content), err
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func fetchUrl(url string) (string, error) {
	// rewrite url, the webcal protocol is just https in disguise
	if strings.HasPrefix(url, "webcal") {
		url = strings.Replace(url, "webcal", "https", 1)
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func fetch(thing string) (string, error) {
	content, err := fetchFile(thing)
	if err == nil {
		// file exists and was read successfully
		return content, nil
	}
	if IsUrl(thing) {
		return fetchUrl(thing)
	}

	// neither
	return "", fmt.Errorf("error reading \"%s\" - %s", thing, err.Error())
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

	title := DEFAULT_TITLE
	if entry.Title != nil {
		title = *entry.Title
	}
	merged += fmt.Sprintf("X-WR-CALNAME:%s\r\n", title)
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

func unite(caldata CalData) func() {
	// closures are awesome!
	return func() {
		for _, datum := range caldata {
			// get merged calendar
			merged, err := fetchAndMerge(datum.Entry)
			if err != nil {
				log.Print(err)
				continue
			}

			// create directory if it doesn't exist
			err = os.MkdirAll(filepath.Dir(datum.File), os.ModePerm)
			if err != nil {
				log.Print(err)
				continue
			}
			// write merged calendar
			err = os.WriteFile(datum.File, []byte(merged), 0666)
			if err != nil {
				log.Print(err)
			}
		}
	}
}
