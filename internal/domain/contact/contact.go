package contact

import "errors"

// Message represents a contact form submission.
type Message struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Validate checks that required fields are present.
func (m *Message) Validate() error {
	if m.Name == "" {
		return errors.New("name is required")
	}
	if m.Email == "" {
		return errors.New("email is required")
	}
	if m.Body == "" {
		return errors.New("message body is required")
	}
	return nil
}

// Sender defines the contract for sending contact messages.
type Sender interface {
	Send(msg *Message) error
}
