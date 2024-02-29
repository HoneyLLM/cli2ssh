package server

import (
	"errors"
	"fmt"
	"net"
	"os"
	"os/exec"
	"syscall"

	"github.com/PeronGH/cli2ssh/internal/env"
	"github.com/PeronGH/cli2ssh/internal/path"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
	"github.com/charmbracelet/wish/activeterm"
	"github.com/charmbracelet/wish/logging"
)

var (
	ErrCommandProviderRequired = errors.New("command provider is required")
)

type CreateServerOptions struct {
	// Required
	CommandProvider func(s ssh.Session) []string

	// Optional
	Host        string
	Port        string
	HostKeyPath string
	EnvProvider func(s ssh.Session) []string
}

func CreateServer(opts CreateServerOptions) (*ssh.Server, error) {
	// Required options
	if opts.CommandProvider == nil {
		return nil, ErrCommandProviderRequired
	}

	// Optional options
	if opts.Host == "" {
		opts.Host = env.Host
	}
	if opts.Port == "" {
		opts.Port = env.Port
	}
	if opts.HostKeyPath == "" {
		if env.HostKeyPath == "" {
			hostKeyPath, err := path.GetDefaultHostKeyPath()
			if err != nil {
				return nil, fmt.Errorf("could not get default host key path: %w", err)
			}
			opts.HostKeyPath = hostKeyPath
		} else {
			opts.HostKeyPath = env.HostKeyPath
		}
	}
	if opts.EnvProvider == nil {
		opts.EnvProvider = func(s ssh.Session) []string {
			return os.Environ()
		}
	}

	return wish.NewServer(
		wish.WithAddress(net.JoinHostPort(opts.Host, opts.Port)),
		wish.WithHostKeyPath(opts.HostKeyPath),
		ssh.AllocatePty(),
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					pty, _, hasPty := s.Pty()
					if !hasPty {
						wish.Fatalln(s, "client has no PTY.")
						next(s)
						return
					}

					command := opts.CommandProvider(s)
					if len(command) == 0 {
						wish.Fatalln(s, "your session has no command to execute.")
						next(s)
						return
					}

					cmd := exec.CommandContext(s.Context(), command[0], command[1:]...)
					cmd.Env = append(opts.EnvProvider(s), fmt.Sprintf("TERM=%s", pty.Term))
					cmd.Stdin = pty.Slave
					cmd.Stdout = pty.Slave
					cmd.Stderr = pty.Slave
					cmd.SysProcAttr = &syscall.SysProcAttr{
						Setctty: true,
						Setsid:  true,
					}

					if err := cmd.Run(); err != nil {
						wish.Fatalln(s, "Could not start command", "error", err)
					}

					next(s)
				}
			},
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
}
