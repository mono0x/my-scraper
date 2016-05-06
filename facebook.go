package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gorilla/feeds"
)

const FacebookApiEndpoint = "https://graph.facebook.com/v2.6/"

type FacebookPosts struct {
	Data []*FacebookPost `json:"data"`
}

// https://developers.facebook.com/docs/graph-api/reference/v2.6/post
type FacebookPost struct {
	Id          string           `json:"id"`
	CreatedTime string           `json:"created_time"`
	From        *FacebookProfile `json:"from"`
	Link        string           `json:"link"`
	Message     string           `json:"message"`
	Picture     string           `json:"picture"`
}

type FacebookProfile struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetPostsFromFacebook(userId string) (*FacebookPosts, error) {
	values := &url.Values{}
	values.Set("access_token", os.Getenv("FACEBOOK_ACCESS_TOKEN"))
	values.Set("fields", "created_time,from,link,message,picture")

	resp, err := http.Get(FacebookApiEndpoint + userId + "/posts?" + values.Encode())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	jsonData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var posts FacebookPosts
	if err := json.Unmarshal(jsonData, &posts); err != nil {
		return nil, err
	}
	return &posts, nil
}

func RenderFacebookFeed(posts *FacebookPosts, userId string) (*feeds.Feed, error) {
	messageReplacer := strings.NewReplacer("\n", "<br />")

	if len(posts.Data) == 0 {
		return nil, errors.New("no posts found")
	}

	items := make([]*feeds.Item, 0, len(posts.Data))
	for _, post := range posts.Data {
		created, err := time.Parse("2006-01-02T15:04:05-0700", post.CreatedTime)
		if err != nil {
			return nil, err
		}

		var title, description string
		if index := strings.Index(post.Message, "\n"); index >= 0 {
			title = post.Message[0:index]
			description = messageReplacer.Replace(post.Message)
			if post.Picture != "" {
				description += fmt.Sprintf(`<br /><img src="%s" />`, post.Picture)
			}
		} else {
			title = post.Message
		}

		items = append(items, &feeds.Item{
			Id:          post.Id,
			Author:      &feeds.Author{Name: post.From.Name},
			Title:       title,
			Description: description,
			Created:     created,
			Link:        &feeds.Link{Href: post.Link},
		})
	}

	feed := &feeds.Feed{
		Title: posts.Data[0].From.Name,
		Link:  &feeds.Link{Href: "https://www.facebook.com/" + userId},
		Items: items,
	}
	return feed, nil
}
