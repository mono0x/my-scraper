package scraper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKittychanInfoSource(t *testing.T) {
	f, err := os.Open("data/www.kittychan.info/information.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	source := NewKittychanInfoSource()
	feed, err := source.ScrapeFromReader(f)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 100, len(feed.Items))
}
