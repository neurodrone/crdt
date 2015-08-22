package crdt

type PNCounter struct {
	pCounter *GCounter
	nCounter *GCounter
}

func NewPNCounter() *PNCounter {
	return &PNCounter{
		pCounter: NewGCounter(),
		nCounter: NewGCounter(),
	}
}

func (pn *PNCounter) Inc() {
	pn.IncVal(1)
}

func (pn *PNCounter) IncVal(incr int) {
	pn.pCounter.IncVal(incr)
}

func (pn *PNCounter) Dec() {
	pn.DecVal(1)
}

func (pn *PNCounter) DecVal(decr int) {
	pn.nCounter.IncVal(decr)
}

func (pn *PNCounter) Count() int {
	return pn.pCounter.Count() - pn.nCounter.Count()
}

func (pn *PNCounter) Merge(pnpn *PNCounter) {
	pn.pCounter.Merge(pnpn.pCounter)
	pn.nCounter.Merge(pnpn.nCounter)
}
