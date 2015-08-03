package crdt

import "testing"

func TestORSetAddContains(t *testing.T) {
	orSet := NewORSet()

	var testValue string = "test-str"
	orSet.Add(testValue)

	if !orSet.Contains(testValue) {
		t.Errorf("Expected set to contain: %v, but not found", testValue)
	}
}

func TestORSetAddRemoveContains(t *testing.T) {
	orSet := NewORSet()

	var testValue string = "test-str"
	orSet.Add(testValue)

	orSet.Remove(testValue)

	if orSet.Contains(testValue) {
		t.Errorf("Expected set to not contain: %v, but found", testValue)
	}
}

func TestORSetAddRemoveAddContains(t *testing.T) {
	orSet := NewORSet()

	var testValue string = "test-str"

	orSet.Add(testValue)
	orSet.Remove(testValue)
	orSet.Add(testValue)

	if !orSet.Contains(testValue) {
		t.Errorf("Expected set to contain: %v, but not found", testValue)
	}
}

func TestORSetAddAddRemoveContains(t *testing.T) {
	orSet := NewORSet()

	var testValue string = "test-str"

	orSet.Add(testValue)
	orSet.Add(testValue)
	orSet.Remove(testValue)

	if orSet.Contains(testValue) {
		t.Errorf("Expected set to not contain: %v, but found", testValue)
	}
}
