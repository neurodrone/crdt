package crdt

import "testing"

func TestTwoPhaseSetAdd(t *testing.T) {
	tpset := NewTwoPhaseSet()
	elem := "some-test-element"

	if tpset.Contains(elem) {
		t.Errorf("set should not contain %q, not yet added", elem)
	}

	tpset.Add(elem)
	if !tpset.Contains(elem) {
		t.Errorf("set should contain %q", elem)
	}
}

func TestTwoPhaseSetAddRemove(t *testing.T) {
	for _, tt := range []struct {
		elem        interface{}
		fnAddRemove func(*TwoPhaseSet, interface{})
	}{
		{
			"test-element",

			// First add the element and remove it. As a result
			// of this the set should not contain that element.
			func(tp *TwoPhaseSet, obj interface{}) {
				tp.Add(obj)
				tp.Remove(obj)
			},
		},
		{
			"test-element",

			// First remove the element and add it. Because of the
			// commutative property of this set the outcome should
			// still be the set that doesn't contain the element.
			func(tp *TwoPhaseSet, obj interface{}) {
				tp.Remove(obj)
				tp.Add(obj)
			},
		},
	} {
		tpset := NewTwoPhaseSet()

		if tpset.Contains(tt.elem) {
			t.Errorf("set should not contain elem %q", tt.elem)
		}

		tt.fnAddRemove(tpset, tt.elem)

		if tpset.Contains(tt.elem) {
			t.Errorf("set should not contain elem %q", tt.elem)
		}
	}
}
