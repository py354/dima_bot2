package vkapi

import "fmt"

func (c *Client) RemoveChatUser(chatID, userID int) {
	c.Request("messages.removeChatUser", fmt.Sprintf("user_id=%d&chat_id=%d", userID, chatID))
}
