package trust

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
)

var approvalTestTime = time.Date(2026, time.July, 23, 10, 0, 0, 0, time.UTC)
var approvalPlanDigest = "sha256:" + strings.Repeat("a", 64)

func TestApprovalV2ContractsAreCanonicalAndMutationDetecting(t *testing.T) {
	left := approvalRequestFixture()
	right := approvalRequestFixture()
	left.EvidenceReferences = []string{"plan:" + approvalPlanDigest, "context:package"}
	right.EvidenceReferences = []string{"context:package", "plan:" + approvalPlanDigest}
	boundLeft, err := BindApprovalRequestDigest(left, "correlation-1")
	if err != nil {
		t.Fatal(err)
	}
	boundRight, err := BindApprovalRequestDigest(right, "correlation-1")
	if err != nil {
		t.Fatal(err)
	}
	if boundLeft.RequestDigest != boundRight.RequestDigest || boundLeft.ID != boundRight.ID {
		t.Fatal("permitted collection ordering changed canonical approval request identity")
	}
	mutated := boundLeft
	mutated.Scope.ProjectID = "other-project"
	assertApprovalFailure(t, mutated.VerifyDigest("correlation-1"), failure.VerificationFailed)

	record, err := BindApprovalRecordV2Integrity(ApprovalRecordV2{
		SchemaVersion: ApprovalRecordV2SchemaVersion, Issuer: "forge.approval.developer-trial",
		RequestID: boundLeft.ID, RequestDigest: boundLeft.RequestDigest, ApproverID: "human-1",
		Scope: boundLeft.Scope, Environment: boundLeft.Environment, Purpose: boundLeft.Purpose,
		Classification: boundLeft.Classification, RiskClass: boundLeft.RiskClass,
		Action: boundLeft.Action, SubjectType: boundLeft.SubjectType, SubjectID: boundLeft.SubjectID,
		SubjectDigest: boundLeft.SubjectDigest, Decision: "approved", DecidedAt: approvalTestTime,
		ExpiresAt: approvalTestTime.Add(10 * time.Minute), Conditions: []string{"require-independent-verification"},
		PolicyDecisionIDs: []string{"decision-eligibility"}, EvidenceReferences: boundLeft.EvidenceReferences,
		EvidenceCorrelationID: "correlation-1",
	}, "correlation-1")
	if err != nil {
		t.Fatal(err)
	}
	expected := ApprovalExpectedFacts{
		RequestID: boundLeft.ID, RequestDigest: boundLeft.RequestDigest, ApproverID: "human-1",
		Scope: boundLeft.Scope, Environment: "developer-trial", Purpose: boundLeft.Purpose,
		Classification: "Internal", RiskClass: "low", Action: "forge.delivery.execute-plan",
		SubjectType: "delivery-plan", SubjectID: "plan-1", SubjectDigest: approvalPlanDigest,
		EvidenceReferences:  boundLeft.EvidenceReferences,
		PermittedConditions: []string{"require-independent-verification"},
	}
	if err := record.Authorize(expected, approvalTestTime.Add(time.Minute), "correlation-1"); err != nil {
		t.Fatal(err)
	}
	crossEnvironment := expected
	crossEnvironment.Environment = "production"
	assertApprovalFailure(t, record.Authorize(crossEnvironment, approvalTestTime.Add(time.Minute), "correlation-1"), failure.Forbidden)
	changedPlan := expected
	changedPlan.SubjectDigest = "sha256:" + strings.Repeat("b", 64)
	assertApprovalFailure(t, record.Authorize(changedPlan, approvalTestTime.Add(time.Minute), "correlation-1"), failure.Forbidden)
	unsupported := record
	unsupported.ID = ""
	unsupported.Integrity = Integrity{}
	unsupported.Conditions = []string{"widen-scope"}
	unsupported, err = BindApprovalRecordV2Integrity(unsupported, "correlation-1")
	if err != nil {
		t.Fatal(err)
	}
	assertApprovalFailure(t, unsupported.Authorize(expected, approvalTestTime.Add(time.Minute), "correlation-1"), failure.Forbidden)
}

func TestApprovalV2StrictDecodingAndV1Separation(t *testing.T) {
	request, err := BindApprovalRequestDigest(approvalRequestFixture(), "correlation-1")
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := json.Marshal(request)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := DecodeApprovalRequest(encoded, "correlation-1"); err != nil {
		t.Fatal(err)
	}
	var raw map[string]any
	if err := json.Unmarshal(encoded, &raw); err != nil {
		t.Fatal(err)
	}
	raw["unexpected"] = true
	unknown, _ := json.Marshal(raw)
	assertApprovalFailure(t, func() error {
		_, decodeErr := DecodeApprovalRequest(unknown, "correlation-1")
		return decodeErr
	}(), failure.InvalidRequest)

	record := ApprovalRecord{
		SchemaVersion: ApprovalRecordSchemaVersion, ID: "approval-v1", ApproverID: "human-1",
		Action: "forge.delivery.execute-plan", SubjectType: "delivery-plan", SubjectDigest: "sha256:plan",
		Decision: "approved", DecidedAt: approvalTestTime, ExpiresAt: approvalTestTime.Add(time.Minute),
		EvidenceCorrelationID: "correlation-1",
	}
	if err := record.Authorize("forge.delivery.execute-plan", "delivery-plan", "sha256:plan", approvalTestTime.Add(time.Second), "correlation-1"); err != nil {
		t.Fatal(err)
	}
	v1JSON, _ := json.Marshal(record)
	assertApprovalFailure(t, func() error {
		_, decodeErr := DecodeApprovalRecordV2(v1JSON, "correlation-1")
		return decodeErr
	}(), failure.InvalidRequest)
}

func BenchmarkApprovalV2Integrity(b *testing.B) {
	request := approvalRequestFixture()
	for i := 0; i < b.N; i++ {
		if _, err := BindApprovalRequestDigest(request, "correlation-1"); err != nil {
			b.Fatal(err)
		}
	}
}

func approvalRequestFixture() ApprovalRequest {
	return ApprovalRequest{
		SchemaVersion: ApprovalRequestSchemaVersion, Issuer: "forge.approval.developer-trial",
		RequestingActorID: "human-1", IntendedApproverID: "human-1",
		Scope:       ApprovalScope{TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1", RepositoryID: "repository-1", DeliveryRunID: "run-1"},
		Environment: "developer-trial", Purpose: "forge.m1.prepare-change", Classification: "Internal",
		RiskClass: "low", Action: "forge.delivery.execute-plan", SubjectType: "delivery-plan",
		SubjectID: "plan-1", SubjectDigest: approvalPlanDigest,
		EvidenceReferences:  []string{"context:package", "plan:" + approvalPlanDigest},
		PermittedConditions: []string{"require-independent-verification"},
		PolicyDecisionIDs:   []string{"decision-context"}, CreatedAt: approvalTestTime,
		ExpiresAt: approvalTestTime.Add(30 * time.Minute), MaxDecisionLifetime: 1800,
		Lifecycle: "open", EvidenceCorrelationID: "correlation-1",
	}
}

func assertApprovalFailure(t *testing.T, err error, code failure.Code) {
	t.Helper()
	var normalized *failure.Error
	if err == nil || !errors.As(err, &normalized) || normalized.Code != code {
		t.Fatalf("error = %v, want %s", err, code)
	}
}
