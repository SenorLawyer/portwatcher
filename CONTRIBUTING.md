# Contributing

Run the checks before opening a pull request:

```bash
bun run check
```

Keep changes focused. For scanner behavior, prefer provider interfaces and tests with fakes over OS-specific assumptions.

Release-significant changes need a `.bumpy/*.md` bump file. Do not edit `CHANGELOG.md` manually.
