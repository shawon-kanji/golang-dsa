package main

import "fmt"

type Node struct {
	Value int32
	Next  *Node
}

type LinkedList struct {
	head *Node
	size int
}

func (ll *LinkedList) InsertAtHead(val int32) {
	newNode := &Node{
		Value: val,
		Next:  ll.head,
	}
	ll.head = newNode
	ll.size++
}

func (ll *LinkedList) InsertAtTail(val int32) {
	newNode := &Node{Value: val}

	if ll.head == nil {
		ll.head = newNode
		ll.size++
		return
	}

	current := ll.head
	for current.Next != nil {
		current = current.Next
	}
	current.Next = newNode
	ll.size++
}

func (ll *LinkedList) Delete(val int32) bool {
	if ll.head == nil {
		return false
	}

	if ll.head.Value == val {
		ll.head = ll.head.Next
		ll.size--
		return true
	}

	current := ll.head
	for current.Next != nil {
		if current.Next.Value == val {
			current.Next = current.Next.Next
			ll.size--
			return true
		}
		current = current.Next
	}
	return false
}

func (ll *LinkedList) Find(val int32) *Node {
	current := ll.head
	for current != nil {
		if current.Value == val {
			return current
		}
		current = current.Next
	}
	return nil
}

func (ll *LinkedList) Traverse() {
	current := ll.head
	for current != nil {
		fmt.Printf("%d -> ", current.Value)
		current = current.Next
	}
	fmt.Println("nil")
}

func (ll *LinkedList) Size() int {
	return ll.size
}

func main() {
	ll := LinkedList{}

	ll.InsertAtHead(50)
	ll.InsertAtHead(555)
	ll.InsertAtTail(100)

	ll.Traverse()
	fmt.Printf("Size: %d\n", ll.Size())
}
