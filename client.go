package aptosindexergrpcgo

import (
	"context"
	"io"
	"sync"
	"sync/atomic"

	v1 "github.com/ice-coldbell/aptos-indexer-grpc-go/aptos/transaction/v1"
	"github.com/ice-coldbell/aptos-indexer-grpc-go/internal/client"
)

type Client[T StreamResponse] struct {
	ctx    context.Context
	cancel context.CancelFunc

	httpAddr string
	grpcAddr string

	grpcClient client.Client[T]
	txChan     chan []*v1.Transaction
	once       sync.Once

	retryCount  atomic.Int64
	lastVersion atomic.Uint64

	wg sync.WaitGroup
}

type StreamResponse = client.StreamResponse
type FullnodeStreamResponse = client.FullnodeStreamResponse
type DataserviceStreamResponse = client.DataserviceStreamResponse

func NewClient[T StreamResponse](
	ctx context.Context,
	httpAddr string,
	grpcAddr string,
	apiKey string,
	useTLS bool,
) (indexer *Client[T], err error) {
	ctx, cancel := context.WithCancel(ctx)
	indexer = &Client[T]{
		ctx:      ctx,
		cancel:   cancel,
		httpAddr: httpAddr,
		grpcAddr: grpcAddr,
		txChan:   make(chan []*v1.Transaction, 512),
	}

	indexer.grpcClient, err = client.NewClient[T](grpcAddr, apiKey, useTLS)
	if err != nil {
		return nil, err
	}
	return indexer, nil
}

func (i *Client[T]) Start(version *uint64) {
	i.once.Do(func() {
		if version != nil {
			i.lastVersion.Store(*version)
		}
		i.wg.Add(1)
		go i.streamLoop()
	})
}

func (i *Client[T]) Stop() {
	i.cancel()
	i.wg.Wait()
}

func (i *Client[T]) Read() ([]*v1.Transaction, error) {
	select {
	case txs, ok := <-i.txChan:
		if !ok {
			return nil, io.EOF
		}
		return txs, nil
	case <-i.ctx.Done():
		return nil, i.ctx.Err()
	}
}
