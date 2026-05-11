# Project Notes

- This is a Go CLI. Keep runtime code in Go.
- Run `go test ./...`, `go vet ./...`, and `go build ./cmd/portwatch` before release work.
- Releases use GoReleaser and `v*` tags.
- Use `patch` for fixes, docs, cleanup, and dependency updates.
- Use `minor` for new user-facing behavior.
- Use `major` for breaking CLI, config, or output changes.
