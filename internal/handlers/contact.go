package handlers

import (
	"bytes"
	"errors"
	"html/template"
	"log/slog"
	"net"
	"net/http"
	"strings"

	"github.com/adrock-miles/about-me/internal/platform/notify"
	"github.com/adrock-miles/about-me/internal/platform/ratelimit"
	"github.com/adrock-miles/about-me/internal/platform/sse"
)

// ContactForm holds the form's editable values plus an optional error
// message. Used both for initial render (zero value) and for re-renders
// after a validation or send failure (values preserved, Error set).
type ContactForm struct {
	Name    string
	Email   string
	Subject string
	Body    string
	Error   string
}

func (f ContactForm) validate() error {
	if strings.TrimSpace(f.Name) == "" {
		return errors.New("Name is required.")
	}
	if strings.TrimSpace(f.Email) == "" {
		return errors.New("Email is required.")
	}
	if strings.TrimSpace(f.Body) == "" {
		return errors.New("Message can't be empty.")
	}
	return nil
}

// Contact handles inbound contact-form submissions as a Datastar BFF.
type Contact struct {
	tmpl    *template.Template
	mailer  *notify.Mailer
	logger  *slog.Logger
	limiter *ratelimit.Limiter
}

// NewContact builds a Contact handler bound to the parsed template set.
// The limiter throttles submissions per client IP and is consulted before
// validation so spammers can't bypass it with malformed payloads.
func NewContact(tmpl *template.Template, mailer *notify.Mailer, logger *slog.Logger, limiter *ratelimit.Limiter) *Contact {
	return &Contact{tmpl: tmpl, mailer: mailer, logger: logger, limiter: limiter}
}

// Submit accepts a urlencoded form POST from Datastar and streams back
// SSE patches that swap the form area in place — either re-rendered with
// errors, or replaced with a thank-you message.
func (c *Contact) Submit(w http.ResponseWriter, r *http.Request) {
	// BFF guard: this endpoint exists only for the Datastar frontend.
	if r.Header.Get("Datastar-Request") == "" {
		http.Error(w, "this endpoint expects a Datastar request", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		c.logger.Warn("parse form", "error", err)
		c.patchForm(w, ContactForm{Error: "Couldn't read your submission — please try again."})
		return
	}

	form := ContactForm{
		Name:    r.FormValue("name"),
		Email:   r.FormValue("email"),
		Subject: r.FormValue("subject"),
		Body:    r.FormValue("body"),
	}

	ip := clientIP(r)
	if !c.limiter.Allow(ip) {
		c.logger.Warn("contact form rate limited", "ip", ip)
		form.Error = "Too many submissions just now — please wait a few minutes before trying again."
		c.patchForm(w, form)
		return
	}

	if err := form.validate(); err != nil {
		form.Error = err.Error()
		c.patchForm(w, form)
		return
	}

	msg := notify.Message{
		Name:    form.Name,
		Email:   form.Email,
		Subject: form.Subject,
		Body:    form.Body,
	}
	if err := c.mailer.SendContact(r.Context(), msg); err != nil {
		c.logger.Error("sending contact email", "error", err)
		form.Error = "Something went wrong on my end — please try again, or email directly."
		c.patchForm(w, form)
		return
	}

	c.patchThanks(w)
}

// patchForm streams the form fragment (preserving values + showing error).
func (c *Contact) patchForm(w http.ResponseWriter, form ContactForm) {
	c.patchFragment(w, "contact-form-area", form)
}

// patchThanks streams the thank-you fragment.
func (c *Contact) patchThanks(w http.ResponseWriter) {
	c.patchFragment(w, "contact-thanks", nil)
}

// clientIP returns the host portion of r.RemoteAddr (which chi.RealIP has
// already resolved from X-Forwarded-For when running behind a proxy).
func clientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func (c *Contact) patchFragment(w http.ResponseWriter, name string, data any) {
	var buf bytes.Buffer
	if err := c.tmpl.ExecuteTemplate(&buf, name, data); err != nil {
		c.logger.Error("rendering fragment", "name", name, "error", err)
		http.Error(w, "render error", http.StatusInternalServerError)
		return
	}
	if err := sse.New(w).PatchElements("", sse.ModeOuter, buf.String()); err != nil {
		c.logger.Error("writing sse", "error", err)
	}
}
