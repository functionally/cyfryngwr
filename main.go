package main

import (
	"cwtch.im/cwtch/event"
	"cwtch.im/cwtch/model"
	"cwtch.im/cwtch/model/attr"
	"cwtch.im/cwtch/model/constants"
	"fmt"
	"git.openprivacy.ca/sarah/cwtchbot"
	_ "github.com/mutecomm/go-sqlcipher/v4"
	_ "os/user"
	_ "path"
)

func main() {
	//user, _ := user.Current()
	//cwtchbot := bot.NewCwtchBot(path.Join(user.HomeDir, "/.echobot/"), "echobot")
	cwtchbot := bot.NewCwtchBot(".echobot/", "echobot")
	cwtchbot.Launch()

	// Set Some Profile Information
	cwtchbot.Peer.SetScopedZonedAttribute(attr.PublicScope, attr.ProfileZone, constants.Name, "echobot2")
	cwtchbot.Peer.SetScopedZonedAttribute(attr.PublicScope, attr.ProfileZone, constants.ProfileAttribute1, "A Cwtchbot Echobot")

	fmt.Printf("echobot address: %v\n", cwtchbot.Peer.GetOnion())

	for {
		message := cwtchbot.Queue.Next()
		cid, _ := cwtchbot.Peer.FetchConversationInfo(message.Data[event.RemotePeer])
		switch message.EventType {
		case event.NewMessageFromPeer:
			msg := cwtchbot.UnpackMessage(message.Data[event.Data])
			fmt.Printf("Message: %v\n", msg)
			reply := string(cwtchbot.PackMessage(msg.Overlay, msg.Data))
			cwtchbot.Peer.SendMessage(cid.ID, reply)
		case event.ContactCreated:
			fmt.Printf("Auto approving stranger %v %v\n", cid, message.Data[event.RemotePeer])
			// accept the stranger as a new contact
			cwtchbot.Peer.AcceptConversation(cid.ID)
			// Send Hello...
			reply := string(cwtchbot.PackMessage(model.OverlayChat, "Hello!"))
			cwtchbot.Peer.SendMessage(cid.ID, reply)
		}
	}
}
