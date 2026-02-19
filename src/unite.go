package main

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

const DEFAULT_TITLE string = "Calunite Calendar"

func fetchFile(name string) (string, error) {
	slog.Debug(fmt.Sprintf("Fetching file %v", name))
	content, err := os.ReadFile(name)
	return string(content), err
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}

func fetchUrl(url string) (string, error) {
	slog.Debug(fmt.Sprintf("Fetching URL %v", url))

	// rewrite url, the webcal protocol is just https in disguise
	if strings.HasPrefix(url, "webcal") {
		url = strings.Replace(url, "webcal", "https", 1)
	}

	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode/100 != 2 {
		return "", fmt.Errorf("Got status code %d for %s", resp.StatusCode, url)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func fetch(thing string) (string, error) {
	if IsUrl(thing) {
		return fetchUrl(thing)
	}

	content, err := fetchFile(thing)
	if err == nil {
		// file exists and was read successfully
		return content, nil
	}

	// neither
	return "", fmt.Errorf("Error reading \"%s\" - %s", thing, err.Error())
}

func (calEntry *CalEntry) formatIfSUMMARY(originalLine string) string {
	if calEntry.EventFormat == nil {
		return originalLine
	}

	// "SUMMARY:title: potential colons" -> ["SUMMARY", "title: potential colons"]
	split := strings.SplitN(originalLine, ":", 2)
	if split[0] == "SUMMARY" {
		return split[0] + ":" + fmt.Sprintf(*calEntry.EventFormat, split[1])
	}

	return originalLine
}

func (calEntry *CalEntry) extractVEVENT(calendar string) string {
	var sbExtracted strings.Builder

	insideEvent := false
	// filter out CR here, allows us to parse LF-only calendars as well
	withoutCR := strings.ReplaceAll(calendar, "\r", "")
	for line := range strings.SplitSeq(withoutCR, "\n") {
		// start of event section
		if strings.HasPrefix(line, "BEGIN:VEVENT") {
			insideEvent = true
		}

		// append line if inside event
		if insideEvent {
			// custom event formatting
			formatted := calEntry.formatIfSUMMARY(line)

			sbExtracted.WriteString(formatted)
			sbExtracted.WriteString("\r\n")
		}

		// end of event section
		if strings.HasPrefix(line, "END:VEVENT") {
			insideEvent = false
		}
	}

	return sbExtracted.String()
}

func (entry *CalEntry) fetchAndMerge() (string, error) {
	var sbMerged strings.Builder

	// header
	sbMerged.WriteString("BEGIN:VCALENDAR\r\n")
	sbMerged.WriteString("VERSION:2.0\r\n")

	title := DEFAULT_TITLE
	if entry.Title != nil {
		title = *entry.Title
	}
	sbMerged.WriteString("X-WR-CALNAME:")
	sbMerged.WriteString(title)
	sbMerged.WriteString("\r\n")

	sbMerged.WriteString("PRODID:-//")
	sbMerged.WriteString(ProdID)
	sbMerged.WriteString("//NONSGML v1.0//EN\r\n")

	for _, thing := range entry.Urls {
		content, err := fetch(thing)
		if err != nil {
			return "", err
		}

		// assume well-formatted RFC 5545
		sbMerged.WriteString(entry.extractVEVENT(content))
	}

	// footer
	sbMerged.WriteString("END:VCALENDAR\r\n")
	return sbMerged.String(), nil
}

func unite(caldata CalData) func() {
	// closures are awesome!
	return func() {
		for _, datum := range caldata {
			slog.Debug(fmt.Sprintf("Creating calendar %s", datum.File))

			// get merged calendar
			merged, err := datum.Entry.fetchAndMerge()
			if err != nil {
				slog.Warn(fmt.Sprintf("Failed to merge calendar %s", datum.File), ErrAttr(err))
				continue
			}

			// create directory if it doesn't exist
			err = os.MkdirAll(filepath.Dir(datum.File), os.ModePerm)
			if err != nil {
				slog.Warn(fmt.Sprintf("Could not create directories for %s", datum.File), ErrAttr(err))
				continue
			}
			// write merged calendar
			err = os.WriteFile(datum.File, []byte(merged), 0666)
			if err != nil {
				slog.Warn(fmt.Sprintf("Could not write %s", datum.File), ErrAttr(err))
			}

			slog.Debug(fmt.Sprintf("Created calendar %s", datum.File))
		}
	}
}
