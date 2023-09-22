package err

import (
	"fmt"
	"strings"
)

// Collection of errors that should not be modified
type ErrorCollection struct {
	Name   string
	Errors []error
}

func (e *ErrorCollection) Error() string {
	var sb strings.Builder
	// lists the errors
	sb.WriteString(fmt.Sprintf("%v had %v errors:", e.Name, len(e.Errors)))
	for _, e := range e.Errors {
		sb.WriteString(fmt.Sprintf("\n\t'%v'", e))
	}
	return sb.String()
}

// Used to catch many errors without making a if-statement for every error
// still a bit of a hazzle for multiple return values
type ErrorCatcher struct {
	ErrorCollection
}

func (e *ErrorCatcher) Catch(err error) {
	if err != nil {
		e.Errors = append(e.Errors, err)
	}
}

func (e *ErrorCatcher) Flush() *ErrorCollection {
	if len(e.Errors) == 0 {
		return nil
	}
	// Copy errors, because original list will be emptied
	reslst := make([]error, len(e.Errors))
	copy(reslst, e.Errors)

	return &ErrorCollection{e.Name, reslst}
}

func Catcher(Name string) *ErrorCatcher {
	return &ErrorCatcher{ErrorCollection{Name, make([]error, 0)}}
}
