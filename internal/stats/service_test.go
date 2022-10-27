package stats_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tofuoverdose/txstat/internal/stats"
)

func TestService_GetAddressWithGreatestExchangeDiff(t *testing.T) {
	expectedAddr := "John Doe"

	f := &MockFetcher{
		ShouldReturnResult: []stats.Transaction{
			{
				SenderAddr:   expectedAddr,
				ReceiverAddr: "Jane",
				Amount:       big.NewInt(666),
				Fee:          big.NewInt(3),
			},
			{
				SenderAddr:   "Alice",
				ReceiverAddr: "Bob",
				Amount:       big.NewInt(5),
				Fee:          big.NewInt(1),
			},
		},
	}

	s := stats.NewService(f)
	addr, err := s.TopExchangeDiffAddress(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, expectedAddr, addr)
}

func TestService_GetAddressWithGreatestExchangeDiff_EmptyResult(t *testing.T) {
	f := &MockFetcher{
		ShouldReturnResult: nil,
	}

	s := stats.NewService(f)
	addr, err := s.TopExchangeDiffAddress(context.Background())
	assert.Empty(t, addr)
	assert.ErrorIs(t, stats.ErrEmptyBlockChain, err)
}

func TestService_GetAddressWithGreatestExchangeDiff_ErrorFromFetcher(t *testing.T) {
	testError := errors.New("error yo")
	f := &MockFetcher{
		ShouldReturnError: testError,
	}

	s := stats.NewService(f)
	addr, err := s.TopExchangeDiffAddress(context.Background())
	assert.Empty(t, addr)
	assert.Error(t, err)
	assert.ErrorIs(t, err, testError)
}

type MockFetcher struct {
	ShouldReturnResult []stats.Transaction
	ShouldReturnError  error
}

func (f *MockFetcher) FetchTransactionsNLastBlocks(_ context.Context, _ uint) (<-chan stats.Transaction, <-chan error) {
	resChan, errChan := make(chan stats.Transaction), make(chan error)

	go func() {
		if f.ShouldReturnError != nil {
			errChan <- f.ShouldReturnError
			close(resChan)
			return
		}
		for _, tx := range f.ShouldReturnResult {
			resChan <- tx
		}
		close(resChan)
	}()

	return resChan, errChan
}
