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

I built envlens to quickly search, view, and copy system and local variables in one place, which makes switching projects and debugging environments much easier.

Works on Windows, macOS, and Linux.

## Features

+ **Unified view**: envlens automatically loads your .env files and merges them with your system variables
+ **Fuzzy search**: Type a few letters of a name or value and find what you need fast
+ **Clipboard integration**: Copy one or several variables easily
+ **Mask private keys**: Hide values when you need to

## Quick Start

### Install with Go

Requires Go 1.21+

```bash
go install github.com/craigf-svg/envlens@latest
```

### Or download the binary

Grab the latest release for your platform from the [Releases page](https://github.com/craigf-svg/envlens/releases).

### Run from anywhere
<details>
<summary>Adding to PATH</summary>

**Go installation**:

- On Windows, head to System Properties → Environment Variables and add `%USERPROFILE%\go\bin` to your PATH.
- On macOS or Linux, put `export PATH=$PATH:$HOME/go/bin` in your `~/.bashrc` or `~/.zshrc`.

**Binary download**:
- On Windows, add the folder where your binary lives to PATH in System Properties → Environment Variables.
- On macOS or Linux, run `chmod +x envlens && sudo mv envlens /usr/local/bin/`.
</details>

Once your PATH is set up, close your terminal and open it again.

To make sure everything worked, type `envlens --version`.

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

