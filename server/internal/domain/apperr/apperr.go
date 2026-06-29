// Package apperr holds the application's sentinel errors — the predefined,
// expected failure modes the services return and the controllers map to HTTP
// statuses. Unexpected/infrastructure failures (raw DB and transaction errors)
// are not enumerated here; they propagate untyped and surface as 500s.
//
// Errors are declared as a string-backed const type so they are true constants
// (immutable, cannot be reassigned at runtime). Match them with errors.Is, which
// works whether the error is returned directly or wrapped with %w.
package apperr

// Error is an immutable sentinel error.
type Error string

func (e Error) Error() string { return string(e) }

const (
	// Expense (mapped to 404/400 by controllers).

	// ErrExpenseNotFound: the expense does not exist or is not owned by the
	// requesting user. → 404.
	ErrExpenseNotFound Error = "expense not found"
	// ErrInvalidStartDate: a recurring schedule's start date is not a valid
	// YYYY-MM-DD date. → 400.
	ErrInvalidStartDate Error = "invalid recurring schedule start date"
	// ErrInvalidInstallmentCount: an expense is split into fewer than 2
	// payments. → 400.
	ErrInvalidInstallmentCount Error = "installment count must be at least 2"
	// ErrInvalidKind: an entry's kind is neither "expense" nor "income". → 400.
	ErrInvalidKind Error = "kind must be 'expense' or 'income'"

	// Auth / user (mapped to 401/409/404 by controllers).

	// ErrInvalidCredentials: unknown username or password mismatch. → 401.
	ErrInvalidCredentials Error = "invalid username or password"
	// ErrUsernameTaken: registration hit the username uniqueness constraint.
	// → 409.
	ErrUsernameTaken Error = "username already taken"
	// ErrUserNotFound: the user does not exist (e.g. delete of a missing/already
	// removed account). → 404.
	ErrUserNotFound Error = "user not found"
)
