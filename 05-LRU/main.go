package main

import "fmt"

type User struct {
	name string
	age  int
}

type Node[T any] struct {
	key   string
	value T
	next  *Node[T]
	prev  *Node[T]
}

type LRU[T any] struct {
	maxSize     int
	currentSize int
	head        *Node[T]
	tail        *Node[T]
	addressMap  map[string]*Node[T]
}

func newLRU[T any](size int) LRU[T] {
	return LRU[T]{
		maxSize:     size,
		currentSize: 0,
		addressMap:  make(map[string]*Node[T]),
	}
}

func (lru *LRU[T]) Put(key string, data T) {
	// Check if key already exists, update value and move to front
	if existingNode, ok := lru.addressMap[key]; ok {
		existingNode.value = data
		lru.moveToFront(existingNode)
		return
	}

	newNode := &Node[T]{
		key:   key,
		value: data,
	}

	if lru.head == nil {
		lru.head = newNode
		lru.tail = newNode
		lru.currentSize++
		lru.addressMap[key] = newNode
	} else if lru.currentSize == lru.maxSize {
		fmt.Println("LRU max capacity.. replacing node")
		delete(lru.addressMap, lru.tail.key)
		lru.tail = lru.tail.prev
		if lru.tail != nil {
			lru.tail.next = nil
		}
		newNode.next = lru.head
		lru.head.prev = newNode
		lru.head = newNode
		lru.addressMap[key] = newNode
	} else {
		newNode.next = lru.head
		lru.head.prev = newNode
		lru.head = newNode
		lru.currentSize++
		lru.addressMap[key] = newNode
	}
}

func (lru *LRU[T]) Get(key string) (T, bool) {
	if node, ok := lru.addressMap[key]; ok {
		lru.moveToFront(node)
		return node.value, true
	}
	var zero T
	return zero, false
}

func (lru *LRU[T]) moveToFront(node *Node[T]) {
	if node == lru.head {
		return
	}

	// Remove node from current position
	if node.prev != nil {
		node.prev.next = node.next
	}
	if node.next != nil {
		node.next.prev = node.prev
	}

	// Update tail if needed
	if node == lru.tail {
		lru.tail = node.prev
	}

	// Move to front
	node.prev = nil
	node.next = lru.head
	lru.head.prev = node
	lru.head = node
}

func main() {
	cache := newLRU[int](4)

	fmt.Println(cache, cache.maxSize, cache.currentSize)

	cache.Put("a", 11)
	cache.Put("b", 22)
	cache.Put("c", 33)
	cache.Put("d", 44)
	cache.Put("e", 55)
	cache.Put("f", 66)
	cache.Put("g", 77)
	cache.Put("h", 5555)

	fmt.Println("\n--- Int Cache ---")
	for p := cache.head; p != nil; p = p.next {
		fmt.Printf("Key: %s, Value: %d\n", p.key, p.value)
	}

	// Test Get method
	if val, ok := cache.Get("h"); ok {
		fmt.Printf("\nGet 'h': %d\n", val)
	}
	if val, ok := cache.Get("z"); ok {
		fmt.Printf("Get 'z': %d\n", val)
	} else {
		fmt.Println("Key 'z' not found")
	}

	fmt.Println("\n--- After Get 'h' (should be at front) ---")
	for p := cache.head; p != nil; p = p.next {
		fmt.Printf("Key: %s, Value: %d\n", p.key, p.value)
	}

	userCache := newLRU[User](4)

	userCache.Put("user1", User{
		name: "Alice Johnson",
		age:  28,
	})
	userCache.Put("user2", User{
		name: "Bob Smith",
		age:  35,
	})
	userCache.Put("user3", User{
		name: "Charlie Brown",
		age:  42,
	})
	userCache.Put("user4", User{
		name: "Diana Prince",
		age:  30,
	})
	userCache.Put("user5", User{
		name: "Edward Norton",
		age:  45,
	})
	userCache.Put("user6", User{
		name: "Fiona Davis",
		age:  27,
	})
	userCache.Put("user7", User{
		name: "George Miller",
		age:  33,
	})

	fmt.Println("\n--- User Cache ---")
	for p := userCache.head; p != nil; p = p.next {
		fmt.Printf("Key: %s, Value: %+v\n", p.key, p.value)
	}

	getUser, _ := userCache.Get("user5")
	fmt.Println(getUser)

	fmt.Println("\n--- User Cache after access ---")
	for p := userCache.head; p != nil; p = p.next {
		fmt.Printf("Key: %s, Value: %+v\n", p.key, p.value)
	}

}
