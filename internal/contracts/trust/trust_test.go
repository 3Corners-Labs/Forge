package trust

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
)

const testCorrelationID = "correlation-1"

var testNow = time.Date(2026, 7, 21, 10, 0, 0, 0, time.UTC)

func TestStrictDecodersAcceptValidContracts(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		value  any
		decode func([]byte, string) error
	}{
		{"policy decision", validPolicyDecision(), func(data []byte, id string) error { _, err := DecodePolicyDecision(data, id); return err }},
		{"execution envelope", validEnvelope(), func(data []byte, id string) error { _, err := DecodeExecutionEnvelope(data, id); return err }},
		{"context package", validContextPackage(), func(data []byte, id string) error { _, err := DecodeContextPackage(data, id); return err }},
		{"approval record", validApprovalRecord(), func(data []byte, id string) error { _, err := DecodeApprovalRecord(data, id); return err }},
		{"evidence record", validEvidenceRecord(), func(data []byte, id string) error { _, err := DecodeEvidenceRecord(data, id); return err }},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			data, err := json.Marshal(test.value)
			if err != nil {
				t.Fatal(err)
			}
			if err := test.decode(data, testCorrelationID); err != nil {
				t.Fatalf("decode error = %v", err)
			}
		})
	}
}

func TestStrictDecoderRejectsUnknownAuthorityAndTrailingData(t *testing.T) {
	t.Parallel()
	contextWithAuthority := []byte(`{
		"schemaVersion":"forge.context-package.v1","id":"context-1",
		"executionEnvelopeId":"envelope-1","objective":"prepare a bounded change",
		"assembledAt":"2026-07-21T10:01:00Z","policyDecisionIds":["policy-1"],
		"items":[{"sourceId":"source-1","sourceRevision":"revision-1","locator":"docs/input.md","classification":"internal","authority":"reference","purpose":"task context","contentDigest":"sha256:context","actions":["publish"]}],
		"manifestDigest":"sha256:manifest"
	}`)
	_, err := DecodeContextPackage(contextWithAuthority, testCorrelationID)
	assertFailureCode(t, err, failure.InvalidRequest)

	policyData, err := json.Marshal(validPolicyDecision())
	if err != nil {
		t.Fatal(err)
	}
	policyData = append(policyData, []byte(` {}`)...)
	_, err = DecodePolicyDecision(policyData, testCorrelationID)
	assertFailureCode(t, err, failure.InvalidRequest)
}

func TestPolicyDecisionAuthorizationIsExactAndFailClosed(t *testing.T) {
	t.Parallel()
	decision := validPolicyDecision()
	if err := decision.Authorize("execution", "subject-1", "prepare-change", "project-1", testNow.Add(2*time.Minute), []string{"record-evidence"}, testCorrelationID); err != nil {
		t.Fatalf("Authorize() error = %v", err)
	}
	assertFailureCode(t, decision.Authorize("execution", "subject-1", "prepare-change", "project-2", testNow.Add(2*time.Minute), []string{"record-evidence"}, testCorrelationID), failure.Forbidden)
	decision.Outcome = "deny"
	assertFailureCode(t, decision.Authorize("execution", "subject-1", "prepare-change", "project-1", testNow.Add(2*time.Minute), []string{"record-evidence"}, testCorrelationID), failure.Forbidden)
	decision.Outcome = "indeterminate"
	assertFailureCode(t, decision.Authorize("execution", "subject-1", "prepare-change", "project-1", testNow.Add(2*time.Minute), []string{"record-evidence"}, testCorrelationID), failure.Indeterminate)
	decision = validPolicyDecision()
	assertFailureCode(t, decision.Authorize("execution", "subject-1", "prepare-change", "project-1", testNow.Add(11*time.Minute), []string{"record-evidence"}, testCorrelationID), failure.Expired)
	assertFailureCode(t, decision.Authorize("execution", "subject-1", "prepare-change", "project-1", testNow.Add(-time.Second), []string{"record-evidence"}, testCorrelationID), failure.Expired)
	assertFailureCode(t, decision.Authorize("execution", "subject-1", "prepare-change", "project-1", testNow.Add(2*time.Minute), nil, testCorrelationID), failure.Forbidden)
}

func TestExecutionEnvelopeAuthorizesOnlyActiveGrantedActions(t *testing.T) {
	t.Parallel()
	envelope := validEnvelope()
	if err := envelope.AuthorizeAction("prepare-change", testNow.Add(2*time.Minute), testCorrelationID); err != nil {
		t.Fatalf("AuthorizeAction() error = %v", err)
	}
	assertFailureCode(t, envelope.AuthorizeAction("publish", testNow.Add(2*time.Minute), testCorrelationID), failure.Forbidden)
	assertFailureCode(t, envelope.AuthorizeAction("prepare-change", testNow.Add(31*time.Minute), testCorrelationID), failure.Expired)
}

