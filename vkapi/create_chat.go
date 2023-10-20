package vkapi

type createChatData struct {
	Response int `json:"response"`
}

func (c *UserClient) CreateChat(title string) int {
	jsonR := c.Request("messages.createChat", "title="+title)
	r := createChatData{}
	err := json.Unmarshal(jsonR, &r)
	CheckError(err)
	return r.Response
}
