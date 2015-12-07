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
	err = json.Unmarshal(jsonData, &events)
	if err != nil {
		t.Fatal(err)
	}

	feed, err := GetSanrioEventsCalendarFromEvents(&events)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 199, len(feed.Items))
}
