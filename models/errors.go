package models

import (
	"fmt"
)

func NewError(errType ErrorType, msg string) *Error {
	return &Error{
		Type:    errType,
		Message: msg,
	}
}

func ConvertError(err error) *Error {
	if err == nil {
		return nil
	}

	modelErr, ok := err.(*Error)
	if !ok {
		modelErr = NewError(ErrorTypeUnknownError, err.Error())
	}
	return modelErr
}

func (err *Error) ToError() error {
	if err == nil {
		return nil
	}
	return err
}

func (err *Error) Equal(other error) bool {
	if e, ok := other.(*Error); ok {
		if err == nil && e != nil {
			return false
		}
		return e.Type == err.Type
	}
	return false
}

func (err *Error) Error() string {
	return err.Message
}

var (
	ErrResourceNotFound = &Error{
		Type:    ErrorTypeResourceNotFound,
		Message: "the requested resource could not be found",
	}

	ErrResourceExists = &Error{
		Type:    ErrorTypeResourceExist,
		Message: "the requested resource already exists",
	}

	ErrResourceConflict = &Error{
		Type:    ErrorTypeResourceConflict,
		Message: "the requested resource is in a conflicting state",
	}

	ErrDeadlock = &Error{
		Type:    ErrorTypeDeadlock,
		Message: "the request failed due to deadlock",
	}

	ErrBadRequest = &Error{
		Type:    ErrorTypeInvalidRequest,
		Message: "the request received is invalid",
	}

	ErrUnknownError = &Error{
		Type:    ErrorTypeUnknownError,
		Message: "the request failed for an unknown reason",
	}

	ErrDeserialize = &Error{
		Type:    ErrorTypeDeserialize,
		Message: "could not deserialize record",
	}

	ErrFailedToOpenEnvelope = &Error{
		Type:    ErrorTypeFailedToOpenEnvelope,
		Message: "could not open envelope",
	}

	ErrGUIDGeneration = &Error{
		Type:    ErrorTypeGUIDGeneration,
		Message: "cannot generate random guid",
	}
)

type ErrInvalidField struct {
	Field string
}

func (err ErrInvalidField) Error() string {
	return "Invalid field: " + err.Field
}

type ErrInvalidModification struct {
	InvalidField string
}

func (err ErrInvalidModification) Error() string {
	return "attempt to make invalid change to field: " + err.InvalidField
}

func NewTaskTransitionError(from, to State) *Error {
	return &Error{
		Type:    ErrorTypeInvalidStateTransition,
		Message: fmt.Sprintf("Cannot transition from %v to %v", from, to),
	}
}

func NewUnrecoverableError(err error) *Error {
	return &Error{
		Type:    ErrorTypeUnrecoverable,
		Message: fmt.Sprint("Unrecoverable Error: ", err),
	}
}
