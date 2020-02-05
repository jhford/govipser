package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

func MustEnv(name string) string {
	if pkg, ok := os.LookupEnv(name); ok {
		return strings.TrimSpace(pkg)
	} else {
		panic(errors.New(name + " env var required"))
	}
}

type SourceFile struct {
	Package string
	Variable string
	Input io.Reader
	Output io.Writer
}

func (s SourceFile) Printf(msg string, args ...interface{}) {
	_, err := fmt.Fprintf(s.Output, msg, args...)
	if err != nil {
		panic(err)
	}
}

func (s SourceFile) Run() {
	s.Printf("package %s\n\n", s.Package)

	s.Printf("var %s = []byte{", s.Variable)

	buf := make([]byte, 4096)

	totalBytes := 0

	for n, err := s.Input.Read(buf); err != io.EOF; n, err = s.Input.Read(buf) {
		for i := 0; i < n; i++ {
			s.Printf("0x%.2X,", buf[i])

			if (totalBytes + i) % 10 == 0 {
				s.Printf("\n  ")
			} else {
				s.Printf(" ")
			}
		}

		totalBytes += n
	}


	s.Printf("\n}\n")
}

func main() {
	s := SourceFile{
		Package: MustEnv("PACKAGE"),
		Variable: MustEnv("VAR"),
		Input: os.Stdin,
		Output: os.Stdout,
	}

	s.Run()
}
