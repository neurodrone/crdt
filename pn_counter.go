package crdt

// PNCounter represents a state-based PN-Counter. It is
// implemented as sets of two G-Counters, one that tracks
// increments while the other decrements.
type PNCounter struct {
	pCounter *GCounter
	nCounter *GCounter
}

// NewPNCounter returns a new *PNCounter with both its
// G-Counters initialized.
func NewPNCounter() *PNCounter {
	return &PNCounter{
		pCounter: NewGCounter(),
		nCounter: NewGCounter(),
	}
}

// Inc monotonically increments the current value of the
// PN-Counter by one.
func (pn *PNCounter) Inc() {
	pn.IncVal(1)
}

// IncVal increments the current value of the PN-Counter
// by the delta incr that is provided. The value of delta
// has to be >= 0. If the value of delta is < 0, then this
// implementation panics.
func (pn *PNCounter) IncVal(incr int) {
	pn.pCounter.IncVal(incr)
}

// Dec monotonically decrements the current value of the
// PN-Counter by one.
func (pn *PNCounter) Dec() {
	pn.DecVal(1)
}

// DecVal adds a decrement to the current value of
// PN-Counter by the value of delta decr. Similar to
// IncVal, the value of decr cannot be less than 0.
func (pn *PNCounter) DecVal(decr int) {
	pn.nCounter.IncVal(decr)
}

// Count returns the current value of the counter. It
// subtracts the value of negative G-Counter from the
// positive grow-only counter and returns the result.
// Because this counter can grow in either direction,
// negative integers as results are possible.
func (pn *PNCounter) Count() int {
	return pn.pCounter.Count() - pn.nCounter.Count()
}

// Merge combines both the PN-Counters and saves the result
// in the invoking counter. Respective G-Counters are merged
// i.e. +ve with +ve and -ve with -ve, but not computation
// is actually performed.
func (pn *PNCounter) Merge(pnpn *PNCounter) {
	pn.pCounter.Merge(pnpn.pCounter)
	pn.nCounter.Merge(pnpn.nCounter)
}
