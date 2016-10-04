package scraper

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetPRTimesFromReader(t *testing.T) {
	f, err := os.Open("data/prtimes.jp/main/action.php")
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

	assert.Equal(t, 40, len(feed.Items))
	assert.Equal(t, "今年のテーマは“ギフト”！サンリオキャラクターが贈る冬のイベント「ピューロウィンターギフト」開催決定！", feed.Items[0].Title)
	assert.Equal(t, "http://prtimes.jp/main/html/rd/p/000000085.000007643.html", feed.Items[0].Link.Href)
	assert.Equal(t, "株式会社サンリオエンターテイメント", feed.Items[0].Author.Name)
	assert.WithinDuration(t, time.Date(2016, 9, 29, 14, 3, 3, 0, loc), feed.Items[0].Created, 0)
}