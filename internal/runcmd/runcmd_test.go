package runcmd_test

import (
	"context"
	"errors"
	"os/exec"
	"testing"
	"time"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure-cli/internal/runcmd"
	"github.com/JosiahWitt/ensure/ensurepkg"
)

func TestRunnerExec(t *testing.T) {
	ensure := ensure.New(t)

	ensure.Run("with valid command execution", func(ensure ensurepkg.Ensure) {
		runner := runcmd.Runner{}
		result, err := runner.Exec(context.Background(), &runcmd.ExecParams{
			PWD:  "/tmp",
			CMD:  "sh",
			Args: []string{"-c", "pwd"},
		})

		ensure(err).IsNotError()
		ensure(result).Equals("/tmp\n")
	})

	ensure.Run("with invalid command", func(ensure ensurepkg.Ensure) {
		runner := runcmd.Runner{}
		result, err := runner.Exec(context.Background(), &runcmd.ExecParams{
			CMD: "this-command-does-not-exist",
		})

		var expectedErr *exec.Error
		ensure(errors.As(err, &expectedErr)).IsTrue()
		ensure(result).IsEmpty()
	})

	ensure.Run("with failing command", func(ensure ensurepkg.Ensure) {
		runner := runcmd.Runner{}
		result, err := runner.Exec(context.Background(), &runcmd.ExecParams{
			CMD:  "sh",
			Args: []string{"-c", "echo 'abc'; exit 1"},
		})

		ensure(err.Error()).Equals("abc\n")
		ensure(result).IsEmpty()
	})

	ensure.Run("with failing command and no output", func(ensure ensurepkg.Ensure) {
		runner := runcmd.Runner{}
		result, err := runner.Exec(context.Background(), &runcmd.ExecParams{
			CMD:  "sh",
			Args: []string{"-c", "exit 1"},
		})

		ensure(err.Error()).Equals("exit status 1")
		ensure(result).IsEmpty()
	})

	ensure.Run("with cancelled context", func(ensure ensurepkg.Ensure) {
		runner := runcmd.Runner{}
		ctx, cancel := context.WithCancel(context.Background())

		go func() {
			// Allow time for command to start before cancelling
			time.Sleep(10 * time.Millisecond)
			cancel()
		}()

		result, err := runner.Exec(ctx, &runcmd.ExecParams{
			CMD:  "sh",
			Args: []string{"-c", "sleep 0.1"},
		})

		ensure(err).IsError(runcmd.ErrProcessTerminated)
		ensure(result).IsEmpty()
	})
}
