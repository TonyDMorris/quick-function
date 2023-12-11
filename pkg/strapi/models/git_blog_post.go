package models

import "time"

type GitBlogPost struct {
	Title         string     `json:"title,omitempty"`
	Description   string     `json:"description,omitempty"`
	Body          string     `json:"body,omitempty"`
	CommitFrom    *time.Time `json:"commit_from,omitempty"`
	CommitTo      *time.Time `json:"commit_to,omitempty"`
	Repository    string     `json:"repository,omitempty"`
	OwnerUsername string     `json:"owner_username,omitempty"`
}
