package crdt

import "github.com/satori/go.uuid"

// GCounter represent a G-counter in CRDT, which is
// a state-based grow-only counter that only supports
// increments.
type GCounter struct {
	// ident provides a unique identity to each replica.
	ident string

	// counter maps identity of each replica to their
	// entry values i.e. the counter value they individually
	// have.
	counter map[string]int
}

// NewGCounter returns a *GCounter by pre-assigning a unique
// identity to it.
func NewGCounter() *GCounter {
	return &GCounter{
		ident:   uuid.NewV4().String(),
		counter: make(map[string]int),
	}
}

// Inc increments the GCounter by the value of 1 everytime it
// is called.
func (g *GCounter) Inc() {
	g.IncVal(1)
}

// IncVal allows passing in an arbitrary delta to increment the
// current value of counter by. Only positive values are accepted.
// If a negative value is provided the implementation will panic.
func (g *GCounter) IncVal(incr int) {
	if incr < 0 {
		panic("cannot decrement a gcounter")
	}
	g.counter[g.ident] += incr
}

// Count returns the total count of this counter across all the
// present replicas.
func (g *GCounter) Count() (total int) {
	for _, val := range g.counter {
		total += val
	}
	return
}

// Merge combines the counter values across multiple replicas.
// The property of idempotency is preserved here across
// multiple merges as when no state is changed across any replicas,
// the result should be exactly the same everytime.
func (g *GCounter) Merge(c *GCounter) {
	for ident, val := range c.counter {
		if v, ok := g.counter[ident]; !ok || v < val {
			g.counter[ident] = val
		}
	}
}
