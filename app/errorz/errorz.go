package errorz

import "fmt"

// ErrorCode is wrapper for our error codes
type ErrorCode string

const (
	// unknown errors
	ERUnknownSqlError ErrorCode = "A001"

	// errors returned by sql queries
	ERNoRecordFound ErrorCode = "A011"
	ERRecordExists  ErrorCode = "A012"
)

// Error
type Error struct {
	code ErrorCode
	Err  error // the original error
}

func (er *Error) Error() string {
	return fmt.Sprintf("Code <%s>: %s", er.code, errorStatus[er.code])
}

func (er *Error) Debug() error {
	return er.Err
}

// New creates a new Error instance
func New(code ErrorCode, err error) *Error {
	return &Error{code, err}
}

var errorStatus = map[ErrorCode]string{
	ERUnknownSqlError: "unknown sql error returned",

	ERRecordExists:  "record already exists",
	ERNoRecordFound: "record does not exist",
}
