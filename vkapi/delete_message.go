package vkapi

import "fmt"

func (c *Client) DeleteMessage(messageID int) {
	c.Request("messages.delete", fmt.Sprintf("message_ids=%d&spam=1&delete_for_all=1", messageID))
}
