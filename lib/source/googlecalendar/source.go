package googlecalendar

import (
	"context"
	"errors"
	"fmt"
	"html"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/feeds"
	scraper "github.com/mono0x/my-scraper/lib"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

const (
	prefix = `https://calendar.google.com/calendar/embed?src=`
)

type source struct {
	httpClient *http.Client
	calendarID string
}

var _ scraper.Source = (*source)(nil)

var (
	descriptionReplacer = strings.NewReplacer("\n", "<br />")
)

func NewSource(c *http.Client, calendarID string) *source {
	return &source{
		httpClient: c,
		calendarID: calendarID,
	}
}

func (s *source) Scrape() (*feeds.Feed, error) {
	events, err := s.fetch()
	if err != nil {
		return nil, err
	}
	return s.render(events)
}

func (s *source) fetch() (*calendar.Events, error) {
	config, err := google.JWTConfigFromJSON(([]byte)(os.Getenv("GOOGLE_CLIENT_CREDENTIALS")), calendar.CalendarReadonlyScope)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, oauth2.HTTPClient, s.httpClient)

	client := config.Client(ctx)

	service, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	timeMin := time.Now().AddDate(0, -3, 0).Format(time.RFC3339)

	events, err := service.Events.List(s.calendarID).MaxResults(2500).OrderBy("updated").SingleEvents(true).TimeMin(timeMin).Do()
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	items := events.Items
	for pageToken := events.NextPageToken; events.NextPageToken != ""; {
		events, err := service.Events.List(s.calendarID).PageToken(pageToken).Do()
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		items = append(items, events.Items...)
		pageToken = events.NextPageToken
	}
	events.Items = items
	return events, nil
}

func (s *source) render(events *calendar.Events) (*feeds.Feed, error) {
	loc, err := time.LoadLocation(events.TimeZone)
	if err != nil {
		return nil, err
	}

	items := make([]*feeds.Item, 0, len(events.Items))
	for _, event := range events.Items {
		if event.Visibility == "private" {
			continue
		}
		if event.Status == "cancelled" {
			continue
		}

		created, err := time.Parse(time.RFC3339, event.Created)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
		updated, err := time.Parse(time.RFC3339, event.Updated)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		var timeZone string
		if event.Start.TimeZone != "" {
			timeZone = event.Start.TimeZone
		} else if events.TimeZone != "" {
			timeZone = events.TimeZone
		}

		link := event.HtmlLink
		if timeZone != "" {
			u, err := url.Parse(link)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}
			query := u.Query()
			query.Set("ctz", timeZone)
			u.RawQuery = query.Encode()
			link = u.String()
		}

		var startLoc *time.Location
		if event.Start.TimeZone != "" {
			var err error
			startLoc, err = time.LoadLocation(event.Start.TimeZone)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}
		} else {
			startLoc = loc
		}

		var endLoc *time.Location
		if event.End.TimeZone != "" {
			var err error
			endLoc, err = time.LoadLocation(event.End.TimeZone)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}
		} else {
			endLoc = loc
		}

		var duration string

		switch {
		case event.Start.Date != "" && event.End.Date != "":
			start, err := time.ParseInLocation("2006-01-02", event.Start.Date, startLoc)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}
			end, err := time.ParseInLocation("2006-01-02", event.End.Date, endLoc)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}
			end = end.AddDate(0, 0, -1)

			if start.Format("2006-01-02") == end.Format("2006-01-02") {
				duration = start.Format("2006-01-02 (Mon)")
			} else {
				duration = start.Format("2006-01-02 (Mon)") + " - " + end.Format("2006-01-02 (Mon)")
			}

		case event.Start.DateTime != "" && event.End.DateTime != "":
			start, err := time.ParseInLocation(time.RFC3339, event.Start.DateTime, startLoc)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
			}
			end, err := time.ParseInLocation(time.RFC3339, event.End.DateTime, endLoc)
			if err != nil {
				return nil, fmt.Errorf("%w", err)
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

		items = append(items, &feeds.Item{
			Id:          event.Id,
			Title:       event.Summary,
			Description: description,
			Link:        &feeds.Link{Href: link},
			Author:      &feeds.Author{Name: event.Creator.DisplayName, Email: event.Creator.Email},
			Created:     created,
			Updated:     updated,
		})
	}

	updated, err := time.Parse(time.RFC3339, events.Updated)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	feed := &feeds.Feed{
		Title:       events.Summary,
		Description: events.Description,
		Link:        &feeds.Link{Href: prefix + s.calendarID},
		Updated:     updated,
		Items:       items,
	}
	return feed, nil
}
