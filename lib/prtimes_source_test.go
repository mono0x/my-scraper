package scraper

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPRTimesFromReader(t *testing.T) {
	f, err := os.Open("data/prtimes.jp/topics/keywords/sanrio.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	source := NewPRTimesSource()
	feed, err := source.ScrapeFromReader(f)
	if err != nil {
		t.Fatal(err)
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 20, len(feed.Items))
	assert.Equal(t, "ディズニー、スヌーピー、サンリオ、ポケモンなどの人気キャラクターを取り入れた最旬トレンドファッション“キャラディネート”を一冊に！", feed.Items[0].Title)
	assert.Equal(t, "/main/html/rd/p/000000431.000005069.html", feed.Items[0].Link.Href)
	assert.WithinDuration(t, time.Date(2016, 9, 16, 19, 10, 39, 0, loc), feed.Items[0].Created, 0)
}
