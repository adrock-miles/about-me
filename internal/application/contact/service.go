package contact

import (
	"github.com/adrock-miles/about-me/internal/domain/contact"
)

// Service handles contact-related business logic.
type Service struct {
	sender contact.Sender
	toAddr string
}

// NewService creates a new contact service.
func NewService(sender contact.Sender, toAddr string) *Service {
	return &Service{sender: sender, toAddr: toAddr}
}

// SendMessage validates and sends a contact message.
func (s *Service) SendMessage(msg *contact.Message) error {
	if err := msg.Validate(); err != nil {
		return err
	}
	return s.sender.Send(msg)
}
