# Installation

Kairo is a single binary with zero external dependencies (aside from Git for sync). It runs on macOS, Linux, and Windows.

## Quick Install

### macOS (Homebrew)
```bash
brew install programmersd21/kairo/kairo
```

### Linux / macOS (Curl)
```bash
curl -fsSL https://raw.githubusercontent.com/programmersd21/kairo/main/scripts/install.sh | bash
```

### Windows (PowerShell)
```powershell
iwr -useb https://raw.githubusercontent.com/programmersd21/kairo/main/scripts/install.ps1 | iex
```

### From Source (Go)
```bash
go install github.com/programmersd21/kairo/cmd/kairo@latest
```

## First Run

After installing, simply run:
```bash
kairo
```

On your first run, Kairo will:
1. Initialize a SQLite database in your application data directory.
2. Create a default `config.toml`.
3. Launch the **Welcome Tour**.

<img src={require('../assets/tasks_list.png').default} alt="Kairo Welcome Tour" />

## Configuration Path

Kairo stores its configuration and database in the following locations:

- **Linux**: `~/.config/kairo/`
- **macOS**: `~/Library/Application Support/kairo/`
- **Windows**: `%APPDATA%\kairo\`

## Verifying Installation

To verify that Kairo is installed correctly, you can check the version:
```bash
kairo --version
```
