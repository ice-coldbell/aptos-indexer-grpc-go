package client

import (
	"context"
	"fmt"

	indexerv1 "github.com/ice-coldbell/aptos-indexer-grpc-go/aptos/indexer/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type rawDataClient struct {
	rawDataClient indexerv1.RawDataClient

	conn   *grpc.ClientConn
	apiKey string
}

func NewRawDataClient(serverAddr string, apiKey string, useTLS bool) (*rawDataClient, error) {
	var creds credentials.TransportCredentials
	if useTLS {
		creds = credentials.NewTLS(nil)
	} else {
		creds = insecure.NewCredentials()
	}

	opts := []grpc.CallOption{
		grpc.MaxCallRecvMsgSize(1024 * 1024 * 1024), // 1GB
		grpc.MaxCallSendMsgSize(1024 * 1024 * 1024), // 1GB
	}

	conn, err := grpc.NewClient(
		serverAddr,
		grpc.WithTransportCredentials(creds),
		grpc.WithDefaultCallOptions(opts...),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to gRPC server: %w", err)
	}

	return &rawDataClient{
		rawDataClient: indexerv1.NewRawDataClient(conn),
		conn:          conn,
		apiKey:        apiKey,
	}, nil
}

func (c *rawDataClient) NewStream(ctx context.Context, startVersion uint64, count uint64) (Stream[*indexerv1.TransactionsResponse], error) {
	ctx = c.createAuthContext(ctx)

	request := &indexerv1.GetTransactionsRequest{
		StartingVersion:   &startVersion,
		TransactionsCount: &count,
	}
	if count == 0 {
		request.TransactionsCount = nil
	}

	stream, err := c.rawDataClient.GetTransactions(ctx, request)
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (c *rawDataClient) Close() error {
	return c.conn.Close()
}

func (c *rawDataClient) createAuthContext(ctx context.Context) context.Context {
	if c.apiKey != "" {
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.apiKey)
	}
	return ctx
}
