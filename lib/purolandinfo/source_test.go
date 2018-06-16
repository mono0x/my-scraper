package purolandinfo

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
