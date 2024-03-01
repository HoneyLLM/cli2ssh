# cli2ssh

Turn any CLI program into a SSH server.

## Installation

```bash
go install github.com/PeronGH/cli2ssh/cmd/cli2ssh@latest
```

## Usage

```bash
# Use default configuration
cli2ssh echo 'Hello, World!'

# Pass custom configuration via environment variables
CLI2SSH_HOST=0.0.0.0 CLI2SSH_PORT=22 CLI2SSH_HOST_KEY_PATH=~/.ssh/id_rsa cli2ssh bash -l
```
