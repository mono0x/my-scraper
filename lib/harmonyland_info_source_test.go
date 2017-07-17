package scraper

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHarmonylandInfoSource(t *testing.T) {
	f, err := os.Open("data/www.harmonyland.jp/welcome.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	source := NewHarmonylandInfoSource()

	feed, err := source.ScrapeFromReader(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 13, len(feed.Items))
	assert.Equal(t, "2017年はシナモン15周年！                     ハーモニーランドはシナモンイベントがいっぱい♪【6/2～7/14】", feed.Items[0].Title)
	assert.Equal(t, "http://www.harmonyland.jp/event/rain/index.html", feed.Items[0].Link.Href)
	assert.NotEqual(t, "", feed.Items[0].Id)
}
