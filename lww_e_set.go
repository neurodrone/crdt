package crdt

import (
	"errors"
	"sync"
	"time"
)

type TSSet struct {
	set map[interface{}][]time.Time
	mu  sync.Mutex
}

func NewTSSet() *TSSet {
	return &TSSet{}
}

func (t *TSSet) Add(value interface{}) {
	t.AddTS(value, time.Now())
}

func (t *TSSet) AddTS(value interface{}, ts ...time.Time) {
	t.mu.Lock()

	timestamps, ok := t.set[value]
	if !ok {
		timestamps = []time.Time{}
	}

	timestamps = append(timestamps, ts...)
	t.set[value] = timestamps

	t.mu.Unlock()
}

func (t *TSSet) Lookup(value interface{}) []time.Time {
	if timestamps, ok := t.set[value]; ok {
		return timestamps
	}
	return nil
}

func (t *TSSet) Elements() map[interface{}][]time.Time {
	return t.set
}

type LWWSet struct {
	addSet *TSSet
	rmSet  *TSSet

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
		addSet: NewTSSet(),
		rmSet:  NewTSSet(),
		bias:   bias,
	}, nil
}

func (s *LWWSet) Add(value interface{}) {
	s.addSet.Add(value)
}

func (s *LWWSet) Remove(value interface{}) {
	s.rmSet.Add(value)
}

func (s *LWWSet) Contains(value interface{}) bool {
	addTSs := s.addSet.Lookup(value)
	rmTSs := s.rmSet.Lookup(value)

	var maxAddTS, maxRmTS time.Time

	for _, ts := range addTSs {
		if ts.After(maxAddTS) {
			maxAddTS = ts
		}
	}

	for _, ts := range rmTSs {
		if ts.After(maxRmTS) {
			maxRmTS = ts
		}
	}

	return maxAddTS.After(maxRmTS)
}

func (s *LWWSet) Merge(r *LWWSet) {
	for value, tss := range r.addSet.Elements() {
		s.addSet.AddTS(value, tss...)
	}

	for value, tss := range r.rmSet.Elements() {
		s.rmSet.AddTS(value, tss...)
	}
}
