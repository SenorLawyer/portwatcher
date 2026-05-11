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
- GoReleaser release flow

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

- Go 1.25+
- Docker, optional for container mapping
- GoReleaser, optional for local release checks

```bash
go test ./...
go vet ./...
go run ./cmd/portwatch
```

## Release

Each PR adds a `.bumpy/<short-description>.md` file with:

```yaml
---
bump: patch
---
Brief changelog description.
```

When PRs merge to `main`, the Version Packages workflow keeps a single `version-packages` PR updated from pending bump files, selects the highest bump level (`major > minor > patch`), and updates `CHANGELOG.md`. Merging that PR creates a `v*` tag and GoReleaser publishes the release.

For a local package check:

```bash
goreleaser release --snapshot --clean
```
