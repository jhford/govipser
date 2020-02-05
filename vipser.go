
package vipser

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
)

func stringify(it interface{}) string {
	switch value := it.(type) {
	case string:
		return value
	case int:
		return strconv.FormatInt(int64(value), 10)
	case int32:
		return strconv.FormatInt(int64(value), 10)
	case int64:
		return strconv.FormatInt(value, 10)
	case uint:
		return strconv.FormatUint(uint64(value), 10)
	case uint32:
		return strconv.FormatUint(uint64(value), 10)
	case uint64:
		return strconv.FormatUint(value, 10)
	case float32:
		return strconv.FormatFloat(float64(value), 'G', -1, 32)
	case float64:
		return strconv.FormatFloat(value, 'G', -1, 64)
	default:
		panic(errors.New("unknown type to stringify"))
	}
}

// Command is any string, integer type or floating point type.
// Pinky swear it is though, or you'll panic!
type Command []interface{}

// String returns a string representation of a string
func (c Command) String() string {
	s := make([]string, len(c))
	for i, command := range c {
		s[i] = stringify(command)
	}
	return strings.Join(s, ",")
}

type Operation struct {
	Commands []Command
	Input    io.Reader
	Output   io.Writer
	Vipser   string
	finished bool
}

func New() *Operation {
	vipser, err := FindVipser()
	if err != nil {
		panic(err)
	}
	return &Operation{
		Commands: make([]Command, 0),
		Vipser:   vipser,
	}
}

func (o *Operation) Resize(x, y int) *Operation {
	o.Commands = append(o.Commands, Command{"RESIZE", x, y})
	return o
}

func (o *Operation) Stretch(x, y int) *Operation {
	o.Commands = append(o.Commands, Command{"STRETCH", x, y})
	return o
}

func (o *Operation) Expand(x, y int) *Operation {
	o.Commands = append(o.Commands, Command{"EXPAND", x, y})
	return o
}

func (o *Operation) Extract(left, top, width, height int) *Operation {
	o.Commands = append(o.Commands, Command{"EXTRACT", left, top, width, height})
	return o
}

type VipserEmbed int

const (
	VipserEmbedWhite = VipserEmbed(1)
	VipserEmbedBlack = VipserEmbed(2)
)

func (o *Operation) Embed(x, y, width, height int, embed VipserEmbed) *Operation {
	var command string
	if embed == VipserEmbedBlack {
		command = "EMBBLK"
	} else if embed == VipserEmbedWhite {
		command = "EMBWHT"
	} else {
		panic("you found a programming error!")
	}

	o.Commands = append(o.Commands, Command{command, x, y, width, height})
	return o
}

func (o *Operation) EmbedWhite(x, y, width, height int) *Operation {
	return o.Embed(x, y, width, height, VipserEmbedWhite)
}

func (o *Operation) EmbedBlack(x, y, width, height int) *Operation {
	return o.Embed(x, y, width, height, VipserEmbedBlack)
}

func (o *Operation) Blur(sigma float64) *Operation {
	o.Commands = append(o.Commands, Command{"BLUR", sigma})
	return o
}

func (o *Operation) Rotate(angle int) *Operation {
	o.Commands = append(o.Commands, Command{"ROTATE", angle})
	return o
}

func (o *Operation) Autorot() *Operation {
	o.Commands = append(o.Commands, Command{"AUTOROT"})
	return o
}

func (o *Operation) Quality(quality int) *Operation {
	o.Commands = append(o.Commands, Command{"QUALITY", quality})
	return o
}

func (o *Operation) Format(format string) *Operation {
	o.Commands = append(o.Commands, Command{"EXPORT", format})
	return o
}

func (o Operation) RenderCommands() []string {
	commands := make([]string, len(o.Commands))
	for i, command := range o.Commands {
		commands[i] = command.String()
	}
	return commands
}

func (o Operation) RunWithContext(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, o.Vipser, o.RenderCommands()...)
	if o.Input == nil {
		return errors.New("must provide input")
	}

	if o.Output == nil {
		return errors.New("must provide output")
	}

	var stderr bytes.Buffer
	cmd.Stdin = o.Input
	cmd.Stdout = o.Output
	cmd.Stderr = &stderr

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()

	if v, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("%s exited %d\n=======\n%s", cmd, v.ExitCode(), stderr.Bytes())
	}

	return nil
}

func (o Operation) Run() error {
	return o.RunWithContext(context.Background())
}

func (o *Operation) Apply(input []byte) ([]byte, error) {
	var output bytes.Buffer
	o.Output = &output
	o.Input = bytes.NewBuffer(input)

	err := o.Run()
	if err != nil {
		return []byte{}, err
	}

	return output.Bytes(), nil
}
