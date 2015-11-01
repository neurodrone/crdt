package crdt

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/benbjohnson/clock"
)

type LWWSet struct {
	addMap map[interface{}]time.Time
	rmMap  map[interface{}]time.Time

	bias BiasType

	clock clock.Clock
}

type BiasType string

const (
	BiasAdd    BiasType = "a"
	BiasRemove BiasType = "r"
)

var (
	ErrNoSuchBias = errors.New("no such bias found")
)

func NewLWWSet() (*LWWSet, error) {
	return NewLWWSetWithBias(BiasAdd)
}

func NewLWWSetWithBias(bias BiasType) (*LWWSet, error) {
	if bias != BiasAdd && bias != BiasRemove {
		return nil, ErrNoSuchBias
	}

	return &LWWSet{
		addMap: make(map[interface{}]time.Time),
		rmMap:  make(map[interface{}]time.Time),
		bias:   bias,
		clock:  clock.New(),
	}, nil
}

func (s *LWWSet) Add(value interface{}) {
	s.addMap[value] = s.clock.Now()
}

func (s *LWWSet) Remove(value interface{}) {
	s.rmMap[value] = s.clock.Now()
}

func (s *LWWSet) Contains(value interface{}) bool {
	addTime, addOk := s.addMap[value]

	// If a value is not present in added set then
	// always return false, irrespective of whether
	// it is present in the removed set.
	if !addOk {
		return false
	}

	rmTime, rmOk := s.rmMap[value]

	// If a value is present in added set and not in remove
	// we should always return true.
	if !rmOk {
		return true
	}

	switch s.bias {
	case BiasAdd:
		return !addTime.Before(rmTime)

	case BiasRemove:
		return rmTime.Before(addTime)
	}

	// This case will almost always never be hit. Usually
	// if an invalid Bias value is provided, it is called
	// at a higher level.
	return false
}

func (s *LWWSet) Merge(r *LWWSet) {
	for value, ts := range r.addMap {
		if t, ok := s.addMap[value]; ok && t.Before(ts) {
			s.addMap[value] = ts
		} else {
			if t.Before(ts) {
				s.addMap[value] = ts
			} else {
				s.addMap[value] = t
			}
		}
	}

	for value, ts := range r.rmMap {
		if t, ok := s.rmMap[value]; ok && t.Before(ts) {
			s.rmMap[value] = ts
		} else {
			if t.Before(ts) {
				s.rmMap[value] = ts
			} else {
				s.rmMap[value] = t
			}
		}
	}
}

type lwwesetJSON struct {
	T string   `json:"type"`
	B string   `json:"bias"`
	E []elJSON `json:"e"`
}

type elJSON struct {
	Elem interface{} `json:"el"`
	TAdd int64       `json:"ta,omitempty"`
	TDel int64       `json:"td,omitempty"`
}

func (s *LWWSet) MarshalJSON() ([]byte, error) {
	l := &lwwesetJSON{
		T: "lww-e-set",
		B: string(s.bias),
		E: make([]elJSON, 0, len(s.addMap)),
	}

	for e, t := range s.addMap {
		el := elJSON{Elem: e, TAdd: t.Unix()}
		if td, ok := s.rmMap[e]; ok {
			el.TDel = td.Unix()
		}

		l.E = append(l.E, el)
	}

	for e, t := range s.rmMap {
		if _, ok := s.addMap[e]; ok {
			continue
		}

		l.E = append(l.E, elJSON{Elem: e, TDel: t.Unix()})
	}

	return json.Marshal(l)
}
