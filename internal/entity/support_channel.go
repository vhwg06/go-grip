package entity

type SupportChannel struct {
	ID          string `json:"id"`
	ChannelType string `json:"channel_type"`
	Label       string `json:"label"`
	Link        string `json:"link"`
	IsEnabled   bool   `json:"is_enabled"`
}
