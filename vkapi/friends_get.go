package vkapi

import "fmt"

type FriendsData struct {
	Count int   `json:"count"`
	Items []int `json:"items"`
}

type Friends struct {
	Response FriendsData `json:"response"`
}

func (c *UserClient) GetFriends(UserID int) []int {
	response := c.Request("friends.get", fmt.Sprintf("user_id=%d&count=10000", UserID))
	r := Friends{}
	err := json.Unmarshal(response, &r)
	CheckError(err)
	return r.Response.Items
}
