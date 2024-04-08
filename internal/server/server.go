package server

import (
	"errors"
	"fmt"
	"io"
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

	PublicKeyAuth func(ctx ssh.Context, key ssh.PublicKey) bool
	PasswordAuth  func(ctx ssh.Context, password string) bool
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
		ssh.PublicKeyAuth(opts.PublicKeyAuth),
		ssh.PasswordAuth(opts.PasswordAuth),
		ssh.AllocatePty(),
		wish.WithMiddleware(
			func(next ssh.Handler) ssh.Handler {
				return func(s ssh.Session) {
					defer next(s)

					cmd := opts.CommandProvider(s)
					if cmd == nil {
						wish.Fatalln(s, "your session has no command to execute.")
						return
					}

					pty, _, hasPty := s.Pty()

					log.Info("Executing command", "remote", s.RemoteAddr(), "command", cmd, "pty", hasPty)
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

						if err := pipeStdio(cmd, s, s, s.Stderr()); err != nil {
							log.Error("Failed to pipe stdio", "remote", s.RemoteAddr(), "error", err)
							wish.Fatalln(s, "Failed to pipe stdio:", err)
							return
						}
					}

					defer s.Exit(cmd.ProcessState.ExitCode())
					if err := cmd.Run(); err != nil {
						if exitErr, ok := err.(*exec.ExitError); ok {
							log.Warn("Command exited with status", "remote", s.RemoteAddr(), "command", cmd, "status", exitErr.ExitCode())
						} else {
							log.Error("Failed to run the command", "remote", s.RemoteAddr(), "command", cmd, "error", err)
							wish.Fatalln(s, "Failed to run the command:", err)
						}
					}
				}
			},
			logging.Middleware(),
		),
	)
}

// workaround for command hanging
func pipeStdio(cmd *exec.Cmd, stdin io.Reader, stdout, stderr io.Writer) error {
	cmdStdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	cmdStdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmdStderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	go io.Copy(cmdStdin, stdin)
	go io.Copy(stdout, cmdStdout)
	go io.Copy(stderr, cmdStderr)

	return nil
}
