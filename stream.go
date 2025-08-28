package aptosindexergrpcgo

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"slices"
	"strconv"
	"time"

	v1 "github.com/ice-coldbell/aptos-indexer-grpc-go/aptos/transaction/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxRetry int64 = 3

func (i *Client[T]) streamLoop() {
	defer close(i.txChan)
	defer i.wg.Done()
	defer i.cancel()

	retry := i.retryCount.Load()
	for {
		select {
		case <-i.ctx.Done():
			return
		default:
			if err := i.startStream(); err != nil && status.Code(err) != codes.Canceled {
				slog.Warn("failed to start stream", "error", err)
				i.retryCount.Add(1)
				if retry >= maxRetry {
					slog.Error("failed to start stream", "error", err, "retry", retry)
					return
				}
				time.Sleep(time.Second * (1 << retry))
				continue
			}
		}
	}
}

func (i *Client[T]) startStream() error {
	latestVersion, err := i.getLatestVersion()
	if err != nil {
		return err
	}

	slog.Info("starting grpc stream", "latest_version", latestVersion)
	stream, err := i.grpcClient.NewStream(i.ctx, latestVersion, 0)
	if err != nil {
		return err
	}

	i.retryCount.Store(0)
	var resp T
	var transactions []*v1.Transaction
	for {
		select {
		case <-i.ctx.Done():
			return i.ctx.Err()
		default:
			resp, err = stream.Recv()
			if err != nil {
				return err
			}

			txs := FilterFast(resp.GetTransactions(), func(tx *v1.Transaction) bool {
				return tx.Type == v1.Transaction_TRANSACTION_TYPE_USER
			})
			if len(txs) == 0 {
				continue // skip data
			}

			transactions = txs
			slices.SortFunc(transactions, func(a *v1.Transaction, b *v1.Transaction) int {
				return int(a.GetVersion() - b.GetVersion())
			})
			slog.Debug("grpc stream transactions", "count", len(transactions))

			i.lastVersion.Store(transactions[len(transactions)-1].GetVersion())

			i.txChan <- transactions
		}
	}
}

func (i *Client[T]) getLatestVersion() (uint64, error) {
	if i.lastVersion.Load() > 0 {
		return i.lastVersion.Load(), nil
	}

	resp, err := http.Get(i.httpAddr + "/v1")
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	var latestVersion struct {
		LedgerVersion string `json:"ledger_version"`
	}
	if err := json.Unmarshal(body, &latestVersion); err != nil {
		return 0, err
	}
	version, err := strconv.ParseUint(latestVersion.LedgerVersion, 10, 64)
	if err != nil {
		return 0, err
	}

	i.lastVersion.Store(version)
	return version, nil
}
