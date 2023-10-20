package vkapi

import "strconv"

type ConvMembers struct {
	Response struct {
		Items []struct {
			MemberId int `json:"member_id"`
		} `json:"items"`
	} `json:"response"`
}

func (c *Client) GetConvMembers(peerID int) []int {
	result := make([]int, 0, 128)
	answer := c.Request("messages.getConversationMembers", "peer_id="+strconv.Itoa(peerID))
	data := ConvMembers{}
	err := json.Unmarshal(answer, &data)
	CheckError(err)
	for _, m := range data.Response.Items {
		result = append(result, m.MemberId)
	}
	return result
}
