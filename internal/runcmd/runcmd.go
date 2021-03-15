package runcmd

import (
	"context"
	"errors"
	"os/exec"
)

var ErrProcessTerminated = errors.New("process was terminated by a signal")

type ExecParams struct {
	PWD  string
	CMD  string
	Args []string
}

type RunnerIface interface {
	Exec(ctx context.Context, params *ExecParams) (string, error)
}

type Runner struct{}

var _ RunnerIface = &Runner{}

// Exec the command defined in the provided params.
func (*Runner) Exec(ctx context.Context, params *ExecParams) (string, error) {
	//nolint:gosec
	c := exec.CommandContext(ctx, params.CMD, params.Args...)
	c.Dir = params.PWD
	out, err := c.CombinedOutput()
	if err != nil {
		return "", handleExecError(out, err)
	}

	return string(out), err
}

func handleExecError(out []byte, err error) error {
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		// Process state exit code of -1 implies that the process was terminated by a signal,
		// thus return a special error so we can detect this case.
		if exitErr.ProcessState.ExitCode() == -1 {
			return ErrProcessTerminated
		}

		if len(out) == 0 {
			return err // Better to have some error message than no error message
		}

		//nolint:goerr113
		return errors.New(string(out))
	}

	return err
}
