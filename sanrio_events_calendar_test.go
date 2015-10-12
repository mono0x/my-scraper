package main

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetSanrioEventsCalendarFromReader(t *testing.T) {
	f, err := os.Open("data/sanrio_events_calendar.ics")
	defer f.Close()
	feed, err := GetSanrioEventsCalendarFromReader(f)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 181, len(feed.Items))
}
