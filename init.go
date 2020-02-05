package vipser

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

var Vipser string

func verifyFile(p string) {
	if fi, err := os.Stat(p); err != nil {
		panic(fmt.Errorf("%s is not a valid executable", p))
	} else if fi.IsDir() {
		panic(fmt.Errorf("%s is a directory", p))
	} else if fi.Mode() & os.FileMode(0111) == 0 {
		panic(fmt.Errorf("%s does not have exec permissions", p))
	}
}

// VipserInit finds the vipser binary and stores it
func VipserInit() {
	if Vipser != "" {
		verifyFile(Vipser)
	} else if v, ok := os.LookupEnv("VIPSER"); ok {
		verifyFile(v)
		Vipser = v
	} else {
		if prog, err := exec.LookPath("vipser"); err != nil {
			panic(err)
		} else if prog == "" {
			verifyFile("vipser")
			Vipser = "vipser"
		} else {
			verifyFile(prog)
			Vipser = prog
		}
	}
	log.Printf("using vipser %s", Vipser)
}

