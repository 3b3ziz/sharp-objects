# Sharp Objects - Development Notes

A TUI (Terminal User Interface) for quickly finding and killing development processes on monitored ports.

## Project Structure

- `main.go` - Main TUI application using Bubble Tea framework
- `process.go` - Process discovery and management functions
- `config.go` - Configuration loading (YAML-based)
- `deploy.sh` - Build and install script

## Key Features

- **Double-tap kill**: Press 'k' once to mark, press 'k' again to kill (cancels on navigation or other keys)
- **Auto-refresh**: Process list refreshes immediately after killing
- **Navigation**: Arrow keys (↑/↓) only
- **Manual refresh**: 'r' key
- **Quit**: 'q' or Ctrl+C

## Architecture

### Model State
- `processes []Process` - Current process list
- `cursor int` - Selected row index
- `pendingKill int` - Index of process marked for kill (-1 = none)
- `error string` - Error message display
- `successMsg string` - Success message (auto-clears after 3s)

### Key Bindings
- `↑/↓` - Navigate
- `k` - First press marks, second press kills
- `r` - Refresh process list
- `q` - Quit

### Process Detection
Uses `lsof` to find processes listening on configured ports. Filters by:
- Target ports (3000, 5173, etc.)
- Target process names (node, bun, deno, wrangler)

## Building & Installing

Run the deploy script:
```bash
./deploy.sh
```

This builds the binary and installs it to `~/.local/bin/sharp-objects`.

## Configuration

Default ports and processes are hardcoded in `process.go`. Can be extended to support YAML config file at `~/.config/sharp-objects/config.yaml`.

## Dependencies

- `github.com/charmbracelet/bubbletea` - TUI framework
- `github.com/charmbracelet/lipgloss` - Styling
- `gopkg.in/yaml.v3` - Config parsing

## Common Tasks

### Adding a new port
Edit `process.go`, add to `TargetPorts` array in `DefaultConfig()`.

### Changing UI colors
Edit style definitions at top of `main.go` (headerStyle, selectedStyle, etc.).

### Modifying kill behavior
See `killMsg` case in `main.go` Update() function.
