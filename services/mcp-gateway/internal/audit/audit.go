// Package audit posts every gateway decision to ClawSecure with a hash chain.
package audit

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"sync"
	"time"
)

type Sink interface {
	Record(event Event) error
}

type Event struct {
	Timestamp time.Time         `json:"ts"`
	Identity  string            `json:"identity"`
	Tool      string            `json:"tool"`
	Args      map[string]any    `json:"args"`
	Decision  string            `json:"decision"` // allow|deny|prompt
	Reason    string            `json:"reason"`
	PrevHash  string            `json:"prev_hash"`
	Hash      string            `json:"hash"`
}

type clawSecureSink struct {
	url      string
	mu       sync.Mutex
	prevHash string
	client   *http.Client
}

func NewClawSecureSink(url string) Sink {
	return &clawSecureSink{
		url:    url,
		client: &http.Client{Timeout: 5 * time.Second},
	}
}

func (s *clawSecureSink) Record(e Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	e.Timestamp = time.Now().UTC()
	e.PrevHash = s.prevHash
	e.Hash = chain(e)
	s.prevHash = e.Hash

	body, _ := json.Marshal(e)
	req, _ := http.NewRequest(http.MethodPost, s.url, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		// Audit must not block the request path. Buffer locally for retry.
		// (Disk-spool implementation lands in Weeks 7-8 of the MVP.)
		return nil
	}
	defer resp.Body.Close()
	return nil
}

func chain(e Event) string {
	canonical, _ := json.Marshal(struct {
		PrevHash string         `json:"prev_hash"`
		Identity string         `json:"identity"`
		Tool     string         `json:"tool"`
		Args     map[string]any `json:"args"`
		Decision string         `json:"decision"`
		Reason   string         `json:"reason"`
	}{e.PrevHash, e.Identity, e.Tool, e.Args, e.Decision, e.Reason})
	sum := sha256.Sum256(canonical)
	return hex.EncodeToString(sum[:])
}
