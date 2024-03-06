package googlecalendar

import (
	"encoding/json"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/api/calendar/v3"
)

func TestNewSource(t *testing.T) {
	source := NewSource(http.DefaultClient)
	assert.Equal(t, http.DefaultClient, source.httpClient)
}

func TestRender(t *testing.T) {
	file, err := os.Open("testdata/sanrio_events_calendar.json")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()

	var events calendar.Events
	if err := json.NewDecoder(file).Decode(&events); err != nil {
		t.Fatal(err)
	}

	source := NewSource(http.DefaultClient)
	feed, err := source.render(&events, "qsqrk2emvnnvu45debac9dugr8@group.calendar.google.com")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 199, len(feed.Items))
}
