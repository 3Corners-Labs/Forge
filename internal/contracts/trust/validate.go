package trust

import (
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
)

func (value PolicyDecision) Validate(correlationID string) error {
	if value.SchemaVersion != PolicyDecisionSchemaVersion {
		return invalid("policy decision has an unsupported schema version", correlationID)
	}
	for _, item := range []string{value.ID, value.DecisionPoint, value.SubjectID, value.Action, value.Resource, value.EvidenceCorrelationID} {
		if item == "" {
			return invalid("policy decision has an empty required field", correlationID)
		}
	}
	if value.EvaluatedAt.IsZero() {
		return invalid("policy decision evaluatedAt is required", correlationID)
	}
	if value.ExpiresAt != nil && !value.ExpiresAt.After(value.EvaluatedAt) {
		return invalid("policy decision expiry must follow evaluation", correlationID)
	}
	if !oneOf(value.DecisionPoint, "retrieval", "secret-resolution", "selection", "execution", "evidence-disclosure", "publication", "release") {
		return invalid("policy decision point is invalid", correlationID)
	}
	if !oneOf(value.Outcome, "permit", "deny", "indeterminate") {
		return invalid("policy decision outcome is invalid", correlationID)
	}
	if err := validateUniqueRequired(value.PolicyVersions, "policy versions", correlationID); err != nil {
		return err
	}
	if value.Obligations == nil {
		return invalid("policy decision obligations are required", correlationID)
	}
	for _, obligation := range value.Obligations {
		typeName, ok := obligation["type"].(string)
		if !ok || typeName == "" {
			return invalid("policy decision obligation type is required", correlationID)
		}
	}
	return validateNonEmpty(value.ReasonCodes, "policy reason codes", correlationID)
}

// Authorize verifies that a decision grants an exact action and that every
// mandatory obligation has already been enforced by the caller. Optional
// advice is not represented in the Policy Decision obligation list.
func (value PolicyDecision) Authorize(decisionPoint, subjectID, action, resource string, at time.Time, enforcedObligationTypes []string, correlationID string) error {
	if err := value.Validate(correlationID); err != nil {
		return err
	}
	if at.IsZero() {
		return invalid("policy authorization time is required", correlationID)
	}
	if at.Before(value.EvaluatedAt) || (value.ExpiresAt != nil && !at.Before(*value.ExpiresAt)) {
		return failure.New(failure.Expired, "policy decision is not active", correlationID, failure.RetryNewAuthority)
	}
	if value.DecisionPoint != decisionPoint || value.SubjectID != subjectID || value.Action != action || value.Resource != resource {
		return failure.New(failure.Forbidden, "policy decision does not authorize the requested action", correlationID, failure.RetryNone)
	}
	switch value.Outcome {
	case "permit":
		for _, obligation := range value.Obligations {
			if !contains(enforcedObligationTypes, obligation["type"].(string)) {
				return failure.New(failure.Forbidden, "policy obligation has not been enforced", correlationID, failure.RetryNone)
			}
		}
		return nil
	case "deny":
		return failure.New(failure.Forbidden, "policy denied the requested action", correlationID, failure.RetryNone)
	default:
		return failure.New(failure.Indeterminate, "policy could not determine authority", correlationID, failure.RetryNone)
	}
}

