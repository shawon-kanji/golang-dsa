package main

import (
	"fmt"
	"time"
)

// NOTE:
// This folder is split across multiple .go files.
// To run the demo, run the whole package (not just this single file):
//   - from repo root: `go run ./20-db`
//   - or from this folder: `go run .`

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

	// Build an SSTable snapshot of current data
	fmt.Println("-------------- Building SSTable ------------")
	sst, err := BuildSSTable(db, "kdb.sst")
	if err != nil {
		fmt.Printf("SSTable build error: %v\n", err)
	} else {
		defer sst.Close()
		fmt.Println("SSTable built at kdb.sst")
	}

	fmt.Println("-------------- Fetching data ------------")
	getStart := time.Now()
	if v, ok := db.Get("user1"); ok {
		fmt.Printf("user1: %s\n", v)
	}
	if v, ok := db.Get("user50"); ok {
		fmt.Printf("user50: %s\n", v)
	}
	fmt.Printf("Get time elapsed: %v\n", time.Since(getStart))

	// Read back via SSTable (if built)
	if err == nil && sst != nil {
		fmt.Println("-------------- SSTable lookups ------------")
		if v, ok := sst.Get("user1"); ok {
			fmt.Printf("sst user1: %s\n", v)
		}
		if v, ok := sst.Get("user50"); ok {
			fmt.Printf("sst user50: %s\n", v)
		}
		if _, ok := sst.Get("missing"); !ok {
			fmt.Println("sst missing: not found")
		}
	}
}
