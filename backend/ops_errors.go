package main

import (
	"errors"
	"fmt"
)

var (
	ErrOpsNotFound   = errors.New("operations record not found")
	ErrOpsConflict   = errors.New("operations revision conflict")
	ErrOpsInvalid    = errors.New("operations request is invalid")
	ErrOpsTransition = errors.New("operations status transition is not allowed")
	ErrOpsPolicy     = errors.New("operations policy rejected the request")
)

type OpsError struct {
	Code      string
	Operation string
	Cause     error
}

func (e *OpsError) Error() string {
	if e.Cause == nil {
		return e.Code + ": " + e.Operation
	}
	return fmt.Sprintf("%s: %s: %v", e.Code, e.Operation, e.Cause)
}
func (e *OpsError) Unwrap() error { return e.Cause }
func wrapOps(code, operation string, cause error) error {
	return &OpsError{Code: code, Operation: operation, Cause: cause}
}
func opsCode(err error) string {
	var typed *OpsError
	if errors.As(err, &typed) {
		if c := opsCodeFromCause(typed); c != "" {
			return c
		}
		return typed.Code
	}
	return opsCodeFromSentinels(err)
}

// opsCodeFromCause classifies an OpsError by first consulting its sentinel
// cause (so a wrapped "not found" stays "not_found") and only falling back to
// the OpsError.Code when the cause carries no recognised sentinel.
func opsCodeFromCause(typed *OpsError) string {
	if typed.Cause != nil {
		if c := opsCodeFromSentinels(typed.Cause); c != "" {
			return c
		}
	}
	switch typed.Code {
	case "not_found", "conflict", "invalid", "transition", "policy":
		return typed.Code
	}
	return ""
}

// opsCodeFromSentinels maps a wrapped sentinel error to its canonical code,
// returning "" when no sentinel is recognised (so the caller can fall back).
func opsCodeFromSentinels(err error) string {
	switch {
	case err == nil:
		return ""
	case errors.Is(err, ErrOpsNotFound):
		return "not_found"
	case errors.Is(err, ErrOpsConflict):
		return "conflict"
	case errors.Is(err, ErrOpsInvalid):
		return "invalid"
	case errors.Is(err, ErrOpsTransition):
		return "transition"
	case errors.Is(err, ErrOpsPolicy):
		return "policy"
	default:
		return "internal"
	}
}
func opsIsNotFound(err error) bool   { return errors.Is(err, ErrOpsNotFound) }
func opsIsConflict(err error) bool   { return errors.Is(err, ErrOpsConflict) }
func opsIsInvalid(err error) bool    { return errors.Is(err, ErrOpsInvalid) }
func opsIsTransition(err error) bool { return errors.Is(err, ErrOpsTransition) }
func opsIsPolicy(err error) bool     { return errors.Is(err, ErrOpsPolicy) }
