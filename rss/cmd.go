package rss

import (
	"github.com/spf13/cobra"
)

func Cmd(user *User) *cobra.Command {
	var rssCmd = &cobra.Command{
		Use:   "rss",
		Short: "Access RSS feeds",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	rssCmd.AddCommand(fetchCmd(user))

	return rssCmd
}

func fetchCmd(user *User) *cobra.Command {
	var count uint
	var cmd = &cobra.Command{
		Use:   "fetch [url]",
		Short: "Fetch items from RSS feed",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			user.Fetch(Url(args[0]), count)
		},
	}
	cmd.Flags().UintVarP(&count, "limit", "l", 1, "Number of items to retrieve.")
	return cmd
}
