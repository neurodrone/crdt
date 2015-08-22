# CRDT [![Build Status](https://travis-ci.org/neurodrone/crdt.svg?branch=master)](https://travis-ci.org/neurodrone/crdt) [![Coverage Status](https://coveralls.io/repos/neurodrone/crdt/badge.svg?branch=master&service=github)](https://coveralls.io/github/neurodrone/crdt?branch=master) [![GoDoc](https://godoc.org/github.com/neurodrone/crdt?status.svg)](https://godoc.org/github.com/neurodrone/crdt) [![](https://img.shields.io/badge/license-MIT-blue.svg)](https://github.com/neurodrone/crdt/blob/master/LICENSE)

This is an implementation of CRDTs in Golang.

### CRDTs Implemented

 * G-Set
 * 2P-Set
 * LWW-e-Set
 * OR-Set

### CRDTs Remaining

#### Counters

 * G-Counter
 * PN-Counter

#### Sets

 * U-Set
 * MC-Set
 * MVRegister
 * Graphs
 * TreeDoc

### TODO

 * Add a persistence layer.
 * Add separate *counters* and *sets* directory.
