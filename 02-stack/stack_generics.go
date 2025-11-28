package main

import (
	"errors"
	"fmt"
)

type Stack[T any] struct {
	dataArray []T
}

func (s *Stack[T]) Push(item T) {
	s.dataArray = append(s.dataArray, item)
}

func (s *Stack[T]) Pop() (T, error) {
	var zero T
	if len(s.dataArray) > 0 {
		item := s.dataArray[len(s.dataArray)-1]
		s.dataArray = s.dataArray[:len(s.dataArray)-1]
		return item, nil
	} else {
		fmt.Println("No Item in stack")
		return zero, errors.New("Empty Stack")
	}
}

func (s *Stack[T]) Peek() (T, error) {
	var zero T
	if len(s.dataArray) > 0 {
		item := s.dataArray[len(s.dataArray)-1]
		return item, nil
	} else {
		fmt.Println("No Item in stack")
		return zero, errors.New("Empty Stack")
	}
}

type User struct {
	Name    string
	Age     int8
	Country string
}

func main() {
	var newStack Stack[int]
	newStack.Push(100)
	newStack.Push(140)
	newStack.Push(10540)
	var val, _ = newStack.Pop()
	fmt.Println(val)
	val, _ = newStack.Pop()
	fmt.Println(val)
	val, _ = newStack.Pop()
	fmt.Println(val)
	val, _ = newStack.Pop()
	fmt.Println(val)
	val, _ = newStack.Pop()
	fmt.Println(val)
	val, _ = newStack.Pop()
	fmt.Println(val)

	var userStack Stack[User]
	userStack.Push(User{Age: 20, Name: "kanji", Country: "Malaysia"})
	userStack.Push(User{Name: "Shawon", Age: 33, Country: "India"})
	user, _ := userStack.Pop()
	fmt.Println(user)
	user, _ = userStack.Pop()
	fmt.Println(user)
	user, _ = userStack.Pop()
	fmt.Println(user)
	user, _ = userStack.Pop()
	fmt.Println(user)
	user, _ = userStack.Pop()
	fmt.Println(user)

	userP := new(User)

	userP.Name = "pointer"
	userP.Age = 50
	userP.Country = "binary"

	fmt.Println(&userP)

}
