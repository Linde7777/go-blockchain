package main

import (
	"fmt"
	"github.com/dgraph-io/badger"
	"github.com/go-redis/redis"
	blockPkg "go-blockchain/block"
	"go-blockchain/blockchain"
	"go-blockchain/utils"
	"log"
	"os"
	"strconv"
	"time"
)

func main() {
	setupLog()
	defer closeLog()

	db := NewStorage(optionRedis)

	chain, err := blockchain.NewBlockChain(db)
	if err != nil {
		log.Println("Failed to create blockchain:", err)
		return
	}

	for i := 0; i < 6; i++ {
		err = chain.AddBlock(fmt.Sprintf("blocked added at %s", time.Now()))
		if err != nil {
			log.Println("Failed to add block:", err)
			return
		}
	}

	lastBlock, err := chain.GetLastBlock()
	if err != nil {
		utils.LogPanic(err)
	}
	iter := blockchain.NewIterator(db, lastBlock.Hash)

	for iter.HasNext() {
		b, err := iter.Next()
		if err != nil {
			log.Println("Failed to get next block:", err)
			return
		}
		log.Println("---------------------")
		log.Printf("block hash: %x\n", b.Hash)
		log.Printf("previous hash: %x\n", b.PrevHash)
		log.Printf("nonce: %d\n", b.Nonce)
		log.Printf("data: %s\n", b.Coinbase)
		log.Printf("proof of work: %s\n", strconv.FormatBool(blockPkg.Validate(b)))
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

func NewStorage(option string) blockchain.Storage {
	switch option {
	case optionBadgerDB:
		badgerDB := NewBadgerDB()
		return blockchain.NewBadgerDBStorage(badgerDB)
	case optionRedis:
		redisDB := NewRedisDB()
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

func NewRedisDB() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	return client
}
