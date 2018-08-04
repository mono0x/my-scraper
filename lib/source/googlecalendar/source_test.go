package googlecalendar

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/mono0x/my-scraper/lib"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/calendar/v3"
)

var _ scraper.Source = (*GoogleCalendarSource)(nil)

func TestSource(t *testing.T) {
	file, err := os.Open("testdata/sanrio_events_calendar.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var events calendar.Events
	if err := json.NewDecoder(file).Decode(&events); err != nil {
		t.Fatal(err)
	}

	source := NewSource("qsqrk2emvnnvu45debac9dugr8@group.calendar.google.com")
	feed, err := source.Render(&events)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 199, len(feed.Items))
}
