package server

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"syscall"

	"github.com/PeronGH/cli2ssh/internal/path"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
	"github.com/charmbracelet/wish"
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
	Host        string
	Port        string
	HostKeyPath string
}

func CreateServer(opts CreateServerOptions) (*ssh.Server, error) {
	// Required options
	if opts.CommandProvider == nil {
		return nil, ErrCommandProviderRequired
	}

	// Optional options
	if opts.Host == "" {
		opts.Host = "localhost"
	}
	if opts.Port == "" {
		opts.Port = "2222"
	}
	if opts.HostKeyPath == "" {
		opts.HostKeyPath = path.GetDefaultHostKeyPath()
	}

	return wish.NewServer(
		wish.WithAddress(net.JoinHostPort(opts.Host, opts.Port)),
		wish.WithHostKeyPath(opts.HostKeyPath),
		ssh.AllocatePty(),
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					pty, _, hasPty := s.Pty()

					cmd := opts.CommandProvider(s)
					if cmd == nil {
						wish.Fatalln(s, "your session has no command to execute.")
						next(s)
						return
					}

					if hasPty {
						cmd.Env = append(cmd.Env, fmt.Sprintf("TERM=%s", pty.Term))
						cmd.Stdin = pty.Slave
						cmd.Stdout = pty.Slave
						cmd.Stderr = pty.Slave
						cmd.SysProcAttr = &syscall.SysProcAttr{
							Setctty: true,
							Setsid:  true,
						}
					} else {
						cmd.Env = append(cmd.Env, "TERM=dumb")
						cmd.Stdin = s
						cmd.Stdout = s
						cmd.Stderr = s.Stderr()
						// TODO: fix command hanging when no pty
					}

					log.Info("Executing command", "command", cmd, "pty", hasPty)
					if err := cmd.Run(); err != nil {
						if exitErr, ok := err.(*exec.ExitError); ok {
							log.Warn("Command exited with status", "command", cmd, "status", exitErr.ExitCode())
						} else {
							log.Error("Failed to run the command", "command", cmd, "error", err)
							wish.Fatalln(s, "Failed to run the command:", err)
						}
					}

					next(s)
				}
			},
			logging.Middleware(),
		),
	)
}
