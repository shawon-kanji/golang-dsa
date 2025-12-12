package main

func HashStringToInt(key string) int {
	hash := 5381
	for _, char := range key {
		hash = ((hash << 5) + hash) + int(char) // hash * 33 + char
	}
	return hash
}
