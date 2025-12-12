package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// WALOperation represents a single operation in the WAL.
// Each entry is stored as one JSON line.
type WALOperation struct {
	Seq       int64  `json:"seq"`
	Operation string `json:"op"` // "PUT" or "DELETE"
	Key       string `json:"key"`
	Value     string `json:"value"`
	Timestamp int64  `json:"timestamp"`
}

// WAL is a simple write-ahead log.
type WAL struct {
	file     *os.File
	mu       sync.Mutex
	seq      int64
	filePath string
}

func NewWAL(filePath string) (*WAL, error) {
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open WAL file: %w", err)
	}

	wal := &WAL{file: file, filePath: filePath}
	wal.seq, _ = wal.getLastSequence()
	return wal, nil
}

func (w *WAL) getLastSequence() (int64, error) {
	file, err := os.Open(w.filePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	var lastSeq int64
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var op WALOperation
		if err := json.Unmarshal(scanner.Bytes(), &op); err == nil {
			if op.Seq > lastSeq {
				lastSeq = op.Seq
			}
		}
	}

	return lastSeq, nil
}

func (w *WAL) Log(operation, key, value string) error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.seq++
	op := WALOperation{
		Seq:       w.seq,
		Operation: operation,
		Key:       key,
		Value:     value,
		Timestamp: time.Now().UnixNano(),
	}

	data, err := json.Marshal(op)
	if err != nil {
		return fmt.Errorf("failed to marshal WAL operation: %w", err)
	}

	if _, err := w.file.Write(append(data, '\n')); err != nil {
		return fmt.Errorf("failed to write to WAL: %w", err)
	}
	if err := w.file.Sync(); err != nil {
		return fmt.Errorf("failed to sync WAL: %w", err)
	}
	return nil
}

func (w *WAL) Recover() ([]WALOperation, error) {
	file, err := os.Open(w.filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to open WAL for recovery: %w", err)
	}
	defer file.Close()

	var operations []WALOperation
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var op WALOperation
		if err := json.Unmarshal(scanner.Bytes(), &op); err != nil {
			continue
		}
		operations = append(operations, op)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading WAL: %w", err)
	}
	return operations, nil
}

func (w *WAL) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		return w.file.Close()
	}
	return nil
}

func (w *WAL) Truncate() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	if w.file != nil {
		if err := w.file.Close(); err != nil {
			return err
		}
	}

	file, err := os.OpenFile(w.filePath, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = file
	w.seq = 0
	return nil
}
