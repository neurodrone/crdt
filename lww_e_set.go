package crdt

import (
	"errors"
	"time"
)

type LWWSet struct {
	addMap map[interface{}]time.Time
	rmMap  map[interface{}]time.Time

	bias BiasType
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
	}, nil
}

func (s *LWWSet) Add(value interface{}) {
	s.addMap[value] = time.Now()
}

func (s *LWWSet) Remove(value interface{}) {
	s.rmMap[value] = time.Now()
}

func (s *LWWSet) Contains(value interface{}) bool {
	addTime, addOk := s.addMap[value]
	if !addOk {
		return false
	}

	rmTime, rmOk := s.rmMap[value]
	if !rmOk {
		return true
	}

	switch s.bias {
	case BiasAdd:
		return addTime.After(rmTime)

	case BiasRemove:
		return rmTime.After(addTime)
	}

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
