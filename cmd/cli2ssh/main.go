package main

import (
	"flag"
	"os"
	"os/exec"

	"github.com/PeronGH/cli2ssh/internal/args"
	"github.com/PeronGH/cli2ssh/internal/server"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
)

func main() {
	var env args.ArrayArg
	flag.Var(&env, "e", "Set environment variables for each SSH session.")
	command := flag.String("c", "", "Set the command to run for each SSH session.")
	useOsEnv := flag.Bool("os-env", false, "Use the OS environment variables for the command.")

	flag.Parse()

	if *command == "" {
		log.Fatal("No command provided.")
	}

	srv, err := server.CreateServer(server.CreateServerOptions{
		CommandProvider: func(s ssh.Session) *exec.Cmd {
			argSession := args.NewSession(s)
			fmtCmd := argSession.FormatArg(*command)
			fmtEnv := argSession.FormatArgs(env)

			cmd := exec.CommandContext(s.Context(), "sh", "-c", fmtCmd)
			if *useOsEnv {
				cmd.Env = os.Environ()
			}
			cmd.Env = append(cmd.Env, fmtEnv...)
			return cmd
		},
	})

	if err != nil {
		log.Fatalf("could not create server: %v", err)
	}

	log.Info("Starting server...", "address", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
