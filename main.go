package main

import (
	blockPkg "go-blockchain/block"
	"go-blockchain/blockchain"
	"io"
	"log"
	"os"
	"strconv"
)

const logFilename = "blockchain.log"

func main() {
	logFile, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			log.Panic(err)
		}
	}(logFile)
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	chain := blockchain.NewBlockChain()
	chain.AddBlock("I am block0")
	chain.AddBlock("I am block1")
	chain.AddBlock("I am block2")
	chain.AddBlock("I am block3")

	for _, b := range chain.Blocks {
		log.Println("---------------------")
		log.Printf("block hash: %x\n", b.Hash)
		log.Printf("previous hash: %x\n", b.PrevHash)
		log.Printf("nonce: %d\n", b.Nonce)
		log.Printf("data: %s\n", b.Coinbase)
		log.Printf("proof of work: %s\n", strconv.FormatBool(blockPkg.Validate(b)))
	}
}
