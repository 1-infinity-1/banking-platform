package grpc

import (
	"github.com/1-infinity-1/banking-platform/internal/account-service/internal/models"
	accountpb "github.com/1-infinity-1/banking-platform/pkg/proto/generated/go/account"
)

//nolint:unused // scaffold: used when implementing TODO handlers
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

//nolint:unused // scaffold: used when implementing TODO handlers
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
