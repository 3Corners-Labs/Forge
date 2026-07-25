package verification

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
)

type verifierCommandFunc func(context.Context, string, ...string) error

func (function verifierCommandFunc) Run(ctx context.Context, name string, arguments ...string) error {
	return function(ctx, name, arguments...)
}

func TestVerifierDockerRunnerUsesSeparateReadOnlyCheckout(t *testing.T) {
	now := time.Now().UTC()
	bundle, err := workercontract.BindBundle(workercontract.CandidateBundle{
		SchemaVersion: workercontract.BundleVersion, SourceRevision: "revision-1", PlanDigest: verifierDigest('a'),
		Files: []workercontract.CandidateFile{{Path: "internal/example.go", Content: []byte("package example\n")}},
	})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := verificationcontract.BindPlan(verificationcontract.AssurancePlan{
		SchemaVersion: verificationcontract.AssurancePlanVersion, CandidateDigest: bundle.Digest,
		PlanDigest: bundle.PlanDigest, RiskClass: "low", Classification: "Internal",
		AffectedCapabilities: []string{"forge.capability.delivery"}, ChangedPaths: []string{"internal/example.go"},
		Controls:   []verificationcontract.Control{{ID: "forge.candidate-integrity", Version: "1", Consequential: true}},
		ProducedBy: "forge.worker.subject.test", Environment: "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "correlation-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	result, err := verificationcontract.BindResult(verificationcontract.Result{
		SchemaVersion: verificationcontract.ResultVersion, AssurancePlanID: plan.ID,
		AssurancePlanDigest: plan.Digest, CandidateDigest: bundle.Digest, Environment: plan.Environment,
		Findings: []verificationcontract.Finding{{
			ControlID: "forge.candidate-integrity", Outcome: "passed", Code: "passed",
			EvidenceDigest: verifierDigest('b'), Consequential: true,
		}},
		Outcome: "passed", Independent: true, VerifiedBy: "forge.verifier.test", CompletedAt: now, EvidenceCorrelationID: plan.EvidenceCorrelationID,
	})
	if err != nil {
		t.Fatal(err)
	}
	var invoked string
	fake := verifierCommandFunc(func(_ context.Context, name string, arguments ...string) error {
		invoked = name + " " + strings.Join(arguments, " ")
		output := verifierMountedSource(arguments, "/output")
		encoded, _ := json.Marshal(result)
		return os.WriteFile(filepath.Join(output, "result.json"), encoded, 0o600)
	})
	runner, err := NewDockerRunner(DockerConfig{
		Image:          "registry.example/forge-verifier@sha256:" + strings.Repeat("c", 64),
		TempRoot:       t.TempDir(),
		Memory:         "256m",
		CPUs:           "1",
		PIDs:           "8",
		MaxOutputBytes: 1 << 20,
	}, fake)
	if err != nil {
		t.Fatal(err)
	}
	got, err := runner.Verify(context.Background(), plan, bundle, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	for _, required := range []string{
		"--network none", "--read-only", "--cap-drop ALL", "no-new-privileges",
		"--user 65532:65532", "GOMAXPROCS=1", "dst=/workspace,readonly", "dst=/input,readonly",
	} {
		if !strings.Contains(invoked, required) {
			t.Fatalf("Docker invocation lacks %q: %s", required, invoked)
		}
	}
	if got.Digest != result.Digest {
		t.Fatalf("result digest=%q want %q", got.Digest, result.Digest)
	}
}

func TestVerifierDockerRunnerAcceptsImmutableLocalImageIDAndRejectsTag(t *testing.T) {
	config := DockerConfig{
		Image: "sha256:" + strings.Repeat("a", 64), TempRoot: t.TempDir(),
		Memory: "256m", CPUs: "1", PIDs: "8", MaxOutputBytes: 1 << 20,
	}
	if _, err := NewDockerRunner(config, verifierCommandFunc(func(context.Context, string, ...string) error { return nil })); err != nil {
		t.Fatalf("immutable local image ID rejected: %v", err)
	}
	config.Image = "forge-verifier:latest"
	if _, err := NewDockerRunner(config, verifierCommandFunc(func(context.Context, string, ...string) error { return nil })); err == nil {
		t.Fatal("mutable Docker tag was admitted")
	}
}

func TestVerifierDockerRunnerExecutesImmutableDeveloperTrialImage(t *testing.T) {
	image := os.Getenv("FORGE_TEST_VERIFIER_IMAGE")
	if image == "" {
		t.Skip("FORGE_TEST_VERIFIER_IMAGE is not configured")
	}
	now := time.Now().UTC()
	bundle, err := workercontract.BindBundle(workercontract.CandidateBundle{
		SchemaVersion: workercontract.BundleVersion, SourceRevision: "revision-integration",
		PlanDigest: verifierDigest('a'),
		Files: []workercontract.CandidateFile{{
			Path: "docs/result.md", Content: []byte("# Result\n"),
		}},
	})
	if err != nil {
		t.Fatal(err)
	}
	plan, err := verificationcontract.BindPlan(verificationcontract.AssurancePlan{
		SchemaVersion: verificationcontract.AssurancePlanVersion, CandidateDigest: bundle.Digest,
		PlanDigest: bundle.PlanDigest, RiskClass: "low", Classification: "Internal",
		AffectedCapabilities: []string{"forge.capability.delivery"}, ChangedPaths: []string{"docs/result.md"},
		Controls: []verificationcontract.Control{
			{ID: "forge.candidate-integrity", Version: "1", Consequential: true},
			{ID: "forge.declared-path-scope", Version: "1", Consequential: true},
		},
		ProducedBy: "forge.worker.subject.test", Environment: "developer-trial", ExpiresAt: now.Add(time.Hour),
		EvidenceCorrelationID: "correlation-integration",
	})
	if err != nil {
		t.Fatal(err)
	}
	runner, err := NewDockerRunner(DockerConfig{
		Image: image, TempRoot: t.TempDir(), Memory: "256m", CPUs: "1",
		PIDs: "8", MaxOutputBytes: 1 << 20,
	}, nil)
	if err != nil {
		t.Fatal(err)
	}
	result, err := runner.Verify(context.Background(), plan, bundle, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != "passed" || !result.Independent {
		t.Fatalf("unexpected verifier result: %#v", result)
	}
}

func verifierMountedSource(arguments []string, destination string) string {
	for index, argument := range arguments {
		if argument != "--mount" || index+1 >= len(arguments) {
			continue
		}
		var source, target string
		for _, field := range strings.Split(arguments[index+1], ",") {
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
	return ""
}

func verifierDigest(character byte) string {
	return "sha256:" + strings.Repeat(string(character), 64)
}
