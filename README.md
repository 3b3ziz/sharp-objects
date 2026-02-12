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

### Pre-built Binaries

Download the latest release for your platform from the [releases page](https://github.com/3b3ziz/sharp-objects/releases).

Available platforms:
- macOS (Intel and Apple Silicon)
- Linux (amd64 and arm64)

Extract the archive and move the binary to a directory in your PATH:

```bash
# Example for macOS/Linux
tar -xzf sharp-objects_*.tar.gz
sudo mv sharp-objects /usr/local/bin/
```

### Build from Source

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

Tired of running `lsof -i :3000` to find the PID, then `kill -9 <pid>` to free up a port?

This tool eliminates the tedium:
- No more memorizing lsof syntax
- No more copying/pasting PIDs
- No more accidentally killing the wrong process
- See ALL your listening ports at a glance
- Double-tap confirmation prevents mistakes
- Smart filtering hides system noise

Instead of this:
```bash
lsof -i :3000                    # Find what's on port 3000
kill -9 12345                    # Kill it
lsof -i TCP -s TCP:LISTEN        # What else is running?
# ... repeat for each port ...
```

Just run `sharp-objects` and navigate with arrow keys.
