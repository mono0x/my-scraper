package harmonylandinfo

import (
	"os"
	"testing"

	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/stretchr/testify/assert"
)

var _ scraper.Source = (*HarmonylandInfoSource)(nil)

func TestSource(t *testing.T) {
	f, err := os.Open("testdata/www.harmonyland.jp/welcome.html")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	source := NewSource()

	feed, err := source.ScrapeFromReader(f)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, 13, len(feed.Items))
	assert.Equal(t, "2017年はシナモン15周年！                     ハーモニーランドはシナモンイベントがいっぱい♪【6/2～7/14】", feed.Items[0].Title)
	assert.Equal(t, "http://www.harmonyland.jp/event/rain/index.html", feed.Items[0].Link.Href)
	assert.NotEqual(t, "", feed.Items[0].Id)
}
