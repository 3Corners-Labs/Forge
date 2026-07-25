package verification

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestAssurancePlanRequiresIndependentConsequentialOracle(t *testing.T) {
	now := time.Now().UTC()
	base := AssurancePlan{
		SchemaVersion:   AssurancePlanVersion,
		CandidateDigest: digest("candidate"), PlanDigest: digest("plan"),
		RiskClass: "low", Classification: "internal",
		AffectedCapabilities: []string{"forge.capability.delivery"},
		ChangedPaths:         []string{"internal/delivery/example.go"},
		Controls:             []Control{{ID: "candidate-test", Version: "1", Consequential: true, CandidateOwned: true}},
		ProducedBy:           "forge.worker.subject.test", Environment: "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr-1",
	}
	if _, err := BindPlan(base); !errors.Is(err, ErrInvalidContract) {
		t.Fatal("candidate-owned control cannot be the sole oracle")
	}
	base.Controls = append(base.Controls, Control{ID: "forge-contract", Version: "1", Consequential: true})
	plan, err := BindPlan(base)
	if err != nil {
		t.Fatal(err)
	}
	if err := plan.Validate(now); err != nil {
		t.Fatal(err)
	}
}

func TestVerificationContractsRejectUnknownFields(t *testing.T) {
	now := time.Now().UTC()
	plan, err := BindPlan(AssurancePlan{
		SchemaVersion: AssurancePlanVersion, CandidateDigest: digest("candidate"), PlanDigest: digest("plan"),
		RiskClass: "low", Classification: "internal", AffectedCapabilities: []string{"forge.capability.delivery"},
		ChangedPaths: []string{"internal/delivery/example.go"},
		Controls:     []Control{{ID: "forge-contract", Version: "1", Consequential: true}},
		ProducedBy:   "forge.worker.subject.test", Environment: "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(plan)
	encoded[len(encoded)-1] = ','
	encoded = append(encoded, []byte(`"candidateCommand":"go test"}`)...)
	if _, err := DecodePlan(encoded); !errors.Is(err, ErrInvalidContract) {
		t.Fatal("unknown Assurance Plan authority must fail strict decoding")
	}
}

// TestProducerCannotVerifyItsOwnCandidate covers OWASP AISVS AC.8.1. The
// invariant is dormant while the verifier is deterministic, because a
// deterministic verifier is never the producing principal. It becomes
// load-bearing as soon as verification is model-assisted, so it is asserted now
// rather than when it is needed.
func TestProducerCannotVerifyItsOwnCandidate(t *testing.T) {
	now := time.Now().UTC()
	const producer = "forge.worker.subject.test"
	plan, err := BindPlan(AssurancePlan{
		SchemaVersion: AssurancePlanVersion, CandidateDigest: digest("candidate"), PlanDigest: digest("plan"),
		RiskClass: "low", Classification: "internal", AffectedCapabilities: []string{"forge.capability.delivery"},
		ChangedPaths: []string{"internal/delivery/example.go"},
		Controls:     []Control{{ID: "forge-contract", Version: "1", Consequential: true}},
		ProducedBy:   producer, Environment: "developer-trial", ExpiresAt: now.Add(time.Hour),
		EvidenceCorrelationID: "corr-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	build := func(verifiedBy string) Result {
		result, err := BindResult(Result{
			SchemaVersion: ResultVersion, AssurancePlanID: plan.ID, AssurancePlanDigest: plan.Digest,
			CandidateDigest: plan.CandidateDigest, Environment: plan.Environment, Outcome: "passed",
			Independent: true, VerifiedBy: verifiedBy, CompletedAt: now,
			EvidenceCorrelationID: plan.EvidenceCorrelationID,
			Findings:              []Finding{{ControlID: "forge-contract", Outcome: "passed", Code: "ok", EvidenceDigest: digest("evidence"), Consequential: true}},
		})
		if err != nil {
			t.Fatal(err)
		}
		return result
	}
	if err := build("forge.verifier.deterministic.v1").Validate(plan); err != nil {
		t.Fatalf("an independent verifier must be accepted: %v", err)
	}
	if err := build(producer).Validate(plan); !errors.Is(err, ErrInvalidContract) {
		t.Fatal("the producing principal must not be able to verify its own candidate")
	}
}

// TestPlanRequiresRecordedProvenance keeps candidate provenance mandatory. It is
// the field the AC.8.1 invariant and the delivery measures both depend on.
func TestPlanRequiresRecordedProvenance(t *testing.T) {
	now := time.Now().UTC()
	if _, err := BindPlan(AssurancePlan{
		SchemaVersion: AssurancePlanVersion, CandidateDigest: digest("candidate"), PlanDigest: digest("plan"),
		RiskClass: "low", Classification: "internal", AffectedCapabilities: []string{"forge.capability.delivery"},
		ChangedPaths: []string{"internal/delivery/example.go"},
		Controls:     []Control{{ID: "forge-contract", Version: "1", Consequential: true}},
		Environment:  "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr-1",
	}); !errors.Is(err, ErrInvalidContract) {
		t.Fatal("a plan without recorded provenance must be refused")
	}
}

func TestPassedResultMustContainIndependentPass(t *testing.T) {
	now := time.Now().UTC()
	plan, err := BindPlan(AssurancePlan{
		SchemaVersion: AssurancePlanVersion, CandidateDigest: digest("candidate"), PlanDigest: digest("plan"),
		RiskClass: "low", Classification: "internal", AffectedCapabilities: []string{"forge.capability.delivery"},
		ChangedPaths: []string{"internal/delivery/example.go"},
		Controls:     []Control{{ID: "forge-contract", Version: "1", Consequential: true}},
		ProducedBy:   "forge.worker.subject.test", Environment: "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := BindResult(Result{
		SchemaVersion: ResultVersion, AssurancePlanID: plan.ID, AssurancePlanDigest: plan.Digest,
		CandidateDigest: plan.CandidateDigest, Environment: plan.Environment, Outcome: "passed",
		Independent: true, VerifiedBy: "forge.verifier.test", CompletedAt: now, EvidenceCorrelationID: plan.EvidenceCorrelationID,
		Findings: []Finding{{ControlID: "forge-contract", Outcome: "passed", Code: "ok", EvidenceDigest: digest("evidence"), Consequential: true}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := result.Validate(plan); err != nil {
		t.Fatal(err)
	}
	result.Findings[0].CandidateOwned = true
	if _, err := BindResult(result); !errors.Is(err, ErrInvalidContract) {
		t.Fatal("self-owned consequential evidence must not pass")
	}
}

func TestResultMustCoverExactAssuranceControls(t *testing.T) {
	now := time.Now().UTC()
	plan, err := BindPlan(AssurancePlan{
		SchemaVersion: AssurancePlanVersion, CandidateDigest: digest("candidate"), PlanDigest: digest("plan"),
		RiskClass: "low", Classification: "internal", AffectedCapabilities: []string{"forge.capability.delivery"},
		ChangedPaths: []string{"internal/delivery/example.go"},
		Controls: []Control{
			{ID: "forge-contract", Version: "1", Consequential: true},
			{ID: "candidate-test", Version: "1", CandidateOwned: true},
		},
		ProducedBy: "forge.worker.subject.test", Environment: "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := BindResult(Result{
		SchemaVersion: ResultVersion, AssurancePlanID: plan.ID, AssurancePlanDigest: plan.Digest,
		CandidateDigest: plan.CandidateDigest, Environment: plan.Environment, Outcome: "passed",
		Independent: true, VerifiedBy: "forge.verifier.test", CompletedAt: now, EvidenceCorrelationID: plan.EvidenceCorrelationID,
		Findings: []Finding{{
			ControlID: "forge-contract", Outcome: "passed", Code: "ok",
			EvidenceDigest: digest("evidence"), Consequential: true,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if !errors.Is(result.Validate(plan), ErrInvalidContract) {
		t.Fatal("omitted Assurance Plan control must fail")
	}
}
