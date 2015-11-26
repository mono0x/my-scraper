package main

import (
	"encoding/xml"
	"github.com/gorilla/feeds"
	"github.com/kennygrant/sanitize"
	"golang.org/x/tools/blog/atom"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
)

const (
	SanrioAlertsUrl = "http://scraper.mono0x.net/sanrio-alerts"
)

func GetSanrioAlerts() (*feeds.Feed, error) {

	urls := []string{
		"https://www.google.com/alerts/feeds/17240735437045332758/1863509270421926440",
		"https://www.google.com/alerts/feeds/17240735437045332758/1863509270421929515",
		"https://www.google.com/alerts/feeds/17240735437045332758/2414106377807123167",
		"https://www.google.com/alerts/feeds/17240735437045332758/2414106377807124539",
		"https://www.google.com/alerts/feeds/17240735437045332758/2414106377807125523",
		"https://www.google.com/alerts/feeds/17240735437045332758/2636887480119177525",
		"https://www.google.com/alerts/feeds/17240735437045332758/2636887480119178148",
		"https://www.google.com/alerts/feeds/17240735437045332758/2636887480119179073",
	}

	atomChan := make(chan *atom.Feed)
	quitChan := make(chan bool)
	errChan := make(chan error)

	go func() {
		var wg sync.WaitGroup
		for _, url := range urls {
			wg.Add(1)
			go func() {
				defer wg.Done()

				resp, err := http.Get(url)
				if err != nil {
					errChan <- err
					return
				}
				defer resp.Body.Close()

				body, err := ioutil.ReadAll(resp.Body)
				if err != nil {
					errChan <- err
					return
				}

				var atom atom.Feed
				err = xml.Unmarshal(body, &atom)
				if err != nil {
					errChan <- err
					return
				}

				atomChan <- &atom
			}()
		}
		wg.Wait()

		quitChan <- true
	}()

	var atoms []*atom.Feed

loop:
	for {
		select {
		case atom := <-atomChan:
			atoms = append(atoms, atom)

		case <-quitChan:
			break loop

		case err := <-errChan:
			return nil, err
		}
	}

	return GetSanrioAlertsFromAtom(atoms)
}

func GetSanrioAlertsFromAtom(atoms []*atom.Feed) (*feeds.Feed, error) {
	var items []*feeds.Item

	urls := map[string]bool{}

	for _, atom := range atoms {
		for _, entry := range atom.Entry {
			if len(entry.Link) == 0 {
				continue
			}

			href, err := url.Parse(entry.Link[0].Href)
			if err != nil {
				return nil, err
			}

			url := href.Query().Get("url")
			if url == "" {
				return nil, err
			}

			if _, ok := urls[url]; ok {
				continue
			}
			urls[url] = true

			title := sanitize.HTML(entry.Title)
			content := sanitize.HTML(entry.Content.Body)

			published, err := time.Parse(time.RFC3339, string(entry.Published))
			if err != nil {
				return nil, err
			}
			updated, err := time.Parse(time.RFC3339, string(entry.Updated))
			if err != nil {
				return nil, err
			}

			items = append(items, &feeds.Item{
				Title:       title,
				Description: content,
				Id:          url,
				Link:        &feeds.Link{Href: url},
				Created:     published,
				Updated:     updated,
			})
		}
	}

	feed := &feeds.Feed{
		Title: "Sanrio Alerts",
		Link:  &feeds.Link{Href: SanrioAlertsUrl},
		Items: items,
	}

	return feed, nil
}
