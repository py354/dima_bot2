package vkapi

import (
	"fmt"
)

type Story struct {
	Response struct {
		Count int `json:"count"`
		Items []struct {
			Date      int  `json:"date"`
			IsDeleted bool `json:"is_deleted"`
		} `json:"items"`
	} `json:"response"`
}

func (c *Client) GetStories(ownerID, id int) (bool, int) {
	response := c.Request("stories.getById", fmt.Sprintf("stories=%d_%d", ownerID, id))

	r := Story{}
	err := json.Unmarshal(response, &r)
	CheckError(err)
	if len(r.Response.Items) == 1 {
		return r.Response.Items[0].IsDeleted, r.Response.Items[0].Date
	} else {
		return true, 0
	}
}
