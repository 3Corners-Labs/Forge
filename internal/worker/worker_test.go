package worker

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	workercontract "github.com/3corners/forge/internal/contracts/worker"
)

func TestPrepareCreatesDigestBoundCandidateWithoutShell(t *testing.T) {
	now := time.Date(2026, 7, 23, 14, 0, 0, 0, time.UTC)
	engine, _ := NewEngine(func() time.Time { return now })
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "internal"), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "internal", "example.txt"), []byte("before"), 0o600); err != nil {
		t.Fatal(err)
	}
	request := boundRequest(t, now, "internal/example.txt", []byte("after"))
	bundle, result, err := engine.Prepare(context.Background(), request, root)
	if err != nil {
		t.Fatal(err)
	}
	if err := bundle.Validate(); err != nil {
		t.Fatal(err)
	}
	if result.CandidateDigest != bundle.Digest || result.ChangedPaths[0] != "internal/example.txt" {
		t.Fatal("result is not bound to the candidate")
	}
}

func TestPrepareRejectsSymlinkEscapeAndChangedSource(t *testing.T) {
	now := time.Now().UTC()
	engine, _ := NewEngine(func() time.Time { return now })
	root := t.TempDir()
	outside := t.TempDir()
	if err := os.Symlink(outside, filepath.Join(root, "internal")); err != nil {
		t.Fatal(err)
	}
	request := boundRequest(t, now, "internal/escape.txt", []byte("no"))
	if _, _, err := engine.Prepare(context.Background(), request, root); !errors.Is(err, ErrWorkspaceBoundary) {
		t.Fatalf("symlink escape must fail, got %v", err)
	}

	root = t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "internal"), 0o700); err != nil {
		t.Fatal(err)
	}
	request = boundRequest(t, now, "internal/example.txt", []byte("after"))
	request.Replacements[0].ExpectedSourceDigest = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	request, _ = workercontract.BindRequest(request)
	if _, _, err := engine.Prepare(context.Background(), request, root); !errors.Is(err, ErrSourceChanged) {
		t.Fatalf("changed source must fail, got %v", err)
	}
}

func boundRequest(t *testing.T, now time.Time, target string, content []byte) workercontract.Request {
	t.Helper()
	request, err := workercontract.BindRequest(workercontract.Request{
		SchemaVersion: workercontract.RequestVersion, EnvelopeID: "env-1", DelegationID: "delegation-1",
		PlanID: "plan-1", PlanDigest: fixedDigest('a'), PlanApprovalRecordID: "approval-1",
		ContextPackageID: "context-1", ContextManifestDigest: fixedDigest('b'), SourceRevision: "commit-1",
		Scope:  workercontract.Scope{TenantID: "tenant", WorkspaceID: "workspace", ProjectID: "project", RepositoryID: "repo", DeliveryRunID: "run"},
		Action: workercontract.PrepareChangeAction, TransformationID: workercontract.TransformationV1,
		TransformationVersion: "1.0.0", AllowedReadPaths: []string{"internal"}, AllowedWritePaths: []string{"internal"},
		DependencyDigests: []string{fixedDigest('c')}, Replacements: []workercontract.FileReplacement{{Path: target, Content: content}},
		Limits:  workercontract.Limits{MaxDurationSeconds: 60, MaxProcesses: 8, MaxOutputBytes: 1024, MaxChangedFiles: 2, MaxAttempts: 2},
		Attempt: 1, ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr",
	})
	if err != nil {
		t.Fatal(err)
	}
	return request
}

func fixedDigest(value byte) string {
	result := make([]byte, 64)
	for index := range result {
		result[index] = value
	}
	return "sha256:" + string(result)
}
