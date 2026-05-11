# Project Notes

- This is a Go CLI. Keep runtime code in Go.
- Run `go test ./...`, `go vet ./...`, and `go build ./cmd/portwatch` before release work.
- Releases use GoReleaser and `v*` tags.
- Use `patch` for fixes, docs, cleanup, and dependency updates.
- Use `minor` for new user-facing behavior.
- Use `major` for breaking CLI, config, or output changes.

## Release Flow

All PRs must include a bump file at `.bumpy/<short-description>.md` (except the auto-generated `version-packages` PR):

```yaml
---
bump: minor
---
Brief changelog description of the change.
```

When PRs merge to `main`, CI automatically keeps a single `version-packages` PR updated from all pending bump files, applies the highest pending bump level (`major` > `minor` > `patch`), and updates `CHANGELOG.md` with all pending bump entries. Merging that PR triggers the actual release via GoReleaser.
