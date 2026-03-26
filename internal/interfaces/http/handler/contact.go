package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adrock-miles/about-me/internal/application/contact"
	contactDomain "github.com/adrock-miles/about-me/internal/domain/contact"
)

// ContactHandler handles contact-related HTTP requests.
type ContactHandler struct {
	service *contact.Service
}

// NewContactHandler creates a new contact handler.
func NewContactHandler(service *contact.Service) *ContactHandler {
	return &ContactHandler{service: service}
}

// SendMessage handles a contact form submission.
func (h *ContactHandler) SendMessage(w http.ResponseWriter, r *http.Request) {
	var msg contactDomain.Message
	if err := json.NewDecoder(r.Body).Decode(&msg); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if err := h.service.SendMessage(&msg); err != nil {
		log.Printf("Error sending message: %v", err)
		http.Error(w, "Failed to send message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "sent"})
}
