package oid

// !!! Status: in progress (draft, not tested, to be continued)

import "sync"

type Tree struct {
	lock sync.RWMutex
	root *Node
}

func NewTree() *Tree {
	return &Tree{root: &Node{}}
}

func (t *Tree) Clear() {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.root.Children = nil
}

func (t *Tree) AddNode(oid []uint32) *Node {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.root.addNode(oid)
}

func (t *Tree) GetNode(oid []uint32) *Node {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.root.getNode(oid)
}

func (t *Tree) GetNextNode(oid1, oid2 []uint32, include1 bool) *Node {
	t.lock.RLock()
	defer t.lock.RUnlock()
	return t.root.getNextNode(oid1, oid2, include1, 1)
}

func (t *Tree) RemoveNode(oid []uint32) bool {
	t.lock.Lock()
	defer t.lock.Unlock()
	return t.root.removeNode(oid, false)
}

func (t *Tree) ForRangeLeaves(do func([]uint32, *Node)) {
	t.lock.RLock()
	defer t.lock.RUnlock()
	t.root.forRangeLeaves(do)
}
