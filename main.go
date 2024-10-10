package main

import (
	"bufio"
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/go-redis/redis"
	"go-blockchain/blockchain"
	"go-blockchain/utils"
	"log"
	"os"
	"strings"
)

func main() {
	setupLog()
	defer closeLog()

	db := NewBlockChainStorage(optionRedis)

	chain, err := blockchain.NewBlockChain(db)
	if err != nil {
		utils.LogPanic(err)
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		printMenu()
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			addBlock(chain, reader)
		case "2":
			showEntireBlockchain(chain)
		case "3":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid option. Please try again.")
		}
	}
}

const logFilename = "blockchain.log"

func setupLog() {
	logFile, err := os.OpenFile(logFilename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}
	log.SetOutput(logFile)
}

func closeLog() {
	if f, ok := log.Writer().(*os.File); ok {
		err := f.Close()
		if err != nil {
			log.Println("Failed to close log file:", err)
		}
	}
}

const (
	optionBadgerDB = "badgerDB"
	optionRedis    = "redis"
)

func NewBlockChainStorage(option string) blockchain.Storage {
	switch option {
	case optionBadgerDB:
		badgerDB := NewBadgerDB()
		return blockchain.NewBadgerDBStorage(badgerDB)
	case optionRedis:
		redisDB := NewRedisClient()
		return blockchain.NewRedisStorage(redisDB)
	default:
		log.Fatal("Invalid storage option")
	}
	return nil
}

const badgerDBPath = "/tmp/badger"

func NewBadgerDB() *badger.DB {
	badgerDB, err := badger.Open(badger.DefaultOptions(badgerDBPath))
	if err != nil {
		log.Fatal(err)
	}
	return badgerDB
}

const redisAddr = "localhost:6379"

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return client
}

func printMenu() {
	fmt.Println("\n--- Blockchain Menu ---")
	fmt.Println("1. Add a new block")
	fmt.Println("2. Show entire blockchain")
	fmt.Println("3. Exit")
	fmt.Print("Enter your choice: ")
}

func addBlock(chain *blockchain.BlockChain, reader *bufio.Reader) {
	fmt.Print("Enter data for the new block: ")
	data, _ := reader.ReadString('\n')
	data = strings.TrimSpace(data)

	err := chain.AddBlock(data)
	if err != nil {
		fmt.Println("Error adding block:", err)
	} else {
		fmt.Println("Block added successfully!")
	}
}

func showEntireBlockchain(chain *blockchain.BlockChain) {
	lastBlock, err := chain.GetLastBlock()
	if err != nil {
		fmt.Println("Error getting last block:", err)
		return
	}

	iter := blockchain.NewIterator(chain.DB, lastBlock.Hash)

	for iter.HasNext() {
		block, err := iter.Next()
		if err != nil {
			fmt.Println("Error getting next block:", err)
			return
		}
		fmt.Println("------------------------")
		fmt.Printf("Hash: %x\n", block.Hash)
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Data: %s\n", block.Coinbase)
		fmt.Printf("Nonce: %d\n", block.Nonce)
	}
}
