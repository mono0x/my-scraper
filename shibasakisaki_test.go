package main

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShibasakisakiFromReader(t *testing.T) {
	f, err := os.Open("data/shibasakisaki.web.fc2.com/index.html")
	defer f.Close()
	if err != nil {
		t.Fatal(err)
	}
	feed, err := GetShibasakisakiFromReader(bufio.NewReader(f))
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, 5, len(feed.Items))
}
