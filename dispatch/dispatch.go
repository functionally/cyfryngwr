package dispatch

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/google/shlex"
	"github.com/spf13/cobra"

	"github.com/functionally/cyfryngwr/rss"
)

var (
	version   = "dev"
	gitCommit = "none"
)

type Handle string

type Request string

type Response string

type Respond func(Response)

type Dispatcher struct {
	Config     map[string]interface{}
	responders map[Handle]Respond
}

func (self Dispatcher) Register(handle Handle, respond Respond) error {
	self.responders[handle] = respond
	return nil
}

func (self Dispatcher) Request(handle Handle, request Request) error {
	respond, exists := self.responders[handle]
	if !exists {
		return fmt.Errorf("Requestor %v not found")
	}
	result, err := Run(string(request))
	if err != nil {
		return err
	}
	respond(Response(result))
	return nil
}

func New(config map[string]interface{}) (*Dispatcher, error) {
	dispatcher := Dispatcher{
		Config:     config,
		responders: make(map[Handle]Respond),
	}
	return &dispatcher, nil
}

func Run(input string) (string, error) {

	var result bytes.Buffer
	var errResult error = nil

	var rootCmd = &cobra.Command{
		Use:   "/",
		Short: "Cyfryngwr agent",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	rootCmd.SetOut(&result)

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Reply with version information",
		Run: func(cmd *cobra.Command, args []string) {
			result.WriteString(fmt.Sprintf("Cyfryngwr %s (%s)", version, gitCommit))
		},
	}
	rootCmd.AddCommand(versionCmd)

	var rssCmd = &cobra.Command{
		Use:   "rss",
		Short: "Access RSS feeds",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	{
		var url string
		var rssFetchCmd = &cobra.Command{
			Use:   "fetch",
			Short: "Fetch items from RSS feed",
			Run: func(cmd *cobra.Command, args []string) {
				feed, err := rss.FetchRSS(url)
				if err != nil {
					errResult = err
				} else {
					result.WriteString(fmt.Sprintf("%s\n%s", feed.Title, feed.Items[0].Title))
				}
			},
		}
		rssFetchCmd.Flags().StringVarP(&url, "url", "u", "", "URL for feed")
		rssCmd.AddCommand(rssFetchCmd)
	}
	rootCmd.AddCommand(rssCmd)

	args, err := shlex.Split(strings.TrimPrefix(input, "/"))
	if err != nil {
		return "", err
	}
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		return "", err
	}

	return result.String(), errResult
}
