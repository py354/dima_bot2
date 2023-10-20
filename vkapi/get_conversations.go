package vkapi

import (
	"fmt"
	"log"
)

type Conversations struct {
	Response struct {
		Count int `json:"count"`
		Items []struct {
			Conversation struct {
				Peer struct {
					Id   int    `json:"id"`
					Type string `json:"type"`
				} `json:"peer"`
			} `json:"conversation"`
		} `json:"items"`
	} `json:"response"`
}

func (c *Client) GetConversations() []int {
	vkids := make([]int, 0)

	offset := 0
	for {
		jsonR := c.Request("messages.getConversations", fmt.Sprintf("count=200&offset=%d", offset))
		offset += 200

		convs := Conversations{}
		err := json.Unmarshal(jsonR, &convs)
		if err != nil {
			log.Println(err)
			continue
		}

		log.Println(convs.Response.Count)
		if len(convs.Response.Items) == 0 {
			break
		}

		for _, conv := range convs.Response.Items {
			if conv.Conversation.Peer.Type == "user" {
				vkids = append(vkids, conv.Conversation.Peer.Id)
			}
		}

	}

	return vkids
}
