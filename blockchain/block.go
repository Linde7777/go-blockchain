package blockchain

type Block struct {
	Hash         string
	PrevHash     string
	Transactions []*Transaction
	Nonce        int
}

func NewGenesisBlock(txs []*Transaction) *Block {
	return Mining("", txs)
}
