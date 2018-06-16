package twitter

import (
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/ChimeraCoder/anaconda"
	scraper "github.com/mono0x/my-scraper/lib"
	"github.com/stretchr/testify/assert"
)

var _ scraper.Source = (*TwitterSource)(nil)

func TestSource(t *testing.T) {
	jsonData, err := ioutil.ReadFile("testdata/api.twitter.com/1.1/statuses/user_timeline.json")
	if err != nil {
		t.Fatal(err)
	}

	var timeline []anaconda.Tweet
	if err := json.Unmarshal(jsonData, &timeline); err != nil {
		t.Fatal(err)
	}

	source := NewSource(725638238291943424)
	feed, err := source.Render(timeline)
	if err != nil {
		t.Fatal(err)
	}

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, "Sanrio Events (@sanrio_events)", feed.Title)
	assert.Equal(t, 20, len(feed.Items))
	assert.Equal(t, "12/24 KT 藤崎本館 (藤崎本館): https://t.co/Gh7HkC6yM2", feed.Items[0].Title)
	assert.Equal(t, "https://twitter.com/sanrio_events/status/805532146110582784", feed.Items[0].Link.Href)
	assert.WithinDuration(t, time.Date(2016, 12, 4, 22, 0, 4, 0, loc), feed.Items[0].Created, 0)
}
