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
			Name:        "magc",
			Description: "Real-time blood glucose tracking with live alerts for friends and family. Integrates with CGMs like Dexcom and LibreLinkUp, plus native iOS, watchOS, and widget companions.",
			Languages:   []string{"Go", "Swift", "Postgres"},
			URL:         "https://magc.app/",
			UpdatedAt:   time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:        "about-me",
			Description: "This site. A Go server rendering HTML templates with Datastar for interactivity and Tailwind v4 alongside an editorial design system.",
			Languages:   []string{"Go", "HTML", "CSS"},
			URL:         "https://github.com/adrock-miles/about-me",
			UpdatedAt:   time.Date(2026, 5, 5, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:        "poe-armory",
			Description: "A Path of Exile character tracker. Imports builds and snapshots equipment, gems, skill trees, and ascendancies over time. Go + SQLite backend, React frontend.",
			Languages:   []string{"TypeScript", "Go"},
			URL:         "https://github.com/adrock-miles/poe-armory",
			UpdatedAt:   time.Date(2026, 3, 31, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:        "go-laserbeak",
			Description: "A Discord LLM bot in Go. Listens for a wake phrase in voice channels, transcribes speech via the OpenAI Whisper API, and replies through any OpenAI-compatible chat model.",
			Languages:   []string{"Go"},
			URL:         "https://github.com/adrock-miles/go-laserbeak",
			UpdatedAt:   time.Date(2026, 3, 2, 0, 0, 0, 0, time.UTC),
		},
		{
			Name:        "greenlight",
			Description: "A Markdown-driven portfolio template. React, TypeScript, Tailwind, and shadcn/ui — no CMS required.",
			Languages:   []string{"TypeScript", "CSS"},
			URL:         "https://github.com/adrock-miles/greenlight",
			UpdatedAt:   time.Date(2026, 2, 28, 0, 0, 0, 0, time.UTC),
		},
	}
}
