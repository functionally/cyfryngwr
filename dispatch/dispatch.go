package dispatch

import (
	"bytes"
	"container/list"
	"fmt"
	"log"
	"strings"
	"sync"

	"github.com/google/shlex"
	"github.com/spf13/cobra"

	"github.com/functionally/cyfryngwr/rss"
)

var (
	version   = "dev"
	gitCommit = "none"
)

type Responder func(string) ()

type Dispatcher struct {
	config     map[string]interface{}
	responders map[string]Responder
	mu sync.Mutex
	shutdown chan any
}

func New(config map[string]interface{}) (*Dispatcher, error) {
	const bufferSize = 10
	dispatcher := Dispatcher{
		config:     config,
		responders: make(map[string]Responder, bufferSize),
		mu: sync.Mutex{},
		shutdown: make(chan any, 1),
	}
	return &dispatcher, nil
}

func (self Dispatcher) Loop() {
	for {
		select {
		case <-self.shutdown:
			return
		}
	}
}

func (self Dispatcher) Online(handle string, respond Responder) {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.responders[handle] = respond
}

func (self Dispatcher) Offline(handle string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	_, exists := self.responders[handle]
	if !exists {
		log.Printf("Offline user not found: %s\n", handle)
	} else {
		delete(self.responders, handle)
	}
}

func (self Dispatcher) Request(handle string, request string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	respond, exists := self.responders[handle]
	if !exists {
		log.Printf("Requesting user not found: %v\n", handle)
	}
	results, err := Run(request)
	if err != nil {
		log.Printf("Command failed: %s\n%s\n", request, err.Error())
		respond(err.Error())
	}
	for _, result := range results {
		respond(result)
	}
}

func (self Dispatcher) Shutdown() {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.shutdown <- nil
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
