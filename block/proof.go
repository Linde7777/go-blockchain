package block

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

const difficulty = 12

func Mining(prevHash []byte, data string) (block *Block) {
	one := big.NewInt(1)
	target := one.Lsh(one, uint(256-difficulty))

	nonce := 0
	var hashInBytes [32]byte
	fmt.Println("mining...")
	for nonce < math.MaxInt64 {
		hashInBytes = sha256.Sum256(combineData(prevHash, data, nonce, difficulty))
		fmt.Printf("trying nonce: %d, mined hash:\r%x\n", nonce, hashInBytes)

		var hashInInt big.Int
		hashInInt.SetBytes(hashInBytes[:])
		if hashInInt.Cmp(target) != -1 {
			nonce += 1
		} else {
			break
		}
	}

	return &Block{
		Hash:     hashInBytes[:],
		PrevHash: prevHash,
		Data:     data,
		Nonce:    nonce,
	}

}

func combineData(prevHash []byte, data string, nonce, difficulty int) []byte {
	return bytes.Join(
		[][]byte{
			prevHash,
			[]byte(data),
			int2Bytes(nonce),
			int2Bytes(difficulty),
		},
		[]byte{},
	)
}

func int2Bytes(num int) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
