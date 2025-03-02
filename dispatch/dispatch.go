package dispatch

import (
	"bytes"
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

type Responder func(string)

type Dispatcher struct {
	config     map[string]interface{}
	responders map[string]Responder
	rssManager *rss.Manager
	rssUsers   map[string]*rss.User
	mu         sync.Mutex
	shutdown   chan any
}

func New(config map[string]interface{}) (*Dispatcher, error) {
	const bufferSize = 10
	dispatcher := Dispatcher{
		config:     config,
		responders: make(map[string]Responder, bufferSize),
		rssManager: rss.New(),
		rssUsers:   make(map[string]*rss.User, bufferSize),
		mu:         sync.Mutex{},
		shutdown:   make(chan any, 1),
	}
	return &dispatcher, nil
}

func (self Dispatcher) Loop() {

	for {
		select {
		case <-self.shutdown:
			for _, user := range self.rssUsers {
				user.Stop()
			}
			return
		}
	}

}

func (self Dispatcher) Online(handle string, respond Responder) {
	self.mu.Lock()
	self.responders[handle] = respond
	rssUser, exists := self.rssUsers[handle]
	if !exists {
		rssUser = rss.MakeUser(self.rssManager)
		self.rssUsers[handle] = rssUser
	}
	self.mu.Unlock()
	go rssUser.Start(respond)
}

func (self Dispatcher) Offline(handle string) {
	self.mu.Lock()
	_, exists := self.responders[handle]
	if !exists {
		log.Printf("Offline user not found: %s\n", handle)
		return
	}
	rssUser, exists := self.rssUsers[handle]
	if !exists {
		log.Printf("Offline RSS user not found: %s\n", handle)
		return
	}
	delete(self.responders, handle)
	self.mu.Unlock()
	rssUser.Stop()
}

func (self Dispatcher) Request(handle string, request string) {
	self.mu.Lock()
	defer self.mu.Unlock()
	respond, exists := self.responders[handle]
	if !exists {
		log.Printf("Requesting user not found: %v\n", handle)
		return
	}
	rssUser, exists := self.rssUsers[handle]
	if !exists {
		log.Printf("Requesting RSS user not found: %v\n", handle)
		return
	}
	Run(rssUser, request, respond)
}

func (self Dispatcher) Shutdown() {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.shutdown <- nil
}

func Run(rssUser *rss.User, input string, respond func(string)) {

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
			respond(fmt.Sprintf("Cyfryngwr %s (%s)", version, gitCommit))
		},
	}
	rootCmd.AddCommand(versionCmd)

	rootCmd.AddCommand(rss.Cmd(rssUser))

	args, err := shlex.Split(strings.TrimPrefix(input, "/"))
	if err != nil {
		respond(err.Error())
		return
	}
	rootCmd.SetArgs(args)
	if err := rootCmd.Execute(); err != nil {
		respond(err.Error())
		return
	}

	if cmdResult.Len() > 0 {
		respond(cmdResult.String())
	}

}
