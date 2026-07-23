package execution

import (
	"testing"
	"time"

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
