package rss

import (
	"bytes"
	"fmt"
	"log"
	"sync"
	"time"
)

type User struct {
	subscriptions map[Url]*time.Time
	showFeed      bool
	showTitle     bool
	showUrl       bool
	showDetail    bool
	manager       *Manager
	mu            sync.Mutex
	stop          chan any
	responses     chan string
}

const bufferSize = 10

func MakeUser(manager *Manager) *User {
	return &User{
		subscriptions: make(map[Url]*time.Time),
		showFeed:      true,
		showTitle:     true,
		showUrl:       true,
		showDetail:    false,
		manager:       manager,
		mu:            sync.Mutex{},
		stop:          make(chan any, 1),
		responses:     make(chan string),
	}
}

func (self User) Start(respond func(string)) {

	tick := time.Tick(60 * time.Second)

	for {
		select {
		case <-self.stop:
			return
		case response := <-self.responses:
			respond(response)
		case <-tick:
			log.Printf("RSS TICK\n")
		}
	}

}

func (self User) Stop() {
	self.responses = make(chan string)
	self.stop <- nil
}

func (self User) Add(url Url) {
	self.mu.Lock()
	defer self.mu.Unlock()
	_, exists := self.subscriptions[url]
	if exists {
		self.responses <- fmt.Sprintf("Already subscribed: %v", url)
		return
	}
	since := time.UnixMilli(0)
	self.subscriptions[url] = &since
}

func (self User) Remove(url Url) {
	self.mu.Lock()
	defer self.mu.Unlock()
	_, exists := self.subscriptions[url]
	if !exists {
		self.responses <- fmt.Sprintf("Not subscribed: %v", url)
		return
	}
	delete(self.subscriptions, url)
}

func (self User) List() {
	self.mu.Lock()
	defer self.mu.Unlock()
	var buffer bytes.Buffer
	for url, _ := range self.subscriptions {
		buffer.WriteString(fmt.Sprintf("%s\n", url))
	}
	self.responses <- buffer.String()
}

func (self User) Clear() {
	self.mu.Lock()
	defer self.mu.Unlock()
	self.subscriptions = make(map[Url]*time.Time, bufferSize)
}

func (self User) Fetch(url Url, count uint) {
	self.mu.Lock()
	defer self.mu.Unlock()
	since := time.UnixMilli(0)
	feed, items, err := self.manager.Fetch(url, &since, count)
	if err != nil {
		self.responses <- fmt.Sprintf("Failed to fetch feed: %s\n,%s", url, err)
		return
	}
	for _, item := range items {
		self.responses <- formatItem(feed, item, self.showFeed, self.showTitle, !self.showDetail, self.showUrl)
	}
}

// func (self User) Show

// func (self User) Import

// func (self User) Export
