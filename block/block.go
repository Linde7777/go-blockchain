package block

import (
	"go-blockchain/blockchain"
)

type Block struct {
	Hash         string
	PrevHash     string
	Transactions []*blockchain.Transaction
	Nonce        int
}

func NewGenesisBlock(txs []*blockchain.Transaction) *Block {
	return Mining("", txs)
}
