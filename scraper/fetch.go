package scraper

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

var ErrNotFound = errors.New("not found")

func Fetch(ctx context.Context, client *http.Client, url string) (io.ReadCloser, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		res.Body.Close()
		if res.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return res.Body, nil
}
