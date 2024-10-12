package blockchain

import (
	"errors"
)

type BlockChain struct {
	LastBlockHash string
	DB            Storage
}

const (
	DBKeyLastBlockHash = "last-block-hash"
)

func NewBlockChain(db Storage, yourAddress string) (*BlockChain, error) {
	chain := &BlockChain{}

	chain.DB = db

	lastBlock, err := chain.GetLastBlock()
	switch {
	// i.e., a blockchain with existing data is already persistent in the storage
	case err == nil:
		chain.LastBlockHash = lastBlock.Hash

	// i.e., there is no blockchain persistent in the storage
	case db.KeyNotFound(err):
		tx, err := NewCoinbaseTX(yourAddress, "")
		if err != nil {
			return nil, err
		}
		txs := []*Transaction{tx}
		genesis := NewGenesisBlock(txs)
		err = db.SetBlock(genesis.Hash, genesis)
		if err != nil {
			return nil, err
		}

		err = db.Set(DBKeyLastBlockHash, genesis.Hash)
		if err != nil {
			return nil, err
		}
		chain.LastBlockHash = genesis.Hash

	case err != nil:
		return nil, err
	}

	return chain, nil
}

func (chain *BlockChain) AddBlock(txs []*Transaction) error {
	lastBlock, err := chain.GetLastBlock()
	if err != nil {
		return err
	}

	newBlock := Mining(lastBlock.Hash, txs)
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

func (chain *BlockChain) GetLastBlock() (*Block, error) {
	lastBlockHash, err := chain.GetLastBlockHash()
	if err != nil {
		return nil, err
	}
	return chain.GetBlock(lastBlockHash)
}

func (chain *BlockChain) GetLastBlockHash() (string, error) {
	lastBlockHash, err := chain.DB.Get(DBKeyLastBlockHash)
	if err != nil {
		return "", err
	}
	return lastBlockHash, nil
}

func (chain *BlockChain) GetBlock(blockHash string) (*Block, error) {
	return chain.DB.GetBlock(blockHash)
}

func (chain *BlockChain) FindUnspentTransactions(address string) ([]*Transaction, error) {
	var unspentTxs []*Transaction

	lastBlockHash, err := chain.GetLastBlockHash()
	if err != nil {
		return nil, err
	}
	iter := NewIterator(chain.DB, lastBlockHash)

	spentTXOs := make(map[string][]int)
	for iter.HasNext() {
		b, err := iter.Next()
		if err != nil {
			return nil, err
		}
		for _, tx := range b.Transactions {

		LabelOutput:
			for outIdx, out := range tx.Outputs {
				if spentTXOs[tx.ID] != nil {
					for _, spentOut := range spentTXOs[tx.ID] {
						if spentOut == outIdx {
							// consider it as a multilevel `continue`
							continue LabelOutput
						}
					}
				}
				// in this simplify implementation, if the output
				// can be unlocked by the given address, it's unspent
				if out.CanUnlock(address) {
					unspentTxs = append(unspentTxs, tx)
				}
			}
			if !tx.IsCoinbase() {
				for _, in := range tx.Inputs {
					if in.CanUnlock(address) {
						spentTXOs[in.ID] = append(spentTXOs[in.ID], in.Out)
					}
				}
			}
		}
	}
	return unspentTxs, nil
}

func (chain *BlockChain) GetUTXO(address string) ([]TXOutput, error) {
	var utxo []TXOutput
	unspentTxs, err := chain.FindUnspentTransactions(address)
	if err != nil {
		return nil, err
	}
	for _, tx := range unspentTxs {
		for _, out := range tx.Outputs {
			utxo = append(utxo, out)
		}
	}
	return utxo, nil
}

func (chain *BlockChain) GetBalance(address string) (amount int, err error) {
	utxo, err := chain.GetUTXO(address)
	if err != nil {
		return 0, err
	}

	for _, out := range utxo {
		amount += out.CoinCount
	}
	return amount, nil
}

var ErrNoSpendableAmount = errors.New("no spendable amount")

// FindSpendableOutputs 's return parameter `unspendOuts` store such k-v: transactionID<->transaction's outputs
func (chain *BlockChain) FindSpendableOutputs(address string, amount int) (
	spendableAmount int, unspentOuts map[string][]int, err error) {

	unspentTxs, err := chain.FindUnspentTransactions(address)
	if err != nil {
		return 0, nil, err
	}
	for _, tx := range unspentTxs {
		for outIdx, out := range tx.Outputs {
			spendableAmount += out.CoinCount
			unspentOuts[tx.ID] = append(unspentOuts[tx.ID], outIdx)
			if spendableAmount >= amount {
				return
			}
		}
	}

	return 0, nil, ErrNoSpendableAmount
}
