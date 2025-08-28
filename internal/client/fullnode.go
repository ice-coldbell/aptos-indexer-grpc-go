package client

import (
	"context"
	"fmt"

	fullnodev1 "github.com/ice-coldbell/aptos-indexer-grpc-go/aptos/fullnode/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
)

type fullnodeClient struct {
	fullnodeClient fullnodev1.FullnodeDataClient

	conn   *grpc.ClientConn
	apiKey string
}

func NewFullnodeClient(serverAddr string, apiKey string, useTLS bool) (*fullnodeClient, error) {
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

	return &fullnodeClient{
		fullnodeClient: fullnodev1.NewFullnodeDataClient(conn),
		conn:           conn,
		apiKey:         apiKey,
	}, nil
}

func (c *fullnodeClient) NewStream(ctx context.Context, startVersion uint64, count uint64) (Stream[*fullnodev1.TransactionsFromNodeResponse], error) {
	ctx = c.createAuthContext(ctx)

	request := &fullnodev1.GetTransactionsFromNodeRequest{
		StartingVersion:   &startVersion,
		TransactionsCount: &count,
	}
	if count == 0 {
		request.TransactionsCount = nil
	}

	stream, err := c.fullnodeClient.GetTransactionsFromNode(ctx, request)
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (c *fullnodeClient) Close() error {
	return c.conn.Close()
}

func (c *fullnodeClient) createAuthContext(ctx context.Context) context.Context {
	if c.apiKey != "" {
		return metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+c.apiKey)
	}
	return ctx
}
