package models

import "time"

type RepositoryConfiguration struct {
	ID             int           `json:"id"`
	LastGeneration *time.Time    `json:"last_generation"`
	Private        bool          `json:"private"`
	Cron           string        `json:"cron"`
	CreatedAt      time.Time     `json:"createdAt"`
	UpdatedAt      time.Time     `json:"updatedAt"`
	Repository     *Repository   `json:"repository"`
	Installation   *Installation `json:"installation"`
}

type Repository struct {
	ID           int       `json:"id"`
	Name         string    `json:"name"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
	FullName     string    `json:"full_name"`
	Private      bool      `json:"private"`
	RepositoryID string    `json:"repository_id"`
}

type Installation struct {
	ID             int       `json:"id"`
	InstallationID string    `json:"installation_id"`
	Username       string    `json:"username"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type Count struct {
	Count int `json:"count"`
}
