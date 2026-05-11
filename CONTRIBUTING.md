# Contributing

Run the checks before opening a pull request:

```bash
go test ./...
go vet ./...
go build ./cmd/portwatch
```

Keep changes focused. For scanner behavior, prefer provider interfaces and tests with fakes over OS-specific assumptions.

## Bump files

For any PR that contains a user-visible change, add a bump file at `.bumpy/<short-description>.md`:

```yaml
---
bump: minor
---
Brief changelog description of the change.
```

Bump levels:
- `patch` — bug fixes, docs, cleanup, dependency updates
- `minor` — new user-facing features or behaviour
- `major` — breaking changes to the CLI, config, or output format

Files inside `.bumpy/` whose names start with `_` (e.g. `_config.json`) are reserved for tooling and are never treated as bump files.

When your PR merges to `main`, the **Version Packages** workflow automatically creates (or updates) a `version-packages` PR that bumps the version and updates `CHANGELOG.md`. Merging that PR triggers a GoReleaser release.

### Optional: VERSION_PAT secret

By default the version-packages PR is opened with the repo's built-in `GITHUB_TOKEN`. GitHub's anti-recursion guard means CI will **not** run automatically on that PR unless you add a fine-grained PAT as a repository secret named `VERSION_PAT` with *Contents (read & write)* and *Pull requests (read & write)* permissions.
