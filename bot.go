package main

import (
	"dima_bot/bot"
	"dima_bot/primitives"
	"dima_bot/store"
	"dima_bot/vkapi"
	"log"
	"math/rand"
	"time"
)

var mainTokens = []string{
	"vk1.a._iN09sGUnHs5Iez4FmMQN_l4mSk9rIWLqyS_uPYbgZ2G9gpl9dumNOCFnGHZ3KH9YPMpqx7SMXe84Zrr24wSZTkw4Cm-XiXoTDBBVobgP-rW_GVeAvmrV_D7mCdvVtWnrUzoBsYbx-WNTsZZuIxu3ZjNQbto9u48RXBC7BE0WU02zKndZ-O0yRdpQJxPTVMAlCNeukpYY2ZyyiGVmwA8pQ",
}

var mainClient = vkapi.NewPool(mainTokens)

func main() {
	rand.Seed(time.Now().Unix())
	store.Connect("danis", "IAmRyk", "dima_bot2", "127.0.0.1")
	Handler()
}

func Handler() {
	inputMessages := make(chan *vkapi.Message, 100)
	lp := vkapi.NewLongpoll(mainTokens[0])
	log.Println("Start Listener")
	go lp.Listen(inputMessages)
	for msg := range inputMessages {
		if msg.FromID < 0 {
			return
		}
		go bot.Handle(&primitives.Message{
			Client:    mainClient,
			LowerText: "",
			Message:   msg,
			Payload:   msg.GetPayload(),
		})
	}
}
