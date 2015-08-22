package crdt

import "testing"

func TestPNCounter(t *testing.T) {
	for _, tt := range []struct {
		incOne, decOne int
		incTwo, decTwo int

		result int
	}{
		{5, 5, 6, 6, 0},
		{5, 6, 7, 8, -2},
		{8, 7, 6, 5, 2},
		{5, 0, 6, 0, 11},
		{0, 5, 0, 6, -11},
	} {
		pOne, pTwo := NewPNCounter(), NewPNCounter()

		for i := 0; i < tt.incOne; i++ {
			pOne.Inc()
		}

		for i := 0; i < tt.decOne; i++ {
			pOne.Dec()
		}

		for i := 0; i < tt.incTwo; i++ {
			pTwo.Inc()
		}

		for i := 0; i < tt.decTwo; i++ {
			pTwo.Dec()
		}

		pOne.Merge(pTwo)

		if pOne.Count() != tt.result {
			t.Errorf("expected the total count to be: %d, actual: %d",
				tt.result,
				pOne.Count())
		}

		pTwo.Merge(pOne)

		if pTwo.Count() != tt.result {
			t.Errorf("expected the total count to be: %d, actual: %d",
				tt.result,
				pTwo.Count())
		}
	}
}

func TestPNCounterInvalidP(t *testing.T) {
	pn := NewPNCounter()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("panic expected here")
		}
	}()

	pn.IncVal(-5)
}

func TestPNCounterInvalidN(t *testing.T) {
	pn := NewPNCounter()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("panic expected here")
		}
	}()

	pn.DecVal(-5)
}
