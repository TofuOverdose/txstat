package fetcher

import (
	"context"
	"fmt"
	"math/big"

	"github.com/tofuoverdose/txstat/internal/txstat/domain/stats"
	"github.com/tofuoverdose/txstat/pkg/getblock/eth"
)

// Fetcher implements stats.TransactionsFetcher
type Fetcher struct {
	Client eth.Client
}

func (f *Fetcher) FetchTransactionsNLastBlocks(ctx context.Context, blocksCount uint) (<-chan stats.Transaction, <-chan error) {
	resChan, errChan := make(chan stats.Transaction), make(chan error)

	go func() {
		var curBlock string
		var err error
		for i := uint(0); i < blocksCount; i++ {
			select {
			case <-ctx.Done():
				close(resChan)
				errChan <- ctx.Err()
				close(errChan)
				return
			default:
				curBlock, err = f.Client.BlockNumber(ctx)
				if err != nil {
					close(resChan)
					errChan <- fmt.Errorf("failed to get number of next block: %w", err)
					close(errChan)
					return
				}

				block, err := f.Client.GetBlockByNumber(ctx, curBlock, true)
				if err != nil {
					close(resChan)
					errChan <- fmt.Errorf("failed to get data for block %s: %w", curBlock, err)
					close(errChan)
					return
				}

				for _, ethTx := range block.Transactions {
					tx, err := transactionFromEth(ethTx)
					if err != nil {
						close(resChan)
						errChan <- fmt.Errorf("failed to create transaction struct: %w", err)
						close(errChan)
						return
					}
					resChan <- *tx
				}

				curBlock = block.ParentHash
			}
		}
		close(errChan)
		close(resChan)
	}()

	return resChan, errChan
}

func transactionFromEth(ethTx eth.Transaction) (*stats.Transaction, error) {
	tx := stats.Transaction{
		SenderAddr:   ethTx.From,
		ReceiverAddr: ethTx.To,
	}

	amount := &big.Int{}
	if _, ok := amount.SetString(ethTx.Value, 0); !ok {
		return nil, fmt.Errorf("failed to parse value %s", ethTx.Value)
	}
	tx.Amount = amount

	fee := &big.Int{}
	if _, ok := fee.SetString(ethTx.GasPrice, 0); !ok {
		return nil, fmt.Errorf("failed to parse gas price %s", ethTx.GasPrice)
	}
	tx.Fee = fee

	return &tx, nil
}
