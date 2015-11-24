package main

import (
	"github.com/gorilla/feeds"
	"google.golang.org/api/calendar/v3"
)

const (
	SanrioEventsCalendarId  = "qsqrk2emvnnvu45debac9dugr8@group.calendar.google.com"
	SanrioEventsCalendarUrl = "https://calendar.google.com/calendar/embed?src=qsqrk2emvnnvu45debac9dugr8@group.calendar.google.com"
)

func GetSanrioEventsCalendar() (*feeds.Feed, error) {
	events, err := GetEventsFromGoogleCalendar(SanrioEventsCalendarId)
	if err != nil {
		return nil, err
	}
	return GetSanrioEventsCalendarFromEvents(events)
}

func GetSanrioEventsCalendarFromEvents(events *calendar.Events) (*feeds.Feed, error) {
	return RenderGoogleCalendarFeed(events, SanrioEventsCalendarUrl)
}
