package models

import "time"

type StrapiWebhookPayload struct {
	Event     string    `json:"event"`
	CreatedAt time.Time `json:"createdAt"`
	Model     string    `json:"model"`
	UID       string    `json:"uid"`
	Entry     struct {
		ID             int         `json:"id"`
		LastGeneration interface{} `json:"last_generation"`
		Private        bool        `json:"private"`
		Cron           string      `json:"cron"`
		CreatedAt      time.Time   `json:"createdAt"`
		UpdatedAt      time.Time   `json:"updatedAt"`
	} `json:"entry"`
}

type RepositoryConfiguration struct {
	ID             int         `json:"id"`
	LastGeneration interface{} `json:"last_generation"`
	Private        bool        `json:"private"`
	Cron           string      `json:"cron"`
	CreatedAt      time.Time   `json:"createdAt"`
	UpdatedAt      time.Time   `json:"updatedAt"`
}
