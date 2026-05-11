# Contributing

Run the checks before opening a pull request:

```bash
go test ./...
go vet ./...
go build ./cmd/portwatch
```

Keep changes focused. For scanner behavior, prefer provider interfaces and tests with fakes over OS-specific assumptions.

Release-significant changes need a `.bumpy/*.md` bump file. Do not edit `CHANGELOG.md` manually.
The release workflow expects `BUMPY_GH_TOKEN` to be configured so generated version PRs and release tags trigger their follow-up workflows.

Use bumpy's package-keyed format:

```md
---
portwatcher: patch
---
Brief changelog entry for users.
```

Use `patch` for fixes, docs, cleanup, and dependency updates; `minor` for new user-facing behavior; and `major` for breaking CLI, config, or output changes.