func TestChildEnvelopeCannotExpandParentAuthority(t *testing.T) {
	t.Parallel()
	parent := validEnvelope()
	parent.Actions = []string{"prepare-change", "verify-change"}
	parent.Classifications = []string{"internal", "public"}
	parent.ToolIDs = []string{"git", "go"}
	child := parent
	child.ID = "envelope-child"
	child.ParentID = parent.ID
	child.IssuedAt = parent.IssuedAt.Add(time.Minute)
	child.ExpiresAt = parent.ExpiresAt.Add(-time.Minute)
	child.Actions = []string{"prepare-change"}
	child.Classifications = []string{"internal"}
	child.ToolIDs = []string{"git"}
	child.Limits.MaxDurationSeconds--
	if err := child.ValidateChildOf(parent, testNow.Add(2*time.Minute), testCorrelationID); err != nil {
		t.Fatalf("ValidateChildOf() error = %v", err)
	}
	tests := []struct {
		name   string
		mutate func(*ExecutionEnvelope)
	}{
		{"action expansion", func(value *ExecutionEnvelope) { value.Actions = append(value.Actions, "publish") }},
		{"scope change", func(value *ExecutionEnvelope) { value.Scope.ProjectID = "project-2" }},
		{"purpose change", func(value *ExecutionEnvelope) { value.Purpose = "publish a change" }},
		{"limit expansion", func(value *ExecutionEnvelope) { value.Limits.MaxConcurrency++ }},
		{"validity expansion", func(value *ExecutionEnvelope) { value.ExpiresAt = parent.ExpiresAt.Add(time.Second) }},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			candidate := child
			test.mutate(&candidate)
			assertFailureCode(t, candidate.ValidateChildOf(parent, testNow.Add(2*time.Minute), testCorrelationID), failure.Forbidden)
		})
	}
	equalChild := parent
	equalChild.ID = "envelope-equal-child"
	equalChild.ParentID = parent.ID
	assertFailureCode(t, equalChild.ValidateChildOf(parent, testNow.Add(2*time.Minute), testCorrelationID), failure.Forbidden)
}

func TestContextPackageIsBoundToOneActiveEnvelope(t *testing.T) {
	t.Parallel()
	envelope := validEnvelope()
	contextPackage := validContextPackage()
	if err := contextPackage.ValidateFor(envelope, testNow.Add(2*time.Minute), testCorrelationID); err != nil {
		t.Fatalf("ValidateFor() error = %v", err)
	}
	contextPackage.ExecutionEnvelopeID = "envelope-2"
	assertFailureCode(t, contextPackage.ValidateFor(envelope, testNow.Add(2*time.Minute), testCorrelationID), failure.Forbidden)
	contextPackage = validContextPackage()
	expiresAfterEnvelope := envelope.ExpiresAt.Add(time.Second)
	contextPackage.ExpiresAt = &expiresAfterEnvelope
	assertFailureCode(t, contextPackage.ValidateFor(envelope, testNow.Add(2*time.Minute), testCorrelationID), failure.Forbidden)
	contextPackage = validContextPackage()
	assertFailureCode(t, contextPackage.ValidateFor(envelope, testNow.Add(30*time.Second), testCorrelationID), failure.Expired)
}

func TestApprovalRequiresExactActionAndDigest(t *testing.T) {
	t.Parallel()
	approval := validApprovalRecord()
	if err := approval.Authorize("publish-pr", "candidate-change", "sha256:candidate", testNow.Add(2*time.Minute), testCorrelationID); err != nil {
		t.Fatalf("Authorize() error = %v", err)
	}
	assertFailureCode(t, approval.Authorize("publish-pr", "candidate-change", "sha256:changed", testNow.Add(2*time.Minute), testCorrelationID), failure.Conflict)
	assertFailureCode(t, approval.Authorize("release", "candidate-change", "sha256:candidate", testNow.Add(2*time.Minute), testCorrelationID), failure.Forbidden)
	approval.Decision = "revoked"
	assertFailureCode(t, approval.Authorize("publish-pr", "candidate-change", "sha256:candidate", testNow.Add(2*time.Minute), testCorrelationID), failure.Forbidden)
	approval = validApprovalRecord()
	assertFailureCode(t, approval.Authorize("publish-pr", "candidate-change", "sha256:candidate", testNow.Add(11*time.Minute), testCorrelationID), failure.Expired)
}

func TestEvidenceRequiresExactSubjectBinding(t *testing.T) {
	t.Parallel()
	record := validEvidenceRecord()
	if err := record.ValidateSubject("candidate-change", "candidate-1", "sha256:candidate", testCorrelationID); err != nil {
		t.Fatalf("ValidateSubject() error = %v", err)
	}
	assertFailureCode(t, record.ValidateSubject("candidate-change", "candidate-1", "sha256:other", testCorrelationID), failure.VerificationFailed)
}

