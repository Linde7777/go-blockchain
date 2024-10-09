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

func Mining(prevHash []byte, coinbase string) (block *Block) {
	one := big.NewInt(1)
	target := one.Lsh(one, uint(targetLength-difficulty))

	nonce := 0
	var hashInBytes [hashBytesLength]byte
	log.Println("-------------")
	log.Println("mining...")
	for nonce < math.MaxInt64 {
		hashInBytes = sha256.Sum256(combineData(prevHash, coinbase, nonce, difficulty))
		log.Printf("trying nonce: %d, mined hash:\r%x\n", nonce, hashInBytes)

		var hashInInt big.Int
		hashInInt.SetBytes(hashInBytes[:])
		if hashInInt.Cmp(target) != -1 {
			nonce += 1
		} else {
			log.Println("nonce found!!!")
			break
		}
	}

	return &Block{
		Hash:     hashInBytes[:],
		PrevHash: prevHash,
		Coinbase: coinbase,
		Nonce:    nonce,
	}

}

func combineData(prevHash []byte, coinbase string, nonce, difficulty int) []byte {
	return bytes.Join(
		[][]byte{
			prevHash,
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
