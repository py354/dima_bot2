package vkapi

import "log"

const (
	KbRed   = "negative"
	KbGreen = "positive"
	KbWhite = "default"
	KbBlue  = "primary"
)

type Payload struct {
	Command string `json:"button"`
}

type Action struct {
	Type    string `json:"type"`
	AppID   string `json:"app_id,omitempty"`
	Payload string `json:"payload"`
	Label   string `json:"label"`
}

type Button struct {
	Action `json:"action"`
	Color  string `json:"color,omitempty"`
}

func (bi *Action) SetPayload(pl string) {
	j, err := json.Marshal(Payload{pl})
	if err != nil {
		panic(err)
	}

	bi.Payload = string(j)
}

func (msg *Message) GetPayload() string {
	if msg.RawPayload == "" {
		return ""
	}

	data := Payload{}
	err := json.Unmarshal([]byte(msg.RawPayload), &data)
	if err != nil {
		return ""
	}
	return data.Command
}

func MakeAppButton(appID, label string) Button {
	return Button{
		Action: Action{
			Type:    "open_app",
			AppID:   appID,
			Payload: "",
			Label:   label,
		},
		Color: "",
	}
}

func MakeButton(name, payload, color string) Button {
	action := Action{
		Type:  "text",
		Label: name,
	}
	action.SetPayload(payload)

	return Button{
		Action: action,
		Color:  color,
	}
}

type Keyboard struct {
	OneTime     bool       `json:"one_time"`
	Inline      bool       `json:"inline,omitempty"`
	ButtonsGrid [][]Button `json:"buttons"`
	Cache       string     `json:"-"`
}

func (kb *Keyboard) String() string {
	if kb.Cache == "" {
		if len(kb.ButtonsGrid) == 0 {
			return ""
		}

		data, err := json.Marshal(kb)
		if err != nil {
			log.Fatalln(data)
		}
		kb.Cache = string(data)
	}
	return kb.Cache
}
