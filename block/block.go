package block

import (
	"fmt"
	"time"
)

type Block struct {
	Hash     string
	PrevHash string
	Coinbase string
	Nonce    int
}

func NewGenesisBlock() *Block {
	return Mining("", fmt.Sprintf("I am Genesis Block, created at %s", time.Now()))
}
