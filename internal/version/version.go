// Package version exposes build-time version metadata. Populated by
// linker flags in the release pipeline:
//
//	go build -ldflags "-X github.com/jimmershere/clawfirm/internal/version.Version=$VERSION ..."
package version

var (
	Version = "dev"
	Commit  = "unknown"
	Date    = "unknown"
)
