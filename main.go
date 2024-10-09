package main

import (
	"fmt"
	"go-blockchain/blockchain"
)

func main() {
	chain := blockchain.NewBlockChain()
	chain.AddBlock("I am block0")
	chain.AddBlock("I am block1")
	chain.AddBlock("I am block2")
	chain.AddBlock("I am block3")

	for _, block := range chain.Blocks {
		fmt.Println("---------------------")
		fmt.Printf("block hash: %x\n", block.Hash)
		fmt.Printf("previous hash: %x\n", block.PrevHash)
		fmt.Printf("nonce: %d\n", block.Nonce)
		fmt.Printf("data: %s\n", block.Data)
	}
}
