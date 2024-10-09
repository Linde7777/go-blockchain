package block

type Block struct {
	Hash     []byte
	PrevHash []byte
	Coinbase string
	Nonce    int
}

func NewGenesisBlock() *Block {
	return &Block{
		Hash:     []byte{},
		PrevHash: []byte{},
		Coinbase: "I am Genesis Block",
		Nonce:    0,
	}
}
