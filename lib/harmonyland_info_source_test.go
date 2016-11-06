package scraper

import (
	"os"
	"testing"
	"time"

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

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 13, len(feed.Items))
	assert.Equal(t, "とびっきりかわいいクリスマスを！ 「ハッピークリスマス」", feed.Items[0].Title)
	assert.Equal(t, "http://www.harmonyland.jp/event/xmas/index.html", feed.Items[0].Link.Href)
	assert.WithinDuration(t, time.Date(2016, 11, 6, 0, 0, 0, 0, loc), feed.Items[0].Created, 0)
}
