package project

import "time"

// Project represents a GitHub project in the domain.
type Project struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Languages   []string  `json:"languages"`
	Topics      []string  `json:"topics"`
	URL         string    `json:"url"`
	Stars       int       `json:"stars"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// Repository defines the contract for fetching projects.
type Repository interface {
	FetchRecentProjects() ([]Project, error)
}
