package memoirsofshibasakisaki

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSource(t *testing.T) {
	f, err := os.Open("testdata/shibasakisaki.web.fc2.com/index.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	source := NewSource()
	feed, err := source.ScrapeFromReader(f)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 2, len(feed.Items))
	assert.Equal(t, "5/14,15　調布観光フェスティバル", feed.Items[0].Title)
	assert.Equal(t, "http://www.csa.gr.jp/enjoy/bussanten.html", feed.Items[0].Link.Href)
}
