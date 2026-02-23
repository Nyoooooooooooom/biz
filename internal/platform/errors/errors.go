package errors

import "fmt"

type Kind string

const (
	KindValidation            Kind = "VALIDATION_ERROR"
	KindNotFound              Kind = "NOT_FOUND"
	KindDependencyUnavailable Kind = "DEPENDENCY_UNAVAILABLE"
	KindConflict              Kind = "CONFLICT"
	KindInternal              Kind = "INTERNAL_ERROR"
)

type Error struct {
	Kind    Kind
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause == nil {
		return e.Message
	}
	return fmt.Sprintf("%s: %v", e.Message, e.Cause)
}

func (e *Error) Unwrap() error { return e.Cause }

func New(kind Kind, message string) *Error {
	return &Error{Kind: kind, Message: message}
}

func Wrap(kind Kind, message string, cause error) *Error {
	return &Error{Kind: kind, Message: message, Cause: cause}
}

func KindOf(err error) Kind {
	if err == nil {
		return ""
	}
	var te *Error
	if ok := As(err, &te); ok {
		return te.Kind
	}
	return KindInternal
}

func As(err error, target interface{}) bool {
	switch t := target.(type) {
	case **Error:
		for err != nil {
			e, ok := err.(*Error)
			if ok {
				*t = e
				return true
			}
			u, ok := err.(interface{ Unwrap() error })
			if !ok {
				return false
			}
			err = u.Unwrap()
		}
	}
	return false
}

func ExitCode(err error) int {
	switch KindOf(err) {
	case KindValidation, KindNotFound:
		return 2
	case KindDependencyUnavailable:
		return 3
	case KindConflict, KindInternal:
		return 4
	default:
		if err == nil {
			return 0
		}
		return 4
	}
}
