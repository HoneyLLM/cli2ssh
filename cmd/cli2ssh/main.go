package main

import (
	"os"
	"os/exec"

	"github.com/PeronGH/cli2ssh/internal/server"
	"github.com/charmbracelet/log"
	"github.com/charmbracelet/ssh"
)

func main() {
	command := os.Args[1:]

	srv, err := server.CreateServer(server.CreateServerOptions{
		CommandProvider: func(s ssh.Session) *exec.Cmd {
			return exec.CommandContext(s.Context(), command[0], command[1:]...)
		},
	})

	if err != nil {
		log.Fatalf("could not create server: %v", err)
	}

	log.Info("Starting server...", "address", srv.Addr)
	log.Fatal(srv.ListenAndServe())
}
