package block

type Block struct {
	Hash     []byte
	PrevHash []byte
	Data     string
	Nonce    int
}

func NewGenesisBlock() *Block {
	return &Block{
		Hash:     []byte{},
		PrevHash: []byte{},
		Data:     "I am Genesis Block",
		Nonce:    0,
	}
}
