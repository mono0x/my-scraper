package scraper

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestInstagramSource(t *testing.T) {
	f, err := os.Open("testdata/www.instagram.com/fukkachan628/index.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	source := NewInstagramSource("fukkachan628")
	feed, err := source.ScrapeFromReader(f)
	if err != nil {
		t.Fatal(err)
	}

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "ふっかちゃん【公式】", feed.Title)
	assert.Equal(t, "https://www.instagram.com/fukkachan628/", feed.Link.Href)
	assert.Equal(t, 12, len(feed.Items))
	assert.Equal(t, "今年の漢字は「金」！Y(o≧ω≦o)Y\n #今年の漢字 #金 #みんなからもらった金メダル #うれしす #ふっかちゃん #過去のお写真", feed.Items[0].Title)
	assert.Equal(t, "http://www.instagram.com/p/BN59EPyhA09/", feed.Items[0].Link.Href)
	assert.WithinDuration(t, time.Date(2016, 12, 12, 5, 34, 41, 0, loc), feed.Items[0].Created, 0)
}
