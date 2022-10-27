package stats

import (
	"context"
	"errors"
	"fmt"
	"math/big"
	"sort"
)

const (
	headBlockCount = 100
)

type Service interface {
	TopExchangeDiffAddress(ctx context.Context) (string, error)
}

type service struct {
	txf TransactionsFetcher
}

func NewService(txf TransactionsFetcher) Service {
	s := service{
		txf: txf,
	}
	return s
}

func (s service) TopExchangeDiffAddress(ctx context.Context) (string, error) {
	// this implementation uses map that accumulates all transactions for each address and then it sorts them
	// it's super slow but straightforward
	// todo rewrite to improve performance

	stats := make(exchangeStats)

	txChan, errChan := s.txf.FetchTransactionsNLastBlocks(ctx, headBlockCount)
fl:
	for {
		select {
		case err, open := <-errChan:
			if !open {
				break fl
			}
			return "", fmt.Errorf("error during block fetch: %w", err)
		case tx := <-txChan:
			stats.AddTransaction(tx)
		}
	}

	sorted := stats.GetSortedByDiffDesc()
	if len(sorted) == 0 {
		return "", ErrEmptyBlockChain
	}

	return sorted[0].Addr, nil
}

// aggregates exchangeStat by address
type exchangeStats map[string]exchangeStat

func (stats exchangeStats) AddTransaction(tx Transaction) {
	// add to receiver
	stats.Add(tx.ReceiverAddr, tx.AmountReceived())
	// subtract from sender
	stats.Sub(tx.SenderAddr, tx.AmountSpent())
}

func (stats exchangeStats) Add(addr string, amount *big.Int) {
	es, ok := stats[addr]
	if !ok {
		es = exchangeStat{Addr: addr, Diff: &big.Int{}}
	}
	es.Diff = es.Diff.Add(es.Diff, amount)
	stats[addr] = es
}

func (stats exchangeStats) Sub(addr string, amount *big.Int) {
	stats.Add(addr, amount.Neg(amount))
}

func (stats exchangeStats) GetSortedByDiffDesc() []exchangeStat {
	ret := make([]exchangeStat, 0, len(stats))
	for _, es := range stats {
		ret = append(ret, es)
	}
	sort.SliceStable(ret, func(i, j int) bool {
		a, b := ret[i].Diff, ret[j].Diff
		return a.Cmp(b) > 0
	})
	return ret
}

// holds diff of all received and sent currency
type exchangeStat struct {
	Diff *big.Int
	Addr string
}

type serviceError struct {
	code string
	text string
}

func (e serviceError) Error() string {
	return e.msg()
}

func (e serviceError) Failed() error {
	return errors.New(e.msg())
}

func (e serviceError) msg() string {
	return fmt.Sprintf("%s %s", e.code, e.text)
}

var (
	ErrEmptyBlockChain = serviceError{"ErrEmptyBlockChain", "blockchain is empty"}
)
