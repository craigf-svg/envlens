package main

type SelectionModel struct {
	variables []string
	cursor    int
	selected  map[int]struct{}
	choices   []string
}
