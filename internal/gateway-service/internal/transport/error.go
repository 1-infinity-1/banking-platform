package transport

import (
	"context"
	"errors"
	"net/http"

	api "github.com/1-infinity-1/banking-platform/internal/gateway-service/api/ogen"
	"github.com/1-infinity-1/banking-platform/internal/gateway-service/internal/models"
	"github.com/ogen-go/ogen/ogenerrors"
)

func (g *GatewayHandler) NewError(_ context.Context, err error) *api.ErrorStatusCode {
	var secErr *ogenerrors.SecurityError
	if errors.As(err, &secErr) {
		var forbiddenErr *models.ForbiddenError
		if errors.As(secErr.Err, &forbiddenErr) {
			return &api.ErrorStatusCode{
				StatusCode: http.StatusForbidden,
				Response:   api.Error{Code: "FORBIDDEN", Message: forbiddenErr.Error()},
			}
		}
		return &api.ErrorStatusCode{
			StatusCode: http.StatusUnauthorized,
			Response:   api.Error{Code: "UNAUTHORIZED", Message: "unauthorized"},
		}
	}

	var unauthorizedErr *models.UnauthorizedError
	if errors.As(err, &unauthorizedErr) {
		return &api.ErrorStatusCode{
			StatusCode: http.StatusUnauthorized,
			Response:   api.Error{Code: "UNAUTHORIZED", Message: unauthorizedErr.Error()},
		}
	}

	var businessErr *models.BusinessError
	if errors.As(err, &businessErr) {
		return &api.ErrorStatusCode{
			StatusCode: http.StatusLocked,
			Response:   api.Error{Code: "LOCKED", Message: businessErr.Error()},
		}
	}

	var conflictErr *models.ConflictError
	if errors.As(err, &conflictErr) {
		return &api.ErrorStatusCode{
			StatusCode: http.StatusConflict,
			Response:   api.Error{Code: "CONFLICT", Message: conflictErr.Error()},
		}
	}

	var validationErr *models.ValidationError
	if errors.As(err, &validationErr) {
		return &api.ErrorStatusCode{
			StatusCode: http.StatusBadRequest,
			Response:   api.Error{Code: "BAD_REQUEST", Message: validationErr.Error()},
		}
	}

	var notFoundErr *models.NotFoundError
	if errors.As(err, &notFoundErr) {
		return &api.ErrorStatusCode{
			StatusCode: http.StatusNotFound,
			Response:   api.Error{Code: "NOT_FOUND", Message: notFoundErr.Error()},
		}
	}

	return &api.ErrorStatusCode{
		StatusCode: http.StatusInternalServerError,
		Response:   api.Error{Code: "INTERNAL_ERROR", Message: "internal server error"},
	}
}
