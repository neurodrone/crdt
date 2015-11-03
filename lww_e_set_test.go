package crdt

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
)

func TestLWWESetAddContains(t *testing.T) {
	lww, err := NewLWWSet()
	if err != nil {
		t.Fatalf("Unexpected error creating lwwset: %s", err)
	}

	testStr := "object1"
	if lww.Contains(testStr) {
		t.Errorf("set should not contain elem: %q", testStr)
	}

	lww.Add(testStr)

	if !lww.Contains(testStr) {
		t.Errorf("Expected set to contain: %v, but not found", testStr)
	}
}

func TestLWWESetAddRemoveContains(t *testing.T) {
	lww, err := NewLWWSet()
	if err != nil {
		t.Fatalf("Unexpected error creating lwwset: %s", err)
	}

	testStr := "object2"
	lww.Add(testStr)
	lww.Remove(testStr)

	if lww.Contains(testStr) {
		t.Errorf("Expected set to not contain: %v, but found", testStr)
	}
}

func TestInvalidBias(t *testing.T) {
	var InvalidBias BiasType = "invalid"

	if _, err := NewLWWSetWithBias(InvalidBias); err != ErrNoSuchBias {
		t.Errorf("error expected here when invalid bias provided: %s", err)
	}

	mock := clock.NewMock()

	lww := &LWWSet{
		addMap: make(map[interface{}]time.Time),
		rmMap:  make(map[interface{}]time.Time),
		bias:   InvalidBias,
		clock:  mock,
	}

	elem := "object1"

	// Remove the element before it is added. Since the time of adding
	// the element is greater than the time to remove it this set should
	// technically comprise of that element. But because the Bias is invalid
	// verify that it should always generate a false response.
	lww.Add(elem)
	mock.Add(-1 * time.Minute)
	lww.Remove(elem)

	if lww.Contains(elem) {
		t.Errorf("set should not contain element and should trigger an invalid case")
	}

}

func TestLWWESetAddRemoveConflict(t *testing.T) {
	for _, tt := range []struct {
		bias       BiasType
		testObject string
		elapsed    time.Duration
		testFn     func(*LWWSet, interface{}) bool
	}{
		{
			BiasAdd,
			"object2",
			0,
			func(l *LWWSet, obj interface{}) bool { return l.Contains(obj) },
		},
		{
			BiasRemove,
			"object3",
			0,
			func(l *LWWSet, obj interface{}) bool { return !l.Contains(obj) },
		},
		{
			BiasAdd,
			"object4",
			1 * time.Minute,
			func(l *LWWSet, obj interface{}) bool { return !l.Contains(obj) },
		},
		{
			BiasAdd,
			"object5",
			-1 * time.Minute,
			func(l *LWWSet, obj interface{}) bool { return l.Contains(obj) },
		},
		{
			BiasRemove,
			"object6",
			1 * time.Minute,
			func(l *LWWSet, obj interface{}) bool { return !l.Contains(obj) },
		},
		{
			BiasRemove,
			"object7",
			-1 * time.Minute,
			func(l *LWWSet, obj interface{}) bool { return l.Contains(obj) },
		},
	} {
		// Create a LWW Set by explicitly setting a bias.
		lww, err := NewLWWSetWithBias(tt.bias)
		if err != nil {
			t.Fatalf("Unexpected error creating lwwset: %s", err)
		}

		// Mock the time so we can time travel forward and back.
		mock := clock.NewMock()
		lww.clock = mock

		// Create an object that will be removed:
		//  a. right the very moment it is added
		//  b. in future
		//  c. in past
		lww.Add(tt.testObject)

		// This will be our time travel tuner for now.
		mock.Add(tt.elapsed)

		lww.Remove(tt.testObject)

		// Verify that the object is correctly present or absent from the LWW set.
		if !tt.testFn(lww, tt.testObject) {
			t.Errorf("value: '%v' in in invalid state in the set when bias: %q", tt.testObject, tt.bias)
		}
	}
}

