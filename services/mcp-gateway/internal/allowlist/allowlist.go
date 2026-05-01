// Package allowlist parses and evaluates per-tool / per-agent allowlist rules.
package allowlist

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config is the on-disk allowlist configuration.
type Config struct {
	Version int          `yaml:"version"`
	Tools   []ToolRule   `yaml:"tools"`
	Agents  []AgentRule  `yaml:"agents"`
	Egress  []EgressRule `yaml:"egress"`
}

type ToolRule struct {
	Name        string   `yaml:"name"`
	Allowed     bool     `yaml:"allowed"`
	RequireTier int      `yaml:"require_governance_tier"` // minimum tier
	Tags        []string `yaml:"tags"`                    // safe|risky|dangerous|read|write|egress
}

type AgentRule struct {
	Identity string   `yaml:"identity"`
	Tools    []string `yaml:"tools"`
	MaxBudget int     `yaml:"max_budget_usd_per_day,omitempty"`
}

type EgressRule struct {
	Tool  string   `yaml:"tool"`
	Hosts []string `yaml:"hosts"` // exact or *.example.com
	Ports []int    `yaml:"ports"`
}

// Load reads YAML config from disk.
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var c Config
	return &c, yaml.Unmarshal(data, &c)
}

// AllowTool returns true iff the named tool is allowed for the given agent at the given tier.
func (c *Config) AllowTool(tool, agent string, tier int) bool {
	for _, t := range c.Tools {
		if t.Name == tool && t.Allowed && tier >= t.RequireTier {
			return true
		}
	}
	return false
}
