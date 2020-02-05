package vipser_test

import (
	"bytes"
	"github.com/jhford/govipser"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"os"
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
		o.Resize(1, 2)
		o.Stretch(3, 4)
		assert.Equal(t, []string{"RESIZE,1,2", "STRETCH,3,4"}, o.RenderCommands())
	})

	t.Run("no commands", func(t *testing.T) {
		o := vipser.New()
		assert.Equal(t, []string{}, o.RenderCommands())
	})
}

func TestOperation_AllCommands(t *testing.T) {
	tests := map[string]func(*vipser.Operation){
		"RESIZE,1,2":  func(o *vipser.Operation) { o.Resize(1, 2) },
		"STRETCH,1,2": func(o *vipser.Operation) { o.Stretch(1, 2) },
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

	o.Resize(1, 2)
	o.Stretch(3, 4)

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

var testfile string = "test.png"
var testinput []byte

func init() {
	if v, ok := os.LookupEnv("TEST_IMAGE"); ok {
		log.Printf("Setting test image path to %s", v)
		testfile = v
	}

	i, err := ioutil.ReadFile(testfile)
	if err != nil {
		panic(err)
	}
	testinput = i
}

func TestOperation_ResizeImage_Readers(t *testing.T) {
	//var output bytes.Buffer
	op := vipser.New()
	op.Resize(2, 2)

	var output bytes.Buffer

	op.Input = bytes.NewBuffer(testinput)

	op.Output = &output

	t.Logf("Vipser: %s", op.Vipser)

	err := op.Run()
	assert.NoError(t, err)
	assert.Greater(t, output.Len(), 0)

}

func TestOperation_ResizeImage_Buffers(t *testing.T) {
	op := vipser.New()
	op.Resize(2, 2)

	output, err := op.Apply(testinput)
	assert.NoError(t, err)
	assert.Greater(t, len(output), 0)
}

func TestOperation_ResizeImage_Files(t *testing.T) {
	op := vipser.New()

	op.Autorot()
	op.Resize(512, 384)
	op.EmbedWhite(512/2, 384/2, 1024, 768)
	op.Format("png")
	op.Extract(90, 90, 1024-180, 768-180)
	op.Blur(3)
	op.Rotate(90)

	in, err := os.Open(testfile)
	assert.NoError(t, err)
	op.Input = in

	outname := ".test_" + t.Name() + ".png"
	out, err := os.Create(outname)
	assert.NoError(t, err)
	op.Output = out

	err = op.Run()
	assert.NoError(t, err)

	fi, err := os.Stat(outname)
	assert.NoError(t, err)
	if err == nil {
		assert.Greater(t, fi.Size(), int64(0))
	}
}

func dobench(b *testing.B, mods func(*vipser.Operation)) {
	var out []byte
	for i := 0; i < b.N; i++ {
		op := vipser.New()

		mods(op)

		_out, err := op.Apply(testinput)
		if err != nil {
			panic(err)
		}

		out = _out
	}

	b.ReportMetric(float64(len(out)), "bytes_out")
	b.ReportMetric(float64(len(testinput)), "bytes_in")
	b.ReportMetric(float64(len(testinput)-len(out)), "bytes_red")
}

func BenchmarkNoCommands(b *testing.B) {
	dobench(b, func(operation *vipser.Operation) {})
}

func BenchmarkResize(b *testing.B) {
	dobench(b, func(operation *vipser.Operation) {
		operation.Resize(200, 200)
	})
}

func BenchmarkBlur2(b *testing.B) {
	dobench(b, func(operation *vipser.Operation) {
		operation.Blur(2)
	})
}

func BenchmarkBlurZero5(b *testing.B) {
	dobench(b, func(operation *vipser.Operation) {
		operation.Blur(0.5)
	})
}
