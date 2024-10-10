package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"math"
	"math/big"
)

const (
	difficulty      = 12
	targetLength    = 256
	hashBytesLength = 32
)

var target = big.NewInt(1).Lsh(big.NewInt(1), uint(targetLength-difficulty))

func Mining(prevHash string, coinbase string) (block *Block) {
	nonce := 0
	var hashInBytes [hashBytesLength]byte
	// log.Println("-------------")
	// log.Println("mining...")
	for nonce < math.MaxInt64 {
		hashInBytes = sha256.Sum256(combineData(prevHash, coinbase, nonce, difficulty))
		// log.Printf("trying nonce: %d, mined hash:\r%x\n", nonce, hashInBytes)

		var hashInInt big.Int
		hashInInt.SetBytes(hashInBytes[:])
		if hashInInt.Cmp(target) != -1 {
			nonce += 1
		} else {
			// log.Println("nonce found!!!")
			break
		}
	}

	return &Block{
		Hash:     string(hashInBytes[:]),
		PrevHash: prevHash,
		Coinbase: coinbase,
		Nonce:    nonce,
	}

}

func Validate(block *Block) bool {
	hashInBytes := sha256.Sum256(combineData(block.PrevHash, block.Coinbase, block.Nonce, difficulty))
	var hashInInt big.Int
	hashInInt.SetBytes(hashInBytes[:])
	return hashInInt.Cmp(target) == -1
}

func combineData(prevHash string, coinbase string, nonce, difficulty int) []byte {
	return bytes.Join(
		[][]byte{
			[]byte(prevHash),
			[]byte(coinbase),
			int2Bytes(nonce),
			int2Bytes(difficulty),
		},
		[]byte{},
	)
}

func int2Bytes(num int) []byte {
	newNum := int64(num)
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, newNum)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
