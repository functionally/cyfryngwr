package dispatch

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"strings"

	"github.com/google/shlex"
	"github.com/spf13/cobra"

	"github.com/functionally/cyfryngwr/rss"
	"github.com/functionally/cyfryngwr/state"
)

var (
	version   = "dev"
	gitCommit = "none"
)

type Dispatcher struct {
	config     map[string]interface{}
	responders map[state.Handle]state.Responder
	online     chan state.Online
	offline    chan state.Offline
	request    chan state.Request
}

func New(config map[string]interface{}) (*Dispatcher, error) {
	dispatcher := Dispatcher{
		config:     config,
		responders: make(map[state.Handle]state.Responder),
		online:     make(chan state.Online),
		offline:    make(chan state.Offline),
		request:    make(chan state.Request),
	}
	return &dispatcher, nil
}

func (self Dispatcher) Loop(shutdown chan any, finished chan any) {
	for {
		select {
		case x := <-self.online:
			self.responders[x.User] = x.Respond
		case x := <-self.offline:
			_, exists := self.responders[x.User]
			if !exists {
				log.Printf("Offline user not found: %s\n", x.User)
			} else {
				delete(self.responders, x.User)
			}
		case x := <-self.request:
			respond, exists := self.responders[x.User]
			if !exists {
				log.Printf("Requesting user not found: %v\n", x.User)
			}
			results, err := Run(x.Text)
			if err != nil {
				log.Printf("Command failed: %s\n%s\n", x.Text, err.Error())
				respond(err.Error())
			}
			for _, result := range results {
				respond(result)
			}
		case <-shutdown:
			finished <- nil
			return
		}
	}
}

func (self Dispatcher) Online(handle string, respond state.Responder) {
	self.online <- state.Online{User: state.Handle(handle), Respond: respond}
}

func (self Dispatcher) Offline(handle string) {
	self.offline <- state.Offline{User: state.Handle(handle)}
}

func (self Dispatcher) Request(handle string, request string) {
	self.request <- state.Request{User: state.Handle(handle), Text: request}
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
