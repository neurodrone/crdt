package crdt

import (
	"bytes"
	"encoding/json"
	"reflect"
	"testing"
)

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

	expectedCount := 1
	if gset.Len() != expectedCount {
		t.Errorf("set expected to contain %d elements, actual: %d", expectedCount, gset.Len())
	}
}

func TestGSetElems(t *testing.T) {
	gset := NewGSet()

	if !reflect.DeepEqual([]interface{}{}, gset.Elems()) {
		t.Errorf("gset should be empty but found: %v", gset.Elems())
	}

	var expectedElems []interface{}
	for _, i := range []int{1, 2, 3} {
		expectedElems = append(expectedElems, i)
		gset.Add(i)
	}

	if !reflect.DeepEqual(expectedElems, gset.Elems()) {
		t.Errorf("expected set to contain: %+v, actual: %+v", expectedElems, gset.Elems())
	}
}

func TestGSetMarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		add      []interface{}
		expected string
	}{
		{[]interface{}{}, `{"type":"g-set","e":[]}`},
		{[]interface{}{1}, `{"type":"g-set","e":[1]}`},
		{[]interface{}{"alpha"}, `{"type":"g-set","e":["alpha"]}`},
	} {

		gset := NewGSet()

		for _, e := range tt.add {
			gset.Add(e)
		}

		out, err := json.Marshal(gset)
		if err != nil {
			t.Fatalf("unexpected error on marshalling gset: %s", err)
		}

		if !bytes.Equal([]byte(tt.expected), out) {
			t.Errorf("expected marshalled output: %q, actual: %q", tt.expected, out)
		}
	}
}
