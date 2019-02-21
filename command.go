package jr

import (
	"bytes"
	"context"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
)

// OSCommand defines the OS command
type OSCommand struct {
	WorkDir string
	Command string
	Args    []string
	Stdout  *bytes.Buffer
	Stderr  *bytes.Buffer
}

// NewOSCommand returns a OSCommand pointer
func NewOSCommand(workDir, command string, args []string, stdout, stderr *bytes.Buffer) *OSCommand {
	return &OSCommand{
		WorkDir: workDir,
		Command: command,
		Args:    args,
		Stdout:  stdout,
		Stderr:  stderr,
	}
}

// Run runs the OS command
func (o *OSCommand) Run(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, o.Command, o.Args...)
	cmd.Stdout = o.Stdout
	cmd.Stderr = o.Stderr
	cmd.Dir = o.WorkDir

	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "%s %s failed", o.Command, strings.Join(o.Args, " "))
	}
	return nil
}
