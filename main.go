package main

import (
	"go-blockchain/blockchain"
	"log"
)

func main() {
	chain := blockchain.NewBlockChain()
	chain.AddBlock("I am block0")
	chain.AddBlock("I am block1")
	chain.AddBlock("I am block2")
	chain.AddBlock("I am block3")

	for _, block := range chain.Blocks {
		log.Println("---------------------")
		log.Printf("block hash: %x\n", block.Hash)
		log.Printf("previous hash: %x\n", block.PrevHash)
		log.Printf("nonce: %d\n", block.Nonce)
		log.Printf("data: %s\n", block.Data)
	}
}
