package crdt

type mainSet map[interface{}]struct{}

type GSet struct {
	mainSet mainSet
}

func NewGSet() *GSet {
	return &GSet{
		mainSet: mainSet{},
	}
}

func (s *GSet) Add(elem interface{}) {
	s.mainSet[elem] = struct{}{}
}

func (s *GSet) Contains(elem interface{}) bool {
	_, ok := s.mainSet[elem]
	return ok
}

var (
	_ Set = &GSet{}
)
