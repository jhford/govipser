//go:generate bash -c "PACKAGE=vipser VAR=VipserLinux64 go run cmd/embed/embed.go < vipser > vipser64.go"

package vipser

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// Just to compile in if needed
const Vipser string = ""

func verifyFile(p string) error {
	if fi, err := os.Stat(p); err != nil {
		return fmt.Errorf("%s is not a valid executable", p)
	} else if fi.IsDir() {
		return fmt.Errorf("%s is a directory", p)
	} else if fi.Mode()&os.FileMode(0111) == 0 {
		return fmt.Errorf("%s does not have exec permissions", p)
	}
	return nil
}

// MakeVipserBin creates a Vipser binary in the same directory if it knows
// how to, otherwise returns an error
func MakeVipserBin() (string, error) {
	switch runtime.GOARCH {
	case "amd64":
		return "vipser-linux-amd64", ioutil.WriteFile("vipser-linux-amd64", vipser64, 0755)
	default:
		return "", fmt.Errorf("cannot create vipser binary for %s", runtime.GOARCH)
	}
}

// VipserInit finds the vipser binary and stores it
func FindVipser() (string, error) {
	if Vipser != "" {
		return Vipser, verifyFile(Vipser)
	} else if v, ok := os.LookupEnv("VIPSER"); ok {
		return v, verifyFile(v)
	} else {
		if prog, _ := exec.LookPath("vipser"); prog == "" {
			vipser, err := filepath.Abs("vipser")
			if err != nil {
				return MakeVipserBin()
			}
			return vipser, verifyFile(vipser)
		} else {
			return prog, verifyFile(prog)
		}
	}
}
