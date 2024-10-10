package blockchain

import (
	"go-blockchain/block"
)

type BlockChain struct {
	LastBlockHash string
	DB            Storage
}

const (
	DBKeyLastBlockHash = "last-block-hash"
)

func NewBlockChain(db Storage) (*BlockChain, error) {
	chain := &BlockChain{}

	chain.DB = db

	lastBlock, err := chain.GetLastBlock()
	switch {
	// i.e., a blockchain with existing data is already persistent in the storage
	case err == nil:
		chain.LastBlockHash = lastBlock.Hash

	// i.e., there is no blockchain persistent in the storage
	case db.KeyNotFound(err):
		genesis := block.NewGenesisBlock()
		err = db.SetBlock(genesis.Hash, genesis)
		if err != nil {
			return nil, err
		}

		err = db.Set(DBKeyLastBlockHash, genesis.Hash)
		if err != nil {
			return nil, err
		}
		chain.LastBlockHash = genesis.Hash

	default:
		return nil, err
	}

	return chain, nil
}

func (chain *BlockChain) AddBlock(coinbase string) error {
	lastBlock, err := chain.GetLastBlock()
	if err != nil {
		return err
	}

	newBlock := block.Mining(lastBlock.Hash, coinbase)
	err = chain.DB.SetBlock(newBlock.Hash, newBlock)
	if err != nil {
		return err
	}

	chain.LastBlockHash = newBlock.Hash
	err = chain.DB.Set(DBKeyLastBlockHash, newBlock.Hash)
	if err != nil {
		return err
	}

	return nil
}

func (chain *BlockChain) GetLastBlock() (*block.Block, error) {
	lastBlock, err := chain.DB.Get(DBKeyLastBlockHash)
	if err != nil {
		return nil, err
	}
	return chain.GetBlock(lastBlock)
}

func (chain *BlockChain) GetBlock(blockHash string) (*block.Block, error) {
	return chain.DB.GetBlock(blockHash)
}
