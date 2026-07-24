package execution

import (
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/trust"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
)

func TestWorkerResultCannotExceedTrustedLimits(t *testing.T) {
	now := time.Now().UTC()
	request := workercontract.Request{
		Limits: workercontract.Limits{
			MaxDurationSeconds: 10, MaxOutputBytes: 100, MaxChangedFiles: 1,
		},
		ExpiresAt: now.Add(time.Minute),
	}
	bundle := workercontract.CandidateBundle{
		Files: []workercontract.CandidateFile{{Path: "internal/example.go"}},
	}
	result := workercontract.Result{
		SchemaVersion: workercontract.ResultVersion, ChangedPaths: []string{"internal/example.go"},
		ResourceUse: workercontract.ResourceUse{
			DurationMillis: 10, OutputBytes: 50, ChangedFiles: 1,
		},
		CleanupOutcome: "completed", CompletedAt: now.Add(time.Second),
	}
	if !validWorkerResult(result, request, bundle, now) {
		t.Fatal("bounded worker result was rejected")
	}
	result.ResourceUse.OutputBytes = 101
	if validWorkerResult(result, request, bundle, now) {
		t.Fatal("oversized worker result must fail")
	}
	result.ResourceUse.OutputBytes = 50
	result.CleanupOutcome = "not-completed"
	if validWorkerResult(result, request, bundle, now) {
		t.Fatal("unproven worker cleanup must fail")
	}
}

func TestProposalMaterializationAuthorityIsBoundToExactProposal(t *testing.T) {
	now := time.Now().UTC()
	expiresAt := now.Add(time.Minute)
	decision := trust.PolicyDecision{
		SchemaVersion: trust.PolicyDecisionSchemaVersion,
		ID:            "decision-1", DecisionPoint: "execution", SubjectID: "agent-1",
		Action:      workercontract.PrepareChangeAction,
		Resource:    "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		EvaluatedAt: now, ExpiresAt: &expiresAt, PolicyVersions: []string{"policy-v1"},
		Outcome: "permit", Obligations: []map[string]any{{"type": "proposal-admitted"}},
		ReasonCodes: []string{"bounded-proposal"}, EvidenceCorrelationID: "correlation-1",
	}
	envelope := trust.ExecutionEnvelope{
		Principal:         trust.EnvelopePrincipal{AgentID: "agent-1"},
		PolicyDecisionIDs: []string{"decision-1"},
	}
	if !validProposalMaterializationAuthority(
		decision, envelope, decision.Resource, now.Add(time.Second), "correlation-1",
	) {
		t.Fatal("exact proposal authority was rejected")
	}
	if validProposalMaterializationAuthority(
		decision, envelope,
		"sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		now.Add(time.Second), "correlation-1",
	) {
		t.Fatal("authority for a different proposal must fail closed")
	}
	envelope.PolicyDecisionIDs = []string{"decision-other"}
	if validProposalMaterializationAuthority(
		decision, envelope, decision.Resource, now.Add(time.Second), "correlation-1",
	) {
		t.Fatal("unreferenced policy authority must fail closed")
	}
}
