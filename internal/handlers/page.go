// Package handlers wires the HTTP endpoints to the page templates and the
// underlying services (mailer, content).
package handlers

import (
	"html/template"
	"log/slog"
	"net/http"

	"github.com/adrock-miles/about-me/internal/content"
)

// HeroMeta is one label/value pair in the hero meta strip.
type HeroMeta struct{ Label, Value string }

// HeroContent holds editorial copy for the hero block.
type HeroContent struct {
	Headline template.HTML
	Tagline  string
	Meta     []HeroMeta
}

// Stat is a label/value pair in the about side rail.
type Stat struct{ Label, Value string }

// AboutContent holds the about block.
type AboutContent struct {
	Paragraphs []string
	Stats      []Stat
}

// ContactLink is one outbound channel link.
type ContactLink struct {
	Label    string
	Value    string
	Href     string
	External bool
}

// ContactContent holds the contact block.
type ContactContent struct {
	Email string
	Links []ContactLink
	Form  ContactForm
}

// PageData is the shape passed to the index template.
type PageData struct {
	Title       string
	Description string
	Name        string
	BuildTag    string
	Hero        HeroContent
	Projects    []content.Project
	About       AboutContent
	Contact     ContactContent
}

// Page renders top-level HTML pages.
type Page struct {
	tmpl   *template.Template
	logger *slog.Logger
}

// NewPage builds a Page handler around a parsed template set.
func NewPage(tmpl *template.Template, logger *slog.Logger) *Page {
	return &Page{tmpl: tmpl, logger: logger}
}

// Index renders the portfolio home page.
func (p *Page) Index(w http.ResponseWriter, r *http.Request) {
	data := PageData{
		Title:       "Adam Miles — Software Engineer",
		Description: "Adam Miles is a software engineer crafting calm, considered tools.",
		Name:        "Adam Miles",
		BuildTag:    "Last updated May 2026",
		Hero: HeroContent{
			Headline: template.HTML("Building for <em>function</em><br>and form."),
			Tagline:  "Adam Miles is a software engineer crafting calm, considered tools — where the seams hold up under pressure and the edges still feel hand-finished.",
			Meta: []HeroMeta{
				{"Currently", "Building independent"},
				{"Based in", "Oklahoma City, OK"},
				{"Available", "Q3 2026"},
			},
		},
		Projects: content.Projects(),
		About: AboutContent{
			Paragraphs: []string{
				"I've spent sixteen years writing software that has to <em>actually work</em> — for people, in the wild, on the worst day of their week. Systems people lean on quietly.",
				"I move between the runtime and the interface comfortably. There's no silver bullet for building systems — but there are plenty of bad decisions you can avoid. The right answer is usually a structural one, not a feature.",
				"Outside of work I read, skate, play guitar, and spend time on anything creative — which is where the love of <em>silhouette and weight</em> in this site comes from.",
			},
			Stats: []Stat{
				{"Years building", "16"},
				{"Languages", "Go · PHP · Swift"},
				{"Currently", "Independent"},
				{"Based in", "Oklahoma City, OK"},
				{"Stack here", "Go · Datastar · Tailwind v4"},
				{"Reading now", "The Name of the Wind"},
			},
		},
		Contact: ContactContent{
			Email: "miles.adrock@gmail.com",
			Links: []ContactLink{
				{Label: "Email", Value: "miles.adrock@gmail.com", Href: "mailto:miles.adrock@gmail.com"},
				{Label: "GitHub", Value: "@adrock-miles", Href: "https://github.com/adrock-miles", External: true},
				{Label: "LinkedIn", Value: "/in/adammiles", Href: "https://www.linkedin.com/in/adammiles", External: true},
				{Label: "Instagram", Value: "@adrock.miles", Href: "https://www.instagram.com/adrock.miles", External: true},
			},
		},
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := p.tmpl.ExecuteTemplate(w, "index.html", data); err != nil {
		p.logger.Error("rendering index template", "error", err)
	}
}
