package main

import "fmt"

func (db *KDB) PrintTree() {
	fmt.Println("-------------- Tree Structure ------------")
	if db.head == nil {
		fmt.Println("(empty tree)")
		return
	}
	db.printNode(db.head, 0, "ROOT")
}

func (db *KDB) printNode(node *Node, level int, position string) {
	if node == nil {
		return
	}

	indent := ""
	for i := 0; i < level; i++ {
		indent += "    "
	}

	fmt.Printf("%s[%s] Node (size=%d):\n", indent, position, node.size)

	if node.size >= 1 {
		val := "<nil>"
		if node.rp1 != nil {
			val = *node.rp1
		}
		fmt.Printf("%s  key1: %d -> \"%s\"\n", indent, node.key1, val)
	}
	if node.size >= 2 {
		val := "<nil>"
		if node.rp2 != nil {
			val = *node.rp2
		}
		fmt.Printf("%s  key2: %d -> \"%s\"\n", indent, node.key2, val)
	}
	if node.size >= 3 {
		val := "<nil>"
		if node.rp3 != nil {
			val = *node.rp3
		}
		fmt.Printf("%s  key3: %d -> \"%s\"\n", indent, node.key3, val)
	}

	if node.cp1 != nil {
		db.printNode(node.cp1, level+1, "CP1")
	}
	if node.cp2 != nil {
		db.printNode(node.cp2, level+1, "CP2")
	}
	if node.cp3 != nil {
		db.printNode(node.cp3, level+1, "CP3")
	}
	if node.cp4 != nil {
		db.printNode(node.cp4, level+1, "CP4")
	}
}
