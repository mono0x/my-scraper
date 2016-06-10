package scraper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPurolandInfoFromDocument(t *testing.T) {
	f, err := os.Open("data/www.puroland.jp/api/live/get_information/index.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	source := NewPurolandInfoSource()
	feed, err := source.ScrapeFromReader(f)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 5, len(feed.Items))
}
