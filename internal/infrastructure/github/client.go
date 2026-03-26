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
	LanguagesURL string  `json:"languages_url"`
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

const maxProjects = 5

// FetchRecentProjects fetches the user's public repos from GitHub.
func (c *Client) FetchRecentProjects() ([]project.Project, error) {
	url := fmt.Sprintf("%s/users/%s/repos?type=owner&sort=updated&per_page=10", c.apiURL, c.username)

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

		languages := c.fetchLanguages(r.LanguagesURL)

		updatedAt, _ := time.Parse(time.RFC3339, r.UpdatedAt)
		projects = append(projects, project.Project{
			Name:        r.Name,
			Description: r.Description,
			Languages:   languages,
			Topics:      r.Topics,
			URL:         r.HTMLURL,
			Stars:       r.Stars,
			UpdatedAt:   updatedAt,
		})
	}

	sort.Slice(projects, func(i, j int) bool {
		return projects[i].UpdatedAt.After(projects[j].UpdatedAt)
	})

	if len(projects) > maxProjects {
		projects = projects[:maxProjects]
	}

	return projects, nil
}

// fetchLanguages calls the repo languages endpoint and returns sorted language names.
func (c *Client) fetchLanguages(url string) []string {
	if url == "" {
		return nil
	}

	resp, err := c.http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil
	}

	// GitHub returns {"Go": 12345, "TypeScript": 6789, ...}
	var langBytes map[string]int
	if err := json.NewDecoder(resp.Body).Decode(&langBytes); err != nil {
		return nil
	}

	// Sort by bytes descending
	type langEntry struct {
		name  string
		bytes int
	}
	entries := make([]langEntry, 0, len(langBytes))
	for name, b := range langBytes {
		entries = append(entries, langEntry{name, b})
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].bytes > entries[j].bytes
	})

	languages := make([]string, len(entries))
	for i, e := range entries {
		languages[i] = e.name
	}
	return languages
}
