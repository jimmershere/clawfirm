// Command clawfirm is the operator CLI for the ClawFirm AI factory.
//
// Usage examples:
//
//	clawfirm init --tier solo
//	clawfirm up
//	clawfirm tier promote 1
//	clawfirm node join <master-ip> <token>
//	clawfirm bundle build --tier enterprise --output clawfirm.airgap.tar.zst
//	clawfirm doctor
//
// See `clawfirm help` for the full command surface.
package main

import (
	"fmt"
	"os"

	"github.com/jimmershere/clawfirm/internal/cli"
	"github.com/jimmershere/clawfirm/internal/version"
)

func main() {
	root := cli.NewRootCommand(version.Version, version.Commit, version.Date)
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, "clawfirm:", err)
		os.Exit(1)
	}
}
