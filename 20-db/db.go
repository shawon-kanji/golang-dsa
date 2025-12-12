package main

import "fmt"

type KDB struct {
	head *Node
	Size int
	wal  *WAL
}

// newKDB creates an in-memory DB without WAL.
func newKDB() *KDB {
	return &KDB{}
}

// newKDBWithWAL creates a DB with WAL and replays existing WAL entries.
func newKDBWithWAL(walPath string) (*KDB, error) {
	db := &KDB{}

	wal, err := NewWAL(walPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create WAL: %w", err)
	}
	db.wal = wal

	if err := db.recoverFromWAL(); err != nil {
		_ = wal.Close()
		return nil, fmt.Errorf("failed to recover from WAL: %w", err)
	}
	return db, nil
}

func (db *KDB) Close() error {
	if db.wal != nil {
		return db.wal.Close()
	}
	return nil
}

func (db *KDB) Checkpoint() error {
	if db.wal != nil {
		return db.wal.Truncate()
	}
	return nil
}

func (db *KDB) checkIsDuplicateKey(key string) bool {
	return db.searchKey(db.head, key) != nil
}

func (db *KDB) searchKey(node *Node, key string) *Node {
	hash := HashStringToInt(key)
	if node == nil {
		return nil
	}

	if node.size >= 1 && node.key1 == hash {
		return node
	}
	if node.size >= 2 && node.key2 == hash {
		return node
	}
	if node.size >= 3 && node.key3 == hash {
		return node
	}

	if node.size >= 1 && hash < node.key1 {
		return db.searchKey(node.cp1, key)
	}
	if node.size >= 2 && hash < node.key2 {
		return db.searchKey(node.cp2, key)
	}
	if node.size >= 3 && hash < node.key3 {
		return db.searchKey(node.cp3, key)
	}
	return db.searchKey(node.cp4, key)
}

func (db *KDB) Put(key string, val string) (bool, string) {
	hashkey := HashStringToInt(key)

	if db.checkIsDuplicateKey(key) {
		return false, val
	}

	if db.wal != nil {
		if err := db.wal.Log("PUT", key, val); err != nil {
			return false, val
		}
	}

	db.insert(hashkey, &val)
	return true, val
}

func (db *KDB) Get(key string) (string, bool) {
	hash := HashStringToInt(key)
	node := db.head

	for node != nil {
		if node.size >= 1 && node.key1 == hash {
			if node.rp1 == nil {
				return "", false
			}
			return *node.rp1, true
		}
		if node.size >= 2 && node.key2 == hash {
			if node.rp2 == nil {
				return "", false
			}
			return *node.rp2, true
		}
		if node.size >= 3 && node.key3 == hash {
			if node.rp3 == nil {
				return "", false
			}
			return *node.rp3, true
		}

		if hash < node.key1 {
			node = node.cp1
		} else if node.size == 1 || hash < node.key2 {
			node = node.cp2
		} else if node.size == 2 || hash < node.key3 {
			node = node.cp3
		} else {
			node = node.cp4
		}
	}

	return "", false
}

func (db *KDB) recoverFromWAL() error {
	if db.wal == nil {
		return nil
	}

	ops, err := db.wal.Recover()
	if err != nil {
		return err
	}

	// Disable WAL during replay to avoid re-logging.
	wal := db.wal
	db.wal = nil
	defer func() { db.wal = wal }()

	for _, op := range ops {
		switch op.Operation {
		case "PUT":
			db.Put(op.Key, op.Value)
		}
	}
	return nil
}
