# envlens

A Terminal UI for viewing, searching, and selecting environment variables, built with [Go](https://github.com/golang/go) and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Usage

Run the program in your terminal:

```bash
go run .
```

If a `.env` file is present in the current directory, it will be loaded automatically.

## Controls

### General

| Key | Action |
| :--- | :--- |
| `q` / `ctrl+c` | Quit the application |
| `↑` / `k` | Move cursor up |
| `↓` / `j` | Move cursor down |
| `space` / `enter` | Select/deselect variable |
| `y` | Copy focused variable to clipboard |
| `Y` | Copy all selected variables to clipboard |

### Modes

| Key | Mode | Description |
| :--- | :--- | :--- |
| `s` | **Search** | Type to filter the variable list. Press `esc` to return. |
| `d` | **Local .env** | Switch to view variables from the local `.env` file. Press `esc` to return. |