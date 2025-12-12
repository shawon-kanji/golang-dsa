package main

import (
	"fmt"
	"time"
)

func main() {
	walPath := "kdb.wal"
	db, err := newKDBWithWAL(walPath)
	if err != nil {
		fmt.Printf("Failed to create database: %v\n", err)
		return
	}
	defer db.Close()

	if db.Size == 0 {
		fmt.Println("-------------- Inserting 50 records ------------")
		insertStart := time.Now()

		for i := 1; i <= 50; i++ {
			key := fmt.Sprintf("user%d", i)
			val := fmt.Sprintf("value-%d", i)
			_, _ = db.Put(key, val)
		}

		fmt.Printf("Insert time elapsed: %v\n", time.Since(insertStart))
	} else {
		fmt.Printf("Database recovered with %d entries from WAL\n", db.Size)
	}

	fmt.Println("-------------- Tree Stats ------------")
	fmt.Printf("Total entries: %d\n", db.Size)

	db.PrintTree()

	fmt.Println("-------------- Fetching data ------------")
	getStart := time.Now()
	if v, ok := db.Get("user1"); ok {
		fmt.Printf("user1: %s\n", v)
	}
	if v, ok := db.Get("user50"); ok {
		fmt.Printf("user50: %s\n", v)
	}
	fmt.Printf("Get time elapsed: %v\n", time.Since(getStart))
}
