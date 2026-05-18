package ledger

import (
	"context"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (c *Client) GetStatement(ctx context.Context, params models.GetStatementParams) (*models.Statement, error) {
	resp, err := c.svc.GetStatement(ctx, toProtoGetStatementRequest(params))
	if err != nil {
		return nil, mapGRPCError(err)
	}
	statement, err := toStatement(resp)
	if err != nil {
		return nil, fmt.Errorf("toStatement: %w", err)
	}
	return statement, nil
}
