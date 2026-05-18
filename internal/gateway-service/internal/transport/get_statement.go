package transport

import (
	"context"
	"fmt"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
)

func (g *GatewayHandler) GetStatement(ctx context.Context, params api.GetStatementParams) (api.GetStatementRes, error) {
	statement, err := g.ledgerSvc.GetStatement(ctx, models.GetStatementParams{
		AccountID: params.AccountID,
		From:      params.From,
		To:        params.To,
	})
	if err != nil {
		return nil, fmt.Errorf("g.ledgerSvc.GetStatement: %w", err)
	}

	res := toAPIStatement(*statement)
	return &res, nil
}
