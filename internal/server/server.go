package server

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"strings"
	"syscall"

	"github.com/PeronGH/cli2ssh/internal/env"
	"github.com/PeronGH/cli2ssh/internal/path"
	"github.com/PeronGH/cli2ssh/internal/set"
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

	// Always use `exec.CommandContext` to create the command.
	CommandProvider func(s ssh.Session) *exec.Cmd

	// Optional
	Host         string
	Port         string
	HostKeyPath  string
	AuthProvider func(s ssh.Session) bool
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
		if env.HostKeyPath != "" {
			opts.HostKeyPath = env.HostKeyPath
		} else {
			hostKeyPath, err := path.GetDefaultHostKeyPath()
			if err != nil {
				return nil, fmt.Errorf("could not get default host key path: %w", err)
			}
			opts.HostKeyPath = hostKeyPath
		}
	}
	if opts.AuthProvider == nil {
		if env.AllowedUsers == "" {
			opts.AuthProvider = func(s ssh.Session) bool {
				return true
			}
		} else {
			allowedUsers := set.NewFromSlice(strings.Split(env.AllowedUsers, ","))
			opts.AuthProvider = func(s ssh.Session) bool {
				return allowedUsers.Has(s.User())
			}
		}
	}

	return wish.NewServer(
		wish.WithAddress(net.JoinHostPort(opts.Host, opts.Port)),
		wish.WithHostKeyPath(opts.HostKeyPath),
		ssh.AllocatePty(),
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					if !opts.AuthProvider(s) {
						wish.Fatalln(s, "you are not allowed to access the service.")
						next(s)
						return
					}

					pty, _, hasPty := s.Pty()
					if !hasPty {
						wish.Fatalln(s, "client has no PTY.")
						next(s)
						return
					}

					cmd := opts.CommandProvider(s)
					if cmd == nil {
						wish.Fatalln(s, "your session has no command to execute.")
						next(s)
						return
					}

					cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", pty.Term))
					cmd.Stdin = pty.Slave
					cmd.Stdout = pty.Slave
					cmd.Stderr = pty.Slave
					cmd.SysProcAttr = &syscall.SysProcAttr{
						Setctty: true,
						Setsid:  true,
					}

					if err := cmd.Run(); err != nil {
						wish.Fatalln(s, "Failed to run the command:", err)
					}

					next(s)
				}
			},
			activeterm.Middleware(),
			logging.Middleware(),
		),
	)
}
