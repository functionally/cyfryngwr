package rss

import (
	"bytes"
	"container/list"
	"fmt"
	"net/http"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

func Cmd(results *list.List, errResult *error) *cobra.Command {
	var rssCmd = &cobra.Command{
		Use:   "rss",
		Short: "Access RSS feeds",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	{
		var count uint
		var source bool
		var title bool
		var brief bool = true
		var url bool
		var rssFetchCmd = &cobra.Command{
			Use:   "fetch [url]",
			Short: "Fetch items from RSS feed",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				feed, err := fetchRSS(args[0])
				if err != nil {
					*errResult = err
				} else {
					fmt.Printf("UPDATED %s\n", feed.Updated)
					for i, item := range feed.Items {
						if uint(i) >= count {
							break
						}
						results.PushBack(formatItem(feed, item, source, title, brief, url))
					}
				}
			},
		}
		rssFetchCmd.Flags().UintVarP(&count, "limit", "l", 1, "Number of items to retrieve.")
		rssFetchCmd.Flags().BoolVarP(&source, "feed", "f", true, "Show the name of the feed.")
		rssFetchCmd.Flags().BoolVarP(&title, "title", "t", true, "Show the title of the item.")
		rssFetchCmd.Flags().BoolVarP(&url, "url", "u", true, "Show the URL for the item.")
		rssCmd.AddCommand(rssFetchCmd)
	}

	return rssCmd
}

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
	p := bluemonday.StrictPolicy()
	var buffer bytes.Buffer
	if source {
		buffer.WriteString(fmt.Sprintf("*%s*\n", feed.Title))
	}
	if title {
		buffer.WriteString(fmt.Sprintf("**%s**\n", item.Title))
	}
	if source || title {
		buffer.WriteString("\n")
	}
	if brief {
		buffer.WriteString(fmt.Sprintf("%s", p.Sanitize(item.Description)))
	} else {
		buffer.WriteString(fmt.Sprintf("%s", p.Sanitize(item.Content)))
	}
	if url {
		buffer.WriteString(fmt.Sprintf("\n\n`%s`", item.Link))
	}
	return buffer.String()
}
