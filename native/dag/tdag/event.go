package tdag

import (
	"github.com/unicornultrafoundation/go-helios/hash"
	"github.com/unicornultrafoundation/go-helios/native/dag"
)

type TestEvent struct {
	dag.MutableBaseEvent
	Name string
}

func (e *TestEvent) AddParent(id hash.Event) {
	parents := e.Parents()
	parents.Add(id)
	e.SetParents(parents)
}
