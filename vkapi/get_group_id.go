package vkapi

import (
	"fmt"
)

type getGroupIDData struct {
	ID int `json:"id"`
}

type getGroupIDResponse struct {
	Response []getGroupIDData `json:"response"`
}

var groupID int

func (c *Client) GetGroupID() int {
	if groupID != 0 {
		return groupID
	}

	jsonR := c.Request("groups.getById", "")
	response := getGroupIDResponse{}
	err := json.Unmarshal(jsonR, &response)
	CheckError(err)

	if len(response.Response) == 0 {
		panic(fmt.Sprintf("GetGroupID() return 0 \n%s\n%#v\n", string(jsonR), c))
	}

	groupID = response.Response[0].ID
	return groupID
}
