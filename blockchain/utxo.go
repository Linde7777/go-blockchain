package blockchain

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
)

type UTXO interface {
	// Basic properties
	GetAmount() int
	GetScriptPubKey() string
	GetTxID() string
	GetOutputIndex() int

	// Spending and validation
	CanBeSpentBy(pubKey string, signature string) bool
	IsSpent() bool
	MarkAsSpent()

	// Serialization
	Serialize() (string, error)
	Deserialize(data string) error

	// Helper methods
	Hash() []byte
	String() string
}

type UTXOImpl struct {
	Amount        int    `json:"amount"`
	ScriptPubKey  string `json:"scriptPubKey"`
	TxID          string `json:"txID"`
	OutputIndex   int    `json:"outputIndex"`
	IsSpentStatus bool   `json:"isSpent"`
}

var _ UTXO = (*UTXOImpl)(nil)

func NewUTXO(amount int, scriptPubKey string, txID string, outputIndex int) *UTXOImpl {
	return &UTXOImpl{
		Amount:        amount,
		ScriptPubKey:  scriptPubKey,
		TxID:          txID,
		OutputIndex:   outputIndex,
		IsSpentStatus: false,
	}
}

func (u *UTXOImpl) GetAmount() int {
	return u.Amount
}

func (u *UTXOImpl) GetScriptPubKey() string {
	return u.ScriptPubKey
}

func (u *UTXOImpl) GetTxID() string {
	return u.TxID
}

func (u *UTXOImpl) GetOutputIndex() int {
	return u.OutputIndex
}

func (u *UTXOImpl) CanBeSpentBy(pubKey string, signature string) bool {
	// In a real implementation, this would involve verifying the signature
	// against the public key and the transaction data.
	// For simplicity, we're just checking if the pubKey matches the ScriptPubKey.
	return u.ScriptPubKey == pubKey
}

func (u *UTXOImpl) IsSpent() bool {
	return u.IsSpentStatus
}

func (u *UTXOImpl) MarkAsSpent() {
	u.IsSpentStatus = true
}

func (u *UTXOImpl) Serialize() (string, error) {
	data, err := json.Marshal(u)
	if err != nil {
		return "", fmt.Errorf("failed to serialize UTXO: %w", err)
	}
	return string(data), nil
}

func (u *UTXOImpl) Deserialize(data string) error {
	err := json.Unmarshal([]byte(data), u)
	if err != nil {
		return fmt.Errorf("failed to deserialize UTXO: %w", err)
	}
	return nil
}

func (u *UTXOImpl) Hash() []byte {
	data := fmt.Sprintf("%s:%d:%d:%s", u.TxID, u.OutputIndex, u.Amount, u.ScriptPubKey)
	hash := sha256.Sum256([]byte(data))
	return hash[:]
}

func (u *UTXOImpl) String() string {
	return fmt.Sprintf("UTXO(TxID: %s, OutputIndex: %d, Amount: %d, ScriptPubKey: %s, IsSpent: %t)",
		u.TxID, u.OutputIndex, u.Amount, u.ScriptPubKey, u.IsSpentStatus)
}
