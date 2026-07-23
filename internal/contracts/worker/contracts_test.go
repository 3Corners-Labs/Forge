package worker

import (
	"encoding/json"
	"errors"
	"testing"
	"time"
)

func TestRequestAndCandidateDigestsFailClosed(t *testing.T) {
	now := time.Date(2026, 7, 23, 14, 0, 0, 0, time.UTC)
	request, err := BindRequest(requestFixture(now))
	if err != nil {
		t.Fatal(err)
	}
	if err := request.Validate(now); err != nil {
		t.Fatal(err)
	}
	mutated := request
	mutated.Replacements[0].Content = []byte("changed")
	if !errors.Is(mutated.Validate(now), ErrDigestMismatch) {
		t.Fatal("mutated request must fail its digest")
	}

	bundle, err := BindBundle(CandidateBundle{
		SchemaVersion: BundleVersion, SourceRevision: "commit-1", PlanDigest: request.PlanDigest,
		Files: []CandidateFile{{Path: "internal/example.txt", Content: []byte("after")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := bundle.Validate(); err != nil {
		t.Fatal(err)
	}
	bundle.Files[0].Content[0] = 'A'
	if !errors.Is(bundle.Validate(), ErrDigestMismatch) {
		t.Fatal("mutated candidate must fail its digest")
	}
}

func TestRequestRejectsTraversalUnknownFieldsAndAuthorityExpansion(t *testing.T) {
	now := time.Now().UTC()
	fixture := requestFixture(now)
	fixture.Replacements[0].Path = "../escape"
	if _, err := BindRequest(fixture); !errors.Is(err, ErrInvalidContract) {
		t.Fatal("path traversal must fail")
	}
	fixture = requestFixture(now)
	fixture.AllowedWritePaths = []string{"docs"}
	if _, err := BindRequest(fixture); !errors.Is(err, ErrInvalidContract) {
		t.Fatal("replacement outside write scope must fail")
	}
	valid, err := BindRequest(requestFixture(now))
	if err != nil {
		t.Fatal(err)
	}
	encoded, _ := json.Marshal(valid)
	encoded[len(encoded)-1] = ','
	encoded = append(encoded, []byte(`"command":"sh"}`)...)
	if _, err := DecodeRequest(encoded); !errors.Is(err, ErrInvalidContract) {
		t.Fatal("unknown command authority must fail strict decoding")
	}
}

func TestWorkerResultRejectsUnknownFields(t *testing.T) {
	encoded := []byte(`{"schemaVersion":"forge.worker-result.v1","candidateDigest":"sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa","command":"sh"}`)
	if _, err := DecodeResult(encoded); !errors.Is(err, ErrInvalidContract) {
		t.Fatal("unknown worker result authority must fail strict decoding")
	}
}

func requestFixture(now time.Time) Request {
	return Request{
		SchemaVersion: RequestVersion, EnvelopeID: "env-1", DelegationID: "delegation-1",
		PlanID: "plan-1", PlanDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		PlanApprovalRecordID: "approval-1", ContextPackageID: "context-1",
		ContextManifestDigest: "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb",
		SourceRevision:        "commit-1", Scope: Scope{TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1", RepositoryID: "repo-1", DeliveryRunID: "run-1"},
		Action: PrepareChangeAction, TransformationID: TransformationV1, TransformationVersion: "1.0.0",
		AllowedReadPaths: []string{"internal"}, AllowedWritePaths: []string{"internal"},
		DependencyDigests: []string{"sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc"},
		Replacements:      []FileReplacement{{Path: "internal/example.txt", Content: []byte("after")}},
		Limits:            Limits{MaxDurationSeconds: 60, MaxProcesses: 8, MaxOutputBytes: 1 << 20, MaxChangedFiles: 4, MaxAttempts: 2},
		Attempt:           1, ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr-1",
	}
}
