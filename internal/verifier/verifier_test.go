package verifier

import (
	"context"
	"testing"
	"time"

	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
)

func TestVerifierUsesIndependentForgeControls(t *testing.T) {
	now := time.Now().UTC()
	bundle, err := workercontract.BindBundle(workercontract.CandidateBundle{
		SchemaVersion: workercontract.BundleVersion, SourceRevision: "commit",
		PlanDigest: fixedDigest('a'), Files: []workercontract.CandidateFile{{
			Path: "internal/example.go", Content: []byte("package example\n"),
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := verificationcontract.BindPlan(verificationcontract.AssurancePlan{
		SchemaVersion:   verificationcontract.AssurancePlanVersion,
		CandidateDigest: bundle.Digest, PlanDigest: bundle.PlanDigest, RiskClass: "low",
		Classification: "internal", AffectedCapabilities: []string{"forge.capability.delivery"},
		ChangedPaths: []string{"internal/example.go"},
		Controls: []verificationcontract.Control{
			{ID: "candidate.tests", Version: "1", Consequential: false, CandidateOwned: true},
			{ID: ControlCandidateIntegrity, Version: "1", Consequential: true},
			{ID: ControlDeclaredPathScope, Version: "1", Consequential: true},
			{ID: ControlGoFormat, Version: "1", Consequential: true},
		},
		Environment: "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr",
	})
	if err != nil {
		t.Fatal(err)
	}
	engine, _ := NewEngine(func() time.Time { return now })
	result, err := engine.Verify(context.Background(), plan, bundle, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != "passed" || !result.Independent {
		t.Fatal("independent verification did not pass")
	}
}

func fixedDigest(value byte) string {
	result := make([]byte, 64)
	for index := range result {
		result[index] = value
	}
	return "sha256:" + string(result)
}
