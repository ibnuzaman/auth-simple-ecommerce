package constants

import "errors"

var (
	RoleCustomer = "Customer"
)
var (
	// ErrServerError      = errors.New("internal server error")
	ErrFailedBadRequest = errors.New("bad request")
	ErrConflict         = errors.New("conflict")

	UniqueViolation = "23505"
)

type ConflictError struct {
	Field string // "email" / "username" / "phone_number"
}

func (e *ConflictError) Error() string        { return "conflict on " + e.Field }
func (e *ConflictError) Is(target error) bool { return target == ErrConflict }
