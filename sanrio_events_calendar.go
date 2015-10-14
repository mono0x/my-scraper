package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/gorilla/feeds"
	"github.com/mvdan/xurls"
	"io"
	"net/http"
	"strings"
	"time"
)

func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t'
}

func scanContentLines(data []byte, atEOF bool) (int, []byte, error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	var (
		advance = 0
		buffer  bytes.Buffer
	)
	for {
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			advance += i + 1
			buffer.Write(dropCR(data[0:i]))

			data = data[i+1:]
			if len(data) > 0 && isSpace(data[0]) {
				j := 1
				for ; j < len(data); j++ {
					if !isSpace(data[j]) {
						break
					}
				}
				advance += j
				data = data[j:]
			} else {
				break
			}
		} else if atEOF {
			advance += len(data)
			buffer.Write(dropCR(data))
		} else {
			return 0, nil, nil
		}
	}
	if advance > 0 {
		return advance, buffer.Bytes(), nil
	}
	return 0, nil, nil
}

type Event struct {
	Summary      string
	Description  string
	Location     string
	Uid          string
	Start        time.Time
	End          time.Time
	Created      time.Time
	LastModified time.Time
}

func parseTime(dateString string, tzId string) (time.Time, error) {
	if tzId == "" {
		location, err := time.LoadLocation("UTC")
		if err != nil {
			return time.Time{}, err
		}
		result, err := time.ParseInLocation("20060102T150405Z", dateString, location)
		if err != nil {
			return time.ParseInLocation("20060102", dateString, location)
		}
		return result, nil
	} else {
		location, err := time.LoadLocation(tzId)
		if err != nil {
			return time.Time{}, err
		}
		result, err := time.ParseInLocation("20060102T150405", dateString, location)
		if err != nil {
			return time.ParseInLocation("20060102", dateString, location)
		}
		return result, nil
	}
}

func unescape(data string) (string, error) {
	tmp := data
	tmp = strings.Replace(tmp, "\\n", "\n", -1)
	tmp = strings.Replace(tmp, "\\", "", -1)
	return tmp, nil
}

func parse(scanner *bufio.Scanner) ([]Event, error) {
	var events []Event

	inEvent := false

	var (
		summary      string
		description  string
		location     string
		uid          string
		start        time.Time
		end          time.Time
		created      time.Time
		lastModified time.Time
	)
	for scanner.Scan() {
		var (
			key        string
			value      string
			err        error
			properties = map[string]string{}
		)

		array := strings.SplitN(scanner.Text(), ":", 2)
		if len(array) != 2 {
			return nil, errors.New("Error parsing a content line")
		}
		value = array[1]
		array = strings.SplitN(array[0], ";", 2)
		key = array[0]
		for i := 1; i < len(array); i++ {
			a := strings.SplitN(array[i], "=", 2)
			if len(array) != 2 {
				return nil, errors.New("Error parsing a property")
			}
			properties[a[0]] = a[1]
		}

		switch {
		case key == "BEGIN" && value == "VEVENT":
			if inEvent {
				return nil, errors.New("Error parsing BEGIN:VEVENT")
			}
			inEvent = true

			summary = ""
			description = ""
			location = ""
			uid = ""
			start = time.Time{}
			end = time.Time{}
			created = time.Time{}
			lastModified = time.Time{}

		case key == "END" && value == "VEVENT":
			if !inEvent {
				return nil, errors.New("Error parsing END:VEVENT")
			}
			inEvent = false

			events = append(events, Event{
				Summary:      summary,
				Description:  description,
				Location:     location,
				Uid:          uid,
				Start:        start,
				End:          end,
				Created:      created,
				LastModified: lastModified,
			})

		case key == "SUMMARY":
			if !inEvent {
				break
			}
			summary, err = unescape(value)
			if err != nil {
				return nil, err
			}
		case key == "DESCRIPTION":
			if !inEvent {
				break
			}
			description, err = unescape(value)
			if err != nil {
				return nil, err
			}
		case key == "LOCATION":
			if !inEvent {
				break
			}
			location, err = unescape(value)
			if err != nil {
				return nil, err
			}
		case key == "UID":
			if !inEvent {
				break
			}
			uid = value
		case key == "DTSTART":
			if !inEvent {
				break
			}
			start, err = parseTime(value, properties["TZID"])
			if err != nil {
				return nil, err
			}
		case key == "DTEND":
			if !inEvent {
				break
			}
			end, err = parseTime(value, properties["TZID"])
			if err != nil {
				return nil, err
			}
		case key == "CREATED":
			if !inEvent {
				break
			}
			created, err = parseTime(value, properties["TZID"])
			if err != nil {
				return nil, err
			}
		case key == "LAST-MODIFIED":
			if !inEvent {
				break
			}
			lastModified, err = parseTime(value, properties["TZID"])
			if err != nil {
				return nil, err
			}
		}

	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	if inEvent {
		return nil, errors.New("END:VEVENT is missing")
	}

	return events, nil
}

func GetSanrioEventsCalendar() (*feeds.Feed, error) {
	resp, err := http.Get("https://www.google.com/calendar/ical/qsqrk2emvnnvu45debac9dugr8%40group.calendar.google.com/public/basic.ics")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return GetSanrioEventsCalendarFromReader(resp.Body)
}

func GetSanrioEventsCalendarFromReader(reader io.Reader) (*feeds.Feed, error) {
	scanner := bufio.NewScanner(reader)
	scanner.Split(scanContentLines)

	events, err := parse(scanner)
	if err != nil {
		return nil, err
	}

	var items []*feeds.Item
	items = make([]*feeds.Item, len(events))
	for i, event := range events {
		start := event.Start.Local()
		end := event.End.Local()
		var duration string
		if start.Format("20060102") == end.Format("20060102") {
			duration = start.Format("2006-01-02 (Mon) 15:04") + " - " + end.Format("15:04")
		} else {
			duration = start.Format("2006-01-02 (Mon) 15:04") + " - " + end.Format("2006-01-02 (Mon) 15:04")
		}
		items[i] = &feeds.Item{
			Title:       event.Summary,
			Description: fmt.Sprintf("Location: %s<br />Duration: %s<br /><br />Description: %s", event.Location, duration, event.Description),
			Link:        &feeds.Link{Href: xurls.Strict.FindString(event.Description)},
			Id:          event.Uid,
			Created:     event.Created,
			Updated:     event.LastModified,
		}
	}

	feed := &feeds.Feed{
		Title: "Sanrio Events",
		Link:  &feeds.Link{Href: "http://ameblo.jp/ohtaket/entry-12059393801.html"},
		Items: items,
	}
	return feed, nil
}
