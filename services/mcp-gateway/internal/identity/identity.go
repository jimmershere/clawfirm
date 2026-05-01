// Package identity verifies OIDC tokens from Authentik and issues
// short-lived per-request tokens for upstream MCP servers (no passthrough).
package identity

import "context"

type Verifier interface {
	Verify(ctx context.Context, bearer string) (*Claims, error)
	MintForTool(ctx context.Context, claims *Claims, tool string) (string, error)
}

type Claims struct {
	Subject string
	Email   string
	Groups  []string
	Tier    int
}

// NewVerifier connects to the OIDC issuer and returns a Verifier.
// MVP scaffold returns a no-op verifier when issuer == "".
func NewVerifier(issuer string) (Verifier, error) {
	if issuer == "" {
		return &noopVerifier{}, nil
	}
	// TODO: actual go-oidc provider setup.
	return &noopVerifier{}, nil
}

type noopVerifier struct{}

func (noopVerifier) Verify(_ context.Context, _ string) (*Claims, error) {
	return &Claims{Subject: "anonymous", Tier: 0}, nil
}
func (noopVerifier) MintForTool(_ context.Context, c *Claims, _ string) (string, error) {
	return "stub-token-for-" + c.Subject, nil
}
