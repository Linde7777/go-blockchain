package blockchain

import "go-blockchain/block"

type BlockChain struct {
	Blocks []*block.Block
}

func NewBlockChain() *BlockChain {
	return &BlockChain{
		Blocks: []*block.Block{block.NewGenesisBlock()},
	}
}

func (chain *BlockChain) AddBlock(data string) {
	lastBlock := chain.Blocks[len(chain.Blocks)-1]
	newBlock := block.Mining(lastBlock.Hash, data)
	chain.Blocks = append(chain.Blocks, newBlock)
}
