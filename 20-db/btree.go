package main

import "sort"

// 4-way search / m=4 (max 3 keys, max 4 children)

type Node struct {
	cp1  *Node
	key1 int
	rp1  *string
	cp2  *Node
	key2 int
	rp2  *string
	cp3  *Node
	key3 int
	rp3  *string
	cp4  *Node
	size int
}

// splitResult holds the result of a node split.
type splitResult struct {
	promoted    bool
	promotedKey int
	promotedVal *string
	newRight    *Node
}

// insert adds a key-value pair into the B-tree.
func (db *KDB) insert(key int, val *string) {
	if db.head == nil {
		db.head = &Node{key1: key, rp1: val, size: 1}
		db.Size++
		return
	}

	result := db.insertKey(db.head, key, val)
	if result.promoted {
		newRoot := &Node{
			key1: result.promotedKey,
			rp1:  result.promotedVal,
			cp1:  db.head,
			cp2:  result.newRight,
			size: 1,
		}
		db.head = newRoot
	}
}

func (db *KDB) insertKeyIntoNode(node *Node, key int, val *string) {
	if node.size == 0 {
		node.key1, node.rp1, node.size = key, val, 1
		return
	}
	if node.size == 1 {
		if key < node.key1 {
			node.key2, node.rp2 = node.key1, node.rp1
			node.key1, node.rp1 = key, val
		} else {
			node.key2, node.rp2 = key, val
		}
		node.size = 2
		return
	}
	if node.size == 2 {
		if key < node.key1 {
			node.key3, node.rp3 = node.key2, node.rp2
			node.key2, node.rp2 = node.key1, node.rp1
			node.key1, node.rp1 = key, val
		} else if key < node.key2 {
			node.key3, node.rp3 = node.key2, node.rp2
			node.key2, node.rp2 = key, val
		} else {
			node.key3, node.rp3 = key, val
		}
		node.size = 3
	}
}

func (db *KDB) isLeaf(node *Node) bool {
	return node.cp1 == nil && node.cp2 == nil && node.cp3 == nil && node.cp4 == nil
}

func (db *KDB) insertKey(node *Node, key int, val *string) splitResult {
	if db.isLeaf(node) {
		if node.size < 3 {
			db.insertKeyIntoNode(node, key, val)
			db.Size++
			return splitResult{}
		}
		db.Size++
		return db.splitNode(node, key, val, nil, nil)
	}

	var child **Node
	var childIdx int
	switch {
	case key < node.key1:
		child, childIdx = &node.cp1, 0
	case node.size == 1 || key < node.key2:
		child, childIdx = &node.cp2, 1
	case node.size == 2 || key < node.key3:
		child, childIdx = &node.cp3, 2
	default:
		child, childIdx = &node.cp4, 3
	}

	if *child == nil {
		*child = &Node{key1: key, rp1: val, size: 1}
		db.Size++
		return splitResult{}
	}

	res := db.insertKey(*child, key, val)
	if !res.promoted {
		return splitResult{}
	}

	if node.size < 3 {
		db.insertPromotedKey(node, res.promotedKey, res.promotedVal, res.newRight, childIdx)
		return splitResult{}
	}
	return db.splitNode(node, res.promotedKey, res.promotedVal, res.newRight, &childIdx)
}

func (db *KDB) insertPromotedKey(node *Node, key int, val *string, newRight *Node, afterChildIdx int) {
	if node.size == 1 {
		if afterChildIdx == 0 {
			node.key2, node.rp2 = node.key1, node.rp1
			node.key1, node.rp1 = key, val
			node.cp3 = node.cp2
			node.cp2 = newRight
		} else {
			node.key2, node.rp2 = key, val
			node.cp3 = newRight
		}
		node.size = 2
		return
	}

	// node.size == 2
	if afterChildIdx == 0 {
		node.key3, node.rp3 = node.key2, node.rp2
		node.key2, node.rp2 = node.key1, node.rp1
		node.key1, node.rp1 = key, val
		node.cp4 = node.cp3
		node.cp3 = node.cp2
		node.cp2 = newRight
	} else if afterChildIdx == 1 {
		node.key3, node.rp3 = node.key2, node.rp2
		node.key2, node.rp2 = key, val
		node.cp4 = node.cp3
		node.cp3 = newRight
	} else {
		node.key3, node.rp3 = key, val
		node.cp4 = newRight
	}
	node.size = 3
}

func (db *KDB) splitNode(node *Node, newKey int, newVal *string, newChild *Node, afterChildIdx *int) splitResult {
	type kv struct {
		k int
		v *string
	}
	items := []kv{{node.key1, node.rp1}, {node.key2, node.rp2}, {node.key3, node.rp3}, {newKey, newVal}}
	sort.Slice(items, func(i, j int) bool { return items[i].k < items[j].k })

	// left keeps items[0]
	node.key1, node.rp1 = items[0].k, items[0].v
	node.key2, node.rp2 = 0, nil
	node.key3, node.rp3 = 0, nil
	node.size = 1

	sibling := &Node{key1: items[2].k, rp1: items[2].v, key2: items[3].k, rp2: items[3].v, size: 2}

	if afterChildIdx != nil {
		children := []*Node{node.cp1, node.cp2, node.cp3, node.cp4, nil}
		idx := *afterChildIdx + 1
		copy(children[idx+1:], children[idx:])
		children[idx] = newChild

		node.cp1, node.cp2, node.cp3, node.cp4 = children[0], children[1], nil, nil
		sibling.cp1, sibling.cp2, sibling.cp3, sibling.cp4 = children[2], children[3], children[4], nil
	} else {
		node.cp1, node.cp2, node.cp3, node.cp4 = nil, nil, nil, nil
	}

	return splitResult{promoted: true, promotedKey: items[1].k, promotedVal: items[1].v, newRight: sibling}
}
