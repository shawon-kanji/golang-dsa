package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// SSTable is a simple immutable sorted-string table stored on disk.
// Keys are hashed ints (using HashStringToInt) and values are plain strings.
type SSTable struct {
	path  string
	file  *os.File
	index map[int]int64 // key hash -> byte offset of line start
}

// ForEachInOrder traverses the B-tree in sorted order and calls fn for each key/value.
func (db *KDB) ForEachInOrder(fn func(key int, val string)) {
	var walk func(n *Node)
	walk = func(n *Node) {
		if n == nil {
			return
		}
		// cp1, key1
		walk(n.cp1)
		if n.rp1 != nil {
			fn(n.key1, *n.rp1)
		}
		// cp2, key2
		walk(n.cp2)
		if n.size >= 2 && n.rp2 != nil {
			fn(n.key2, *n.rp2)
		}
		// cp3, key3
		walk(n.cp3)
		if n.size >= 3 && n.rp3 != nil {
			fn(n.key3, *n.rp3)
		}
		// cp4
		walk(n.cp4)
	}
	walk(db.head)
}

// BuildSSTable creates an SSTable file from the current DB contents.
// The file is rewritten (truncated) each time.
func BuildSSTable(db *KDB, path string) (*SSTable, error) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("open sstable: %w", err)
	}
	w := bufio.NewWriter(f)
	index := make(map[int]int64)
	var pos int64

	db.ForEachInOrder(func(key int, val string) {
		line := fmt.Sprintf("%d\t%s\n", key, val)
		index[key] = pos
		_, _ = w.WriteString(line)
		pos += int64(len(line))
	})

	if err := w.Flush(); err != nil {
		f.Close()
		return nil, fmt.Errorf("flush sstable: %w", err)
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return nil, fmt.Errorf("sync sstable: %w", err)
	}

	return &SSTable{path: path, file: f, index: index}, nil
}

// LoadSSTable opens an existing SSTable and builds its in-memory index.
func LoadSSTable(path string) (*SSTable, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open sstable: %w", err)
	}

	index := make(map[int]int64)
	scanner := bufio.NewScanner(f)
	var offset int64
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) != 2 {
			offset += int64(len(line) + 1)
			continue
		}
		k, err := strconv.Atoi(parts[0])
		if err == nil {
			index[k] = offset
		}
		offset += int64(len(line) + 1) // +1 for newline
	}
	if err := scanner.Err(); err != nil {
		f.Close()
		return nil, fmt.Errorf("scan sstable: %w", err)
	}

	// Reset file pointer for reads
	if _, err := f.Seek(0, 0); err != nil {
		f.Close()
		return nil, fmt.Errorf("seek sstable: %w", err)
	}

	return &SSTable{path: path, file: f, index: index}, nil
}

// Get looks up a string key (hashed) in the SSTable.
func (s *SSTable) Get(key string) (string, bool) {
	h := HashStringToInt(key)
	off, ok := s.index[h]
	if !ok {
		return "", false
	}
	if _, err := s.file.Seek(off, 0); err != nil {
		return "", false
	}
	reader := bufio.NewReader(s.file)
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", false
	}
	line = strings.TrimRight(line, "\n")
	parts := strings.SplitN(line, "\t", 2)
	if len(parts) != 2 {
		return "", false
	}
	return parts[1], true
}

// Close releases the SSTable file handle.
func (s *SSTable) Close() error {
	if s.file != nil {
		return s.file.Close()
	}
	return nil
}
