package block

type Block struct {
	Hash     []byte
	PrevHash []byte
	Data     string
	Nonce    int
}
