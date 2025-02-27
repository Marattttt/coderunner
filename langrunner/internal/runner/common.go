package runner

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"time"
)

type RunResult struct {
	ExitCode   int           `json:"exitCode"`
	Stdout     string        `json:"stdout"`
	Stderr     string        `json:"stderr"`
	TimeTook   time.Duration `json:"timeTook"`
	TimeTookMs int64         `json:"timeTookMs"`
	Extra      *any          `json:"extra,omitempty"`
}

func writeFile(path string, code []byte) (*os.File, error) {
	f, err := os.OpenFile(path, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return nil, fmt.Errorf("opening: %w", err)
	}

	if _, err = f.Write(code); err != nil {
		return nil, fmt.Errorf("writing: %w", err)
	}

	return f, nil
}

type commandOutput struct {
	Stdout   string
	Stderr   string
	ExitCode int
}

func readStdoutStderr(ctx context.Context, env *SecureEnv) (*commandOutput, error) {
	var outBuf, errBuf bytes.Buffer
	var wg sync.WaitGroup
	var stdoutErr, stderrErr error

	wg.Add(2)
	go func() {
		defer wg.Done()
		_, stdoutErr = io.Copy(&outBuf, env.StdOutPipe)
	}()
	go func() {
		defer wg.Done()
		_, stderrErr = io.Copy(&errBuf, env.StdErrPipe)
	}()

	cmdErrCh := make(chan error, 1)
	go func() {
		wg.Wait()
		err := env.Cmd.Wait()
		if err != nil {
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				cmdErrCh <- nil
				return
			}
		}
		cmdErrCh <- err
	}()

	// Wait for either context cancellation or command completion
	var cmdErr error
	select {
	case cmdErr = <-cmdErrCh:
		break
	case <-ctx.Done():
		// Try to kill the process if context is cancelled
		_ = env.Cmd.Process.Kill()
		return nil, fmt.Errorf("ctx cancelled: %w", ctx.Err())
	}

	if stdoutErr != nil {
		return nil, fmt.Errorf("reading stdout: %w", stdoutErr)
	}
	if stderrErr != nil {
		return nil, fmt.Errorf("reading stderr: %w", stderrErr)
	}

	if cmdErr != nil {
		cmdErr = fmt.Errorf("waiting for cmd execution: %w", cmdErr)
	}

	return &commandOutput{
		Stdout:   outBuf.String(),
		Stderr:   errBuf.String(),
		ExitCode: env.Cmd.ProcessState.ExitCode(),
	}, cmdErr
}
