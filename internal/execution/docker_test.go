package execution

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	workercontract "github.com/3corners/forge/internal/contracts/worker"
)

type commandRunnerFunc func(context.Context, string, ...string) error

func (function commandRunnerFunc) Run(ctx context.Context, name string, arguments ...string) error {
	return function(ctx, name, arguments...)
}

func TestDockerRunnerEnforcesIsolationAndReadsClosedOutputs(t *testing.T) {
	now := time.Now().UTC()
	request, err := workercontract.BindRequest(workercontract.Request{
		SchemaVersion: workercontract.RequestVersion, EnvelopeID: "envelope-1", DelegationID: "delegation-1",
		PlanID: "plan-1", PlanDigest: testDigest('a'), PlanApprovalRecordID: "approval-1",
		ContextPackageID: "context-1", ContextManifestDigest: testDigest('b'), SourceRevision: "revision-1",
		Scope: workercontract.Scope{
			TenantID: "tenant-1", WorkspaceID: "workspace-1", ProjectID: "project-1",
			RepositoryID: "repository-1", DeliveryRunID: "run-1",
		},
		Action: workercontract.PrepareChangeAction, TransformationID: workercontract.TransformationV1,
		TransformationVersion: "1", AllowedReadPaths: []string{"internal/example.txt"},
		AllowedWritePaths: []string{"internal/example.txt"}, DependencyDigests: []string{testDigest('c')},
		Replacements: []workercontract.FileReplacement{{Path: "internal/example.txt", Content: []byte("after")}},
		Limits: workercontract.Limits{
			MaxDurationSeconds: 60, MaxProcesses: 4, MaxOutputBytes: 1 << 20, MaxChangedFiles: 1, MaxAttempts: 1,
		},
		Attempt: 1, ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "correlation-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	bundle, err := workercontract.BindBundle(workercontract.CandidateBundle{
		SchemaVersion: workercontract.BundleVersion, SourceRevision: request.SourceRevision,
		PlanDigest: request.PlanDigest, Files: []workercontract.CandidateFile{{Path: "internal/example.txt", Content: []byte("after")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	var invoked string
	fake := commandRunnerFunc(func(_ context.Context, name string, arguments ...string) error {
		invoked = name + " " + strings.Join(arguments, " ")
		output := mountedSource(arguments, "/output")
		bundleBytes, _ := json.Marshal(bundle)
		resultBytes, _ := json.Marshal(workercontract.Result{
			SchemaVersion: workercontract.ResultVersion, RequestID: request.ID, RequestDigest: request.Digest,
			CandidateDigest: bundle.Digest, TreeDigest: workercontract.CandidateTreeDigest(bundle),
			ChangedPaths: []string{"internal/example.txt"}, BundleDigest: bundle.Digest,
			ResourceUse: workercontract.ResourceUse{ChangedFiles: 1}, CleanupOutcome: "not-completed",
			EvidenceReferences: []string{}, CompletedAt: now,
		})
		if err := os.WriteFile(filepath.Join(output, "bundle.json"), bundleBytes, 0o600); err != nil {
			return err
		}
		return os.WriteFile(filepath.Join(output, "result.json"), resultBytes, 0o600)
	})
	runner, err := NewDockerRunner(DockerConfig{
		Image:    "registry.example/forge-worker@sha256:" + strings.Repeat("d", 64),
		TempRoot: t.TempDir(), Memory: "256m", CPUs: "1",
	}, fake)
	if err != nil {
		t.Fatal(err)
	}
	workspace := t.TempDir()
	gotBundle, gotResult, err := runner.Run(context.Background(), request, workspace)
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"--network none", "--read-only", "--cap-drop ALL", "no-new-privileges",
		"--user 65532:65532", "--pids-limit 4", "GOMAXPROCS=1", "dst=/input,readonly",
	} {
		if !strings.Contains(invoked, required) {
			t.Fatalf("Docker invocation lacks %q: %s", required, invoked)
		}
	}
	if gotBundle.Digest != bundle.Digest || gotResult.CleanupOutcome != "completed" {
		t.Fatalf("unexpected outputs: bundle=%q cleanup=%q", gotBundle.Digest, gotResult.CleanupOutcome)
	}
}

func TestDockerRunnerAcceptsImmutableLocalImageIDAndRejectsTag(t *testing.T) {
	config := DockerConfig{
		Image: "sha256:" + strings.Repeat("a", 64), TempRoot: t.TempDir(), Memory: "256m", CPUs: "1",
	}
	if _, err := NewDockerRunner(config, commandRunnerFunc(func(context.Context, string, ...string) error { return nil })); err != nil {
		t.Fatalf("immutable local image ID rejected: %v", err)
	}
	config.Image = "forge-worker:latest"
	if _, err := NewDockerRunner(config, commandRunnerFunc(func(context.Context, string, ...string) error { return nil })); err == nil {
		t.Fatal("mutable Docker tag was admitted")
	}
}

func TestDockerRunnerExecutesImmutableDeveloperTrialImage(t *testing.T) {
	image := os.Getenv("FORGE_TEST_WORKER_IMAGE")
	if image == "" {
		t.Skip("FORGE_TEST_WORKER_IMAGE is not configured")
	}
	now := time.Now().UTC()
	request, err := workercontract.BindRequest(workercontract.Request{
		SchemaVersion: workercontract.RequestVersion, EnvelopeID: "envelope-integration",
		DelegationID: "delegation-integration", PlanID: "plan-integration",
		PlanDigest: testDigest('a'), PlanApprovalRecordID: "approval-integration",
		ContextPackageID: "context-integration", ContextManifestDigest: testDigest('b'),
		SourceRevision: "revision-integration",
		Scope: workercontract.Scope{
			TenantID: "tenant", WorkspaceID: "workspace", ProjectID: "project",
			RepositoryID: "repository", DeliveryRunID: "run",
		},
		Action: workercontract.PrepareChangeAction, TransformationID: workercontract.TransformationV1,
		TransformationVersion: "1", AllowedReadPaths: []string{"docs/result.md"},
		AllowedWritePaths: []string{"docs/result.md"}, DependencyDigests: []string{testDigest('c')},
		Replacements: []workercontract.FileReplacement{{Path: "docs/result.md", Content: []byte("# Result\n")}},
		Limits: workercontract.Limits{
			MaxDurationSeconds: 60, MaxProcesses: 8, MaxOutputBytes: 1 << 20,
			MaxChangedFiles: 1, MaxAttempts: 1,
		},
		Attempt: 1, ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "correlation-integration",
	})
	if err != nil {
		t.Fatal(err)
	}
	root := t.TempDir()
	workspace := filepath.Join(root, "workspace")
	if err := os.Mkdir(workspace, 0o777); err != nil {
		t.Fatal(err)
	}
	runner, err := NewDockerRunner(DockerConfig{
		Image: image, TempRoot: root, Memory: "256m", CPUs: "1",
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	bundle, result, err := runner.Run(context.Background(), request, workspace)
	if err != nil {
		t.Fatal(err)
	}
	if bundle.Digest == "" || result.CandidateDigest != bundle.Digest {
		t.Fatalf("unbound runtime result: %#v %#v", bundle, result)
	}
}

func mountedSource(arguments []string, destination string) string {
	for index, argument := range arguments {
		if argument == "--mount" && index+1 < len(arguments) {
			fields := strings.Split(arguments[index+1], ",")
			source, target := "", ""
			for _, field := range fields {
				if strings.HasPrefix(field, "src=") {
					source = strings.TrimPrefix(field, "src=")
				}
				if strings.HasPrefix(field, "dst=") {
					target = strings.TrimPrefix(field, "dst=")
				}
			}
			if target == destination {
				return source
			}
		}
	}
	return ""
}

func testDigest(character byte) string {
	return "sha256:" + strings.Repeat(string(character), 64)
}
