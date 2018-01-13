package charactershow

import (
	"os"
	"testing"

	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
)

func TestSource(t *testing.T) {
	f, err := os.Open("testdata/charactershow.jp/index.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	doc, err := goquery.NewDocumentFromReader(f)
	if err != nil {
		t.Fatal(err)
	}
	source := NewSource()
	feed, err := source.ScrapeFromDocument(doc)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 115, len(feed.Items))
}
