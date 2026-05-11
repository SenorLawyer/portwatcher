# Project Notes

- This is a Go CLI. Keep runtime code in Go.
- Run `go test ./...`, `go vet ./...`, and `go build ./cmd/portwatch` before release work.
- Releases use GoReleaser and `v*` tags.
- Use `patch` for fixes, docs, cleanup, and dependency updates.
- Use `minor` for new user-facing behavior.
- Use `major` for breaking CLI, config, or output changes.

## Release Flow

PRs that contain user-visible changes must include a bump file at `.bumpy/<short-description>.md`:

```yaml
---
bump: minor
---
Brief changelog description of the change.
```

When a PR merges to `main`, CI automatically creates (or updates) a `version-packages` PR that applies the version bump and updates `CHANGELOG.md`. Merging that PR triggers the actual release via GoReleaser.
