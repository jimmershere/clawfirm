// Package tier defines the three ClawFirm deployment tiers.
package tier

import "fmt"

type Tier int

const (
	Solo Tier = iota
	SMB
	Enterprise
)

func (t Tier) String() string {
	switch t {
	case Solo:
		return "solo"
	case SMB:
		return "smb"
	case Enterprise:
		return "enterprise"
	default:
		return "unknown"
	}
}

// Parse converts a CLI flag value into a Tier.
func Parse(s string) (Tier, error) {
	switch s {
	case "solo":
		return Solo, nil
	case "smb":
		return SMB, nil
	case "enterprise":
		return Enterprise, nil
	default:
		return 0, fmt.Errorf("unknown tier %q (expected solo|smb|enterprise)", s)
	}
}
