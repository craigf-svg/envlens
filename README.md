# envlens

A Terminal UI for viewing, searching, and selecting environment variables, built with [Go](https://github.com/golang/go) and [Bubble Tea](https://github.com/charmbracelet/bubbletea).

<p align="center">
  <img width="100%" src="https://github.com/user-attachments/assets/ea7b42b7-61ea-4700-9c79-2cd5aa26ad0a" alt="User copying and searching environment variables in the terminal" />
</p>

## Motivation

I was frustrated with the way Windows handles environment variables. I had to use System Properties or PowerShell to inspect or copy a single variable, and there was no easy way to see how local `.env` files overlapped with system variables.

I built envlens to quickly search, view, and copy system and local variables in one place, which makes switching projects and debugging environments fast and painless.

## Usage

Start from root:

```bash
go run .
```

Any `.env` file in the current directory will be loaded automatically.
