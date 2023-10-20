package vkapi

import (
	"fmt"
	"net/url"
	"strconv"
	"time"
)

func (c *Client) SendSticker(peerID int, stickerID int) []byte {
	params := fmt.Sprintf("peer_id=%d&sticker_id=%d&random_id=&", peerID, stickerID)
	return c.Request("messages.send", params)
}

func (c *Client) sendBroadcastStickers(userIDs string, stickerID int) []byte {
	params := fmt.Sprintf("user_ids=%s&sticker_id=%d&random_id=&", userIDs, stickerID)
	return c.Request("messages.send", params)
}

func (c *Client) sendMessage(dst string, message, keyboard, attachment string) []byte {
	params := dst + "&random_id=&disable_mentions=1&dont_parse_links=1&"

	if message != "" {
		params += "message=" + url.QueryEscape(message) + "&"
	}

	if keyboard != "" {
		params += "keyboard=" + keyboard + "&"
	}

	if attachment != "" {
		params += "attachment=" + attachment + "&"
	}

	return c.Request("messages.send", params)
}

func (c *Client) SendMessage(peerID int, message, keyboard, attachment string) []byte {
	return c.sendMessage("peer_id="+strconv.Itoa(peerID), message, keyboard, attachment)
}

func (c *Client) broadcast(userIDs, message, keyboard, attachment string) []byte {
	return c.sendMessage("user_ids="+userIDs, message, keyboard, attachment)
}

const maxBroadcast = 100

func (c *Client) Broadcast(userIDs []int, message, keyboard, attachment string) {
	offset := 0
	for offset != len(userIDs) {
		stringUserIDs := ""
		passed := len(userIDs) - offset
		if passed > maxBroadcast {
			passed = maxBroadcast
		}

		for i := offset; i < passed+offset; i++ {
			stringUserIDs += strconv.Itoa(userIDs[i]) + ","
		}

		offset += passed
		c.broadcast(stringUserIDs[:len(stringUserIDs)-1], message, keyboard, attachment)
	}
}

func (c *Client) SlowBroadcast(userIDs []int, message, keyboard, attachment string, sleep time.Duration) {
	offset := 0
	for offset != len(userIDs) {
		stringUserIDs := ""
		passed := len(userIDs) - offset
		if passed > maxBroadcast {
			passed = maxBroadcast
		}

		for i := offset; i < passed+offset; i++ {
			stringUserIDs += strconv.Itoa(userIDs[i]) + ","
		}

		time.Sleep(sleep)
		offset += passed
		c.broadcast(stringUserIDs[:len(stringUserIDs)-1], message, keyboard, attachment)
	}
}

func (c *Client) SendBroadcastStickers(userIDs []int, stickerID int) [][]byte {
	result := make([][]byte, 0, 100)

	offset := 0
	for offset != len(userIDs) {
		stringUserIDs := ""
		passed := len(userIDs) - offset
		if passed > maxBroadcast {
			passed = maxBroadcast
		}

		for i := offset; i < passed+offset; i++ {
			stringUserIDs += strconv.Itoa(userIDs[i]) + ","
		}

		offset += passed
		result = append(result, c.sendBroadcastStickers(stringUserIDs[:len(stringUserIDs)-1], stickerID))
	}

	return result
}
