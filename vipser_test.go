package vipser_test

import (
	"bytes"
	"github.com/jhford/govipser"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCommand_String(t *testing.T) {
	assert.Equal(t, "CMD,-1,-32,-64,1,32,64,0.1,0.1", vipser.Command{
		"CMD",
		int(-1),
		int32(-32),
		int64(-64),
		uint(1),
		uint32(32),
		uint64(64),
		float32(0.1),
		float64(0.1),
	}.String())
}

func TestOperation_RenderCommands(t *testing.T) {
	t.Run("with commands", func(t *testing.T) {
		o := vipser.New()
		o.Resize(1,2)
		o.Stretch(3,4)
		assert.Equal(t, []string{"RESIZE,1,2", "STRETCH,3,4"}, o.RenderCommands())
	})

	t.Run("no commands", func(t *testing.T) {
		o := vipser.New()
		assert.Equal(t, []string{}, o.RenderCommands())
	})
}

func TestOperation_AllCommands(t *testing.T) {
	tests := map[string]func(*vipser.Operation){
		"RESIZE,1,2": func(o *vipser.Operation) { o.Resize(1,2) },
		"STRETCH,1,2": func(o *vipser.Operation) { o.Stretch(1,2) },
	}

	for expected, doer := range tests {
		expected := expected
		doer := doer

		t.Run(expected, func(t *testing.T) {
			o := vipser.New()
			doer(o)
			assert.Equal(t, expected, o.RenderCommands()[0])
		})
	}
}

func TestOperation_RunCorrectArgs(t *testing.T) {
	var output bytes.Buffer
	o := vipser.New()
	o.Vipser = "echo"
	o.Input = &bytes.Buffer{}
	o.Output = &output

	o.Resize(1,2)
	o.Stretch(3,4)

	err := o.Run()
	assert.NoError(t, err)

	assert.Equal(t, "RESIZE,1,2 STRETCH,3,4\n", output.String())
}

func TestOperation_RunCorrectInput(t *testing.T) {
	var output bytes.Buffer
	o := vipser.New()
	o.Vipser = "cat"
	o.Input = bytes.NewBufferString("Hello!")
	o.Output = &output

	err := o.Run()
	assert.NoError(t, err)

	assert.Equal(t, []byte("Hello!"), output.Bytes())
}

