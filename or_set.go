package crdt

import "github.com/satori/go.uuid"

type ORSet struct {
	addMap map[interface{}]map[string]struct{}
	rmMap  map[interface{}]map[string]struct{}
}

func NewORSet() *ORSet {
	return &ORSet{
		addMap: make(map[interface{}]map[string]struct{}),
		rmMap:  make(map[interface{}]map[string]struct{}),
	}
}

func (o *ORSet) Add(value interface{}) {
	if m, ok := o.addMap[value]; ok {
		m[uuid.NewV4().String()] = struct{}{}
		o.addMap[value] = m
		return
	}

	m := make(map[string]struct{})

	m[uuid.NewV4().String()] = struct{}{}
	o.addMap[value] = m
}

func (o *ORSet) Remove(value interface{}) {
	r, ok := o.rmMap[value]
	if !ok {
		r = make(map[string]struct{})
	}

	if m, ok := o.addMap[value]; ok {
		for uid, _ := range m {
			r[uid] = struct{}{}
		}
	}

	o.rmMap[value] = r
}

func (o *ORSet) Contains(value interface{}) bool {
	addMap, ok := o.addMap[value]
	if !ok {
		return false
	}

	rmMap, ok := o.rmMap[value]
	if !ok {
		return true
	}

	for uid, _ := range addMap {
		if _, ok := rmMap[uid]; !ok {
			return true
		}
	}

	return false
}

func (o *ORSet) Merge(r *ORSet) {
	for value, m := range r.addMap {
		addMap, ok := o.addMap[value]
		if ok {
			for uid, _ := range m {
				addMap[uid] = struct{}{}
			}

			continue
		}

		o.addMap[value] = m
	}

	for value, m := range r.rmMap {
		rmMap, ok := o.rmMap[value]
		if ok {
			for uid, _ := range m {
				rmMap[uid] = struct{}{}
			}

			continue
		}

		o.rmMap[value] = m
	}
}
