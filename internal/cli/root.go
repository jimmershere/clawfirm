// Package cli wires together the clawfirm CLI subcommands.
package cli

import (
	"github.com/spf13/cobra"
)

// NewRootCommand returns the top-level `clawfirm` command with all subcommands attached.
func NewRootCommand(ver, commit, date string) *cobra.Command {
	root := &cobra.Command{
		Use:           "clawfirm",
		Short:         "ClawFirm — an AI factory in a box",
		Long:          longDescription,
		Version:       ver + " (" + commit + " @ " + date + ")",
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	// Global flags
	root.PersistentFlags().String("config", "", "config file (default ~/.clawfirm/config.yaml)")
	root.PersistentFlags().Bool("verbose", false, "verbose logging")

	// Subcommands
	root.AddCommand(
		newInitCommand(),
		newUpCommand(),
		newDownCommand(),
		newStatusCommand(),
		newDoctorCommand(),
		newTierCommand(),
		newNodeCommand(),
		newBundleCommand(),
		newUpgradeCommand(),
		newSecretCommand(),
		newApprovalCommand(),
		newSSOCommand(),
		newNetCommand(),
		newBackupCommand(),
		newLogsCommand(),
	)

	return root
}

const longDescription = `ClawFirm composes:
  - clawrails  - inference router + cost layer
  - clawboot   - bare-metal installer
  - clawsecure - approvals + policy + audit
  ...with a curated upstream stack (Ollama, vLLM, Dify, n8n, LangGraph,
  OpenHands, OpenBao, Authentik, Headscale, Langfuse) into a tier-aware
  AI platform that you can run on a laptop, a 3-node SMB cluster, or a
  full Enterprise HA topology with air-gap support.

  Defaults are secure: chat-only at first launch, sandbox-first, default-deny
  egress, signed skills, mandatory audit log. Promote up the four-tier
  governance ladder when you're ready.

  See https://docs.clawfirm.io for the full guide.`
