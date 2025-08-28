package aptosindexergrpcgo

import (
	"context"
	"log/slog"
	"testing"

	fullnodev1 "github.com/ice-coldbell/aptos-indexer-grpc-go/aptos/fullnode/v1"
	indexerv1 "github.com/ice-coldbell/aptos-indexer-grpc-go/aptos/indexer/v1"
)

// Using Your Own Fullnode.
func TestFullnodeClient(t *testing.T) {
	client, err := NewClient[*fullnodev1.TransactionsFromNodeResponse](
		context.Background(),
		"http://localhost:8080",
		"localhost:50051",
		"",
		false,
	)
	if err != nil {
		t.Fatalf("failed to create indexer: %v", err)
	}
	defer client.Stop()

	client.Start(nil)

	for {
		txs, err := client.Read()
		if err != nil {
			t.Fatalf("failed to read txs: %v", err)
		}
		slog.Info("received txs", "count", len(txs))
	}
}

// Using Aptos Labs Dataservice.
func TestDataserviceClient(t *testing.T) {
	client, err := NewClient[*indexerv1.TransactionsResponse](
		context.Background(),
		"https://api.mainnet.aptoslabs.com",
		"grpc.mainnet.aptoslabs.com:443",
		"YOUR_API_KET",
		true,
	)
	if err != nil {
		t.Fatalf("failed to create indexer: %v", err)
	}
	defer client.Stop()

	client.Start(nil)

	for {
		txs, err := client.Read()
		if err != nil {
			t.Fatalf("failed to read txs: %v", err)
		}
		slog.Info("received txs", "count", len(txs))
	}
}
