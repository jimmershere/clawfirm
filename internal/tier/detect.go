package tier

import (
	"os"
	"os/exec"
	"runtime"
)

// Capabilities describes what the host can run. Populated by Detect().
type Capabilities struct {
	NumCPU          int
	GoArch          string
	GoOS            string
	HasNvidiaGPU    bool // populated via nvidia-smi probe
	HasKVM          bool // /dev/kvm
	IsAppleSilicon  bool
}

// Detect returns a best-effort summary of the host's capabilities. The CLI
// uses this to recommend a tier and a profile; the operator still chooses.
func Detect() Capabilities {
	return Capabilities{
		NumCPU:         runtime.NumCPU(),
		GoArch:         runtime.GOARCH,
		GoOS:           runtime.GOOS,
		HasNvidiaGPU:   probeNvidia(),
		HasKVM:         probeKVM(),
		IsAppleSilicon: runtime.GOOS == "darwin" && runtime.GOARCH == "arm64",
	}
}

func probeNvidia() bool {
	path, err := exec.LookPath("nvidia-smi")
	if err != nil {
		return false
	}
	cmd := exec.Command(path, "-L")
	if err := cmd.Run(); err == nil {
		return true
	}
	return false
}

func probeKVM() bool {
	_, err := os.Stat("/dev/kvm")
	return err == nil
}
