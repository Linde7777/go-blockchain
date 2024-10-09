package blockchain

import "go-blockchain/block"

type BlockChain struct {
	blocks []*block.Block
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		blocks: []*block.Block{block.NewGenesisBlock()},
	}
}

func (chain *BlockChain) AddBlock(data string) {
	lastBlock := chain.blocks[len(chain.blocks)-1]
	newBlock := block.Mining(lastBlock.Hash, data)
	chain.blocks = append(chain.blocks, newBlock)
}
