// Package installer contains the per-tier deploy drivers (Compose for solo,
// k3s for SMB, k0s+Flux for Enterprise) behind a single Driver interface.
package installer

import (
	"context"
	"fmt"

	"github.com/jimmershere/clawfirm/internal/seed"
	"github.com/jimmershere/clawfirm/internal/tier"
)

// Driver is the interface every per-tier installer implements.
type Driver interface {
	Up(ctx context.Context, cfg *seed.Seed) error
	Down(ctx context.Context, cfg *seed.Seed) error
	Status(ctx context.Context, cfg *seed.Seed) error
}

// For returns the right Driver for the given tier.
func For(t tier.Tier) (Driver, error) {
	switch t {
	case tier.Solo:
		return &composeDriver{}, nil
	case tier.SMB:
		return &k3sDriver{}, nil
	case tier.Enterprise:
		return &k0sDriver{}, nil
	}
	return nil, fmt.Errorf("no driver for tier %s", t)
}

// Doctor runs cross-tier sanity checks (Docker socket, port conflicts, disk space, etc.)
func Doctor(ctx context.Context) error {
	// TODO: implement. Stubbed for the MVP scaffold.
	fmt.Println("clawfirm doctor: scaffold; see scripts/doctor.sh for the current set of checks")
	return nil
}

// PromoteTier moves the system to the target governance tier (0..3).
// Calls into ClawSecure to update the active rule bundle.
func PromoteTier(ctx context.Context, target string) error {
	fmt.Printf("clawfirm tier promote %s: scaffold (would POST to clawsecure /api/bundles/clawfirm-<tier>/apply)\n", target)
	return nil
}

// PanicDemote forces the system to Tier 0 immediately and drains the approval queue.
func PanicDemote(ctx context.Context) error {
	fmt.Println("clawfirm tier panic: scaffold (would force-cancel in-flight tool calls and apply clawfirm-tier0 bundle)")
	return nil
}
