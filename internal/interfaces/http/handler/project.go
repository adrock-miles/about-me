package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/adrock-miles/about-me/internal/application/project"
)

// ProjectHandler handles project-related HTTP requests.
type ProjectHandler struct {
	service *project.Service
}

// NewProjectHandler creates a new project handler.
func NewProjectHandler(service *project.Service) *ProjectHandler {
	return &ProjectHandler{service: service}
}

// GetProjects returns the list of recent projects.
func (h *ProjectHandler) GetProjects(w http.ResponseWriter, r *http.Request) {
	projects, err := h.service.GetRecentProjects()
	if err != nil {
		log.Printf("Error fetching projects: %v", err)
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}
