package controllers

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/mislavperi/jafa/server/internal/domain/apperr"
	"github.com/mislavperi/jafa/server/internal/server/httperr"
	"github.com/mislavperi/jafa/server/internal/server/middleware"
)

// requireUser returns the authenticated user's ID, or writes a 401 and reports
// false when the request is unauthenticated. Handlers should return early when
// ok is false.
func requireUser(ctx *gin.Context) (int64, bool) {
	uid, ok := middleware.CurrentUserID(ctx)
	if !ok {
		httperr.Unauthorized(ctx, "authentication required")
		return 0, false
	}
	return uid, true
}

// respondExpenseError maps an ExpenseService error to its HTTP response and
// reports whether it handled the error. A nil error is not handled. This keeps
// the not-found/validation/internal mapping in one place across handlers.
func respondExpenseError(ctx *gin.Context, err error) bool {
	switch {
	case err == nil:
		return false
	case errors.Is(err, apperr.ErrExpenseNotFound):
		httperr.NotFound(ctx, "expense not found")
	case errors.Is(err, apperr.ErrInvalidStartDate),
		errors.Is(err, apperr.ErrInvalidInstallmentCount),
		errors.Is(err, apperr.ErrInvalidKind):
		httperr.BadRequest(ctx, err.Error(), err)
	default:
		httperr.Internal(ctx, err)
	}
	return true
}
