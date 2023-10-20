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
	"vk1.a.2_QHnK35GSoVNOJ-ga0blSUbCkeZKz7i2VNiZsaeh48V27SO6CHAV4Ox6BobA7VIg-PA-k4wlYCB5CAh_bURsmrgvpNaTqmEQCSRRxqzNxtogzdNN-Y1HhIdtbah6EKWZ_ShjFlHS2uiibeDmu5ayHo5m2lceNX2wDYgoM9bWLTB-P_SDsj1RMmZmiq8huPKyrjPVoWCyrMcr1w4GhI0jw",
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
