package github

import (
	"encoding/json"
	"net/http"
)

// Client is a GitHub http client
type Client struct {
	client *http.Client
	user   string
	repo   string
}

type gitHubVersion struct {
	TagVersion string `json:"tag_name"`
}

// NewClient creates a new GitHub client
func NewClient(user, repo string) *Client {
	return &Client{&http.Client{}, user, repo}
}

// LatestVersion returns the last release version
func (c *Client) LatestVersion() (version string, err error) {
	url := "https://github.com/" + c.user + "/" + c.repo + "/releases/latest"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	lastVersion := gitHubVersion{}
	json.NewDecoder(resp.Body).Decode(&lastVersion)
	return lastVersion.TagVersion, nil
}
