package purolandinfo

import (
	"os"
	"testing"

	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/stretchr/testify/assert"
)

var _ scraper.Source = (*PurolandInfoSource)(nil)

func TestSource(t *testing.T) {
	f, err := os.Open("testdata/www.puroland.jp/api/live/get_information/index.json")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()

	source := NewSource()
	feed, err := source.ScrapeFromReader(f)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 5, len(feed.Items))
}
