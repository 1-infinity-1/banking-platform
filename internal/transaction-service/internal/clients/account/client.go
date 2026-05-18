package account

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/transaction-service/internal/models"
	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc"
)

type Client struct {
	conn accountpb.AccountServiceClient
}

func NewClient(grpcConn *grpc.ClientConn) *Client {
	return &Client{conn: accountpb.NewAccountServiceClient(grpcConn)}
}

func (c *Client) Debit(ctx context.Context, req models.DebitRequest) (*models.DebitResult, error) {
	resp, err := c.conn.Debit(ctx, &accountpb.DebitRequest{
		AccountId:      req.AccountID.String(),
		Amount:         req.Amount.String(),
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, fmt.Errorf("conn.Debit: %w", err)
	}

	balanceAfter, err := decimal.NewFromString(resp.GetBalanceAfter())
	if err != nil {
		return nil, fmt.Errorf("decimal.NewFromString(balance_after=%q): %w", resp.GetBalanceAfter(), err)
	}

	return &models.DebitResult{
		AccountID:    resp.GetAccountId(),
		BalanceAfter: balanceAfter,
	}, nil
}

func (c *Client) Credit(ctx context.Context, req models.CreditRequest) (*models.CreditResult, error) {
	resp, err := c.conn.Credit(ctx, &accountpb.CreditRequest{
		AccountId:      req.AccountID.String(),
		Amount:         req.Amount.String(),
		IdempotencyKey: req.IdempotencyKey,
	})
	if err != nil {
		return nil, fmt.Errorf("conn.Credit: %w", err)
	}

	balanceAfter, err := decimal.NewFromString(resp.GetBalanceAfter())
	if err != nil {
		return nil, fmt.Errorf("decimal.NewFromString(balance_after=%q): %w", resp.GetBalanceAfter(), err)
	}

	return &models.CreditResult{
		AccountID:    resp.GetAccountId(),
		BalanceAfter: balanceAfter,
	}, nil
}
