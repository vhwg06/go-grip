package notification

import "time"

type Notification struct {
	Channel string `json:"channel"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

type UserNotification struct {
	ID         int64     `json:"id"`
	UserID     string    `json:"user_id"`
	Type       string    `json:"type"`
	TitleKey   string    `json:"title_key"`
	ContentKey string    `json:"content_key"`
	Data       string    `json:"data"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

type BroadcastMessage struct {
	ID         int64     `json:"id"`
	TitleKey   string    `json:"title_key"`
	ContentKey string    `json:"content_key"`
	Data       string    `json:"data"`
	Sender     string    `json:"sender"`
	CreatedAt  time.Time `json:"created_at"`
}

type BroadcastRead struct {
	ID        int64     `json:"id"`
	MessageID int64     `json:"message_id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type AdminMessage struct {
	ID          int64     `json:"id"`
	TargetType  string    `json:"target_type"`
	TargetValue string    `json:"target_value"`
	Title       string    `json:"title"`
	Body        string    `json:"body"`
	Sender      string    `json:"sender"`
	CreatedAt   time.Time `json:"created_at"`
}
