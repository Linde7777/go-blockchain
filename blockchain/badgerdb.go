package blockchain

import (
	"errors"
	"github.com/dgraph-io/badger"
	"go-blockchain/block"
)

type BadgerDB struct {
	db *badger.DB
}

var _ Storage = &BadgerDB{}

func NewBadgerDBStorage(db *badger.DB) *BadgerDB {
	return &BadgerDB{db: db}
}

func (b *BadgerDB) Get(key string) (string, error) {
	var data []byte
	err := b.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		data, err = item.ValueCopy(nil)
		return err
	})
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (b *BadgerDB) GetBlock(key string) (*block.Block, error) {
	data, err := b.Get(key)
	if err != nil {
		return nil, err
	}
	return str2Block(data)
}

func (b *BadgerDB) Set(key, value string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
}

func (b *BadgerDB) SetBlock(key string, block *block.Block) error {
	serBlock, err := block2Str(block)
	if err != nil {
		return err
	}
	return b.Set(key, serBlock)
}

func (b *BadgerDB) Delete(key string) error {
	return b.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

func (b *BadgerDB) KeyNotFound(err error) bool {
	return errors.Is(err, badger.ErrKeyNotFound)
}
