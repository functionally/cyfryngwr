package rss

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"

	"cwtch.im/cwtch/model"
	"git.openprivacy.ca/sarah/cwtchbot"
)

func Message(cwtchbot *bot.CwtchBot, url string) (string, error) {
	feed, err := fetchRSS(url)
	if err != nil {
		log.Fatalf("Error fetching RSS: %v", err)
	}
	s := fmt.Sprintf("%s\n%s", feed.Title, feed.Items[0].Title)
	reply := string(cwtchbot.PackMessage(model.OverlayChat, s))
	return reply, nil

}

func fetchRSS(url string) (*gofeed.Feed, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 response: %d", resp.StatusCode)
	}
	parser := gofeed.NewParser()
	feed, err := parser.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to parse feed: %w", err)
	}
	return feed, nil
}
