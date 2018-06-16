package kittychaninfo

import (
	"os"
	"testing"
	"time"

	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/stretchr/testify/assert"
)

var _ scraper.Source = (*KittychanInfoSource)(nil)

func TestSource(t *testing.T) {
	f, err := os.Open("testdata/www.kittychan.info/information.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	source := NewSource()
	feed, err := source.ScrapeFromReader(f)
	if err != nil {
		t.Fatal(err)
	}

	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 100, len(feed.Items))
	assert.Equal(t, "多摩センターイルミネーション\u3000今年も開催～キティちゃんのイルミネーション＆オープニングセレモニーでキティちゃんのショー開催＆キティちゃん達のパレードも～", feed.Items[0].Title)
	assert.Equal(t, "http://www.tamacenter-cm.com/illumi/", feed.Items[0].Link.Href)
	assert.WithinDuration(t, time.Date(2016, 11, 6, 0, 0, 0, 0, loc), feed.Items[0].Created, 0)

	assert.Equal(t, "「Hello Kitty Japanダイバーシティ東京店」　リニューアルオープン", feed.Items[99].Title)
	assert.Equal(t, "http://www.sanrio.co.jp/news/kt-gu-hkj-odaiba-renewal-20160921/", feed.Items[99].Link.Href)
	assert.WithinDuration(t, time.Date(2016, 9, 21, 0, 0, 0, 0, loc), feed.Items[99].Created, 0)
}
