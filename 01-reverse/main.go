package main

import (
	"fmt"
	util "golangdsa/01-reverse/util"
)

func main() {
	input := "Hello, 世界"
	result := util.ReverseString(input)
	fmt.Println("Input:", input)
	fmt.Println("Output:", result)
}
