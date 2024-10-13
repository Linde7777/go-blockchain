package blockchain

import (
	"bytes"
	"encoding/gob"
	"strings"
)

type Storage interface {
	Get(key string) (string, error)
	GetBlock(key string) (*Block, error)
	Set(key, value string) error
	SetBlock(key string, value *Block) error
	Delete(key string) error
	KeyNotFound(err error) bool
}

func block2Str(b *Block) (string, error) {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)
	err := encoder.Encode(b)
	if err != nil {
		return "", err
	}
	return res.String(), nil
}

func str2Block(data string) (*Block, error) {
	var b Block
	decoder := gob.NewDecoder(strings.NewReader(data))
	err := decoder.Decode(&b)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

// DBIterator iterate the db like iterating a
// linked-list(prevHash field in block.Block).
// It can iterate till the Genesis block.
type DBIterator interface {
	HasNext() bool
	Next() (*Block, error)
}

type DBIteratorV1 struct {
	nextHash string
	db       Storage
}

var _ DBIterator = &DBIteratorV1{}

func NewIterator(db Storage, lastBlockHash string) *DBIteratorV1 {
	return &DBIteratorV1{
		nextHash: lastBlockHash,
		db:       db,
	}
}

func (iter *DBIteratorV1) HasNext() bool {
	// when iter.nextHash will be "" ? check the code in Next()
	return iter.nextHash != ""
}

func (iter *DBIteratorV1) Next() (*Block, error) {
	if !iter.HasNext() {
		return nil, nil
	}

	b, err := iter.db.GetBlock(iter.nextHash)
	switch {
	// this case is related to the HasNext()
	case err != nil && iter.db.KeyNotFound(err):
		b.PrevHash = ""
	case err != nil:
		return nil, err
	}

	iter.nextHash = b.PrevHash
	return b, nil
}
