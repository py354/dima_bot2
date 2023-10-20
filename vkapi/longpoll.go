package vkapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type LongpollEvent struct {
	ID          int          `json:"id"`
	FromID      int          `json:"from_id"`
	PostID      int          `json:"post_id"`
	Out         int          `json:"out"`
	PeerID      int          `json:"peer_id"`
	Text        string       `json:"text"`
	Date        int          `json:"date"`
	Payload     string       `json:"payload"`
	Attachments []Attachment `json:"attachments"`
	Action      LongAction   `json:"action"`
	Ref         string       `json:"ref"`
}

type LongAction struct {
	Type     string `json:"type"`
	MemberID int    `json:"member_id"`
	Text     string `json:"text"`
}

type longpollUpdate struct {
	LongpollEvent `json:"object"`
	Type          string `json:"type"`
	GroupId       int    `json:"group_id"`
}

type longpollResponse struct {
	Ts      string           `json:"ts"`
	Updates []longpollUpdate `json:"updates"`
	Failed  int              `json:"failed"`
}

type getLongpollServerResponse struct {
	Response getLongPollServerData `json:"response"`
}

func NewLongpoll(token string) *Longpoll {
	return &Longpoll{
		Client: NewClient(token),
	}
}

func (lp *Longpoll) initVKParams() {
	jsonR := lp.Request("groups.getLongPollServer", "group_id="+strconv.Itoa(lp.GetGroupID()))
	response := getLongpollServerResponse{}
	//log.Println(string(jsonR))
	err := json.Unmarshal(jsonR, &response)
	CheckError(err)

	lp.Key = response.Response.Key
	lp.Server = response.Response.Server
	lp.TS = response.Response.Ts
}

func (lp *Longpoll) getEvents() (longpollResponse, error) {
	url := fmt.Sprintf("%s?act=a_check&key=%s&ts=%s&wait=25", lp.Server, lp.Key, lp.TS)
	r, err := http.Get(url)
	if err != nil {
		return longpollResponse{}, err
	}

	defer r.Body.Close()

	answer, err := ioutil.ReadAll(r.Body)
	CheckError(err)

	response := longpollResponse{}
	//log.Println(string(answer))
	err = json.Unmarshal(answer, &response)
	CheckError(err)
	return response, nil
}

// Listen support only "message_new" event from users
func (lp *Longpoll) Listen(inputMessages chan<- *Message) {
	//log.Println("Start Listener")
	lp.initVKParams()

	for {
		response, err := lp.getEvents()
		if response.Failed != 0 || response.Ts == "" || err != nil {
			lp.initVKParams()
			time.Sleep(5 * time.Second)
			continue
		}

		lp.TS = response.Ts
		for _, event := range response.Updates {
			if event.FromID < 0 {
				continue
			}
			if event.Type == "message_new" && event.Out != 1 {
				message := Message{
					FromID:      event.FromID,
					PeerID:      event.PeerID,
					Date:        event.Date,
					Text:        event.Text,
					RawPayload:  event.Payload,
					Attachments: event.Attachments,
					Action:      event.Action,
					Ref:         event.Ref,
					Timestamp:   time.Now(),
				}

				if event.FromID == event.PeerID {
					message.ChatID = 0
				} else {
					message.ChatID = event.PeerID - 2000000000
				}

				if lp.IsClosed {
					return
				}
				inputMessages <- &message
			}
		}
	}
}

type getLongPollServerData struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	Ts     string `json:"ts"`
}

type Longpoll struct {
	*Client
	IsClosed bool
	Key      string
	Server   string
	TS       string
}

type Market struct {
	Id int `json:"id"`
}

type Attachment struct {
	Type   string `json:"type"`
	Market Market `json:"market"`
}

type Message struct {
	FromID      int
	PeerID      int
	Date        int
	Text        string
	RawPayload  string
	ChatID      int
	Attachments []Attachment
	Action      LongAction
	Ref         string
	Timestamp   time.Time
}

type Comment struct {
	ID     int
	FromID int
	PostID int
	Text   string
}


