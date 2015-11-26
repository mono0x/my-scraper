package main

import (
	"errors"
	"fmt"
	"github.com/gorilla/feeds"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"html"
	"io/ioutil"
	"strings"
	"time"
)

func GetEventsFromGoogleCalendar(calendarId string) (*calendar.Events, error) {
	json, err := ioutil.ReadFile("google_client_credentials.json")
	if err != nil {
		return nil, err
	}

	config, err := google.JWTConfigFromJSON(json, calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, err
	}

	client := config.Client(oauth2.NoContext)

	service, err := calendar.New(client)
	if err != nil {
		return nil, err
	}

	timeMin := time.Now().AddDate(0, -3, 0).Format(time.RFC3339)

	events, err := service.Events.List(calendarId).MaxResults(2500).OrderBy("updated").SingleEvents(true).TimeMin(timeMin).Do()
	if err != nil {
		return nil, err
	}

	items := events.Items
	for pageToken := events.NextPageToken; events.NextPageToken != ""; {
		events, err := service.Events.List(calendarId).PageToken(pageToken).Do()
		if err != nil {
			return nil, err
		}
		items = append(items, events.Items...)
		pageToken = events.NextPageToken
	}
	events.Items = items
	return events, nil
}

func RenderGoogleCalendarFeed(events *calendar.Events, href string) (*feeds.Feed, error) {
	descriptionReplacer := strings.NewReplacer("\n", "<br />")

	var items []*feeds.Item
	items = make([]*feeds.Item, len(events.Items))
	for i, event := range events.Items {
		created, err := time.Parse(time.RFC3339, event.Created)
		if err != nil {
			return nil, err
		}
		updated, err := time.Parse(time.RFC3339, event.Updated)
		if err != nil {
			return nil, err
		}

		var duration string

		switch {
		case event.Start.Date != "" && event.End.Date != "":
			startLoc, err := time.LoadLocation(event.Start.TimeZone)
			if err != nil {
				return nil, err
			}
			start, err := time.ParseInLocation("2006-01-02", event.Start.Date, startLoc)
			if err != nil {
				return nil, err
			}
			endLoc, err := time.LoadLocation(event.Start.TimeZone)
			if err != nil {
				return nil, err
			}
			end, err := time.ParseInLocation("2006-01-02", event.End.Date, endLoc)
			if err != nil {
				return nil, err
			}

			if start.Format("2006-01-02") == end.Format("2006-01-02") {
				duration = start.Format("2006-01-02 (Mon)")
			} else {
				duration = start.Format("2006-01-02 (Mon)") + " - " + end.Format("2006-01-02 (Mon)")
			}

		case event.Start.DateTime != "" && event.End.DateTime != "":
			startLoc, err := time.LoadLocation(event.Start.TimeZone)
			if err != nil {
				return nil, err
			}
			start, err := time.ParseInLocation(time.RFC3339, event.Start.DateTime, startLoc)
			if err != nil {
				return nil, err
			}
			endLoc, err := time.LoadLocation(event.Start.TimeZone)
			if err != nil {
				return nil, err
			}
			end, err := time.ParseInLocation(time.RFC3339, event.End.DateTime, endLoc)
			if err != nil {
				return nil, err
			}

			if start.Format("2006-01-02") == end.Format("2006-01-02") {
				duration = start.Format("2006-01-02 (Mon) 15:04") + " - " + end.Format("15:04")
			} else {
				duration = start.Format("2006-01-02 (Mon) 15:04") + " - " + end.Format("2006-01-02 (Mon) 15:04")
			}

		default:
			return nil, errors.New("must not happen")
		}

		var description string
		if event.Location != "" {
			description += fmt.Sprintf("Location: %s<br />", html.EscapeString(event.Location))
		}
		description += fmt.Sprintf("Duration: %s<br /><br />", html.EscapeString(duration))
		description += descriptionReplacer.Replace(html.EscapeString(event.Description))

		items[i] = &feeds.Item{
			Id:          event.Etag,
			Title:       event.Summary,
			Description: description,
			Link:        &feeds.Link{Href: event.HtmlLink},
			Author:      &feeds.Author{Name: event.Creator.DisplayName, Email: event.Creator.Email},
			Created:     created,
			Updated:     updated,
		}
	}

	updated, err := time.Parse(time.RFC3339, events.Updated)
	if err != nil {
		return nil, err
	}

	feed := &feeds.Feed{
		Id:          events.Etag,
		Title:       events.Summary,
		Description: events.Description,
		Link:        &feeds.Link{Href: href},
		Updated:     updated,
		Items:       items,
	}
	return feed, nil
}
