package vkapi

import (
	"fmt"
	"net/url"
)

func (c *UserClient) CreateBoardComment(groupID, topicID int, message string) string {
	params := fmt.Sprintf("group_id=%d&topic_id=%d&from_group=1&message=%s", groupID, topicID, url.QueryEscape(message))
	return string(c.Request("board.createComment", params))
}

func (c *UserClient) CloseBoard(groupID, topicID int) string {
	params := fmt.Sprintf("group_id=%d&topic_id=%d", groupID, topicID)
	return string(c.Request("board.closeTopic", params))
}

func (c *UserClient) OpenBoard(groupID, topicID int) string {
	params := fmt.Sprintf("group_id=%d&topic_id=%d", groupID, topicID)
	return string(c.Request("board.openTopic", params))
}
