package types

import (
	"time"
)

type RunMetadata struct {
	ID               int64         `json:"id"`
	Name             string        `json:"name"`
	NodeID           string        `json:"node_id"`
	HeadBranch       string        `json:"head_branch"`
	HeadSha          string        `json:"head_sha"`
	Path             string        `json:"path"`
	DisplayTitle     string        `json:"display_title"`
	RunNumber        int           `json:"run_number"`
	Event            string        `json:"event"`
	Status           string        `json:"status"`
	Conclusion       string        `json:"conclusion"`
	WorkflowID       int           `json:"workflow_id"`
	CheckSuiteID     int64         `json:"check_suite_id"`
	CheckSuiteNodeID string        `json:"check_suite_node_id"`
	URL              string        `json:"url"`
	HTMLURL          string        `json:"html_url"`
	PullRequests     []interface{} `json:"pull_requests"`
	CreatedAt        time.Time     `json:"created_at"`
	UpdatedAt        time.Time     `json:"updated_at"`
	Actor            struct {
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
	} `json:"actor"`
	RunAttempt          int           `json:"run_attempt"`
	ReferencedWorkflows []interface{} `json:"referenced_workflows"`
	RunStartedAt        time.Time     `json:"run_started_at"`
	TriggeringActor     struct {
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
	} `json:"triggering_actor"`
	JobsURL            string      `json:"jobs_url"`
	LogsURL            string      `json:"logs_url"`
	CheckSuiteURL      string      `json:"check_suite_url"`
	ArtifactsURL       string      `json:"artifacts_url"`
	CancelURL          string      `json:"cancel_url"`
	RerunURL           string      `json:"rerun_url"`
	PreviousAttemptURL interface{} `json:"previous_attempt_url"`
	WorkflowURL        string      `json:"workflow_url"`
	HeadCommit         struct {
		ID        string    `json:"id"`
		TreeID    string    `json:"tree_id"`
		Message   string    `json:"message"`
		Timestamp time.Time `json:"timestamp"`
		Author    User      `json:"author"`
		Committer User      `json:"committer"`
	} `json:"head_commit"`
	Repository     Repo     `json:"repository"`
	HeadRepository HeadRepo `json:"head_repository"`
}
