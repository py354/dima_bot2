package primitives

import (
	"dima_bot/store"
	"dima_bot/vkapi"
)

type Message struct {
	User        *store.User
	Client      vkapi.ClientsPool
	AdminClient *vkapi.UserClient
	LowerText   string
	*vkapi.Message
	Payload string
}

func (msg *Message) Answer(text string, kb *vkapi.Keyboard) {
	kbs := ""
	if kb != nil {
		kbs = kb.String()
	}
	go msg.Client.SendMessage(msg.PeerID, text, kbs, "")
}

type Handler func(msg *Message) bool // функция обрабатывает инфу
type Presenter func(msg *Message)    // функция лишь показывает информацию

type AnswerNode struct {
	Performance Presenter // функция когда человек переходит на эту ноду впервые
	Default     Handler   // функция для обработки сообщения игрока
	Parent      int       // индекс предыдущей ноды
}

var nodes = make([]*AnswerNode, 60)

func RegisterNode(id int) int {
	nodes[id] = &AnswerNode{}
	return id
}

func GetNode(i int) *AnswerNode {
	return nodes[i]
}

func (msg *Message) Trans(nodeI int) {
	node := GetNode(nodeI)
	msg.User.SetNode(nodeI)
	if node.Performance != nil {
		node.Performance(msg)
	}
}

func (msg *Message) GetNode() *AnswerNode {
	return GetNode(msg.User.Node)
}