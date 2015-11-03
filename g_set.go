package crdt

import "encoding/json"

type mainSet map[interface{}]struct{}

var (
	// GSet should implement the set interface.
	_ Set = &GSet{}
)

// Gset is a grow-only set.
type GSet struct {
	mainSet mainSet
}

// NewGSet returns an instance of GSet.
func NewGSet() *GSet {
	return &GSet{
		mainSet: mainSet{},
	}
}

// Add lets you add an element to grow-only set.
func (g *GSet) Add(elem interface{}) {
	g.mainSet[elem] = struct{}{}
}

// Contains returns true if an element exists within the
// set or false otherwise.
func (g *GSet) Contains(elem interface{}) bool {
	_, ok := g.mainSet[elem]
	return ok
}

// Len returns the no. of elements present within GSet.
func (g *GSet) Len() int {
	return len(g.mainSet)
}

// Elems returns all the elements present in the set.
func (g *GSet) Elems() []interface{} {
	elems := make([]interface{}, 0, len(g.mainSet))

	for elem := range g.mainSet {
		elems = append(elems, elem)
	}

	return elems
}

type gsetJSON struct {
	T string        `json:"type"`
	E []interface{} `json:"e"`
}

// MarshalJSON will be used to generate a serialized output
// of a given GSet.
func (g *GSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&gsetJSON{
		T: "g-set",
		E: g.Elems(),
	})
}
