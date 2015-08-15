package crdt

import (
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
