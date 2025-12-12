package main

type SelectionModel struct {
	variables []string
	cursor    int
	selected  map[int]struct{}
	hidden    map[int]struct{}
	choices   []string
}
