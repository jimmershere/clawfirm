package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/jimmershere/clawfirm/internal/installer"
	"github.com/jimmershere/clawfirm/internal/seed"
	"github.com/jimmershere/clawfirm/internal/tier"
)

// ---------- init ----------

func newInitCommand() *cobra.Command {
	var (
		tierName string
		profile  string
		domain   string
	)
	cmd := &cobra.Command{
		Use:   "init",
		Short: "Initialize a new ClawFirm install (writes ~/.clawfirm/config.yaml and a tier-appropriate seed)",
		RunE: func(cmd *cobra.Command, args []string) error {
			t, err := tier.Parse(tierName)
			if err != nil {
				return fmt.Errorf("invalid --tier: %w", err)
			}
			s, err := seed.Generate(t, profile, domain)
			if err != nil {
				return fmt.Errorf("seed generation: %w", err)
			}
			return seed.Write(s)
		},
	}
	cmd.Flags().StringVar(&tierName, "tier", "solo", "deployment tier: solo | smb | enterprise")
	cmd.Flags().StringVar(&profile, "profile", "default", "build profile: default | permissive-only | air-gap-fips | minimal")
	cmd.Flags().StringVar(&domain, "domain", "", "FQDN for ingress (SMB+ tiers)")
	return cmd
}

// ---------- up / down / status / doctor ----------

func newUpCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Bring up the ClawFirm stack for the configured tier",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := seed.Load()
			if err != nil {
				return err
			}
			drv, err := installer.For(cfg.Tier)
			if err != nil {
				return err
			}
			return drv.Up(cmd.Context(), cfg)
		},
	}
}

func newDownCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "down",
		Short: "Stop ClawFirm services (state preserved)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := seed.Load()
			if err != nil {
				return err
			}
			drv, err := installer.For(cfg.Tier)
			if err != nil {
				return err
			}
			return drv.Down(cmd.Context(), cfg)
		},
	}
}

func newStatusCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show service health, tier, and active governance",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := seed.Load()
			if err != nil {
				return err
			}
			drv, err := installer.For(cfg.Tier)
			if err != nil {
				return err
			}
			return drv.Status(cmd.Context(), cfg)
		},
	}
}

func newDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Diagnose common install / runtime problems",
		RunE: func(cmd *cobra.Command, args []string) error {
			return installer.Doctor(cmd.Context())
		},
	}
}

// ---------- tier promote/demote ----------

func newTierCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tier",
		Short: "Inspect or change the governance tier (0-3)",
	}
	cmd.AddCommand(
		&cobra.Command{
			Use:   "promote [n]",
			Short: "Promote to governance tier n (0-3)",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return installer.PromoteTier(cmd.Context(), args[0])
			},
		},
		&cobra.Command{
			Use:   "panic",
			Short: "Emergency demote to Tier 0 (chat-only); cancels in-flight tool calls",
			RunE: func(cmd *cobra.Command, args []string) error {
				return installer.PanicDemote(cmd.Context())
			},
		},
	)
	return cmd
}

// ---------- node, bundle, upgrade, secret, approval, sso, net, backup, logs ----------
// These are scaffolds for the 90-day MVP; bodies land sprint-by-sprint.

func newNodeCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "node", Short: "Cluster node management (SMB/Enterprise)"}
	cmd.AddCommand(
		stub("token", "Print the join token for new worker nodes"),
		stub("join", "Join this host to an existing cluster"),
		stub("list", "List nodes in the cluster"),
	)
	return cmd
}

func newBundleCommand() *cobra.Command {
	cmd := &cobra.Command{Use: "bundle", Short: "Build, verify, and apply air-gapped bundles"}
	cmd.AddCommand(
		stub("build", "Build a self-contained air-gap bundle (--tier, --profile, --models, --skills, --output)"),
		stub("verify", "Verify Sigstore signatures on a bundle"),
	)
	return cmd
}

func newUpgradeCommand() *cobra.Command {
	return stub("upgrade", "Pull latest signed images, verify, and rolling-restart")
}
func newSecretCommand() *cobra.Command {
	return stub("secret", "Manage credentials in OpenBao (set/get/rotate/list)")
}
func newApprovalCommand() *cobra.Command {
	return stub("approval", "Configure approval channels (Slack, Discord, Telegram, web)")
}
func newSSOCommand() *cobra.Command {
	return stub("sso", "Configure SSO via Authentik (Okta, Azure AD, Google Workspace, Keycloak)")
}
func newNetCommand() *cobra.Command {
	return stub("net", "Manage zero-trust networking (Headscale + Tailscale)")
}
func newBackupCommand() *cobra.Command {
	return stub("backup", "Backup and restore ClawFirm state to S3-compatible storage")
}
func newLogsCommand() *cobra.Command {
	return stub("logs", "Tail logs from one or more services")
}

func stub(name, short string) *cobra.Command {
	return &cobra.Command{
		Use:   name,
		Short: short,
		RunE: func(cmd *cobra.Command, args []string) error {
			return fmt.Errorf("`%s`: not yet implemented in MVP scaffold; see ROADMAP.md", name)
		},
	}
}
