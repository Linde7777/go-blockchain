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
	"strconv"
	"strings"
)

func main() {
	setupLog()
	defer closeLog()

	LoadConfig()

	db := NewBlockChainStorage(AppConfig.StorageOption)

	chain, err := blockchain.NewBlockChain(db, AppConfig.UserAddress)
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
			newTransaction(chain, reader)
		case "2":
			showEntireBlockchain(chain)
		case "3":
			getBalance(chain, reader)
		case "4":
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
	optionBadgerDB = "badgerdb"
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

func NewBadgerDB() *badger.DB {
	badgerDB, err := badger.Open(badger.DefaultOptions(AppConfig.BadgerDBPath))
	if err != nil {
		log.Fatal(err)
	}
	return badgerDB
}

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: AppConfig.RedisAddr,
	})
	return client
}

func printMenu() {
	fmt.Println("\n--- Blockchain Menu ---")
	fmt.Println("1. Add a new transaction")
	fmt.Println("2. Show entire blockchain")
	fmt.Println("3. Get balance")
	fmt.Println("4. Exit")
	fmt.Print("Enter your choice: ")
}

func newTransaction(chain *blockchain.BlockChain, reader *bufio.Reader) {
	fmt.Print("Enter sender: ")
	sender, _ := reader.ReadString('\n')
	sender = strings.TrimSpace(sender)

	fmt.Print("Enter recipient: ")
	recipient, _ := reader.ReadString('\n')
	recipient = strings.TrimSpace(recipient)

	fmt.Print("Enter amount: ")
	amount, _ := reader.ReadString('\n')
	amount = strings.TrimSpace(amount)

	amountInt, err := strconv.Atoi(amount)
	if err != nil {
		fmt.Println("Error converting amount to integer:", err)
		return
	}

	transaction, err := blockchain.NewTransaction(sender, recipient, amountInt, chain)
	if err != nil {
		fmt.Println("Error creating transaction:", err)
		return
	}

	err = chain.AddBlock([]*blockchain.Transaction{transaction})
	if err != nil {
		fmt.Println("Error adding block:", err)
		return
	}

	fmt.Println("Transaction added successfully!")
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
		fmt.Printf("Transactions: %s\n", block.Transactions)
		fmt.Printf("Nonce: %d\n", block.Nonce)
		fmt.Printf("Proof of Work: %s\n", strconv.FormatBool(block.Validate()))
	}
}

func getBalance(chain *blockchain.BlockChain, reader *bufio.Reader) {
	fmt.Print("Enter address: ")
	address, _ := reader.ReadString('\n')
	address = strings.TrimSpace(address)

	balance, err := chain.GetBalance(address)
	if err != nil {
		fmt.Println("Error getting balance:", err)
		return
	}

	fmt.Printf("Balance of %s: %d coins\n", address, balance)
}
