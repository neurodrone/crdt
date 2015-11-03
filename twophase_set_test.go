package crdt

import (
	"bytes"
	"encoding/json"
	"testing"
)

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

func TestTwoPhaseSetMarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		add, rm  []interface{}
		expected string
	}{
		{[]interface{}{}, []interface{}{}, `{"type":"2p-set","a":[],"r":[]}`},
		{[]interface{}{"alpha"}, []interface{}{}, `{"type":"2p-set","a":["alpha"],"r":[]}`},
		{[]interface{}{}, []interface{}{"beta"}, `{"type":"2p-set","a":[],"r":["beta"]}`},
		{[]interface{}{"alpha"}, []interface{}{"beta"}, `{"type":"2p-set","a":["alpha"],"r":["beta"]}`},
		{[]interface{}{"alpha"}, []interface{}{"alpha"}, `{"type":"2p-set","a":["alpha"],"r":["alpha"]}`},

		{[]interface{}{1}, []interface{}{}, `{"type":"2p-set","a":[1],"r":[]}`},
		{[]interface{}{}, []interface{}{2}, `{"type":"2p-set","a":[],"r":[2]}`},
		{[]interface{}{1}, []interface{}{2}, `{"type":"2p-set","a":[1],"r":[2]}`},
		{[]interface{}{1}, []interface{}{1}, `{"type":"2p-set","a":[1],"r":[1]}`},
	} {
		tpset := NewTwoPhaseSet()

		for _, e := range tt.add {
			tpset.Add(e)
		}

		for _, e := range tt.rm {
			tpset.Remove(e)
		}

		out, err := json.Marshal(tpset)
		if err != nil {
			t.Fatalf("unexpected error marshalling tpset: %s", err)
		}

		if !bytes.Equal([]byte(tt.expected), out) {
			t.Errorf("expected marshalled bytes: %q, actual: %q", tt.expected, out)
		}
	}
}
