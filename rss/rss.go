package rss

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/mmcdole/gofeed"
	"github.com/spf13/cobra"
)

func Cmd(result *bytes.Buffer, errResult *error) *cobra.Command {
	var rssCmd = &cobra.Command{
		Use:   "rss",
		Short: "Access RSS feeds",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

		var url string
		var rssFetchCmd = &cobra.Command{
			Use:   "fetch",
			Short: "Fetch items from RSS feed",
			Run: func(cmd *cobra.Command, args []string) {
				feed, err := fetchRSS(url)
				if err != nil {
					*errResult = err
				} else {
					result.WriteString(fmt.Sprintf("%s\n%s", feed.Title, feed.Items[0].Title))
				}
			},
    }
		rssFetchCmd.Flags().StringVarP(&url, "url", "u", "", "URL for feed")

		rssCmd.AddCommand(rssFetchCmd)

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

// EXAMPLE https://haskellweekly.news/podcast.rss
