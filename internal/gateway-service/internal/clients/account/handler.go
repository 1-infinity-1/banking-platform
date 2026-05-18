package account

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (c *Client) CreateAccount(ctx context.Context, params models.CreateAccountParams) (*models.Account, error) {
	resp, err := c.svc.CreateAccount(ctx, toProtoCreateAccountRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	account, err := toAccount(resp)
	if err != nil {
		return nil, fmt.Errorf("toAccount: %w", err)
	}
	return account, nil
}

func (c *Client) GetUserAccounts(ctx context.Context, params models.GetUserAccountsParams) ([]models.Account, error) {
	resp, err := c.svc.GetUserAccounts(ctx, toProtoGetUserAccountsRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	accounts := make([]models.Account, 0, len(resp.GetAccounts()))
	for _, a := range resp.GetAccounts() {
		acc, mappingErr := toAccount(a)
		if mappingErr != nil {
			return nil, fmt.Errorf("toAccount: %w", mappingErr)
		}
		accounts = append(accounts, *acc)
	}
	return accounts, nil
}

func (c *Client) GetAccount(ctx context.Context, params models.GetAccountParams) (*models.Account, error) {
	resp, err := c.svc.GetAccount(ctx, toProtoGetAccountRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	account, err := toAccount(resp)
	if err != nil {
		return nil, fmt.Errorf("toAccount: %w", err)
	}
	return account, nil
}

func (c *Client) GetBalance(ctx context.Context, params models.GetBalanceParams) (*models.Balance, error) {
	resp, err := c.svc.GetBalance(ctx, toProtoGetBalanceRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	balance, err := toBalance(resp)
	if err != nil {
		return nil, fmt.Errorf("toBalance: %w", err)
	}
	return balance, nil
}

func (c *Client) UpdateStatus(ctx context.Context, params models.UpdateAccountStatusParams) (*models.Account, error) {
	resp, err := c.svc.UpdateStatus(ctx, toProtoUpdateStatusRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	account, err := toAccount(resp.GetAccount())
	if err != nil {
		return nil, fmt.Errorf("toAccount: %w", err)
	}
	return account, nil
}
