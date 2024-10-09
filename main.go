package main

import (
	"go-blockchain/blockchain"
	"io"
	"log"
	"os"
)

func main() {
	logFile, err := os.OpenFile("blockchain.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer logFile.Close()
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

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
