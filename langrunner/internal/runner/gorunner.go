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

//go:embed gorunner_init.sh
var GoInitScript []byte

//go:embed gorunner_run.sh
var GoRunScript []byte

type GoRunner struct {
	Conf    config.GoRunnerConfig
	Manager *RuntimeManager
	Logger  *slog.Logger
}

func (g *GoRunner) RunCode(ctx context.Context, code []byte) (*RunResult, error) {
	g.Logger.Debug(
		"Using scripts",
		slog.String("init", string(GoInitScript)),
		slog.String("run", string(GoRunScript)),
	)

	u, err := g.Manager.AcquireUserWait(ctx)
	if err != nil {
		return nil, fmt.Errorf("aqcuiring user: %w", err)
	}
	defer g.Manager.ReleaseUser(u)

	if err = g.initCurrenntDir(ctx, u, code); err != nil {
		return nil, err
	}

	res, err := g.runCurrentDir(ctx, u)
	if err != nil {
		return nil, fmt.Errorf("running: %w", err)
	}

	return res, nil
}

func (g *GoRunner) initCurrenntDir(ctx context.Context, u *user.User, code []byte) error {
	env, err := g.Manager.CreateSafeEnv(ctx, u, "bash")
	if err != nil {
		return fmt.Errorf("getting a safe env: %w", err)
	}

	_, err = env.StdInPipe.Write(GoInitScript)
	if err != nil {
		return fmt.Errorf("writing go init script to bash: %w", err)
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

	mainPath := path.Join(env.User.HomeDir, "main.go")
	main, err := writeFile(mainPath, code)
	if err != nil {
		return fmt.Errorf("writing main.go: %w", err)
	}
	defer main.Close()

	g.Logger.Info("Wrote main file", slog.String("path", mainPath))

	return nil
}

func (g *GoRunner) runCurrentDir(ctx context.Context, u *user.User) (*RunResult, error) {
	env, err := g.Manager.CreateSafeEnv(ctx, u, "bash")
	if err != nil {
		return nil, fmt.Errorf("getting a safe env: %w", err)
	}

	_, err = env.StdInPipe.Write(GoRunScript)
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
		Stdout:     []byte(output.Stdout),
		Stderr:     []byte(output.Stderr),
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
