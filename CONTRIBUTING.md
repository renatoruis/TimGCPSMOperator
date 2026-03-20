# Contributing to TIMGCPSMOPERATOR

Thanks for your interest in contributing.

## Code of conduct

By participating, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## How to contribute

1. **Issues** — open an issue to discuss bugs or features before large changes.
2. **Pull requests** — fork, branch from `main`, keep changes focused, and reference the issue when applicable.
3. **Tests** — run `make test` (or `go test ./...`) and ensure `go fmt` / `go vet` pass.

## Development

```bash
go mod download
make build
make test
```

CRDs are generated from Go types via `make manifests` (requires `controller-gen`).

## Commits

Clear, short commit messages are preferred (e.g. conventional style: `fix:`, `feat:`, `docs:`).

## License

By contributing, you agree that your contributions are licensed under the **Apache License 2.0**, the same as this project.
