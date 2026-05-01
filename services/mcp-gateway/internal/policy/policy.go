// Package policy is the request decision engine for the MCP gateway.
package policy

import (
	"encoding/json"
	"net/http"

	"github.com/jimmershere/clawfirm/services/mcp-gateway/internal/allowlist"
	"github.com/jimmershere/clawfirm/services/mcp-gateway/internal/audit"
	"github.com/jimmershere/clawfirm/services/mcp-gateway/internal/identity"
)

type Engine struct {
	allow   *allowlist.Config
	id      identity.Verifier
	auditor audit.Sink
}

func New(a *allowlist.Config, id identity.Verifier, s audit.Sink) *Engine {
	return &Engine{allow: a, id: id, auditor: s}
}

// Handle is the HTTP handler that fronts every tool call.
func (e *Engine) Handle() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Tool string         `json:"tool"`
			Args map[string]any `json:"args"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad request", http.StatusBadRequest)
			return
		}

		claims, err := e.id.Verify(r.Context(), bearer(r))
		if err != nil {
			http.Error(w, "unauthenticated", http.StatusUnauthorized)
			return
		}

		decision := "deny"
		reason := "tool not in allowlist"
		if e.allow.AllowTool(req.Tool, claims.Subject, claims.Tier) {
			decision = "allow"
			reason = "ok"
		}

		_ = e.auditor.Record(audit.Event{
			Identity: claims.Subject,
			Tool:     req.Tool,
			Args:     req.Args,
			Decision: decision,
			Reason:   reason,
		})

		if decision != "allow" {
			http.Error(w, reason, http.StatusForbidden)
			return
		}

		// Mint a short-lived per-tool token (NOT the user's bearer).
		_, _ = e.id.MintForTool(r.Context(), claims, req.Tool)

		// Forward to the upstream MCP server. Wiring lands in Weeks 7-8.
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok","note":"upstream MCP forward not wired in MVP scaffold"}`))
	})
}

func bearer(r *http.Request) string {
	h := r.Header.Get("Authorization")
	if len(h) > 7 && h[:7] == "Bearer " {
		return h[7:]
	}
	return ""
}
