package installer

import (
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/jimmershere/clawfirm/internal/seed"
)

// composeDriver brings up the Solo tier via `docker compose`.
type composeDriver struct{}

func (d *composeDriver) Up(ctx context.Context, cfg *seed.Seed) error {
	args := []string{"compose", "-f", "deploy/solo/docker-compose.yml"}
	if cfg.Inference.VLLMEnabled {
		args = append(args, "-f", "deploy/solo/docker-compose.gpu.yml")
	}
	args = append(args, "up", "-d")
	return run(ctx, "docker", args...)
}

func (d *composeDriver) Down(ctx context.Context, cfg *seed.Seed) error {
	args := []string{"compose", "-f", "deploy/solo/docker-compose.yml"}
	if cfg.Inference.VLLMEnabled {
		args = append(args, "-f", "deploy/solo/docker-compose.gpu.yml")
	}
	args = append(args, "down")
	return run(ctx, "docker", args...)
}

func (d *composeDriver) Status(ctx context.Context, cfg *seed.Seed) error {
	fmt.Printf("ClawFirm Solo  tier=%s  governance=Tier %d  sandbox=%s\n",
		cfg.TierName, cfg.Governance.DefaultTier, cfg.Sandbox)
	return run(ctx, "docker", "compose", "-f", "deploy/solo/docker-compose.yml", "ps")
}

// run executes a shell command, streaming stdout/stderr to the user.
func run(ctx context.Context, name string, args ...string) error {
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
