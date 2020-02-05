package vipser

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Just to compile in if needed
const Vipser string = ""

func verifyFile(p string) error {
	if fi, err := os.Stat(p); err != nil {
		return fmt.Errorf("%s is not a valid executable", p)
	} else if fi.IsDir() {
		return fmt.Errorf("%s is a directory", p)
	} else if fi.Mode() & os.FileMode(0111) == 0 {
		return fmt.Errorf("%s does not have exec permissions", p)
	}
	return nil
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
				return "", err
			}
			return vipser, verifyFile(vipser)
		} else {
			return prog, verifyFile(prog)
		}
	}
}

