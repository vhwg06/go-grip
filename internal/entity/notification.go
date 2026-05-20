package entity

type Notification struct {
	Channel string `json:"channel"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}
