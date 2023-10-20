package vkapi

import "fmt"

type HistoryItem struct {
	Date   int `json:"date"`
	FromID int `json:"from_id"`
	ID     int `json:"id"`
}

type History struct {
	Items []HistoryItem `json:"items"`
}

type HistoryResponse struct {
	Response History `json:"response"`
}

func (c *Client) GetMessagesHistory(peerID, offset, count int) []HistoryItem {
	params := fmt.Sprintf("offset=%d&peer_id=%d&count=%d", offset, peerID, count)
	response := c.Request("messages.getHistory", params)
	data := HistoryResponse{}
	err := json.Unmarshal(response, &data)
	CheckError(err)
	return data.Response.Items
}

func (c *Client) GetMessageByConversationMessageID(peerID, conversationMessageID int) {

}
