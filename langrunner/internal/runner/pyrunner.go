package runner

import (
	"context"
	_ "embed"
	"fmt"
	"log/slog"
	"os/user"
	"path"
	"time"

	"github.com/Marattttt/coderunner/langrunner/internal/config"
)

//go:embed pyrunner_init.sh
var PyInitScript []byte

//go:embed pyrunner_run.sh
var PyRunScript []byte

type PyRunner struct {
	Conf    config.PyRunnerConfig
	Manager *RuntimeManager
	Logger  *slog.Logger
}

func NewPyRunner(logger *slog.Logger, conf config.PyRunnerConfig, manager *RuntimeManager) *PyRunner {
	return &PyRunner{
		Conf:    conf,
		Manager: manager,
		Logger:  logger,
	}
}

func (pr *PyRunner) RunCode(ctx context.Context, code []byte) (*RunResult, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Duration(pr.Manager.conf.ExecTimeoutMs*1e6))
	defer cancel()

	pr.Logger.Debug(
		"Using scripts",
		slog.String("init", string(PyInitScript)),
		slog.String("run", string(PyRunScript)),
	)

	u, err := pr.Manager.AcquireUserWait(ctx)
	if err != nil {
		return nil, fmt.Errorf("aqcuiring user: %w", err)
	}
	defer pr.Manager.ReleaseUser(u)

	if err = pr.initCurrenntDir(ctx, u, code); err != nil {
		return nil, err
	}

	res, err := pr.runCurrentDir(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("running: %w", err)
	}

	return res, nil
}

func (g *PyRunner) initCurrenntDir(ctx context.Context, u *user.User, code []byte) error {
	env, err := g.Manager.CreateSafeEnv(ctx, u, "bash")
	if err != nil {
		return fmt.Errorf("getting a safe env: %w", err)
	}

	_, err = env.StdInPipe.Write(PyInitScript)
	if err != nil {
		return fmt.Errorf("writing python init script to bash: %w", err)
	}
	_, err = env.StdInPipe.Write([]byte("\nexit\n"))
	if err != nil {
		return fmt.Errorf("writing 'exit' to bash: %w", err)
	}

	output, err := readStdoutStderr(ctx, env)
	if err != nil {
		g.Logger.Error("Error running init script",
			slog.Any("output", output),
		)
		return fmt.Errorf("running init: %w", err)
	}

	g.Logger.Info("Ran init script")
	g.Logger.Debug("Ran init script", slog.Any("output", output))

	mainPath := path.Join(env.User.HomeDir, "main.py")
	main, err := writeFile(mainPath, code)
	if err != nil {
		return fmt.Errorf("writing main.py: %w", err)
	}
	defer main.Close()

	g.Logger.Info("Wrote main file", slog.String("path", mainPath))

	return nil
}

func (g *PyRunner) runCurrentDir(ctx context.Context, u *user.User) (*RunResult, error) {
	env, err := g.Manager.CreateSafeEnv(ctx, u, "bash")
	if err != nil {
		return nil, fmt.Errorf("getting a safe env: %w", err)
	}

	_, err = env.StdInPipe.Write(PyRunScript)
	if err != nil {
		return nil, fmt.Errorf("writing go run script to bash: %w", err)
	}

	startAt := time.Now()

	output, err := readStdoutStderr(ctx, env)
	if err != nil {
		return nil, fmt.Errorf("Reading stdout, stdin: %w", err)
	}

	end := time.Now()

	res := RunResult{
		ExitCode:   output.ExitCode,
		Stdout:     output.Stdout,
		Stderr:     output.Stderr,
		TimeTook:   end.Sub(startAt),
		TimeTookMs: end.Sub(startAt).Milliseconds(),
	}

	g.Logger.Info("Ran user code",
		slog.Int("exitCode", res.ExitCode),
		slog.Int64("timeTook", res.TimeTookMs),
	)

	g.Logger.Debug("Ran user code",
		slog.String("stdout", output.Stdout),
		slog.String("stderr", output.Stderr),
		slog.Int("exitCode", res.ExitCode),
		slog.Int64("timeTook", res.TimeTookMs),
	)

	return &res, nil
}
