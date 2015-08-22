package crdt

import "testing"

func TestGCounter(t *testing.T) {
	for _, tt := range []struct {
		incsOne int
		incsTwo int
		result  int
	}{
		{5, 10, 15},
		{10, 5, 15},
		{100, 100, 200},
		{1, 2, 3},
	} {
		gOne, gTwo := NewGCounter(), NewGCounter()

		for i := 0; i < tt.incsOne; i++ {
			gOne.Inc()
		}

		for i := 0; i < tt.incsTwo; i++ {
			gTwo.Inc()
		}

		gOne.Merge(gTwo)

		if gOne.Count() != tt.result {
			t.Errorf("expected total count to be: %d, actual: %d",
				tt.result,
				gOne.Count())
		}

		gTwo.Merge(gOne)

		if gTwo.Count() != tt.result {
			t.Errorf("expected total count to be: %d, actual: %d",
				tt.result,
				gTwo.Count())
		}
	}
}

func TestGCounterInvalidInput(t *testing.T) {
	gc := NewGCounter()

	defer func() {
		if r := recover(); r == nil {
			t.Fatalf("panic expected here")
		}
	}()

	gc.IncVal(-5)
}
