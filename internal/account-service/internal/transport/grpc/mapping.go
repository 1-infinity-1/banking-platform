package grpc

import (
	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func toProtoAccountStatus(s models.AccountStatus) accountpb.AccountStatus {
	switch s {
	case models.AccountStatusActive:
		return accountpb.AccountStatus_ACCOUNT_STATUS_ACTIVE
	case models.AccountStatusBlocked:
		return accountpb.AccountStatus_ACCOUNT_STATUS_BLOCKED
	case models.AccountStatusClosed:
		return accountpb.AccountStatus_ACCOUNT_STATUS_CLOSED
	case models.AccountStatusUnspecified:
		return accountpb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED
	}
	return accountpb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED
}

func fromProtoAccountStatus(s accountpb.AccountStatus) models.AccountStatus {
	switch s {
	case accountpb.AccountStatus_ACCOUNT_STATUS_UNSPECIFIED:
		return models.AccountStatusUnspecified
	case accountpb.AccountStatus_ACCOUNT_STATUS_ACTIVE:
		return models.AccountStatusActive
	case accountpb.AccountStatus_ACCOUNT_STATUS_BLOCKED:
		return models.AccountStatusBlocked
	case accountpb.AccountStatus_ACCOUNT_STATUS_CLOSED:
		return models.AccountStatusClosed
	}
	return models.AccountStatusUnspecified
}

func toProtoAccount(a *models.Account) *accountpb.Account {
	return &accountpb.Account{
		Id:        a.PublicID.String(),
		UserId:    a.UserID.String(),
		Currency:  a.Currency,
		Balance:   a.Balance.String(),
		Status:    toProtoAccountStatus(a.Status),
		CreatedAt: timestamppb.New(a.CreatedAt),
		UpdatedAt: timestamppb.New(a.UpdatedAt),
	}
}

func toProtoAccountsList(accs []*models.Account) *accountpb.AccountsList {
	out := &accountpb.AccountsList{
		Accounts: make([]*accountpb.Account, 0, len(accs)),
	}
	for _, a := range accs {
		out.Accounts = append(out.Accounts, toProtoAccount(a))
	}
	return out
}

func toProtoBalance(b *models.Balance) *accountpb.Balance {
	return &accountpb.Balance{
		AccountId: b.AccountID,
		Amount:    b.Amount.String(),
		Currency:  b.Currency,
	}
}

func toProtoUpdateStatusResponse(a *models.Account) *accountpb.UpdateStatusResponse {
	return &accountpb.UpdateStatusResponse{Account: toProtoAccount(a)}
}

func toProtoDebitResponse(r *models.DebitResult) *accountpb.DebitResponse {
	return &accountpb.DebitResponse{
		AccountId:    r.AccountID,
		BalanceAfter: r.BalanceAfter.String(),
	}
}

func toProtoCreditResponse(r *models.CreditResult) *accountpb.CreditResponse {
	return &accountpb.CreditResponse{
		AccountId:    r.AccountID,
		BalanceAfter: r.BalanceAfter.String(),
	}
}
