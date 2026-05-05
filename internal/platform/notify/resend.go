package notify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const resendAPIURL = "https://api.resend.com/emails"

// ResendSender posts emails to the Resend HTTP API.
type ResendSender struct {
	apiKey string
	client *http.Client
}

// NewResendSender builds a ResendSender with a sane default timeout.
func NewResendSender(apiKey string) *ResendSender {
	return &ResendSender{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}
}

type resendRequest struct {
	From    string   `json:"from"`
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Text    string   `json:"text"`
}

// SendEmail posts a single message to the Resend API.
func (r *ResendSender) SendEmail(ctx context.Context, from, to, subject, body string) error {
	payload, err := json.Marshal(resendRequest{
		From:    from,
		To:      []string{to},
		Subject: subject,
		Text:    body,
	})
	if err != nil {
		return fmt.Errorf("resend: marshal: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, resendAPIURL, bytes.NewReader(payload))
	if err != nil {
		return fmt.Errorf("resend: build request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+r.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return fmt.Errorf("resend: do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("resend: status %d: %s", resp.StatusCode, string(body))
	}
	return nil
}
