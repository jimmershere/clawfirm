package adapters

import (
	"encoding/json"
	"net/http"
)

// DifyHandler is the HTTP endpoint that a Dify human-input node POSTs to.
// Returns the same payload shape as Dify's native human-input node.
func DifyHandler(b Backend) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req ApprovalRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		req.Source = "dify"
		dec, err := b.Submit(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		_ = json.NewEncoder(w).Encode(dec)
	})
}
