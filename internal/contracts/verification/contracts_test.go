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
		Environment:          "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr-1",
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
		Environment:  "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr-1",
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

func TestPassedResultMustContainIndependentPass(t *testing.T) {
	now := time.Now().UTC()
	plan, err := BindPlan(AssurancePlan{
		SchemaVersion: AssurancePlanVersion, CandidateDigest: digest("candidate"), PlanDigest: digest("plan"),
		RiskClass: "low", Classification: "internal", AffectedCapabilities: []string{"forge.capability.delivery"},
		ChangedPaths: []string{"internal/delivery/example.go"},
		Controls:     []Control{{ID: "forge-contract", Version: "1", Consequential: true}},
		Environment:  "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := BindResult(Result{
		SchemaVersion: ResultVersion, AssurancePlanID: plan.ID, AssurancePlanDigest: plan.Digest,
		CandidateDigest: plan.CandidateDigest, Environment: plan.Environment, Outcome: "passed",
		Independent: true, CompletedAt: now, EvidenceCorrelationID: plan.EvidenceCorrelationID,
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
		Environment: "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := BindResult(Result{
		SchemaVersion: ResultVersion, AssurancePlanID: plan.ID, AssurancePlanDigest: plan.Digest,
		CandidateDigest: plan.CandidateDigest, Environment: plan.Environment, Outcome: "passed",
		Independent: true, CompletedAt: now, EvidenceCorrelationID: plan.EvidenceCorrelationID,
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
