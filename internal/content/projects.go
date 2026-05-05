// Package content holds the static editorial data the portfolio renders —
// the project list, hero copy, about copy, and contact channels.
package content

import "time"

// Project is one entry in the "selected work" list.
type Project struct {
	Name        string
	Description string
	Languages   []string
	URL         string
	UpdatedAt   time.Time
}

// Projects returns the curated project list, ordered by recency.
func Projects() []Project {
	return []Project{
		{
			Name:        "about-me",
			Description: "This site. A Go server rendering HTML templates with Datastar for interactivity and Tailwind v4 alongside an editorial design system.",
			Languages:   []string{"Go", "HTML", "CSS"},
			URL:         "https://github.com/adrock-miles/about-me",
			UpdatedAt:   time.Date(2026, 5, 4, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:        "poe-armory",
			Description: "A build and item planning toolkit for Path of Exile.",
			Languages:   []string{"TypeScript"},
			URL:         "https://github.com/adrock-miles/poe-armory",
			UpdatedAt:   time.Date(2026, 3, 24, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:        "go-laserbeak",
			Description: "A fast, concurrent web scraping framework written in Go.",
			Languages:   []string{"Go"},
			URL:         "https://github.com/adrock-miles/go-laserbeak",
			UpdatedAt:   time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:        "greenlight",
			Description: "A lightweight deployment approval and release gating tool.",
			Languages:   []string{"TypeScript"},
			URL:         "https://github.com/adrock-miles/greenlight",
			UpdatedAt:   time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:        "system-showcase",
			Description: "Interactive documentation site for system design patterns and architectures.",
			Languages:   []string{"TypeScript", "CSS"},
			URL:         "https://github.com/adrock-miles/system-showcase",
			UpdatedAt:   time.Date(2026, 2, 19, 0, 0, 0, 0, time.UTC),
		},
	}
}
