package stats

import (
	"context"
)

// TransactionsFetcher fetches provider-specific blockchain data and returns it in form of Transaction
type TransactionsFetcher interface {
	// FetchTransactionsNLastBlocks returns channel with Transaction from N last blocks
	FetchTransactionsNLastBlocks(ctx context.Context, blocksCount uint) (<-chan Transaction, <-chan error)
}
