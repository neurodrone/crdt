package crdt

import "testing"

func TestGSetAddContains(t *testing.T) {
	gset := NewGSet()

	elem := "some-test-element"
	if gset.Contains(elem) {
		t.Errorf("set should not contain %q element", elem)
	}

	gset.Add(elem)
	if !gset.Contains(elem) {
		t.Errorf("set should contain %q element", elem)
	}
}
