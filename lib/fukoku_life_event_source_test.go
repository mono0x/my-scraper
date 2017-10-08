package scraper

import (
	"bufio"
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestFukokuLifeEventSource(t *testing.T) {
	f, err := os.Open("testdata/act.fukoku-life.co.jp/event/index.php")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	source := NewFukokuLifeEventSource()
	feed, err := source.ScrapeFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 6, len(feed.Items))
	assert.Equal(t, "https://act.fukoku-life.co.jp/event/detail.php?id=415", feed.Items[0].Link.Href)
	assert.Equal(t, "フコク生命 ピラティス体験イベント（in京浜）", feed.Items[0].Title)
	assert.Contains(t, feed.Items[0].Description, "2016年10月12日（水）")
	assert.Contains(t, feed.Items[0].Description, "川崎市川崎区貝塚1-1-3 川崎フコク生命ビル４Ｆ 特設会場")
	assert.Equal(t, "https://act.fukoku-life.co.jp/event/detail.php?id=435", feed.Items[5].Link.Href)
}
