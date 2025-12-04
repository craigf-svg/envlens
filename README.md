# envlens
![Go](https://img.shields.io/badge/Go-273849?style=for-the-badge&logo=go&logoColor=64b5f6)
![Bubble Tea](https://img.shields.io/badge/Bubble_Tea-273849?style=for-the-badge&logo=data:image/svg+xml;base64,PHN2ZyB3aWR0aD0iMjQiIGhlaWdodD0iMjQiIHZpZXdCb3g9IjAgMCAyNCAyNCIgZmlsbD0ibm9uZSIgeG1sbnM9Imh0dHA6Ly93d3cudzMub3JnLzIwMDAvc3ZnIj4KICA8Y2lyY2xlIGN4PSIxMiIgY3k9IjYiIHI9IjUiIGZpbGw9IiM2NGI1ZjYiLz4KICA8Y2lyY2xlIGN4PSI3IiBjeT0iMTQuNSIgcj0iMy41IiBmaWxsPSIjNjRiNWY2Ii8+CiAgPGNpcmNsZSBjeD0iMTYiIGN5PSIxNy41IiByPSIzLjUiIGZpbGw9IiM2NGI1ZjYiLz4KPC9zdmc+)

> Inspect, search, and copy environment variables from your terminal.

A cross-platform Terminal UI built with [Go](https://github.com/golang/go) and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

<p align="center">
  <img width="100%" src="https://github.com/user-attachments/assets/ea7b42b7-61ea-4700-9c79-2cd5aa26ad0a" alt="User copying and searching environment variables in the terminal" />
</p>

## Motivation

I was frustrated with the way Windows handles environment variables. I had to use System Properties or PowerShell to inspect or copy a single variable, and there was no easy way to see how local `.env` files overlapped with system variables.

I built envlens to quickly search, view, and copy system and local variables in one place, which makes switching projects and debugging environments fast and painless.

Works on Windows, macOS, and Linux.

## Features

+ **Unified view**: Automatically loads `.env` files and merges them with system variables
+ **Fuzzy search**: Quickly filter by variable name or value
+ **Clipboard integration**: Copy single or select multiple variables
+ **Mask private keys**: Toggle visibility to hide keys

## Quick Start

### Install with Go

Requires Go 1.21+

```bash
go install github.com/craigf-svg/envlens@latest
```

### Or download the binary

Grab the latest release for your platform from the [Releases page](https://github.com/craigf-svg/envlens/releases).

### Launch

```bash
envlens
```

Any `.env` file in the current directory will be loaded automatically.

## Usage

### Keybindings

| Key | Action |
|-----|--------|
| `↑/↓` or `k/j` | Navigate up/down |
| `Enter` or `Space` | Select/deselect variable |
| `y` | Copy current variable to clipboard |
| `Y` | Copy selected variables to clipboard |
| `Tab` | Toggle value visibility (mask/unmask) |
| `s` | Enter search mode |
| `d` | View local `.env` file |
| `q` or `Ctrl+C` | Quit |

### Flags

```bash
envlens --demo    # Run with demo data (no real env vars)
```

## Contributing

```bash
git clone https://github.com/craigf-svg/envlens.git
cd envlens
go run .
```

### Submit a pull request

If you'd like to contribute, please fork the repository and open a pull request against `master`.

