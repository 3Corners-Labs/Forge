package verification

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPostgresVerificationContractsSurviveRestart(t *testing.T) {
	dsn := os.Getenv("FORGE_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("FORGE_TEST_POSTGRES_DSN is not configured")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	store, _ := NewPostgresStore(db)
	if err := store.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	plan, err := verificationcontract.BindPlan(verificationcontract.AssurancePlan{
		SchemaVersion:   verificationcontract.AssurancePlanVersion,
		CandidateDigest: verifierDigest('d'), PlanDigest: verifierDigest('e'),
		RiskClass: "low", Classification: "Internal",
		AffectedCapabilities: []string{"forge.capability.delivery"}, ChangedPaths: []string{"internal/example.go"},
		Controls:    []verificationcontract.Control{{ID: "forge.candidate-integrity", Version: "1", Consequential: true}},
		Environment: "developer-trial", ExpiresAt: now.Add(time.Hour),
		EvidenceCorrelationID: "correlation-" + now.Format(time.RFC3339Nano),
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := verificationcontract.BindResult(verificationcontract.Result{
		SchemaVersion: verificationcontract.ResultVersion, AssurancePlanID: plan.ID,
		AssurancePlanDigest: plan.Digest, CandidateDigest: plan.CandidateDigest,
		Environment: plan.Environment, Outcome: "passed", Independent: true, CompletedAt: now,
		EvidenceCorrelationID: plan.EvidenceCorrelationID,
		Findings: []verificationcontract.Finding{{
			ControlID: "forge.candidate-integrity", Outcome: "passed", Code: "passed",
			EvidenceDigest: verifierDigest('f'), Consequential: true,
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := store.SavePlan(ctx, plan); err != nil {
		t.Fatal(err)
	}
	if err := store.SaveResult(ctx, result); err != nil {
		t.Fatal(err)
	}
	restarted, _ := NewPostgresStore(db)
	storedPlan, found, err := restarted.ResolvePlan(ctx, plan.ID)
	if err != nil || !found || storedPlan.Digest != plan.Digest {
		t.Fatalf("plan found=%t digest=%q err=%v", found, storedPlan.Digest, err)
	}
	storedResult, found, err := restarted.ResolveResult(ctx, result.ID)
	if err != nil || !found || storedResult.Digest != result.Digest {
		t.Fatalf("result found=%t digest=%q err=%v", found, storedResult.Digest, err)
	}
}
