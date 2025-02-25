package main


import (
	"fmt"

	"cwtch.im/cwtch/event"
	"cwtch.im/cwtch/model"
	"cwtch.im/cwtch/model/attr"
	"cwtch.im/cwtch/model/constants"
	"git.openprivacy.ca/sarah/cwtchbot"
	_ "github.com/mutecomm/go-sqlcipher/v4"

	"github.com/functionally/cyfryngwr/dispatch"
	"github.com/functionally/cyfryngwr/rss"
)

func main() {

	cwtchbot := bot.NewCwtchBot(".cyfryngwr/", "cyfryngwr")
	cwtchbot.Launch()

	cwtchbot.Peer.SetScopedZonedAttribute(attr.PublicScope, attr.ProfileZone, constants.Name, "cyfryngwr")
	cwtchbot.Peer.SetScopedZonedAttribute(attr.PublicScope, attr.ProfileZone, constants.ProfileAttribute1, "Cyfryngwr, a cwtch agent")

	fmt.Printf("echobot address: %v\n", cwtchbot.Peer.GetOnion())

	for {
		message := cwtchbot.Queue.Next()
		cid, _ := cwtchbot.Peer.FetchConversationInfo(message.Data[event.RemotePeer])
		switch message.EventType {
		case event.NewMessageFromPeer:
			msg := cwtchbot.UnpackMessage(message.Data[event.Data])
			fmt.Printf("Message: %v\n", msg)
			result, _ := dispatch.Run(msg.Data)
	                reply := string(cwtchbot.PackMessage(model.OverlayChat, result))
			cwtchbot.Peer.SendMessage(cid.ID, reply)
	//reply := string(cwtchbot.PackMessage(msg.Overlay, msg.Data))
			reply, _ = rss.Message(cwtchbot, "https://haskellweekly.news/podcast.rss")
			cwtchbot.Peer.SendMessage(cid.ID, reply)
		case event.ContactCreated:
			fmt.Printf("Auto approving stranger %v %v\n", cid, message.Data[event.RemotePeer])
			cwtchbot.Peer.AcceptConversation(cid.ID)
			reply := string(cwtchbot.PackMessage(model.OverlayChat, "Hello!"))
			cwtchbot.Peer.SendMessage(cid.ID, reply)
		}
	}
}
