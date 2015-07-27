package crdt

type TwoPhaseSet struct {
	addSet *GSet
	rmSet  *GSet
}

func NewTwoPhaseSet() *TwoPhaseSet {
	return &TwoPhaseSet{
		addSet: NewGSet(),
		rmSet:  NewGSet(),
	}
}

func (s *TwoPhaseSet) Add(elem interface{}) {
	s.addSet.Add(elem)
}

func (s *TwoPhaseSet) Remove(elem interface{}) {
	s.rmSet.Add(elem)
}

func (s *TwoPhaseSet) Contains(elem interface{}) bool {
	return s.addSet.Contains(elem) && !s.rmSet.Contains(elem)
}

var (
	_ Set = &TwoPhaseSet{}
)
