package project

import (
	"time"

	domain "github.com/adrock-miles/about-me/internal/domain/project"
)

// StaticRepository returns a hardcoded list of recent projects.
type StaticRepository struct{}

// NewStaticRepository creates a new static project repository.
func NewStaticRepository() *StaticRepository {
	return &StaticRepository{}
}

// FetchRecentProjects returns the hardcoded project list.
func (r *StaticRepository) FetchRecentProjects() ([]domain.Project, error) {
	return []domain.Project{
		{
			Name:        "poe-armory",
			Description: "A build and item planning toolkit for Path of Exile.",
			Languages:   []string{"TypeScript"},
			Topics:      nil,
			URL:         "https://github.com/adrock-miles/poe-armory",
			Stars:       0,
			UpdatedAt:   time.Date(2026, 3, 24, 20, 18, 18, 0, time.UTC),
		},
		{
			Name:        "go-laserbeak",
			Description: "A fast, concurrent web scraping framework written in Go.",
			Languages:   []string{"Go", "Makefile", "Dockerfile"},
			Topics:      nil,
			URL:         "https://github.com/adrock-miles/go-laserbeak",
			Stars:       1,
			UpdatedAt:   time.Date(2026, 3, 2, 17, 59, 10, 0, time.UTC),
		},
		{
			Name:        "greenlight",
			Description: "A lightweight deployment approval and release gating tool.",
			Languages:   []string{"TypeScript"},
			Topics:      nil,
			URL:         "https://github.com/adrock-miles/greenlight",
			Stars:       0,
			UpdatedAt:   time.Date(2026, 2, 28, 23, 19, 19, 0, time.UTC),
		},
		{
			Name:        "system-showcase",
			Description: "Interactive documentation site for system design patterns and architectures.",
			Languages:   []string{"TypeScript", "CSS"},
			Topics:      nil,
			URL:         "https://github.com/adrock-miles/system-showcase",
			Stars:       0,
			UpdatedAt:   time.Date(2026, 2, 19, 4, 42, 36, 0, time.UTC),
		},
		{
			Name:        "clawd",
			Description: "Clawing back your time with simple task management.",
			Languages:   []string{"Go"},
			Topics:      nil,
			URL:         "https://github.com/adrock-miles/clawd",
			Stars:       0,
			UpdatedAt:   time.Date(2026, 2, 17, 0, 12, 8, 0, time.UTC),
		},
	}, nil
}
