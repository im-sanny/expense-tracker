package apperrors

import "errors"

var (
	ErrInvalidMin          = errors.New("invalid min value")
	ErrInvalidMax          = errors.New("invalid max value")
	ErrMinGraterThanMax    = errors.New("min cannot be grater than max")
	ErrInvalidFromDate     = errors.New("invalid 'from' date format")
	ErrInvalidToDate       = errors.New("invalid 'to' date format")
	ErrFromDateAfterTo     = errors.New("'from' date must be before 'to' date")
	ErrInvalidPage         = errors.New("invalid page number")
	ErrInvalidLimit        = errors.New("invalid limit value")
	ErrFailedToGetExpenses = errors.New("failed to retrieve expenses")
	ErrFailedToCount       = errors.New("failed to count expenses")
	ErrDatabase            = errors.New("database error")
	ErrInternalServer      = errors.New("internal server error")
	ErrNotFound            = errors.New("resource not found")
	ErrBadRequest          = errors.New("bad request")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrForbidden           = errors.New("forbidden")
	ErrTimeOut             = errors.New("request timeout")
)

func IsValidationError(err error) bool {
	return errors.Is(err, ErrInvalidMax) ||
		errors.Is(err, ErrMinGraterThanMax) ||
		errors.Is(err, ErrInvalidFromDate) ||
		errors.Is(err, ErrInvalidToDate) ||
		errors.Is(err, ErrFromDateAfterTo) ||
		errors.Is(err, ErrInvalidPage) ||
		errors.Is(err, ErrInvalidLimit)
}

func IsClientError(err error) bool {
	return IsValidationError(err) ||
		errors.Is(err, ErrNotFound) ||
		errors.Is(err, ErrUnauthorized) ||
		errors.Is(err, ErrForbidden) ||
		errors.Is(err, ErrBadRequest)
}
