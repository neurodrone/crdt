# CRDT [![Build Status](https://travis-ci.org/neurodrone/crdt.svg?branch=master)](https://travis-ci.org/neurodrone/crdt)

This is an implementation of CRDTs in Golang.

### CRDTs Implemented

 * G-Set
 * 2P-Set
 * LWW-e-Set

### CRDTs Remaining

#### Counters

 * G-Counter
 * PN-Counter

#### Sets

 * U-Set
 * OR-Set
  - AWORSET
  - RWORSET
 * MC-Set
 * MVRegister
 * Graphs
 * TreeDoc

### TODO

 * Add a persistence layer.
 * Add separate *counters* and *sets* directory.
