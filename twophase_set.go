package crdt

import "encoding/json"

var (
	// TwoPhaseSet should implement Set.
	_ Set = &TwoPhaseSet{}
)

// TwoPhaseSet supports both addition and removal of
// elements to set.
type TwoPhaseSet struct {
	addSet *GSet
	rmSet  *GSet
}

// NewTwoPhaseSet returns a new instance of TwoPhaseSet.
func NewTwoPhaseSet() *TwoPhaseSet {
	return &TwoPhaseSet{
		addSet: NewGSet(),
		rmSet:  NewGSet(),
	}
}

// Add inserts element into the TwoPhaseSet.
func (t *TwoPhaseSet) Add(elem interface{}) {
	t.addSet.Add(elem)
}

// Remove deletes the element from the set.
func (t *TwoPhaseSet) Remove(elem interface{}) {
	t.rmSet.Add(elem)
}

// Contains returns true if the set contains the element.
// The set is said to contain the element if it is present
// in the add-set and not in the remove-set.
func (t *TwoPhaseSet) Contains(elem interface{}) bool {
	return t.addSet.Contains(elem) && !t.rmSet.Contains(elem)
}

type tpsetJSON struct {
	T string        `json:"type"`
	A []interface{} `json:"a"`
	R []interface{} `json:"r"`
}

// MarshalJSON serializes TwoPhaseSet into an JSON array.
func (t *TwoPhaseSet) MarshalJSON() ([]byte, error) {
	return json.Marshal(&tpsetJSON{
		T: "2p-set",
		A: t.addSet.Elems(),
		R: t.rmSet.Elems(),
	})
}
