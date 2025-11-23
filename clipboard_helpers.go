package main

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
)

func copySelectedVarsToClipboard(selected map[int]struct{}, variables []string) (string, error) {
	if len(selected) == 0 {
		return "No variables selected to copy", nil
	}

	selectedVars := make([]string, 0, len(selected))
	for idx := range selected {
		selectedVars = append(selectedVars, variables[idx])
	}
	v := strings.Join(selectedVars, "\n")
	if err := clipboard.WriteAll(v); err != nil {
		return "", err
	}

	if len(selectedVars) == 1 {
		return "Copied 1 variable to clipboard", nil
	}
	return fmt.Sprintf("Copied %d variables to clipboard", len(selectedVars)), nil
}
