package models

import "time"

type StrapiWebhookPayload struct {
	Event     string      `json:"event"`
	CreatedAt time.Time   `json:"createdAt"`
	Model     string      `json:"model"`
	UID       string      `json:"uid"`
	Entry     interface{} `json:"entry"`
}

type Carrier struct {
	Data interface{} `json:"data"`
}
