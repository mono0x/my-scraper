package sanrionewsrelease

import (
	"os"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/stretchr/testify/assert"
)

var _ scraper.Source = (*SanrioNewsReleaseSource)(nil)

func TestSource(t *testing.T) {
	f, err := os.Open("testdata/www.sanrio.co.jp/corporate/release/index.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Fatal(err)
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	source := NewSource()
	feed, err := source.ScrapeFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 51, len(feed.Items))
	assert.Equal(t, "ぐでぐでやる気のない「ぐでたま」のイベント九州初上陸！ 夏休み企画 「ぐでたま in ふくおか」 7月21日(木)〜 福岡パルコ & sanrio vivitix 天神地下街店にて開催 (PDF)", feed.Items[0].Title)
	assert.Equal(t, "http://www.sanrio.co.jp/wp-content/uploads/2015/05/20160708-1.pdf", feed.Items[0].Link.Href)
	assert.WithinDuration(t, time.Date(2016, 7, 8, 0, 0, 0, 0, loc), feed.Items[0].Created, 0)
	assert.Equal(t, "2016年バレンタイン向けスペシャルギフト「GODIVA &ハローキティ」・「GODIVA &マイメロディ」1月6日（水）よりサンリオ限定販売", feed.Items[50].Title)
	assert.Equal(t, "http://www.sanrio.co.jp/corporate/release/y2016/d0106/", feed.Items[50].Link.Href)
	assert.WithinDuration(t, time.Date(2016, 1, 6, 0, 0, 0, 0, loc), feed.Items[50].Created, 0)
}
