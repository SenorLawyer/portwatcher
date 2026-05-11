# Contributing

Run the checks before opening a pull request:

```bash
go test ./...
go vet ./...
go build ./cmd/portwatch
```

Keep changes focused. For scanner behavior, prefer provider interfaces and tests with fakes over OS-specific assumptions.
