package stats

import "math/big"

// Transaction represents a provider-independent data about transaction in blockchain
type Transaction struct {
	SenderAddr   string
	ReceiverAddr string
	Amount       *big.Int
	Fee          *big.Int
}

func (tx Transaction) AmountReceived() *big.Int {
	return tx.Amount
}

func (tx Transaction) AmountSpent() *big.Int {
	sum := &big.Int{}
	sum.Add(tx.Amount, tx.Fee)
	return sum
}
