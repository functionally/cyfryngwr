package rss

import (
	"container/list"

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
