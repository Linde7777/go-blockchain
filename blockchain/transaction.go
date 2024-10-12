package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
)

type TXInput struct {
	ID  string // The ID of the transaction containing the output we're spending
	Out int    // The index of the output in the referenced transaction

	// The signature to unlock the output,
	// in this simple implementation,
	// it's just the sender's address
	Signature string
}

func (in *TXInput) CanUnlock(data string) bool {
	return in.Signature == data
}

// -------------------------------------------------------
// -------------------------------------------------------
// -------------------------------------------------------

type TXOutput struct {
	CoinCount int

	// The public key of the recipient, receiver should
	// use his privateKey to unlock. In this simple implementation,
	// it's just an address, we just do a string comparison
	PublicKey string
}

func (out *TXOutput) CanUnlock(data string) bool {
	return out.PublicKey == data
}

// -------------------------------------------------------
// -------------------------------------------------------
// -------------------------------------------------------

type Transaction struct {
	ID      string
	Inputs  []TXInput
	Outputs []TXOutput
}

func NewTransaction(from, to string, amount int, chain *BlockChain) (*Transaction, error) {
	spendableAmount, unspentOuts, err := chain.FindSpendableOutputs(from, amount)
	if err != nil {
		return nil, err
	}

	var inputs []TXInput
	for txID, outs := range unspentOuts {
		for _, out := range outs {
			input := TXInput{
				ID:        txID,
				Out:       out,
				Signature: from,
			}
			inputs = append(inputs, input)
		}
	}

	var outputs []TXOutput
	outputs = append(outputs, TXOutput{amount, to})

	// i.e. if there is a change, return it to the sender
	if spendableAmount > amount {
		o := TXOutput{CoinCount: spendableAmount - amount, PublicKey: from}
		outputs = append(outputs, o)
	}

	tx := &Transaction{Inputs: inputs, Outputs: outputs}
	err = tx.SetID()
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// NewCoinbaseTX create a special transaction which will create 100 coin
func NewCoinbaseTX(to, data string) (*Transaction, error) {
	inputs := []TXInput{
		{
			// coinbaseTX have no previous transaction,
			// so the field `ID` and `Out` are empty
			ID:        "",
			Out:       -1,
			Signature: data,
		},
	}

	outputs := []TXOutput{
		{
			CoinCount: 100,
			PublicKey: to,
		},
	}

	tx := &Transaction{
		ID:      "",
		Inputs:  inputs,
		Outputs: outputs,
	}
	err := tx.SetID()
	return tx, err
}

func (tx *Transaction) SetID() error {
	var encoded bytes.Buffer
	var hash [32]byte

	encode := gob.NewEncoder(&encoded)
	err := encode.Encode(tx)
	if err != nil {
		return err
	}

	hash = sha256.Sum256(encoded.Bytes())
	tx.ID = string(hash[:])
	return nil
}

func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 &&
		len(tx.Outputs) == 1 &&
		tx.Inputs[0].ID == ""
}
