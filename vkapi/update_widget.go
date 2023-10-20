package vkapi

import (
	"fmt"
	"net/url"
)

func (c *Client) UpdateWidget(code, wType string) string {
	params := fmt.Sprintf("code=%s&type=%s", url.QueryEscape(code), url.QueryEscape(wType))
	return string(c.Request("appWidgets.update", params))
}
