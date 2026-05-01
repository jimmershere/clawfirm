package installer

import (
	"context"
	"fmt"

	"github.com/jimmershere/clawfirm/internal/seed"
)

// k3sDriver applies the SMB Kustomize overlay to a local or remote k3s cluster.
type k3sDriver struct{}

func (d *k3sDriver) Up(ctx context.Context, cfg *seed.Seed) error {
	// Stage 1 (MVP): pure kubectl apply -k.
	// Stage 2: replace with controller-runtime client to do dependency
	// ordering, sealed-secret unsealing, and Helm chart application.
	return run(ctx, "kubectl", "apply", "-k", "deploy/smb/k3s")
}

func (d *k3sDriver) Down(ctx context.Context, cfg *seed.Seed) error {
	return run(ctx, "kubectl", "delete", "-k", "deploy/smb/k3s", "--ignore-not-found=true")
}

func (d *k3sDriver) Status(ctx context.Context, cfg *seed.Seed) error {
	fmt.Printf("ClawFirm SMB  tier=%s  governance=Tier %d  sandbox=%s\n",
		cfg.TierName, cfg.Governance.DefaultTier, cfg.Sandbox)
	return run(ctx, "kubectl", "-n", "clawfirm", "get", "pods,svc,sts")
}

// k0sDriver bootstraps the Enterprise tier via k0sctl + Flux GitOps.
type k0sDriver struct{}

func (d *k0sDriver) Up(ctx context.Context, cfg *seed.Seed) error {
	// MVP scaffold: emit instructions, defer to k0sctl + flux bootstrap manually.
	fmt.Println(`Enterprise Up is a multi-step operation:

  1. Edit deploy/enterprise/k0s/k0sctl.yaml.example for your topology.
  2. k0sctl apply --config deploy/enterprise/k0s/k0sctl.yaml
  3. flux bootstrap git --url=ssh://git@... --branch=main --path=clusters/clawfirm

  This will be wrapped into 'clawfirm up' in the Q3-26 milestone.`)
	return nil
}

func (d *k0sDriver) Down(ctx context.Context, cfg *seed.Seed) error {
	return fmt.Errorf("Enterprise teardown is destructive and intentionally not automated; see docs/operations/teardown.md")
}

func (d *k0sDriver) Status(ctx context.Context, cfg *seed.Seed) error {
	return run(ctx, "kubectl", "-n", "clawfirm", "get", "pods,svc,sts,hpa")
}
