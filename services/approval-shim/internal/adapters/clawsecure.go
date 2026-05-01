// Package adapters implements per-framework approval translators.
package adapters

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

// Backend is the abstract approval backend (today: ClawSecure REST).
type Backend interface {
	Submit(req ApprovalRequest) (ApprovalDecision, error)
}

type ApprovalRequest struct {
	Source   string         `json:"source"`   // "langgraph"|"n8n"|"dify"|"openhands"|"openclaw"
	Identity string         `json:"identity"`
	Tool     string         `json:"tool"`
	Args     map[string]any `json:"args"`
	Reason   string         `json:"reason"`
}

type ApprovalDecision struct {
	Approved bool   `json:"approved"`
	Reason   string `json:"reason"`
	BlockingTimeoutSec int `json:"blocking_timeout_sec"`
}

type clawSecureBackend struct {
	baseURL string
	client  *http.Client
}

func NewClawSecure(url string) Backend {
	return &clawSecureBackend{
		baseURL: url,
		client:  &http.Client{Timeout: 30 * time.Second},
	}
}

func (b *clawSecureBackend) Submit(req ApprovalRequest) (ApprovalDecision, error) {
	body, _ := json.Marshal(req)
	resp, err := b.client.Post(b.baseURL+"/api/approvals", "application/json", bytes.NewReader(body))
	if err != nil {
		return ApprovalDecision{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return ApprovalDecision{}, errors.New("clawsecure rejected")
	}
	var d ApprovalDecision
	_ = json.NewDecoder(resp.Body).Decode(&d)
	return d, nil
}
