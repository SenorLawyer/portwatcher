# Contributing

Run the checks before opening a pull request:

```bash
go test ./...
go vet ./...
go build ./cmd/portwatch
```

Keep changes focused. For scanner behavior, prefer provider interfaces and tests with fakes over OS-specific assumptions.

## Bump files

All PRs must include a bump file at `.bumpy/<short-description>.md` (except the auto-generated `version-packages` PR):

```yaml
---
bump: minor
---
Brief changelog description of the change.
```

Bump levels:
- `patch` — bug fixes, docs, cleanup, dependency updates
- `minor` — new user-facing features or behavior
- `major` — breaking changes to the CLI, config, or output format

Files inside `.bumpy/` whose names start with `_` (e.g. `_config.json`) are reserved for tooling and are never treated as bump files.

When PRs merge to `main`, the **Version Packages** workflow automatically keeps a single `version-packages` PR updated with all pending bump files, picks the highest bump level across them (`major` > `minor` > `patch`), and updates `CHANGELOG.md` with all pending bump entries. Merging that PR triggers a GoReleaser release.

### Optional: VERSION_PAT secret

By default the version-packages PR is opened with the repo's built-in `GITHUB_TOKEN`. GitHub's anti-recursion guard means CI will **not** run automatically on that PR unless you add a fine-grained PAT as a repository secret named `VERSION_PAT` with *Contents (read & write)* and *Pull requests (read & write)* permissions.
