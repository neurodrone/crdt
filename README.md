# CRDT [![Build Status](https://travis-ci.org/neurodrone/crdt.svg?branch=master)](https://travis-ci.org/neurodrone/crdt) [![Coverage Status](https://coveralls.io/repos/neurodrone/crdt/badge.svg?branch=master&service=github)](https://coveralls.io/github/neurodrone/crdt?branch=master) [![GoDoc](https://godoc.org/github.com/neurodrone/crdt?status.svg)](https://godoc.org/github.com/neurodrone/crdt) [![](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/neurodrone/crdt/blob/master/LICENSE) [![Report card](http://goreportcard.com/badge/neurodrone/crdt)](http://goreportcard.com/report/neurodrone/crdt)

This is an implementation of [Convergent and Commutative Replicated Data Types](https://hal.inria.fr/inria-00555588/document) in [Go](https://golang.org/).

The following state-based counters and sets have currently been
implemented.

## Counters

### G-Counter

A grow-only counter (G-Counter) can only increase in one direction. The increment
operation increases the value of current replica by 1. The merge
operation combines values from distinct replicas by taking the maximum
of each replica.

```go
gcounter := crdt.NewGCounter()

// We can increase the counter monotonically by 1.
gcounter.Inc()

// Twice.
gcounter.Inc()

// Or we can pass in an arbitrary delta to apply as an increment.
gcounter.IncVal(2)

// Should print '4' as the result.
fmt.Println(gcounter.Count())
```

### PN-Counter

A positive-negative counter (PN-Counter) is a CRDT that can both increase or
decrease and converge correctly in the light of commutative
operations. Both `.Inc()` and `.Dec()` operations are allowed and thus
negative values are possible as a result.

```go
pncounter := crdt.NewPNCounter()

// We can increase the counter by 1.
pncounter.Inc()

// Or more.
pncounter.Inc(100)

// And similarly decrease its value by 1.
pncounter.Dec()

// Or more.
pncounter.DecVal(100)

// End result should equal '0' here.
fmt.Println(pncounter.Count())
```

## Sets

### G-Set

A grow-only (G-Set) set to which element/s can be added to. Removing element
from the set is not possible.

```go
obj := "dummy-object"
gset := crdt.NewGSet()

gset.Add(obj)

// Should always print 'true' as `obj` exists in the g-set.
fmt.Println(gset.Contains(obj))
```

### 2P-Set

Two-phase set (2P-Set) allows both additions and removals to the set.
Internally it comprises of two G-Sets, one to keep track of additions
and the other for removals.

```go
obj := "dummy-object"

ppset := crdt.NewTwoPhaseSet()

ppset.Add(obj)

// Remove the object that we just added to the set, emptying it.
ppset.Remove(obj)

// Should return 'false' as the obj doesn't exist within the set.
fmt.Println(ppset.Contains(obj))
```

### LWW-e-Set

Last-write-wins element set (LWW-e-Set) keeps track of element additions
and removals but with respect to the timestamp that is attached to each
element. Timestamps should be unique and have ordering properties.

```go
obj := "dummy-object"
lwwset := crdt.NewLWWSet()

// Here, we remove the object first before we add it in. For a
// 2P-set the object would be deemed absent from the set. But for
// a LWW-set the object should be present because `.Add()` follows
// `.Remove()`.
lwwset.Remove(obj); lwwset.Add(obj)

// This should print 'true' because of the above.
fmt.Println(lwwset.Contains(obj))
```

### OR-Set

An OR-Set (Observed-Removed-Set) allows deletion and addition of
elements similar to LWW-e-Set, but does not surface only the most recent one. Additions are uniquely tracked via tags and an element is considered member of the set if the deleted set consists of all the tags present within additions.

```go
// object 1 == object 2
obj1, obj2 := "dummy-object", "dummy-object"

orset := crdt.NewORSet()

orset.Add(obj1); orset.Add(obj2)

// Removing any one of the above two objects should remove both
// because they contain the same value.
orset.Remove(obj1)

// Should return 'false'.
fmt.Println(orset.Contains(obj2))
```
