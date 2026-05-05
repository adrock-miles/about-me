package handlers

import "net/http"

// Health responds with a 200 so PaaS healthchecks (Railway, Fly) can probe.
func Health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
