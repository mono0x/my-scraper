package main

import (
	"github.com/gorilla/feeds"
	"google.golang.org/api/calendar/v3"
)

const (
	GotouchiCharaCalendarId  = "gnr0r3kevuuv3j0q6q25gj4hks@group.calendar.google.com"
	GotouchiCharaCalendarUrl = "https://calendar.google.com/calendar/embed?src=gnr0r3kevuuv3j0q6q25gj4hks@group.calendar.google.com"
)

func GetGotouchiCharaCalendar() (*feeds.Feed, error) {
	events, err := GetEventsFromGoogleCalendar(GotouchiCharaCalendarId)
	if err != nil {
		return nil, err
	}
	return GetSanrioEventsCalendarFromEvents(events)
}

func GetGotouchiCharaCalendarFromEvents(events *calendar.Events) (*feeds.Feed, error) {
	return RenderGoogleCalendarFeed(events, GotouchiCharaCalendarUrl)
}
