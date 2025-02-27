package main

import (
	"fmt"
	"log"

	"cwtch.im/cwtch/event"
	"cwtch.im/cwtch/model"
	"cwtch.im/cwtch/model/attr"
	"cwtch.im/cwtch/model/constants"
	"git.openprivacy.ca/sarah/cwtchbot"

	_ "github.com/mutecomm/go-sqlcipher/v4"

	"github.com/functionally/cyfryngwr/dispatch"
)

func main() {

	cwtchbot := bot.NewCwtchBot(".cyfryngwr/", "cyfryngwr")
	cwtchbot.Launch()
	cwtchbot.Peer.SetScopedZonedAttribute(attr.PublicScope, attr.ProfileZone, constants.Name, "cyfryngwr")
	cwtchbot.Peer.SetScopedZonedAttribute(attr.PublicScope, attr.ProfileZone, constants.ProfileAttribute1, "Cyfryngwr, a cwtch agent")
	log.Printf("Cyfryngwr address: %v\n", cwtchbot.Peer.GetOnion())

	for {
		message := cwtchbot.Queue.Next()
		cid, err := cwtchbot.Peer.FetchConversationInfo(message.Data[event.RemotePeer])
		if err != nil {
			log.Printf("Failed to fetch message: %v\n", err)
		} else {
			switch message.EventType {
			case event.NewMessageFromPeer:
				msg := cwtchbot.UnpackMessage(message.Data[event.Data])
				log.Printf("Received message: %v\n", msg.Data)
				result, err := dispatch.Run(msg.Data)
				if err != nil {
					log.Printf("Failed to process message: %v\n", err)
					result = err.Error()
				}
				err = sendMessage(cwtchbot, cid, result)
				if err != nil {
					log.Printf("Failed to send message: %v\n", err)
				}
			case event.ContactCreated:
				fmt.Printf("Auto approving stranger %v %v\n", cid, message.Data[event.RemotePeer])
				cwtchbot.Peer.AcceptConversation(cid.ID)
				sendMessage(cwtchbot, cid, "Hello!")
			}
		}
	}
}

func sendMessage(cwtchbot *bot.CwtchBot, cid *model.Conversation, text string) error {
	msg := string(cwtchbot.PackMessage(model.OverlayChat, text))
	_, err := cwtchbot.Peer.SendMessage(cid.ID, msg)
	return err
}
