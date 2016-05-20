package main

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/calendar/v3"
)

func TestGetSanrioEventsCalendarFromReader(t *testing.T) {
	jsonData, err := ioutil.ReadFile("data/sanrio_events_calendar.json")
	if err != nil {
		t.Fatal(err)
	}

	var events calendar.Events
	if err := json.Unmarshal(jsonData, &events); err != nil {
		t.Fatal(err)
	}

	source := NewSanrioEventsCalendarGoogleCalendarSource()
	feed, err := source.Render(&events)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 199, len(feed.Items))
}
