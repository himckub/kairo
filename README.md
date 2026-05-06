<div align="center">

<img src="screenshots/logo.png" alt="Kairo" width="80" />

# Kairo

**The terminal task manager for developers who live in their editor.**

No browser tabs. No subscriptions. No mouse. Just your tasks — exactly where your brain already is.

<br/>

[![Release](https://img.shields.io/github/v/release/programmersd21/kairo?style=for-the-badge&logo=github&color=7c3aed)](https://github.com/programmersd21/kairo/releases)
[![CI](https://img.shields.io/github/actions/workflow/status/programmersd21/kairo/ci.yml?branch=main&style=for-the-badge&logo=github&color=2563eb)](https://github.com/programmersd21/kairo/actions)
[![Go Report](https://img.shields.io/badge/go%20report-A%2B-brightgreen?style=for-the-badge&logo=go&logoColor=white)](https://goreportcard.com/report/github.com/programmersd21/kairo)
[![Downloads](https://img.shields.io/github/downloads/programmersd21/kairo/total?style=for-the-badge&logo=github)](https://github.com/programmersd21/kairo/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-f59e0b?style=for-the-badge)](https://opensource.org/licenses/MIT)

<br/>

![Kairo Demo](screenshots/demo.gif)

</div>

---

## Why Kairo?

Most task managers make you context-switch out of your flow. Kairo doesn't.

| Pain point | What Kairo does |
|---|---|
| GUI apps break your focus | Lives entirely in your terminal |
| Cloud tools own your data | Everything local, stored in SQLite |
| Plain-text tools lack structure | Full tagging, filtering, and fuzzy search |
| Legacy TUIs feel dated | Modern, animated, keyboard-first UX |

---

## Quick Start

**macOS (Homebrew)**
```bash
brew install programmersd21/kairo/kairo
```

**Linux / macOS**
```bash
curl -fsSL https://raw.githubusercontent.com/programmersd21/kairo/main/scripts/install.sh | bash
```

**Windows (PowerShell)**
```powershell
iwr -useb https://raw.githubusercontent.com/programmersd21/kairo/main/scripts/install.ps1 | iex
```

**Go**
```bash
go install github.com/programmersd21/kairo/cmd/kairo@latest
```

Then just run:
```bash
kairo
```

Press `n` to create your first task. `ctrl+s` to save. That's it.

> Works best on Alacritty. Some terminals may have rendering quirks — see [#16](https://github.com/programmersd21/kairo/issues/16).

---

## Features

![Kairo Home Screen](screenshots/home_screen.png)

### ⚡ Genuinely Fast
Sub-millisecond fuzzy search. Vim bindings (`j/k/gg/G`). Natural language deadlines like `tomorrow 10am` or `next friday`. Full keyboard control — you never touch the mouse.

### 🗂 Nested Tasks & Hierarchy
Organize work into deep hierarchies. Nest tasks via the **Parent** field in the editor, collapse/expand with `Space`, and export/import with full structure preserved — across JSON, CSV, Markdown, and plain text.

### 🔁 Recurring Tasks
Tasks reappear automatically on a schedule. Weekly (`mon,wed,fri`) or monthly (`15`). When completed, Kairo generates the next instance immediately with a smart due-date preview.

### 🔒 Your Data, Locally
SQLite with WAL mode. Fully offline. Optional Git-backed sync — no backend, no account, no lock-in. Export to JSON, CSV, Markdown, or plain text on demand.

### 🧭 Interactive Stats Dashboard
Press `s` to open a next-gen "Command Center". Visualize your **Productivity DNA**, track real-time momentum, and get behavioral insights like "You complete 73% more tasks at night". Fully animated, keyboard-driven, and deeply insightful.

### 🤖 AI — Optional, Never Intrusive
Gemini integration (`gemini-3.1-flash-lite-preview` / `gemini-2.5-flash-lite` / `gemini-2.0-flash-lite`). Toggle with `ctrl+a`. Create and manage complex recurring tasks with natural language. Invisible until you need it.

### 🎨 Beautiful by Default
32 built-in themes with edge-to-edge background coverage. Live switching with `t`. Bento-style layout. Real-time Markdown preview (`ctrl+p`). Cinematic create/complete/delete animations — or disable them entirely in `config.toml`.

### 🧩 Extensible to the Core
A Lua plugin system hooks into task events. A headless CLI API enables full scripting. An MCP server opens Kairo to AI agents — with complete support for recurring schedules and nested hierarchies.

---

## Keyboard Shortcuts

| Key | Action |
|---|---|
| `n` | New task |
| `e` | Edit task |
| `z` | Complete task |
| `d` | Delete task |
| `Space` | Collapse / expand subtasks |
| `s` | Stats dashboard |
| `f` | Filter by tag |
| `t` | Switch theme |
| `ctrl+p` | Command palette / Markdown preview |
| `ctrl+a` | AI panel |
| `ctrl+s` | Settings |
| `x` | Import / Export |
| `?` | Help |
| `ctrl+d` | Welcome tour |

<div align="center">
  <img src="screenshots/new_task.png" width="30%" />
  <img src="screenshots/filter_tags.png" width="30%" />
  <img src="screenshots/help_menu.png" width="30%" />
  <img src="screenshots/settings_menu.png" width="30%" />
  <img src="screenshots/theme_menu.png" width="30%" />
  <img src="screenshots/dashboard.png" width="30%" />
</div>

---

## CLI Automation

Kairo exposes a full CLI API for scripting and CI pipelines, with complete support for `parent_id` and `collapsed` state:

```bash
# Create a task
kairo api create --title "Finish report" --priority 1

# List by tag
kairo api list --tag work

# Mark complete
kairo api update --id <id> --status done

# Export everything
kairo export --format markdown
```

---

## Lua Plugin System

```lua
local plugin = {
    id = "my-plugin",
    name = "My Plugin",
    version = "1.0.0"
}

kairo.on("task_create", function(event)
    kairo.notify("New task: " .. event.task.title)
end)

return plugin
```

Browse [sample plugins →](https://github.com/programmersd21/kairo/tree/main/plugins)

---

## Architecture

```
Input  (CLI · TUI · Lua · AI)
       ↓
Task Service  (single source of truth)
       ↓
SQLite (WAL)  +  optional Git sync
       ↓
Bubble Tea TUI  (instant rendering)
```

**Stack:** Bubble Tea · Lip Gloss · SQLite (WAL) · GopherLua · Gemini API · Git

---

## Everything Included

| Feature | Status |
|---|---|
| Local-first SQLite storage | ✅ |
| Nested tasks & folders | ✅ |
| 32 themes, live switching | ✅ |
| Keyboard-only workflow | ✅ |
| Recurring tasks | ✅ |
| Git sync (no backend) | ✅ |
| Lua plugin system | ✅ |
| CLI automation API | ✅ |
| AI assistant (optional) | ✅ |
| MCP server | ✅ |
| Free & open source | ✅ |

---

## Configuration

Kairo can be configured via `config.toml` in your application data directory.

### Task List
You can customize the fields shown on the right side of the task list:

```toml
[list.order]
right = ["tags", "due", "priority"]
```

Valid values for `right` are: `tags`, `due`, `priority`.

### Task Fields
*   **wait_until**: Hide a task from the task list until the specified datetime. If the task is recurring, new instances are not generated/shown until `wait_until` has passed. Format: `yyyy-MM-dd HH:mm`.
*   **until**: Stop generating new recurring instances after the specified datetime. Existing instances may remain visible. Format: `yyyy-MM-dd HH:mm`.

Auto-generated on first run at:

- **Linux:** `~/.config/kairo/config.toml`
- **macOS:** `~/Library/Application Support/kairo/config.toml`
- **Windows:** `%APPDATA%\kairo\config.toml`

| Option | Description | Default |
|---|---|---|
| `theme` | UI theme name | `catppuccin` |
| `vim_mode` | Vim keybindings | `false` |
| `show_help` | Help footer | `true` |
| `show_id` | Task IDs in detail view | `true` |
| `animations` | UI animations | `true` |
| `rainbow` | Animated rainbow logo | `false` |

Prefer not to edit files? `ctrl+s` opens the in-app settings menu.

---

## Roadmap

- Encrypted multi-workspace support
- Event-sourced sync engine
- Sandboxed plugin environment
- Smart task suggestions
- Plugin marketplace
- Streaming performance optimizations

---

## Star History

<a href="https://www.star-history.com/?repos=programmersd21%2Fkairo&type=date&legend=top-left">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="https://api.star-history.com/chart?repos=programmersd21/kairo&type=date&theme=dark&legend=top-left" />
    <source media="(prefers-color-scheme: light)" srcset="https://api.star-history.com/chart?repos=programmersd21/kairo&type=date&legend=top-left" />
    <img alt="Star History Chart" src="https://api.star-history.com/chart?repos=programmersd21/kairo&type=date&legend=top-left" />
  </picture>
</a>

---

## Contributing

PRs are welcome — especially for themes, plugins, performance, and docs. If something bugs you, fix it.

Huge thanks to [@Tornado300](https://github.com/Tornado300) for key bug fixes that made Kairo better for everyone.

---

<div align="center">

**If Kairo saves you time, a ⭐ helps other developers find it.**

<br/>

*Built for the terminal. Built for focus. Built for you.*

</div>
