package vkapi

import "fmt"

type ResolveDomain struct {
	Response struct {
		Type     string `json:"type"`
		ObjectID int    `json:"object_id"`
	} `json:"response"`
}

func (c *Client) IsGroupDomain(domain string) bool {
	answer := c.Request("utils.resolveScreenName", fmt.Sprintf("screen_name=%s", domain))
	data := ResolveDomain{}
	err := json.Unmarshal(answer, &data)
	if err != nil {
		return false
	}
	return data.Response.Type == "group"
}

func (c *Client) ResolveDomain(domain string) (objectType string, objectID int) {
	answer := c.Request("utils.resolveScreenName", fmt.Sprintf("screen_name=%s", domain))
	data := ResolveDomain{}
	err := json.Unmarshal(answer, &data)
	if err != nil {
		return "", 0
	}
	return data.Response.Type, data.Response.ObjectID
}