func TestLWWESetMerge(t *testing.T) {
	type addRm struct {
		op string
		d  time.Duration
	}

	var addOp, rmOp string = "add", "remove"

	for _, tt := range []struct {
		mapOne, mapTwo map[string]addRm
		valid, invalid map[string]struct{}
	}{
		{
			map[string]addRm{
				"object1": addRm{addOp, 1 * time.Minute},
				"object2": addRm{addOp, 2 * time.Minute},
			},
			map[string]addRm{
				"object1": addRm{rmOp, 2 * time.Minute},
				"object2": addRm{rmOp, 2 * time.Minute},
			},
			map[string]struct{}{
				"object2": struct{}{},
			},
			map[string]struct{}{
				"object1": struct{}{},
			},
		},
		{
			map[string]addRm{
				"object1": addRm{addOp, 1 * time.Minute},
				"object2": addRm{rmOp, 2 * time.Minute},
			},
			map[string]addRm{
				"object3": addRm{addOp, 1 * time.Minute},
				"object4": addRm{rmOp, 2 * time.Minute},
			},
			map[string]struct{}{
				"object1": struct{}{},
				"object3": struct{}{},
			},
			map[string]struct{}{
				"object2": struct{}{},
				"object4": struct{}{},
			},
		},
		{
			map[string]addRm{
				"object1": addRm{addOp, 1 * time.Minute},
				"object2": addRm{addOp, 3 * time.Minute},
			},
			map[string]addRm{
				"object1": addRm{addOp, 2 * time.Minute},
				"object2": addRm{addOp, 2 * time.Minute},
			},
			map[string]struct{}{
				"object1": struct{}{},
				"object2": struct{}{},
			},
			map[string]struct{}{},
		},
		{
			map[string]addRm{
				"object1": addRm{rmOp, 1 * time.Minute},
				"object2": addRm{rmOp, 3 * time.Minute},
			},
			map[string]addRm{
				"object1": addRm{rmOp, 2 * time.Minute},
				"object2": addRm{rmOp, 2 * time.Minute},
			},
			map[string]struct{}{},
			map[string]struct{}{
				"object1": struct{}{},
				"object2": struct{}{},
			},
		},
	} {
		mock1, mock2 := clock.NewMock(), clock.NewMock()

		lww1, err := NewLWWSet()
		if err != nil {
			t.Fatalf("unable to initialize lww set: %s", err)
		}
		lww1.clock = mock1

		lww2, err := NewLWWSet()
		if err != nil {
			t.Fatalf("unable to initialize lww set: %s", err)
		}
		lww2.clock = mock2

		var totalDuration time.Duration

		for obj, addrm := range tt.mapOne {
			curTime := addrm.d - totalDuration

			totalDuration += curTime
			mock1.Add(curTime)

			switch addrm.op {
			case addOp:
				lww1.Add(obj)
			case rmOp:
				lww1.Remove(obj)
			}
		}

		totalDuration = 0 * time.Second

		for obj, addrm := range tt.mapTwo {
			curTime := addrm.d - totalDuration

			totalDuration += curTime
			mock2.Add(curTime)

			switch addrm.op {
			case addOp:
				lww2.Add(obj)
			case rmOp:
				lww2.Remove(obj)
			}
		}

		lww1.Merge(lww2)

		for obj := range tt.valid {
			if !lww1.Contains(obj) {
				t.Errorf("expected merged set to contain: %q", obj)
			}
		}

		for obj := range tt.invalid {
			if lww1.Contains(obj) {
				t.Errorf("expected merged set to not contain: %q", obj)
			}
		}
	}
}

func TestLWWESetMarshalJSON(t *testing.T) {
	for _, tt := range []struct {
		add, rm  map[interface{}]time.Duration
		bias     BiasType
		expected string
	}{
		{
			map[interface{}]time.Duration{},
			map[interface{}]time.Duration{},
			BiasAdd,
			`{"type":"lww-e-set","bias":"a","e":[]}`,
		},
		{
			map[interface{}]time.Duration{},
			map[interface{}]time.Duration{},
			BiasRemove,
			`{"type":"lww-e-set","bias":"r","e":[]}`,
		},
		{
			map[interface{}]time.Duration{
				"object1": 1 * time.Second,
			},
			map[interface{}]time.Duration{
				"object2": 1 * time.Second,
			},
			BiasAdd,
			`{"type":"lww-e-set","bias":"a","e":[{"el":"object1","ta":1},{"el":"object2","td":1}]}`,
		},
		{
			map[interface{}]time.Duration{
				"object1": 1 * time.Second,
			},
			map[interface{}]time.Duration{
				"object2": 1 * time.Second,
			},
			BiasRemove,
			`{"type":"lww-e-set","bias":"r","e":[{"el":"object1","ta":1},{"el":"object2","td":1}]}`,
		},
		{
			map[interface{}]time.Duration{
				"object1": 1 * time.Second,
			},
			map[interface{}]time.Duration{},
			BiasAdd,
			`{"type":"lww-e-set","bias":"a","e":[{"el":"object1","ta":1}]}`,
		},
		{
			map[interface{}]time.Duration{
				"object1": 1 * time.Second,
			},
			map[interface{}]time.Duration{},
			BiasRemove,
			`{"type":"lww-e-set","bias":"r","e":[{"el":"object1","ta":1}]}`,
		},
		{
			map[interface{}]time.Duration{},
			map[interface{}]time.Duration{
				"object1": 1 * time.Second,
			},
			BiasAdd,
			`{"type":"lww-e-set","bias":"a","e":[{"el":"object1","td":1}]}`,
		},
		{
			map[interface{}]time.Duration{},
			map[interface{}]time.Duration{
				"object1": 1 * time.Second,
			},
			BiasRemove,
			`{"type":"lww-e-set","bias":"r","e":[{"el":"object1","td":1}]}`,
		},
		{
			map[interface{}]time.Duration{
				"object1": 1 * time.Second,
			},
			map[interface{}]time.Duration{
				"object1": 2 * time.Second,
			},
			BiasAdd,
			`{"type":"lww-e-set","bias":"a","e":[{"el":"object1","ta":1,"td":2}]}`,
		},
		{
			map[interface{}]time.Duration{
				"object1": 1 * time.Second,
			},
			map[interface{}]time.Duration{
				"object1": 2 * time.Second,
			},
			BiasRemove,
			`{"type":"lww-e-set","bias":"r","e":[{"el":"object1","ta":1,"td":2}]}`,
		},
	} {
		lww, err := NewLWWSetWithBias(tt.bias)
		if err != nil {
			t.Fatalf("unexpected error creating new LWW-E-Set: %s", err)
		}

		mock := clock.NewMock()
		lww.clock = mock

		for e, d := range tt.add {
			mock.Add(d)
			lww.Add(e)

			mock = clock.NewMock()
			lww.clock = mock
		}

		for e, d := range tt.rm {
			mock.Add(d)
			lww.Remove(e)

			mock = clock.NewMock()
			lww.clock = mock
		}

		out, err := json.Marshal(lww)
		if err != nil {
			t.Fatalf("unexpected error on marshalling: %s", err)
		}

		if !bytes.Equal([]byte(tt.expected), out) {
			t.Errorf("expected marshalling array to be: %q, actual: %q", tt.expected, out)
		}
	}
}
