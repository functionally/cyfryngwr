package dispatch

import (
	"bytes"
	"container/list"
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
	results, err := Run(string(request))
	if err != nil {
		return err
	}
	for _, result := range results {
		respond(Response(result))
	}
	return nil
}

func New(config map[string]interface{}) (*Dispatcher, error) {
	dispatcher := Dispatcher{
		Config:     config,
		responders: make(map[Handle]Respond),
	}
	return &dispatcher, nil
}

func Run(input string) ([]string, error) {

	var results = list.New()
	var errResult error = nil

	var rootCmd = &cobra.Command{
		Use:   "/",
		Short: "Cyfryngwr agent",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	var cmdResult bytes.Buffer
	rootCmd.SetOut(&cmdResult)

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Reply with version information",
		Run: func(cmd *cobra.Command, args []string) {
			results.PushBack(fmt.Sprintf("Cyfryngwr %s (%s)", version, gitCommit))
		},
	}
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(rss.Cmd(results, &errResult))

	args, err := shlex.Split(strings.TrimPrefix(input, "/"))
	if err != nil {
		return nil, err
	}
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		return nil, err
	}

	if cmdResult.Len() > 0 {
		results.PushFront(cmdResult.String())
	}
	result := make([]string, results.Len())
	for i := 0; results.Len() > 0; i += 1 {
		x := results.Front()
		result[i] = x.Value.(string)
		results.Remove(x)
	}
	return result, errResult
}
