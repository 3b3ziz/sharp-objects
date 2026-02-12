# Sharp Objects

A simple TUI (Terminal User Interface) for monitoring and killing development processes running on specific ports.

> Named after the Gillian Flynn novel/HBO series. Because it can kill... processes.

## Features

- 🔍 Show ALL listening TCP ports (not just specific dev ports)
- 📊 Display port, process name, PID, and working directory
- 💻 Display full command line with smart path cleaning
- 🎭 Toggle between dev-only and all processes (press `f`)
- ⚡ Quick kill functionality with double-tap confirmation
- 🔄 Manual refresh to update process list
- ⚙️  Smart filtering of system services by default

## Installation

```bash
go build -o sharp-objects
```

Or install it globally:

```bash
go install
```

## Usage

Simply run:

```bash
./sharp-objects
```

### Keyboard Controls

- `↑/↓` - Navigate through processes
- `k` - Kill selected process (press twice to confirm)
- `r` - Refresh process list
- `f` - Toggle filter (dev only ↔ all)
- `q` or `Ctrl+C` - Quit

## Configuration

**Note:** Sharp Objects now shows ALL listening TCP ports by default, so the configuration file is optional and mainly used for identifying "dev" processes when filtering.

You can optionally create a config file at `~/.config/sharp-objects/config.yaml`:

```yaml
# Process names to consider as "dev" processes (for filtering)
processes:
  - node
  - bun
  - deno
  - wrangler
```

### Default Configuration

If no config file exists, these process names are used for dev filtering:
- **Processes:** node, bun, deno, wrangler

When the filter is active (default), system services are hidden. Press `f` to toggle between dev-only and all processes.

## How It Works

The tool:
1. Uses `lsof` to find ALL listening TCP processes
2. Displays cleaned command lines (with smart path truncation)
3. Shows port, PID, full command, and working directory in an interactive table
4. Filters system services by default (toggle with `f`)
5. Double-tap `k` for safe process killing
6. Auto-refreshes after killing a process

## Requirements

- Go 1.24+
- `lsof` command (standard on macOS and Linux)

## Why?

When developing with multiple servers (Next.js, Vite, Wrangler, Expo, etc.), ports can get stuck or you might forget what's running where. This tool gives you a quick overview of ALL listening processes, with smart filtering to focus on dev processes by default, and easy cleanup with double-tap confirmation.
