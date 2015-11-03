package crdt

import (
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
	for _, tt := range []struct {
		add []interface{}
	}{
		{[]interface{}{}},
		{[]interface{}{1}},
		{[]interface{}{1, 2, 3}},
		{[]interface{}{1, 100, 1000, -1}},
		{[]interface{}{"alpha"}},
		{[]interface{}{"alpha", "beta"}},
		{[]interface{}{"alpha", "beta", 1, 2}},
	} {
		gset := NewGSet()

		expectedElems := map[interface{}]struct{}{}
		for _, i := range tt.add {
			expectedElems[i] = struct{}{}
			gset.Add(i)
		}

		actualElems := map[interface{}]struct{}{}
		for _, i := range gset.Elems() {
			actualElems[i] = struct{}{}
		}

		if !reflect.DeepEqual(expectedElems, actualElems) {
			t.Errorf("expected set to contain: %v, actual: %v", expectedElems, actualElems)
		}
	}
}

func TestGSetMarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		add      []interface{}
		expected string
	}{
		{[]interface{}{}, `{"type":"g-set","e":[]}`},
		{[]interface{}{1}, `{"type":"g-set","e":[1]}`},
		{[]interface{}{1, 2, 3}, `{"type":"g-set","e":[3,2,1]}`},
		{[]interface{}{1, 2, 3}, `{"type":"g-set","e":[1,2,3]}`},
		{[]interface{}{"alpha"}, `{"type":"g-set","e":["alpha"]}`},
		{[]interface{}{"alpha", "beta", "gamma"}, `{"type":"g-set","e":["alpha","beta","gamma"]}`},
		{[]interface{}{"alpha", 1, "beta", 2}, `{"type":"g-set","e":[1,2,"alpha","beta"]}`},
	} {

		gset := NewGSet()

		for _, e := range tt.add {
			gset.Add(e)
		}

		out, err := json.Marshal(gset)
		if err != nil {
			t.Fatalf("unexpected error on marshalling gset: %s", err)
		}

		a := struct {
			E []interface{} `json:"e"`
		}{}

		if err = json.Unmarshal(out, &a); err != nil {
			t.Fatalf("unexpected error on unmarshalling serialized %q: %s", tt.expected, err)
		}

		expectedMap := map[interface{}]struct{}{}
		for _, i := range a.E {
			expectedMap[i] = struct{}{}
		}

		if err = json.Unmarshal([]byte(tt.expected), &a); err != nil {
			t.Fatalf("unexpected error on unmarshalling serialized %q: %s", tt.expected, err)
		}

		actualMap := map[interface{}]struct{}{}
		for _, i := range a.E {
			actualMap[i] = struct{}{}
		}

		if !reflect.DeepEqual(expectedMap, actualMap) {
			t.Errorf("expected set to contain: %v, actual: %v", expectedMap, actualMap)
		}
	}
}
