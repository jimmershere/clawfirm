// Package seed handles the on-disk config (~/.clawfirm/config.yaml) which
// drives every other subcommand. The format intentionally mirrors the
// SeedFile.md format from the clawboot repo so the same seed can be consumed
// by the bare-metal first-boot installer and the in-cluster CLI.
package seed

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/jimmershere/clawfirm/internal/tier"
)

// Seed is the on-disk representation of a ClawFirm install.
type Seed struct {
	Version    int          `yaml:"version"`
	Tier       tier.Tier    `yaml:"-"` // serialized as TierName
	TierName   string       `yaml:"tier"`
	Profile    string       `yaml:"profile"`
	Domain     string       `yaml:"domain,omitempty"`
	Governance Governance   `yaml:"governance"`
	Inference  Inference    `yaml:"inference"`
	Sandbox    string       `yaml:"sandbox"`
	Network    Network      `yaml:"network"`
	Skills     SkillsConfig `yaml:"skills"`
}

type Governance struct {
	DefaultTier int  `yaml:"default_tier"` // 0..3
	PanicHotkey bool `yaml:"panic_hotkey_enabled"`
}

type Inference struct {
	OllamaEnabled bool     `yaml:"ollama_enabled"`
	VLLMEnabled   bool     `yaml:"vllm_enabled"`
	DefaultModels []string `yaml:"default_models"`
	FrontierAllow []string `yaml:"frontier_allow"` // empty = deny-by-default
}

type Network struct {
	BindAddress  string `yaml:"bind_address"`     // 127.0.0.1 by default
	OverlayProvider string `yaml:"overlay_provider"` // "headscale" | "tailscale" | "none"
}

type SkillsConfig struct {
	RegistryURL string `yaml:"registry_url"`
	RequireSig  bool   `yaml:"require_signature"`
}

// Generate produces a tier-appropriate Seed with secure defaults.
func Generate(t tier.Tier, profile, domain string) (*Seed, error) {
	s := &Seed{
		Version:  1,
		Tier:     t,
		TierName: t.String(),
		Profile:  profile,
		Domain:   domain,
		Governance: Governance{
			DefaultTier: defaultGovTier(t),
			PanicHotkey: true,
		},
		Inference: Inference{
			OllamaEnabled: true,
			VLLMEnabled:   t != tier.Solo,
			DefaultModels: []string{"qwen3.5-4b"},
			FrontierAllow: nil, // deny-by-default
		},
		Sandbox: defaultSandbox(t),
		Network: Network{
			BindAddress:     "127.0.0.1",
			OverlayProvider: "none",
		},
		Skills: SkillsConfig{
			RegistryURL: "https://registry.clawfirm.io",
			RequireSig:  true,
		},
	}
	return s, nil
}

func defaultGovTier(t tier.Tier) int {
	switch t {
	case tier.Solo:
		return 0 // chat-only at first run
	case tier.SMB:
		return 1 // smart-approval
	case tier.Enterprise:
		return 1 // smart-approval base; Tier 2 enforced for write/egress via policy
	}
	return 0
}

func defaultSandbox(t tier.Tier) string {
	switch t {
	case tier.Solo:
		return "difysandbox+rootless"
	case tier.SMB:
		return "difysandbox+gvisor"
	case tier.Enterprise:
		return "difysandbox+gvisor+firecracker"
	}
	return "difysandbox+rootless"
}

// configPath returns ~/.clawfirm/config.yaml.
func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".clawfirm", "config.yaml"), nil
}

// Write persists a Seed to disk with mode 0600.
func Write(s *Seed) error {
	path, err := configPath()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	data, err := yaml.Marshal(s)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0o600); err != nil {
		return err
	}
	fmt.Println("wrote", path)
	return nil
}

// Load reads ~/.clawfirm/config.yaml.
func Load() (*Seed, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("no config at %s — did you run `clawfirm init`? (%w)", path, err)
	}
	var s Seed
	if err := yaml.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	t, err := tier.Parse(s.TierName)
	if err != nil {
		return nil, err
	}
	s.Tier = t
	return &s, nil
}
