package rss

import (
	"bytes"
	"fmt"
	"html"
	"net/http"
	"regexp"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
)

func fetchRSS(url string) (*gofeed.Feed, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch feed:\n%w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Received non-200 response: %d", resp.StatusCode)
	}
	parser := gofeed.NewParser()
	feed, err := parser.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse feed:\n%w", err)
	}
	return feed, nil
}

func formatItem(feed *gofeed.Feed, item *gofeed.Item, source bool, title bool, brief bool, url bool) string {
	re := regexp.MustCompile(`\n{3,}`)
	p := bluemonday.StrictPolicy()
	var buffer bytes.Buffer
	if source {
		buffer.WriteString(fmt.Sprintf("*%s*\n\n", feed.Title))
	}
	if title {
		buffer.WriteString(fmt.Sprintf("**%s**\n\n", item.Title))
	}
	if brief {
		buffer.WriteString(fmt.Sprintf("%s", p.Sanitize(item.Description)))
	} else {
		buffer.WriteString(fmt.Sprintf("%s", p.Sanitize(item.Content)))
	}
	if url {
		buffer.WriteString(fmt.Sprintf("\n\n`%s`", item.Link))
	}
	return html.UnescapeString(re.ReplaceAllString(buffer.String(), "\n\n"))
}
