package models

type AppInstallation struct {
	ID                     int         `json:"id"`
	Account                Account     `json:"account"`
	AccessTokensURL        string      `json:"access_tokens_url"`
	RepositoriesURL        string      `json:"repositories_url"`
	HTMLURL                string      `json:"html_url"`
	AppID                  int         `json:"app_id"`
	TargetID               int         `json:"target_id"`
	TargetType             string      `json:"target_type"`
	Permissions            Permissions `json:"permissions"`
	Events                 []string    `json:"events"`
	SingleFileName         string      `json:"single_file_name"`
	HasMultipleSingleFiles bool        `json:"has_multiple_single_files"`
	SingleFilePaths        []string    `json:"single_file_paths"`
	RepositorySelection    string      `json:"repository_selection"`
	CreatedAt              string      `json:"created_at"`
	UpdatedAt              string      `json:"updated_at"`
	AppSlug                string      `json:"app_slug"`
	SuspendedAt            interface{} `json:"suspended_at"`
	SuspendedBy            interface{} `json:"suspended_by"`
}
type Account struct {
	Login             string `json:"login"`
	ID                int    `json:"id"`
	NodeID            string `json:"node_id"`
	AvatarURL         string `json:"avatar_url"`
	GravatarID        string `json:"gravatar_id"`
	URL               string `json:"url"`
	HTMLURL           string `json:"html_url"`
	FollowersURL      string `json:"followers_url"`
	FollowingURL      string `json:"following_url"`
	GistsURL          string `json:"gists_url"`
	StarredURL        string `json:"starred_url"`
	SubscriptionsURL  string `json:"subscriptions_url"`
	OrganizationsURL  string `json:"organizations_url"`
	ReposURL          string `json:"repos_url"`
	EventsURL         string `json:"events_url"`
	ReceivedEventsURL string `json:"received_events_url"`
	Type              string `json:"type"`
	SiteAdmin         bool   `json:"site_admin"`
}
type Permissions struct {
	Checks   string `json:"checks"`
	Metadata string `json:"metadata"`
	Contents string `json:"contents"`
}