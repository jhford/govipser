package vipser

import (
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
	Input io.Reader
	Output io.Writer
	Vipser string
}

func New() *Operation {
	return &Operation{
		Commands: make([]Command, 0),
		Vipser: Vipser,
	}
}

func (o *Operation) Resize(x,y int) *Operation {
	o.Commands = append(o.Commands, Command{"RESIZE", x, y})
	return o
}

func (o *Operation) Stretch(x,y int) *Operation {
	o.Commands = append(o.Commands, Command{"STRETCH", x, y})
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
	cmd.Stdin = o.Input
	cmd.Stdout = o.Output

	err := cmd.Start()
	if err != nil {
		return err
	}

	err = cmd.Wait()

	if v, ok := err.(*exec.ExitError); ok {
		return fmt.Errorf("%s exited %d\n=======\n%s", cmd, v.ExitCode(), v.Stderr)
	}

	return nil
}

func (o Operation) Run() error {
	return o.RunWithContext(context.Background())
}

