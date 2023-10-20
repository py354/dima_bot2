package vkapi

import (
	url2 "net/url"
	"strings"
	"time"
)

type getShortLinkData struct {
	ShortURL string `json:"short_url"`
}

type getShortLinkResponse struct {
	Response getShortLinkData `json:"response"`
}

func (c *Client) GetShortLink(url string) string {
	var result string
	tries := 0
	for result == "" && tries < 5 {
		result = c.getShortLink(url)
		if result == "" {
			tries += 1
			time.Sleep(200 * time.Millisecond)
		}
	}
	return strings.Trim(result, "https://")
}

func (c *Client) getShortLink(url string) string {
	jsonR := c.Request("utils.getShortLink", "url="+url2.QueryEscape(url))
	response := getShortLinkResponse{}
	err := json.Unmarshal(jsonR, &response)
	if err != nil {
		return ""
	}
	return response.Response.ShortURL
}
