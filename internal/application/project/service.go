package project

import (
	"github.com/adrock-miles/about-me/internal/domain/project"
)

// Service handles project-related business logic.
type Service struct {
	repo project.Repository
}

// NewService creates a new project service.
func NewService(repo project.Repository) *Service {
	return &Service{repo: repo}
}

// GetRecentProjects returns recently updated GitHub projects.
func (s *Service) GetRecentProjects() ([]project.Project, error) {
	return s.repo.FetchRecentProjects()
}
