package bot

import (
	"dima_bot/primitives"
	"dima_bot/store"
	"log"
	"strings"
	"sync"
)

var msgMutex = make(map[int]*sync.Mutex) // msg.FromID

func Handle(msg *primitives.Message) {
	vkid := msg.FromID
	if msgMutex[vkid] == nil {
		msgMutex[vkid] = &sync.Mutex{}
	}

	msgMutex[vkid].Lock()
	defer msgMutex[msg.FromID].Unlock()

	u := store.GetUser(vkid)
	msg.User = u
	msg.LowerText = strings.ToLower(msg.Text)

	log.Printf("%#v\n", msg.User)
	log.Printf("%#v\n", msg)

	if u.Node == 0 || msg.LowerText == "заново" {
		u.SetUserName(msg.Client.GetUserName(u.VKID))
		msg.Trans(Node1)
	} else {
		node := msg.GetNode()
		if !node.Default(msg) {
			msg.Answer("Я тебя не понимаю, используй кнопки", nil)
		}
	}
}

var Node1 = primitives.RegisterNode(1)
var Node2 = primitives.RegisterNode(2)
var Node3 = primitives.RegisterNode(3)
var Node4 = primitives.RegisterNode(4)
var Node5 = primitives.RegisterNode(5)
var Node6 = primitives.RegisterNode(6)
var Node7 = primitives.RegisterNode(7)
var Node8 = primitives.RegisterNode(8)
