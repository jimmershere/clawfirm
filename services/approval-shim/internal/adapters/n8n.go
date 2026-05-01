package adapters

import (
	"encoding/json"
	"net/http"
)

// N8NHandler returns an http.Handler that n8n's HTTP Request node calls
// from inside an Approve workflow step.
func N8NHandler(b Backend) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ApprovalRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.Source = "n8n"
		dec, err := b.Submit(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		_ = json.NewEncoder(w).Encode(dec)
	})
}
