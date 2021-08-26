package main

/*
General utility types, functions, and constants used in the program.
*/
type matchedLine struct {
	Text  string
	Index int
}

type ByIndex []matchedLine

func (a ByIndex) Len() int {
	return len(a)
}

func (a ByIndex) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a ByIndex) Less(i, j int) bool {
	return a[i].Index < a[j].Index
}

type SearchMode int

const (
	String SearchMode = iota
	Regex
)