func (value ExecutionEnvelope) Validate(correlationID string) error {
	if value.SchemaVersion != ExecutionEnvelopeSchemaVersion {
		return invalid("execution envelope has an unsupported schema version", correlationID)
	}
	for _, item := range []string{
		value.ID, value.Issuer, value.Scope.TenantID, value.Scope.WorkspaceID,
		value.Scope.ProjectID, value.Scope.DeliveryRunID, value.Principal.SubjectID,
		value.Principal.AgentID, value.Purpose, value.EvidenceCorrelationID,
	} {
		if item == "" {
			return invalid("execution envelope has an empty required field", correlationID)
		}
	}
	if value.IssuedAt.IsZero() || value.ExpiresAt.IsZero() || !value.ExpiresAt.After(value.IssuedAt) {
		return invalid("execution envelope has an invalid validity window", correlationID)
	}
	if value.RiskClass != "" && !oneOf(value.RiskClass, "low", "moderate", "high", "critical") {
		return invalid("execution envelope risk class is invalid", correlationID)
	}
	if err := validateUniqueRequired(value.Actions, "execution actions", correlationID); err != nil {
		return err
	}
	if err := validateUniqueRequired(value.PolicyDecisionIDs, "execution policy decisions", correlationID); err != nil {
		return err
	}
	for _, field := range []struct {
		label string
		items []string
	}{
		{"delegation chain", value.Principal.DelegationChain},
		{"classifications", value.Classifications},
		{"destinations", value.Destinations},
		{"model IDs", value.ModelIDs},
		{"tool IDs", value.ToolIDs},
		{"secret reference IDs", value.SecretReferenceIDs},
		{"approval record IDs", value.ApprovalRecordIDs},
	} {
		if err := validateUnique(field.items, field.label, correlationID); err != nil {
			return err
		}
	}
	if value.Principal.DelegationChain == nil {
		return invalid("execution envelope delegation chain is required", correlationID)
	}
	if value.Limits.MaxDurationSeconds < 1 || value.Limits.MaxCostMinorUnits < 0 || value.Limits.MaxInputTokens < 0 || value.Limits.MaxOutputTokens < 0 || value.Limits.MaxConcurrency < 1 {
		return invalid("execution envelope limits are invalid", correlationID)
	}
	if value.Integrity != nil {
		if err := validateIntegrity(*value.Integrity, correlationID); err != nil {
			return err
		}
	}
	return nil
}

func (value ExecutionEnvelope) AuthorizeAction(action string, at time.Time, correlationID string) error {
	if value.Integrity != nil {
		if err := value.VerifyIntegrity(correlationID); err != nil {
			return err
		}
	} else if err := value.Validate(correlationID); err != nil {
		return err
	}
	if at.IsZero() {
		return invalid("execution authorization time is required", correlationID)
	}
	if at.Before(value.IssuedAt) || !at.Before(value.ExpiresAt) {
		return failure.New(failure.Expired, "execution envelope is not active", correlationID, failure.RetryNewAuthority)
	}
	if !contains(value.Actions, action) {
		return failure.New(failure.Forbidden, "execution envelope does not grant the requested action", correlationID, failure.RetryNone)
	}
	return nil
}

func (value ExecutionEnvelope) ValidateChildOf(parent ExecutionEnvelope, at time.Time, correlationID string) error {
	if err := value.Validate(correlationID); err != nil {
		return err
	}
	if err := parent.Validate(correlationID); err != nil {
		return err
	}
	if value.ParentID != parent.ID {
		return failure.New(failure.Forbidden, "child envelope does not reference the parent", correlationID, failure.RetryNone)
	}
	if value.Scope != parent.Scope || value.Principal.SubjectID != parent.Principal.SubjectID {
		return failure.New(failure.Forbidden, "child envelope changes authoritative scope", correlationID, failure.RetryNone)
	}
	if value.Purpose != parent.Purpose || value.RiskClass != parent.RiskClass {
		return failure.New(failure.Forbidden, "child envelope changes parent purpose or risk class", correlationID, failure.RetryNone)
	}
	if value.IssuedAt.Before(parent.IssuedAt) || value.ExpiresAt.After(parent.ExpiresAt) {
		return failure.New(failure.Forbidden, "child envelope expands the parent validity window", correlationID, failure.RetryNone)
	}
	strictlyNarrower := value.IssuedAt.After(parent.IssuedAt) || value.ExpiresAt.Before(parent.ExpiresAt)
	for _, field := range []struct {
		label  string
		child  []string
		parent []string
	}{
		{"actions", value.Actions, parent.Actions},
		{"classifications", value.Classifications, parent.Classifications},
		{"destinations", value.Destinations, parent.Destinations},
		{"models", value.ModelIDs, parent.ModelIDs},
		{"tools", value.ToolIDs, parent.ToolIDs},
		{"secret references", value.SecretReferenceIDs, parent.SecretReferenceIDs},
	} {
		if !subset(field.child, field.parent) {
			return failure.New(failure.Forbidden, "child envelope expands parent "+field.label, correlationID, failure.RetryNone)
		}
		if len(field.child) < len(field.parent) {
			strictlyNarrower = true
		}
	}
	if value.Limits.MaxDurationSeconds > parent.Limits.MaxDurationSeconds ||
		value.Limits.MaxCostMinorUnits > parent.Limits.MaxCostMinorUnits ||
		value.Limits.MaxInputTokens > parent.Limits.MaxInputTokens ||
		value.Limits.MaxOutputTokens > parent.Limits.MaxOutputTokens ||
		value.Limits.MaxConcurrency > parent.Limits.MaxConcurrency {
		return failure.New(failure.Forbidden, "child envelope expands parent limits", correlationID, failure.RetryNone)
	}
	if value.Limits.MaxDurationSeconds < parent.Limits.MaxDurationSeconds ||
		value.Limits.MaxCostMinorUnits < parent.Limits.MaxCostMinorUnits ||
		value.Limits.MaxInputTokens < parent.Limits.MaxInputTokens ||
		value.Limits.MaxOutputTokens < parent.Limits.MaxOutputTokens ||
		value.Limits.MaxConcurrency < parent.Limits.MaxConcurrency {
		strictlyNarrower = true
	}
	if !strictlyNarrower {
		return failure.New(failure.Forbidden, "child envelope is not strictly narrower than its parent", correlationID, failure.RetryNone)
	}
	if at.IsZero() || at.Before(value.IssuedAt) || !at.Before(value.ExpiresAt) || at.Before(parent.IssuedAt) || !at.Before(parent.ExpiresAt) {
		return failure.New(failure.Expired, "parent or child envelope is not active", correlationID, failure.RetryNewAuthority)
	}
	return nil
}

