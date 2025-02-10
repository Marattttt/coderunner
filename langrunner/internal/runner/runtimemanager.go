package runner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os/exec"
	"os/user"
	"slices"
	"strconv"
	"syscall"
	"time"

	"github.com/Marattttt/coderunner/langrunner/internal/config"
)

var ErrNoFreeUsers = fmt.Errorf("no free users")

type SecureEnv struct {
	User       *user.User
	StdInPipe  io.WriteCloser
	StdOutPipe io.ReadCloser
	StdErrPipe io.ReadCloser
	Cmd        *exec.Cmd
}

type RuntimeManager struct {
	conf      *config.RunnerConig
	logger    *slog.Logger
	busyUsers []*user.User
}

func NewRuntimeManager(conf *config.RunnerConig, logger *slog.Logger) *RuntimeManager {
	return &RuntimeManager{
		conf:      conf,
		logger:    logger,
		busyUsers: make([]*user.User, 0),
	}
}

// A blocking call that waits for an available user environment
func (r *RuntimeManager) AcquireUserWait(ctx context.Context) (*user.User, error) {
	envCh := make(chan *user.User)
	errCh := make(chan error)
	go func() {
		for {
			u, err := r.AcquireUser()
			if errors.Is(err, ErrNoFreeUsers) {
				time.Sleep(time.Millisecond * 500)
				continue
			}
			if err != nil {
				errCh <- err
				close(envCh)
				return
			}
			envCh <- u
			close(errCh)
			return
		}
	}()

	select {
	case env := <-envCh:
		return env, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()

	}
}

// Starts a new process info with specified commands as the same or a different user,
// depending on configuration
//
// Works only when executed in a Unix-like OS
func (r *RuntimeManager) CreateSafeEnv(ctx context.Context, u *user.User, name string, args ...string) (*SecureEnv, error) {
	uid, err := strconv.ParseInt(u.Uid, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("parsing uid %s: %w", u.Uid, err)
	}
	gid, err := strconv.ParseInt(u.Gid, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("parsing gid %s: %w", u.Gid, err)
	}

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = nil
	cmd.Dir = u.HomeDir
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	// Works only on *nix-liike systems
	cmd.SysProcAttr.Credential = &syscall.Credential{
		Uid: uint32(uid),
		Gid: uint32(gid),
	}

	stdinPipe, err := cmd.StdinPipe()
	if err != nil {
		return nil, fmt.Errorf("getting stdin pipe: %w", err)
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("gettoutg stdout pipe: %w", err)
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return nil, fmt.Errorf("getterrg stderr pipe: %w", err)
	}

	err = cmd.Start()
	if err != nil {
		return nil, fmt.Errorf("starting command: %w", err)
	}

	r.logger.Info(
		"Initialiized a secure environment",
		slog.String("user", u.Username),
		slog.String("dir", u.HomeDir),
	)

	return &SecureEnv{User: u, Cmd: cmd, StdInPipe: stdinPipe, StdOutPipe: stdoutPipe, StdErrPipe: stderrPipe}, nil
}

func (r *RuntimeManager) AcquireUser() (*user.User, error) {
	if r.conf.SafeUserNames == nil {
		current, err := user.Current()
		if err != nil {
			return nil, fmt.Errorf("Looking up current user: %w", err)
		}
		return current, nil
	}

	var username *string
	for _, uname := range r.conf.SafeUserNames {
		isFree := false
		for _, busy := range r.busyUsers {
			if busy.Name == uname {
				isFree = false
				break
			}
		}

		if !isFree {
			username = &uname
			break
		}
	}

	if username == nil {
		return nil, ErrNoFreeUsers
	}

	u, err := user.Lookup(*username)
	if err != nil {
		return nil, fmt.Errorf("user not found: %s", *username)
	}

	r.busyUsers = append(r.busyUsers, u)

	return u, nil
}

func (r *RuntimeManager) ReleaseUser(u *user.User) {
	i := slices.Index(r.busyUsers, u)
	if i == -1 {
		return
	}

	r.busyUsers[i] = r.busyUsers[len(r.busyUsers)-1]
	r.busyUsers = r.busyUsers[:len(r.busyUsers)-2]
}
