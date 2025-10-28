# Sharp Objects

A simple TUI (Terminal User Interface) for monitoring and killing development processes running on specific ports.

> Named after the Gillian Flynn novel/HBO series. Because it can kill... processes.

## Features

- 🔍 Monitor node/bun/deno processes listening on ports
- 📊 Display port, process name, PID, and working directory
- 🎯 Focus on common dev ports (3000, 8000, 8787, etc.)
- ⚡ Quick kill functionality with confirmation
- 🔄 Manual refresh to update process list
- ⚙️  Configurable ports and process names

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

- `↑/↓` or `j/k` - Navigate through processes
- `K` (Shift+k) - Kill selected process (with confirmation)
- `r` - Refresh process list
- `q` or `Ctrl+C` - Quit

## Configuration

Create a config file at `~/.config/sharp-objects/config.yaml` to customize:

```yaml
# Ports to monitor
ports:
  - 3000
  - 3001
  - 8000
  - 8080
  - 8787
  - 5173

# Process names to monitor
processes:
  - node
  - bun
  - deno
  - wrangler
```

### Default Configuration

If no config file exists, these defaults are used:

- **Ports:** 3000, 3001, 3002, 4321, 5000, 5173, 5174, 8000, 8001, 8080, 8081, 8787, 8788, 9229, 5432, 1313, 7878, 7272
- **Processes:** node, bun, deno, wrangler

The default port list is based on commonly used dev server ports including Next.js, Vite, Astro, Wrangler, Node debugger, PostgreSQL, and more.

## How It Works

The tool:
1. Uses `lsof` to find listening TCP processes
2. Filters for configured ports and process names
3. Displays them in an interactive table
4. Shows the working directory where each process is running from
5. Allows you to quickly kill processes that are blocking ports

## Requirements

- Go 1.24+
- `lsof` command (standard on macOS and Linux)

## Why?

When developing with multiple servers (Next.js, Vite, Wrangler, Expo, etc.), ports can get stuck or you might forget what's running where. This tool gives you a quick overview and easy cleanup.
