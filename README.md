# cli2ssh

Turn any CLI program into a SSH server.

## Installation

```bash
go install github.com/PeronGH/cli2ssh/cmd/cli2ssh@latest
```

## Usage

```bash
# Check usage
cli2ssh -h

# Example: echo the username
cli2ssh -c 'echo Hello, {{ .User }}.'
```

## TODO

- [ ] Authentication
- [ ] Add tests
- [ ] Integrate with GitHub Actions
