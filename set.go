package crdt

type Set interface {
	Add(interface{})
	Contains(interface{}) bool
}
