package vipser_test

import (
	"bytes"
	vipser "github.com/jhford/govipser"
	"github.com/stretchr/testify/assert"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestMakeVipserBin(t *testing.T) {
	if runtime.GOARCH != "amd64" {
		t.Skip("skipping MakeVipserBin")
	}
	path, err := vipser.MakeVipserBin()
	assert.NoError(t, err)
	assert.Contains(t, path, "vipser")
	assert.FileExists(t, path)

	_, err = filepath.Abs(path)
	assert.NoError(t, err)

	var out bytes.Buffer
	cmd := exec.Command("cat")
	cmd.Stdin = bytes.NewBuffer(testinput)
	cmd.Stdout = &out

	err = cmd.Run()
	assert.NoError(t, err)
	assert.Equal(t, testinput, out.Bytes())
}
