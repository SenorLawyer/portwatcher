# Project Notes

- This is a Go CLI. Keep runtime code in Go.
- Run `bun run check` before release work.
- Bumpy is only for release bookkeeping.
- GoReleaser builds binaries from `v*` tags.
- Release-significant changes need a `.bumpy/*.md` bump file.
- Use `patch` for fixes, docs, cleanup, and dependency updates.
- Use `minor` for new user-facing behavior.
- Use `major` for breaking CLI, config, or output changes.
- Do not edit `CHANGELOG.md` manually.
