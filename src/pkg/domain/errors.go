package domain

import "fmt"

// define custom errors
var (
	ErrSimpleExample = fmt.Errorf("just a simple example error without details")
)

type ErrorPlaceholderMissing struct {
	MissingPlaceholder string
}

func (e ErrorPlaceholderMissing) Error() string {
	return fmt.Sprintf("placeholder is missing: %s", e.MissingPlaceholder)
}
