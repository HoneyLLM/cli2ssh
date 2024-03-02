# cli2ssh

Turn any CLI program into a SSH server.

## Installation

```bash
go install github.com/PeronGH/cli2ssh/cmd/cli2ssh@latest
```

## Usage

```bash
# Check usage
cli2ssh --help

# Basic example: echo the username
cli2ssh -c 'echo Hello, {{ .User }}.'

# More practical example: serve oterm publicly
cli2ssh -h 0.0.0.0 -e 'OTERM_DATA_DIR=userdata/{{ .User }}' -c $(which oterm)
```

## Use Cases

- Share a CLI program with someone who doesn't have it installed.
- Publicly host a TUI program, allowing it to be accessed like a web page.
- Let me know if you have any other ideas!

```bash

## TODO

- [ ] Authentication
- [ ] Add tests
- [ ] Integrate with GitHub Actions
