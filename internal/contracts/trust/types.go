// Package trust defines the versioned data contracts shared across Forge trust boundaries.
// It owns no policy, context, approval, evidence, or delivery state.
package trust

import "time"

const (
	PolicyDecisionSchemaVersion    = "forge.policy-decision.v1"
	ExecutionEnvelopeSchemaVersion = "forge.execution-envelope.v1"
	ContextPackageSchemaVersion    = "forge.context-package.v1"
	ApprovalRecordSchemaVersion    = "forge.approval-record.v1"
	EvidenceRecordSchemaVersion    = "forge.evidence-record.v1"
)

type PolicyDecision struct {
	SchemaVersion         string           `json:"schemaVersion"`
	ID                    string           `json:"id"`
	DecisionPoint         string           `json:"decisionPoint"`
	SubjectID             string           `json:"subjectId"`
	Action                string           `json:"action"`
	Resource              string           `json:"resource"`
	EvaluatedAt           time.Time        `json:"evaluatedAt"`
	ExpiresAt             *time.Time       `json:"expiresAt,omitempty"`
	PolicyVersions        []string         `json:"policyVersions"`
	Outcome               string           `json:"outcome"`
	Obligations           []map[string]any `json:"obligations"`
	ReasonCodes           []string         `json:"reasonCodes,omitempty"`
	EvidenceCorrelationID string           `json:"evidenceCorrelationId"`
}

type ExecutionEnvelope struct {
	SchemaVersion         string            `json:"schemaVersion"`
	ID                    string            `json:"id"`
	ParentID              string            `json:"parentId,omitempty"`
	Issuer                string            `json:"issuer"`
	IssuedAt              time.Time         `json:"issuedAt"`
	ExpiresAt             time.Time         `json:"expiresAt"`
	Scope                 ExecutionScope    `json:"scope"`
	Principal             EnvelopePrincipal `json:"principal"`
	Purpose               string            `json:"purpose"`
	RiskClass             string            `json:"riskClass,omitempty"`
	Actions               []string          `json:"actions"`
	Classifications       []string          `json:"classifications,omitempty"`
	Destinations          []string          `json:"destinations,omitempty"`
	ModelIDs              []string          `json:"modelIds,omitempty"`
	ToolIDs               []string          `json:"toolIds,omitempty"`
	SecretReferenceIDs    []string          `json:"secretReferenceIds,omitempty"`
	PolicyDecisionIDs     []string          `json:"policyDecisionIds"`
	ApprovalRecordIDs     []string          `json:"approvalRecordIds,omitempty"`
	Limits                ExecutionLimits   `json:"limits"`
	IdempotencyKey        string            `json:"idempotencyKey,omitempty"`
	EvidenceCorrelationID string            `json:"evidenceCorrelationId"`
	Integrity             *Integrity        `json:"integrity,omitempty"`
}

type ExecutionScope struct {
	TenantID      string `json:"tenantId"`
	WorkspaceID   string `json:"workspaceId"`
	ProjectID     string `json:"projectId"`
	DeliveryRunID string `json:"deliveryRunId"`
}

type EnvelopePrincipal struct {
	SubjectID       string   `json:"subjectId"`
	AgentID         string   `json:"agentId"`
	DelegationChain []string `json:"delegationChain"`
}

type ExecutionLimits struct {
	MaxDurationSeconds int64 `json:"maxDurationSeconds"`
	MaxCostMinorUnits  int64 `json:"maxCostMinorUnits"`
	MaxInputTokens     int64 `json:"maxInputTokens,omitempty"`
	MaxOutputTokens    int64 `json:"maxOutputTokens,omitempty"`
	MaxConcurrency     int64 `json:"maxConcurrency"`
}

type Integrity struct {
	Algorithm          string `json:"algorithm"`
	Digest             string `json:"digest"`
	SignatureReference string `json:"signatureReference,omitempty"`
}

type ContextPackage struct {
	SchemaVersion       string        `json:"schemaVersion"`
	ID                  string        `json:"id"`
	ExecutionEnvelopeID string        `json:"executionEnvelopeId"`
	Objective           string        `json:"objective"`
	AssembledAt         time.Time     `json:"assembledAt"`
	ExpiresAt           *time.Time    `json:"expiresAt,omitempty"`
	PolicyDecisionIDs   []string      `json:"policyDecisionIds"`
	Items               []ContextItem `json:"items"`
	Omissions           []string      `json:"omissions,omitempty"`
	Conflicts           []string      `json:"conflicts,omitempty"`
	ManifestDigest      string        `json:"manifestDigest"`
}

type ContextItem struct {
	SourceID       string     `json:"sourceId"`
	SourceRevision string     `json:"sourceRevision"`
	Locator        string     `json:"locator"`
	Classification string     `json:"classification"`
	Authority      string     `json:"authority"`
	Purpose        string     `json:"purpose"`
	FreshnessAt    *time.Time `json:"freshnessAt,omitempty"`
	Confidence     *float64   `json:"confidence,omitempty"`
	ContentDigest  string     `json:"contentDigest"`
}

type ApprovalRecord struct {
	SchemaVersion         string    `json:"schemaVersion"`
	ID                    string    `json:"id"`
	ApproverID            string    `json:"approverId"`
	DelegationReference   string    `json:"delegationReference,omitempty"`
	Action                string    `json:"action"`
	SubjectType           string    `json:"subjectType"`
	SubjectDigest         string    `json:"subjectDigest"`
	Decision              string    `json:"decision"`
	DecidedAt             time.Time `json:"decidedAt"`
	ExpiresAt             time.Time `json:"expiresAt"`
	Conditions            []string  `json:"conditions,omitempty"`
	EvidenceCorrelationID string    `json:"evidenceCorrelationId"`
}

type EvidenceRecord struct {
	SchemaVersion  string            `json:"schemaVersion"`
	ID             string            `json:"id"`
	Type           string            `json:"type"`
	RecordedAt     time.Time         `json:"recordedAt"`
	ActorID        string            `json:"actorId"`
	Scope          EvidenceScope     `json:"scope"`
	Subject        EvidenceSubject   `json:"subject"`
	References     []string          `json:"references"`
	Outcome        string            `json:"outcome,omitempty"`
	Classification string            `json:"classification"`
	RetentionClass string            `json:"retentionClass"`
	Integrity      EvidenceIntegrity `json:"integrity"`
}

type EvidenceScope struct {
	TenantID      string `json:"tenantId"`
	WorkspaceID   string `json:"workspaceId"`
	ProjectID     string `json:"projectId"`
	DeliveryRunID string `json:"deliveryRunId,omitempty"`
}

type EvidenceSubject struct {
	Type   string `json:"type"`
	ID     string `json:"id"`
	Digest string `json:"digest"`
}

type EvidenceIntegrity struct {
	Algorithm            string `json:"algorithm"`
	Digest               string `json:"digest"`
	PreviousRecordDigest string `json:"previousRecordDigest,omitempty"`
	SignatureReference   string `json:"signatureReference,omitempty"`
}
