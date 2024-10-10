package block

type Block struct {
	Hash     []byte
	PrevHash []byte
	Coinbase string
	Nonce    int
}

func NewGenesisBlock() *Block {
	return Mining([]byte{}, "I am Genesis Block")
}
