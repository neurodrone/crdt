package crdt

import "github.com/satori/go.uuid"

type GCounter struct {
	ident   string
	counter map[string]int
}

func NewGCounter() *GCounter {
	return &GCounter{
		ident:   uuid.NewV4().String(),
		counter: make(map[string]int),
	}
}

func (g *GCounter) Inc() {
	g.IncVal(1)
}

func (g *GCounter) IncVal(incr int) {
	if incr < 0 {
		panic("cannot decrement a gcounter")
	}
	g.counter[g.ident] += incr
}

func (g *GCounter) Count() (total int) {
	for _, val := range g.counter {
		total += val
	}
	return
}

func (g *GCounter) Merge(c *GCounter) {
	for ident, val := range c.counter {
		if v, ok := g.counter[ident]; !ok || v < val {
			g.counter[ident] = val
		}
	}
}
