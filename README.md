# portwatch

`portwatch` is a fast terminal UI for local ports and the processes that own them. It shows open ports, process metadata, commands, Docker/container mappings when available, quick actions, and a local history of port changes.

The module is `github.com/SenorLawyer/portwatcher`; the binary is `portwatch`.

## Features

- htop-inspired Bubble Tea interface for local TCP/UDP ports
- process name, PID, user, and command enrichment
- optional Docker host-port to container mapping
- keyboard actions for kill, copy, filter, open, refresh, and history
- JSON output for scripting
- JSONL history store with retention
- GoReleaser config and bumpy release flow

## Install

```bash
go install github.com/SenorLawyer/portwatcher/cmd/portwatch@latest
```

## Usage

```bash
portwatch
portwatch list --json
portwatch history --json
portwatch version
```

Key bindings:

| Key | Action |
| --- | --- |
| `/` | Filter rows |
| `h` | Toggle history |
| `r` | Refresh now |
| `k` | Interrupt selected process |
| `K` | Force kill selected process |
| `c` | Copy selected row |
| `p` | Copy selected port |
| `o` | Open `http://localhost:<port>` |
| `q` / `Ctrl+C` | Quit |

## Development

Prerequisites:

- Go 1.22+
- Docker, optional for container mapping
- bumpy, optional for releases
- GoReleaser, optional for local release checks

```bash
go test ./...
go run ./cmd/portwatch
```

## Release

Releases are managed with bumpy. Release-significant PRs should include a `.bumpy/*.md` bump file.

```bash
bunx bumpy add --packages "portwatcher:patch" --message "Fixed local port display." --name "fix-port-display"
bunx bumpy status
```

Version PRs are opened from `main`. When a `v*` release tag is pushed, GoReleaser builds the binaries.

For a local package check:

```bash
goreleaser release --snapshot --clean
```
