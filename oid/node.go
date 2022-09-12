package oid

// !!! Status: in progress (draft, not tested, to be continued)

import "sort"

type Node struct {
	Subid    uint32
	Parent   *Node
	Children []*Node
	Data     interface{}
}

// searchShild searches for child with the specified subid in the
// sorted slice of node's children. If the child exists, it returns
// its index and ok = true. Otherwise, it returns ok = false and
// the index to insert a child.
func (node *Node) searchChild(subid uint32) (ind int, ok bool) {
	ind = sort.Search(len(node.Children), func(i int) bool {
		return node.Children[i].Subid >= subid
	})
	ok = ind < len(node.Children) && node.Children[ind].Subid == subid
	return
}

// addChild adds new child with the specified subid to the sorted
// slice of node's children. If the child already exists, it
// returns nil. Otherwise, it returns the added child.
func (node *Node) addChild(subid uint32) *Node {
	if ind, ok := node.searchChild(subid); !ok {
		child := &Node{Subid: subid, Parent: node}
		if ind == len(node.Children) {
			node.Children = append(node.Children, child)
		} else {
			node.Children = append(node.Children[:ind+1], node.Children[ind:]...)
			node.Children[ind] = child
		}
		return child
	}
	return nil
}

// getChild returns node's child with the specified subid. If the
// child does not exist, it returns nil.
func (node *Node) getChild(subid uint32) *Node {
	if ind, ok := node.searchChild(subid); ok {
		return node.Children[ind]
	}
	return nil
}

// removeChild removes child with the specified subid from the
// slice of node's children. It returns false if there is no
// such a child, and true otherwise.
func (node *Node) removeChild(subid uint32) bool {
	if ind, ok := node.searchChild(subid); ok {
		node.Children = append(node.Children[:ind], node.Children[ind+1:]...)
		return true
	}
	return false
}

// addNode
func (node *Node) addNode(oid []uint32) *Node {
	i := 0
	for ; i < len(oid); i++ {
		if ind, ok := node.searchChild(oid[i]); ok {
			node = node.Children[ind]
			continue
		}
		break
	}
	if i == len(oid) {
		return nil
	}
	for _, subid := range oid[i:] {
		node = node.addChild(subid)
	}
	return node
}

// getNode
func (node *Node) getNode(oid []uint32) *Node {
	for _, subid := range oid {
		if node = node.getChild(subid); node == nil {
			break
		}
	}
	return node
}

// prune
func (node *Node) prune() {
	for len(node.Children) == 0 && node.Parent != nil {
		node.Parent.removeChild(node.Subid)
		node = node.Parent
	}
}

// removeNode
func (node *Node) removeNode(oid []uint32, prune bool) bool {
	if node = node.getNode(oid); node == nil || node.Parent == nil {
		return false
	}
	node.Parent.removeChild(node.Subid)
	if prune {
		node.Parent.prune()
	}
	return true
}

// getNextNode returns repeatition'th leaf node with oid such that
// (see specification of SNMP GetNext and GetBulk commands):
//		oid1 <= oid         if include1 == true  && oid2 == nil
//		oid1 <  oid         if include1 == false && oid2 == nil
//		oid1 <= oid < oid2  if include1 == true  && oid2 != nil
//		oid1 <  oid < oid2  if include1 == false && oid2 != nil
// If there is no such node or repeatition < 1, it returns nil.
// Note that empty oid < not empty oid, so in case oid1 == nil,
// it returns leaf node with the least oid.
func (node *Node) getNextNode(oid1, oid2 []uint32, include1 bool, repeatition int16) (resnode *Node) {
	if repeatition < 1 {
		return
	}

	var iterate func(*Node, int, bool, bool)
	iterate = func(node *Node, ind int, eq1, eq2 bool) {
		if len(node.Children) == 0 {
			if (!eq1 || include1) && !eq2 {
				if repeatition == 1 {
					resnode = node
				} else {
					repeatition--
				}
			}
			return
		}
		for _, child := range node.Children {
			if eq1 && ind < len(oid1) && oid1[ind] > child.Subid {
				continue
			}
			if eq2 && (ind >= len(oid2) || oid2[ind] < child.Subid) {
				continue
			}
			if ind >= len(oid1) {
				if ind >= len(oid2) {
					iterate(child, ind, false, false)
				} else {
					iterate(child, ind+1, false, eq2 && child.Subid == oid2[ind])
				}
			} else {
				if ind >= len(oid2) {
					iterate(child, ind+1, eq1 && child.Subid == oid1[ind], false)
				} else {
					iterate(child, ind+1, eq1 && child.Subid == oid1[ind], eq2 && child.Subid == oid2[ind])
				}
			}
			if resnode != nil {
				return
			}
		}
	}

	iterate(node, 0, true, len(oid2) > 0)
	return
}

func (node *Node) forRangeLeaves(do func([]uint32, *Node)) {
	var iterate func([]uint32, []*Node)
	iterate = func(parentOid []uint32, nodes []*Node) {
		for _, node := range nodes {
			nodeOid := Cat(parentOid, node.Subid)
			if len(node.Children) == 0 {
				do(nodeOid, node)
			} else {
				iterate(nodeOid, node.Children)
			}
		}
	}
	iterate(nil, node.Children)
}
