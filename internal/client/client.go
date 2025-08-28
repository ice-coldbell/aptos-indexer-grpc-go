package client

import (
	"context"
	"fmt"

	fullnodev1 "github.com/ice-coldbell/aptos-indexer-grpc-go/aptos/fullnode/v1"
	indexerv1 "github.com/ice-coldbell/aptos-indexer-grpc-go/aptos/indexer/v1"
	v1 "github.com/ice-coldbell/aptos-indexer-grpc-go/aptos/transaction/v1"
)

type Client[T StreamResponse] interface {
	NewStream(ctx context.Context, startVersion uint64, count uint64) (Stream[T], error)
	Close() error
}

type Stream[T StreamResponse] interface {
	Recv() (T, error)
}

type StreamResponse interface {
	FullnodeStreamResponse | DataserviceStreamResponse

	GetTransactions() []*v1.Transaction
}
type FullnodeStreamResponse = *fullnodev1.TransactionsFromNodeResponse
type DataserviceStreamResponse = *indexerv1.TransactionsResponse

func NewClient[T StreamResponse](serverAddr string, apiKey string, useTLS bool) (Client[T], error) {
	var zero T
	switch any(zero).(type) {
	case *fullnodev1.TransactionsFromNodeResponse:
		client, err := NewFullnodeClient(serverAddr, apiKey, useTLS)
		if err != nil {
			return nil, err
		}
		return any(client).(Client[T]), nil
	case *indexerv1.TransactionsResponse:
		client, err := NewRawDataClient(serverAddr, apiKey, useTLS)
		if err != nil {
			return nil, err
		}
		return any(client).(Client[T]), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", zero)
	}
}
