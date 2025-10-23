package haste_router

import "strings"

const pathSeparater = "/"

type HasteRouterPathTraveller struct {
	Paths []string
	Index int
}

func NewHasteRouterPathTraveller(path string) *HasteRouterPathTraveller {
	path = strings.TrimPrefix(path, pathSeparater)
	path = strings.TrimSuffix(path, pathSeparater)
	paths := strings.Split(path, pathSeparater)
	traveller := &HasteRouterPathTraveller{
		Paths: paths,
		Index: 0,
	}
	return traveller
}

func (traveller *HasteRouterPathTraveller) HasNext() bool {
	return traveller.Index < len(traveller.Paths)
}

func (traveller *HasteRouterPathTraveller) Remaining() int {
	return len(traveller.Paths) - traveller.Index
}

func (traveller *HasteRouterPathTraveller) Peek() (string, bool) {
	if traveller.Index < len(traveller.Paths) {
		return traveller.Paths[traveller.Index], true
	}
	return "", false
}

func (traveller *HasteRouterPathTraveller) Next() (string, bool) {
	currentIndex := traveller.Index
	if currentIndex < len(traveller.Paths) {
		traveller.Index++
		return traveller.Paths[currentIndex], true
	}
	return "", false
}
