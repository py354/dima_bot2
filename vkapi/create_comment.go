package vkapi

import "fmt"

func (c *Client) CreateComment(ownerID, postID, fromGroup, replyToComment int, message string) string {
	temp := "owner_id=%d&post_id=%d&from_group=%d&message=%s&reply_to_comment=%d"
	params := fmt.Sprintf(temp, ownerID, postID, fromGroup, message, replyToComment)
	return string(c.Request("wall.createComment", params))
}