func (value ContextPackage) Validate(correlationID string) error {
	if value.SchemaVersion != ContextPackageSchemaVersion {
		return invalid("context package has an unsupported schema version", correlationID)
	}
	for _, item := range []string{value.ID, value.ExecutionEnvelopeID, value.Objective, value.ManifestDigest} {
		if item == "" {
			return invalid("context package has an empty required field", correlationID)
		}
	}
	if value.AssembledAt.IsZero() {
		return invalid("context package assembledAt is required", correlationID)
	}
	if value.ExpiresAt != nil && !value.ExpiresAt.After(value.AssembledAt) {
		return invalid("context package expiry must follow assembly", correlationID)
	}
	if err := validateUniqueRequired(value.PolicyDecisionIDs, "context policy decisions", correlationID); err != nil {
		return err
	}
	if value.Items == nil {
		return invalid("context package items are required", correlationID)
	}
	for _, item := range value.Items {
		for _, required := range []string{item.SourceID, item.SourceRevision, item.Locator, item.Classification, item.Authority, item.Purpose, item.ContentDigest} {
			if required == "" {
				return invalid("context item has an empty required field", correlationID)
			}
		}
		if !oneOf(item.Authority, "none", "reference", "accepted-knowledge", "task-instruction") {
			return invalid("context item authority is invalid", correlationID)
		}
		if item.FreshnessAt != nil && item.FreshnessAt.IsZero() {
			return invalid("context item freshness is invalid", correlationID)
		}
		if item.Confidence != nil && (*item.Confidence < 0 || *item.Confidence > 1) {
			return invalid("context item confidence is invalid", correlationID)
		}
	}
	if err := validateNonEmpty(value.Omissions, "context omissions", correlationID); err != nil {
		return err
	}
	return validateNonEmpty(value.Conflicts, "context conflicts", correlationID)
}

func (value ContextPackage) ValidateFor(envelope ExecutionEnvelope, at time.Time, correlationID string) error {
	if err := value.Validate(correlationID); err != nil {
		return err
	}
	if err := envelope.Validate(correlationID); err != nil {
		return err
	}
	if value.ExecutionEnvelopeID != envelope.ID {
		return failure.New(failure.Forbidden, "context package belongs to another execution envelope", correlationID, failure.RetryNone)
	}
	if value.AssembledAt.Before(envelope.IssuedAt) || !value.AssembledAt.Before(envelope.ExpiresAt) {
		return failure.New(failure.Forbidden, "context package was assembled outside envelope validity", correlationID, failure.RetryNone)
	}
	if value.ExpiresAt != nil && value.ExpiresAt.After(envelope.ExpiresAt) {
		return failure.New(failure.Forbidden, "context package outlives its execution envelope", correlationID, failure.RetryNone)
	}
	if at.IsZero() || at.Before(value.AssembledAt) || !at.Before(envelope.ExpiresAt) || (value.ExpiresAt != nil && !at.Before(*value.ExpiresAt)) {
		return failure.New(failure.Expired, "context authority is not active", correlationID, failure.RetryNewAuthority)
	}
	return nil
}

