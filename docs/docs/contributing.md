# Contributing

We welcome contributions of all kinds! Whether you're fixing a bug, adding a theme, or improving documentation, your help makes Kairo better for everyone.

## Development Setup

Kairo is written in Go. You'll need Go 1.21 or later to build it from source.

1. **Clone the repository**:
   ```bash
   git clone https://github.com/programmersd21/kairo.git
   cd kairo
   ```

2. **Build the binary**:
   ```bash
   go build -o kairo ./cmd/kairo
   ```

3. **Run tests**:
   ```bash
   go test ./...
   ```

## Architecture Overview

Kairo follows a clean architecture pattern:

- `cmd/kairo`: Entry point and CLI command parsing.
- `internal/core`: Core domain models (Task, Filter, Project).
- `internal/service`: Business logic and service orchestration.
- `internal/storage`: Persistence layer (SQLite).
- `internal/ui`: Bubble Tea components and TUI logic.
- `internal/lua`: Lua engine and plugin bindings.
- `internal/api`: External CLI and MCP API.

## Ways to Contribute

- **Themes**: Add a new theme to `internal/ui/theme/theme.go`.
- **Plugins**: Create and share Lua plugins in the `plugins/` directory.
- **Documentation**: Improve these docs by editing files in the `docs/` folder.
- **Bug Fixes**: Check the GitHub issues and submit a PR.

## Pull Request Process

1. Fork the repo and create your branch from `main`.
2. If you've added code that should be tested, add tests.
3. Ensure the test suite passes.
4. Issue that pull request!

## Code of Conduct

Please be respectful and helpful. We follow the [Contributor Covenant Code of Conduct](https://github.com/programmersd21/kairo/blob/main/CODE_OF_CONDUCT.md).