func TestContractValidationRejectsInvalidStructure(t *testing.T) {
	t.Parallel()
	decision := validPolicyDecision()
	decision.PolicyVersions = []string{"policy-v1", "policy-v1"}
	assertFailureCode(t, decision.Validate(testCorrelationID), failure.InvalidRequest)
	envelope := validEnvelope()
	envelope.Principal.DelegationChain = nil
	assertFailureCode(t, envelope.Validate(testCorrelationID), failure.InvalidRequest)
	contextPackage := validContextPackage()
	contextPackage.Items[0].Authority = "instruction"
	assertFailureCode(t, contextPackage.Validate(testCorrelationID), failure.InvalidRequest)
}

func validPolicyDecision() PolicyDecision {
	expiresAt := testNow.Add(10 * time.Minute)
	return PolicyDecision{SchemaVersion: PolicyDecisionSchemaVersion, ID: "policy-1", DecisionPoint: "execution", SubjectID: "subject-1", Action: "prepare-change", Resource: "project-1", EvaluatedAt: testNow, ExpiresAt: &expiresAt, PolicyVersions: []string{"policy-v1"}, Outcome: "permit", Obligations: []map[string]any{{"type": "record-evidence"}}, EvidenceCorrelationID: testCorrelationID}
}

func validEnvelope() ExecutionEnvelope {
	return ExecutionEnvelope{SchemaVersion: ExecutionEnvelopeSchemaVersion, ID: "envelope-1", Issuer: "policy-service", IssuedAt: testNow, ExpiresAt: testNow.Add(30 * time.Minute), Scope: ExecutionScope{TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1", DeliveryRunID: "run-1"}, Principal: EnvelopePrincipal{SubjectID: "subject-1", AgentID: "agent-1", DelegationChain: []string{}}, Purpose: "prepare a bounded change", Actions: []string{"prepare-change"}, Classifications: []string{"internal"}, ToolIDs: []string{"git"}, PolicyDecisionIDs: []string{"policy-1"}, Limits: ExecutionLimits{MaxDurationSeconds: 600, MaxCostMinorUnits: 1000, MaxInputTokens: 1000, MaxOutputTokens: 1000, MaxConcurrency: 1}, EvidenceCorrelationID: testCorrelationID}
}

func validContextPackage() ContextPackage {
	expiresAt := testNow.Add(20 * time.Minute)
	return ContextPackage{SchemaVersion: ContextPackageSchemaVersion, ID: "context-1", ExecutionEnvelopeID: "envelope-1", Objective: "prepare a bounded change", AssembledAt: testNow.Add(time.Minute), ExpiresAt: &expiresAt, PolicyDecisionIDs: []string{"policy-1"}, Items: []ContextItem{{SourceID: "source-1", SourceRevision: "revision-1", Locator: "docs/input.md", Classification: "internal", Authority: "reference", Purpose: "task context", ContentDigest: "sha256:context"}}, ManifestDigest: "sha256:manifest"}
}

func validApprovalRecord() ApprovalRecord {
	return ApprovalRecord{SchemaVersion: ApprovalRecordSchemaVersion, ID: "approval-1", ApproverID: "human-1", Action: "publish-pr", SubjectType: "candidate-change", SubjectDigest: "sha256:candidate", Decision: "approved", DecidedAt: testNow, ExpiresAt: testNow.Add(10 * time.Minute), EvidenceCorrelationID: testCorrelationID}
}

func validEvidenceRecord() EvidenceRecord {
	return EvidenceRecord{SchemaVersion: EvidenceRecordSchemaVersion, ID: "evidence-1", Type: "verification", RecordedAt: testNow, ActorID: "verifier-1", Scope: EvidenceScope{TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1", DeliveryRunID: "run-1"}, Subject: EvidenceSubject{Type: "candidate-change", ID: "candidate-1", Digest: "sha256:candidate"}, References: []string{}, Classification: "internal", RetentionClass: "delivery-evidence", Integrity: EvidenceIntegrity{Algorithm: "sha256", Digest: "sha256:evidence"}}
}

func assertFailureCode(t *testing.T, err error, want failure.Code) {
	t.Helper()
	if err == nil {
		t.Fatalf("error = nil, want code %q", want)
	}
	var normalized *failure.Error
	if !errors.As(err, &normalized) {
		t.Fatalf("error type = %T, want *failure.Error", err)
	}
	if normalized.Code != want {
		t.Fatalf("error code = %q, want %q (error: %v)", normalized.Code, want, err)
	}
	if validateErr := normalized.Validate(); validateErr != nil {
		t.Fatalf("normalized error validation = %v", validateErr)
	}
}