func (value ApprovalRecord) Validate(correlationID string) error {
	if value.SchemaVersion != ApprovalRecordSchemaVersion {
		return invalid("approval record has an unsupported schema version", correlationID)
	}
	for _, item := range []string{value.ID, value.ApproverID, value.Action, value.SubjectType, value.SubjectDigest, value.Decision, value.EvidenceCorrelationID} {
		if item == "" {
			return invalid("approval record has an empty required field", correlationID)
		}
	}
	if !oneOf(value.Decision, "approved", "rejected", "revoked") {
		return invalid("approval decision is invalid", correlationID)
	}
	if value.DecidedAt.IsZero() || value.ExpiresAt.IsZero() || !value.ExpiresAt.After(value.DecidedAt) {
		return invalid("approval record has an invalid validity window", correlationID)
	}
	return validateNonEmpty(value.Conditions, "approval conditions", correlationID)
}

func (value ApprovalRecord) Authorize(action, subjectType, subjectDigest string, at time.Time, correlationID string) error {
	if err := value.Validate(correlationID); err != nil {
		return err
	}
	if at.IsZero() || at.Before(value.DecidedAt) || !at.Before(value.ExpiresAt) {
		return failure.New(failure.Expired, "approval is not active", correlationID, failure.RetryNewAuthority)
	}
	if value.Decision != "approved" {
		return failure.New(failure.Forbidden, "approval does not permit the action", correlationID, failure.RetryNone)
	}
	if value.Action != action || value.SubjectType != subjectType {
		return failure.New(failure.Forbidden, "approval is bound to another action or subject", correlationID, failure.RetryNone)
	}
	if value.SubjectDigest != subjectDigest {
		return failure.New(failure.Conflict, "approval subject digest does not match", correlationID, failure.RetryReconcile)
	}
	return nil
}

func (value EvidenceRecord) Validate(correlationID string) error {
	if value.SchemaVersion != EvidenceRecordSchemaVersion {
		return invalid("evidence record has an unsupported schema version", correlationID)
	}
	for _, item := range []string{
		value.ID, value.Type, value.ActorID, value.Scope.TenantID, value.Scope.WorkspaceID,
		value.Scope.ProjectID, value.Subject.Type, value.Subject.ID, value.Subject.Digest,
		value.Classification, value.RetentionClass,
	} {
		if item == "" {
			return invalid("evidence record has an empty required field", correlationID)
		}
	}
	if value.RecordedAt.IsZero() {
		return invalid("evidence record recordedAt is required", correlationID)
	}
	if value.References == nil {
		return invalid("evidence references are required", correlationID)
	}
	if err := validateNonEmpty(value.References, "evidence references", correlationID); err != nil {
		return err
	}
	if value.Integrity.Algorithm == "" || value.Integrity.Digest == "" {
		return invalid("evidence integrity algorithm and digest are required", correlationID)
	}
	return nil
}

func (value EvidenceRecord) ValidateSubject(subjectType, subjectID, subjectDigest, correlationID string) error {
	if err := value.Validate(correlationID); err != nil {
		return err
	}
	if value.Subject.Type != subjectType || value.Subject.ID != subjectID || value.Subject.Digest != subjectDigest {
		return failure.New(failure.VerificationFailed, "evidence subject does not match", correlationID, failure.RetryNone)
	}
	return nil
}

func validateIntegrity(value Integrity, correlationID string) error {
	if value.Algorithm == "" || value.Digest == "" {
		return invalid("integrity algorithm and digest are required", correlationID)
	}
	return nil
}

func validateUniqueRequired(items []string, label, correlationID string) error {
	if len(items) == 0 {
		return invalid(label+" are required", correlationID)
	}
	return validateUnique(items, label, correlationID)
}

func validateUnique(items []string, label, correlationID string) error {
	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		if item == "" {
			return invalid(label+" contain an empty value", correlationID)
		}
		if _, exists := seen[item]; exists {
			return invalid(label+" contain a duplicate value", correlationID)
		}
		seen[item] = struct{}{}
	}
	return nil
}

func validateNonEmpty(items []string, label, correlationID string) error {
	for _, item := range items {
		if item == "" {
			return invalid(label+" contain an empty value", correlationID)
		}
	}
	return nil
}

func invalid(message, correlationID string) error {
	return failure.New(failure.InvalidRequest, message, correlationID, failure.RetryNone)
}

func oneOf(value string, allowed ...string) bool {
	return contains(allowed, value)
}

func contains(items []string, wanted string) bool {
	for _, item := range items {
		if item == wanted {
			return true
		}
	}
	return false
}

func subset(child, parent []string) bool {
	for _, item := range child {
		if !contains(parent, item) {
			return false
		}
	}
	return true
}
