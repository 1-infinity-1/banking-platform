package account

import (
	"errors"
	"fmt"

	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func toProtoCreateAccountRequest(p models.CreateAccountParams) *accountpb.CreateAccountRequest {
	return &accountpb.CreateAccountRequest{
		UserId:   p.UserID.String(),
		Currency: p.Currency,
	}
}

func toProtoGetUserAccountsRequest(p models.GetUserAccountsParams) *accountpb.GetUserAccountsRequest {
	return &accountpb.GetUserAccountsRequest{UserId: p.UserID.String()}
}

func toProtoGetAccountRequest(p models.GetAccountParams) *accountpb.GetAccountRequest {
	return &accountpb.GetAccountRequest{AccountId: p.AccountID.String()}
}

func toProtoGetBalanceRequest(p models.GetBalanceParams) *accountpb.GetBalanceRequest {
	return &accountpb.GetBalanceRequest{AccountId: p.AccountID.String()}
}

func toProtoUpdateStatusRequest(p models.UpdateAccountStatusParams) *accountpb.UpdateStatusRequest {
	return &accountpb.UpdateStatusRequest{
		AccountId: p.AccountID.String(),
		Status:    toProtoAccountStatus(p.Status),
	}
}

func toAccount(a *accountpb.Account) (*models.Account, error) {
	if a == nil {
		return nil, errors.New("empty account")
	}
	id, err := uuid.Parse(a.GetId())
	if err != nil {
		return nil, fmt.Errorf("parse account id: %w", err)
	}
	userID, err := uuid.Parse(a.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("parse account user_id: %w", err)
	}
	return &models.Account{
		ID:        id,
		UserID:    userID,
		Currency:  a.GetCurrency(),
		Balance:   a.GetBalance(),
		Status:    toAccountStatus(a.GetStatus()),
		CreatedAt: a.GetCreatedAt().AsTime(),
		UpdatedAt: a.GetUpdatedAt().AsTime(),
	}, nil
}

func toBalance(b *accountpb.Balance) (*models.Balance, error) {
	if b == nil {
		return nil, errors.New("empty balance")
	}
	id, err := uuid.Parse(b.GetAccountId())
	if err != nil {
		return nil, fmt.Errorf("parse balance account_id: %w", err)
	}
	return &models.Balance{
		AccountID: id,
		Amount:    b.GetAmount(),
		Currency:  b.GetCurrency(),
	}, nil
}

func toAccountStatus(s accountpb.AccountStatus) models.AccountStatus {
	switch s { //nolint:exhaustive // unspecified handled by default
	case accountpb.AccountStatus_ACCOUNT_STATUS_ACTIVE:
		return models.AccountStatusActive
	case accountpb.AccountStatus_ACCOUNT_STATUS_BLOCKED:
		return models.AccountStatusBlocked
	case accountpb.AccountStatus_ACCOUNT_STATUS_CLOSED:
		return models.AccountStatusClosed
	default:
		return models.AccountStatusUnspecified
	}
}

func toProtoAccountStatus(s models.AccountStatus) accountpb.AccountStatus {
	switch s { //nolint:exhaustive // unspecified handled by default
	case models.AccountStatusActive:
		return accountpb.AccountStatus_ACCOUNT_STATUS_ACTIVE
	case models.AccountStatusBlocked:
		return accountpb.AccountStatus_ACCOUNT_STATUS_BLOCKED
	case models.AccountStatusClosed:
		return accountpb.AccountStatus_ACCOUNT_STATUS_CLOSED
	default:
		return accountpb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED
	}
}

func mapGRPCError(err error) error {
	st, ok := status.FromError(err)
	if !ok {
		return fmt.Errorf("unexpected gRPC error: %w", err)
	}
	switch st.Code() { //nolint:exhaustive // only meaningful codes handled; default covers the rest
	case codes.NotFound:
		return models.NewNotFoundError(st.Message(), err)
	case codes.InvalidArgument:
		return models.NewValidationError("", st.Message(), err)
	case codes.AlreadyExists:
		return models.NewConflictError(st.Message(), err)
	case codes.Unauthenticated:
		return models.NewUnauthorizedError(st.Message(), err)
	default:
		return fmt.Errorf("account service error: %w", err)
	}
}
