package cwtch

import (
	"fmt"
	"log"
	"os"

	"cwtch.im/cwtch/event"
	"cwtch.im/cwtch/model"
	"cwtch.im/cwtch/model/attr"
	"cwtch.im/cwtch/model/constants"
	"git.openprivacy.ca/sarah/cwtchbot"

	_ "github.com/mutecomm/go-sqlcipher/v4"

	"github.com/functionally/cyfryngwr/dispatch"
)

func Connect(folder string, name string, description string) *bot.CwtchBot {
	cwtchbot := bot.NewCwtchBot(folder, name)
	cwtchbot.Launch()
	cwtchbot.Peer.SetScopedZonedAttribute(attr.PublicScope, attr.ProfileZone, constants.Name, name)
	cwtchbot.Peer.SetScopedZonedAttribute(attr.PublicScope, attr.ProfileZone, constants.ProfileAttribute1, description)
	log.Printf("Cwtch address for %v: %v\n", name, cwtchbot.Peer.GetOnion())
	return cwtchbot
}

func Loop(dispatcher *dispatch.Dispatcher, cwtchbot *bot.CwtchBot, stop chan os.Signal) () {
	shutdown := make(chan any, 1)
	finished := make(chan any, 1)
	go dispatcher.Loop(shutdown, finished)
	for {
		select {
		case message := <-cwtchbot.Queue.OutChan():
			cid, err := cwtchbot.Peer.FetchConversationInfo(message.Data[event.RemotePeer])
			if err != nil {
				log.Printf("Failed to fetch conversation:\n%v\n", err)
			} else {
				handle := cid.Handle
				switch message.EventType {
				case event.PeerStateChange:
					state, exists := message.Data[event.ConnectionState]
					if exists {
						if state == "Authenticated" {
							dispatcher.Online(handle, func(response string) {
								text := string(response)
								if len(text) > 7000 {
									text = text[:7000]
								}
								err = sendMessage(cwtchbot, cid, text)
								if err != nil {
									log.Printf("Failed to send message:\n%v\n", err)
								}
							})
						} else if state == "Disconnected" {
							dispatcher.Offline(handle)
						}
					}
				case event.NewMessageFromPeer:
					msg := cwtchbot.UnpackMessage(message.Data[event.Data])
					text := msg.Data
					if len(text) > 0 && text[0] == '/' {
						dispatcher.Request(handle, text)
					}
				case event.ContactCreated:
					fmt.Printf("Auto approving stranger %v %v\n", cid, message.Data[event.RemotePeer])
					cwtchbot.Peer.AcceptConversation(cid.ID)
					sendMessage(cwtchbot, cid, "Hello!")
				}
			}
		case s := <-stop:
			fmt.Printf("Stopping for signal: %o\n", s)
			shutdown <- true
			<-finished
			return
		}
	}
}

func sendMessage(cwtchbot *bot.CwtchBot, cid *model.Conversation, text string) error {
	msg := string(cwtchbot.PackMessage(model.OverlayChat, text))
	_, err := cwtchbot.Peer.SendMessage(cid.ID, msg)
	return err
}
