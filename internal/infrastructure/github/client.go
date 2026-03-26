package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"time"

	"github.com/adrock-miles/about-me/internal/domain/project"
)

type repoResponse struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Language    string   `json:"language"`
	Topics      []string `json:"topics"`
	HTMLURL     string   `json:"html_url"`
	Stars       int      `json:"stargazers_count"`
	Fork        bool     `json:"fork"`
	UpdatedAt   string   `json:"updated_at"`
}

// Client interacts with the GitHub API.
type Client struct {
	username string
	apiURL   string
	http     *http.Client
}

// NewClient creates a new GitHub API client.
func NewClient(username, apiURL string) *Client {
	return &Client{
		username: username,
		apiURL:   apiURL,
		http:     &http.Client{Timeout: 10 * time.Second},
	}
}

// FetchRecentProjects fetches the user's public repos from GitHub.
func (c *Client) FetchRecentProjects() ([]project.Project, error) {
	url := fmt.Sprintf("%s/users/%s/repos?type=owner&sort=updated&per_page=20", c.apiURL, c.username)

	resp, err := c.http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch repos: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("github API returned status %d", resp.StatusCode)
	}

	var repos []repoResponse
	if err := json.NewDecoder(resp.Body).Decode(&repos); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	var projects []project.Project
	for _, r := range repos {
		if r.Fork {
			continue
		}
		updatedAt, _ := time.Parse(time.RFC3339, r.UpdatedAt)
		projects = append(projects, project.Project{
			Name:        r.Name,
			Description: r.Description,
			Language:    r.Language,
			Topics:      r.Topics,
			URL:         r.HTMLURL,
			Stars:       r.Stars,
			UpdatedAt:   updatedAt,
		})
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].UpdatedAt.After(projects[j].UpdatedAt)
	})

	return projects, nil
}
