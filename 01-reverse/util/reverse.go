package util

import "fmt"

func ReverseString(inputString string) string {
	runeSlices := []rune(inputString)
	left := 0
	right := len(runeSlices) - 1
	for left < right {
		replace := runeSlices[left]
		runeSlices[left] = runeSlices[right]
		runeSlices[right] = replace
		left++
		right--
	}

	outputString := string(runeSlices)
	fmt.Println("reverse : ", outputString)
	fmt.Println(1 << 5)
	return outputString

}
