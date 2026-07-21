// Package failure defines the disclosure-safe error contract shared by Forge capabilities.
package failure

import (
	"errors"
	"fmt"
)

const SchemaVersion = "forge.normalized-error.v1"

// Code is a stable machine-readable failure category.
type Code string

const (
	InvalidRequest         Code = "invalid_request"
	Unauthenticated        Code = "unauthenticated"
	Forbidden              Code = "forbidden"
	Indeterminate          Code = "indeterminate"
	Conflict               Code = "conflict"
	Expired                Code = "expired"
	BudgetExhausted        Code = "budget_exhausted"
	DependencyUnavailable  Code = "dependency_unavailable"
	VerificationFailed     Code = "verification_failed"
	EvidenceUnavailable    Code = "evidence_unavailable"
	ManualRecoveryRequired Code = "manual_recovery_required"
)

// Retry classifies the only safe next action for a failure.
type Retry string

const (
	RetryNone           Retry = "none"
	RetryReauthenticate Retry = "reauthenticate"
	RetryReconcile      Retry = "reconcile"
	RetryNewAuthority   Retry = "new-authority"
	RetryBounded        Retry = "bounded"
	RetryManual         Retry = "manual"
)

// Error is safe to return across a public capability boundary. Details must
// already have passed disclosure policy before they are attached.
type Error struct {
	SchemaVersion string         `json:"schemaVersion"`
	Code          Code           `json:"code"`
	SafeMessage   string         `json:"safeMessage"`
	CorrelationID string         `json:"correlationId"`
	Retry         Retry          `json:"retry"`
	Details       map[string]any `json:"details,omitempty"`
}

func New(code Code, safeMessage, correlationID string, retry Retry) *Error {
	if correlationID == "" {
		correlationID = "forge.validation.uncorrelated"
	}
	return &Error{
		SchemaVersion: SchemaVersion,
		Code:          code,
		SafeMessage:   safeMessage,
		CorrelationID: correlationID,
		Retry:         retry,
	}
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("%s: %s", e.Code, e.SafeMessage)
}

// Validate checks the normalized error itself. It never inspects Details,
// whose disclosure safety is the responsibility of the producing capability.
func (e *Error) Validate() error {
	if e == nil {
		return errors.New("normalized error is nil")
	}
	if e.SchemaVersion != SchemaVersion {
		return fmt.Errorf("unsupported normalized error schema version %q", e.SchemaVersion)
	}
	if !validCode(e.Code) {
		return fmt.Errorf("unknown normalized error code %q", e.Code)
	}
	if e.SafeMessage == "" {
		return errors.New("normalized error safe message is required")
	}
	if e.CorrelationID == "" {
		return errors.New("normalized error correlation ID is required")
	}
	if !validRetry(e.Retry) {
		return fmt.Errorf("unknown normalized error retry class %q", e.Retry)
	}
	return nil
}

func validCode(code Code) bool {
	switch code {
	case InvalidRequest, Unauthenticated, Forbidden, Indeterminate, Conflict,
		Expired, BudgetExhausted, DependencyUnavailable, VerificationFailed,
		EvidenceUnavailable, ManualRecoveryRequired:
		return true
	default:
		return false
	}
}

func validRetry(retry Retry) bool {
	switch retry {
	case RetryNone, RetryReauthenticate, RetryReconcile, RetryNewAuthority,
		RetryBounded, RetryManual:
		return true
	default:
		return false
	}
}
