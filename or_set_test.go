package crdt

import "testing"

func TestORSetAddContains(t *testing.T) {
	orSet := NewORSet()

	var testValue string = "object"

	if orSet.Contains(testValue) {
		t.Errorf("Expected set to not contain: %v, but found", testValue)
	}

	orSet.Add(testValue)

	if !orSet.Contains(testValue) {
		t.Errorf("Expected set to contain: %v, but not found", testValue)
	}
}

func TestORSetAddRemoveContains(t *testing.T) {
	orSet := NewORSet()

	var testValue string = "object"
	orSet.Add(testValue)

	orSet.Remove(testValue)

	if orSet.Contains(testValue) {
		t.Errorf("Expected set to not contain: %v, but found", testValue)
	}
}

func TestORSetAddRemoveAddContains(t *testing.T) {
	orSet := NewORSet()

	var testValue string = "object"

	orSet.Add(testValue)
	orSet.Remove(testValue)
	orSet.Add(testValue)

	if !orSet.Contains(testValue) {
		t.Errorf("Expected set to contain: %v, but not found", testValue)
	}
}

func TestORSetAddAddRemoveContains(t *testing.T) {
	orSet := NewORSet()

	var testValue string = "object"

	orSet.Add(testValue)
	orSet.Add(testValue)
	orSet.Remove(testValue)

	if orSet.Contains(testValue) {
		t.Errorf("Expected set to not contain: %v, but found", testValue)
	}
}

func TestORSetMerge(t *testing.T) {
	type addRm struct {
		addSet []string
		rmSet  []string
	}

	for _, tt := range []struct {
		setOne  addRm
		setTwo  addRm
		valid   map[string]struct{}
		invalid map[string]struct{}
	}{
		{
			addRm{[]string{"object1"}, []string{}},
			addRm{[]string{}, []string{"object1"}},
			map[string]struct{}{
				"object1": struct{}{},
			},
			map[string]struct{}{},
		},
		{
			addRm{[]string{}, []string{"object1"}},
			addRm{[]string{"object1"}, []string{}},
			map[string]struct{}{
				"object1": struct{}{},
			},
			map[string]struct{}{},
		},
		{
			addRm{[]string{"object1"}, []string{"object1"}},
			addRm{[]string{}, []string{}},
			map[string]struct{}{},
			map[string]struct{}{
				"object1": struct{}{},
			},
		},
		{
			addRm{[]string{}, []string{}},
			addRm{[]string{"object1"}, []string{"object1"}},
			map[string]struct{}{},
			map[string]struct{}{
				"object1": struct{}{},
			},
		},
		{
			addRm{[]string{"object2"}, []string{"object1"}},
			addRm{[]string{"object1"}, []string{"object2"}},
			map[string]struct{}{
				"object1": struct{}{},
				"object2": struct{}{},
			},
			map[string]struct{}{},
		},
		{
			addRm{[]string{"object2", "object1"}, []string{"object1"}},
			addRm{[]string{"object1", "object2"}, []string{"object2"}},
			map[string]struct{}{
				"object1": struct{}{},
				"object2": struct{}{},
			},
			map[string]struct{}{},
		},
		{
			addRm{[]string{"object2", "object1"}, []string{"object1", "object2"}},
			addRm{[]string{"object1", "object2"}, []string{"object2", "object1"}},
			map[string]struct{}{},
			map[string]struct{}{
				"object1": struct{}{},
				"object2": struct{}{},
			},
		},
	} {
		orset1, orset2 := NewORSet(), NewORSet()

		for _, add := range tt.setOne.addSet {
			orset1.Add(add)
		}

		for _, rm := range tt.setOne.rmSet {
			orset1.Remove(rm)
		}

		for _, add := range tt.setTwo.addSet {
			orset2.Add(add)
		}

		for _, rm := range tt.setTwo.rmSet {
			orset2.Remove(rm)
		}

		orset1.Merge(orset2)

		for obj, _ := range tt.valid {
			if !orset1.Contains(obj) {
				t.Errorf("expected set to contain: %v", obj)
			}
		}

		for obj, _ := range tt.invalid {
			if orset1.Contains(obj) {
				t.Errorf("expected set to not contain: %v", obj)
			}
		}
	}
}
