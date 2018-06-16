package yuyakekoyakenews

import (
	"os"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/stretchr/testify/assert"
)

var _ scraper.Source = (*YuyakekoyakeNewsSource)(nil)

func TestSource(t *testing.T) {
	f, err := os.Open("testdata/yuyakekoyake.jp/news/index.php")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Fatal(err)
	}

	source := NewSource()
	feed, err := source.ScrapeFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 15, len(feed.Items))
	assert.Equal(t, "http://yuyakekoyake.jp/news/news_detail.php?id=news59f16132092f0", feed.Items[0].Link.Href)
	assert.Equal(t, "イルミネーションの実施につきまして", feed.Items[0].Title)
	assert.WithinDuration(t, time.Date(2017, 10, 26, 0, 0, 0, 0, loc), feed.Items[0].Created, 0)
}
