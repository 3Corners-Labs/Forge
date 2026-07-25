package verifier

import (
	"context"
	"errors"
	"os"
	"path/filepath"
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
		ProducedBy: "forge.worker.subject.test", Environment: "developer-trial", ExpiresAt: now.Add(time.Hour), EvidenceCorrelationID: "corr",
	})
	if err != nil {
		t.Fatal(err)
	}
	engine, _ := NewEngine(func() time.Time { return now }, DefaultIdentity)
	result, err := engine.Verify(context.Background(), plan, bundle, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != "passed" || !result.Independent {
		t.Fatal("independent verification did not pass")
	}
}

// TestVerifierRefusesPlanWithoutApplicableOracle covers the case that motivated
// these controls: a documentation candidate whose only independent controls are
// integrity preconditions incapable of failing. Certifying it would be an empty
// assurance, so verification fails closed instead.
func TestVerifierRefusesPlanWithoutApplicableOracle(t *testing.T) {
	workspace := t.TempDir()
	bundle := mustBundle(t, "docs/notes.md", "# notes\n")
	plan := mustPlan(t, bundle, []verificationcontract.Control{
		{ID: ControlCandidateIntegrity, Version: "1", Consequential: true},
		{ID: ControlDeclaredPathScope, Version: "1", Consequential: true},
	})
	engine, _ := NewEngine(time.Now().UTC, DefaultIdentity)
	if _, err := engine.Verify(context.Background(), plan, bundle, workspace); !errors.Is(err, ErrNoApplicableOracle) {
		t.Fatalf("expected the plan to be refused, got %v", err)
	}
}

func TestDocumentMetadataControlFailsWithoutLifecycleStatus(t *testing.T) {
	for name, content := range map[string]string{
		"missing status": "# Title\n\nOwner: Forge\n",
		"invalid status": "# Title\n\nStatus: Informative\nOwner: Forge\n",
		"missing owner":  "# Title\n\nStatus: Accepted\n",
	} {
		t.Run(name, func(t *testing.T) {
			bundle := mustBundle(t, "docs/governance/example.md", content)
			plan := mustPlan(t, bundle, []verificationcontract.Control{
				{ID: ControlDocumentMetadata, Version: "1", Consequential: true},
			})
			engine, _ := NewEngine(time.Now().UTC, DefaultIdentity)
			result, err := engine.Verify(context.Background(), plan, bundle, t.TempDir())
			if err != nil {
				t.Fatal(err)
			}
			if result.Outcome != "failed" {
				t.Fatalf("expected failure, got %q", result.Outcome)
			}
		})
	}
}

func TestDocumentMetadataControlAcceptsDeclaredDocument(t *testing.T) {
	bundle := mustBundle(t, "docs/governance/example.md", "# Title\n\nStatus: Accepted\nOwner: Forge Architecture\n")
	plan := mustPlan(t, bundle, []verificationcontract.Control{
		{ID: ControlDocumentMetadata, Version: "1", Consequential: true},
	})
	engine, _ := NewEngine(time.Now().UTC, DefaultIdentity)
	result, err := engine.Verify(context.Background(), plan, bundle, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != "passed" {
		t.Fatalf("expected pass, got %q", result.Outcome)
	}
}

func TestMarkdownLinksControlFailsOnUnresolvableLink(t *testing.T) {
	bundle := mustBundle(t, "docs/notes.md", "See [the plan](does-not-exist.md).\n")
	plan := mustPlan(t, bundle, []verificationcontract.Control{
		{ID: ControlMarkdownLinks, Version: "1", Consequential: true},
	})
	engine, _ := NewEngine(time.Now().UTC, DefaultIdentity)
	result, err := engine.Verify(context.Background(), plan, bundle, t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != "failed" {
		t.Fatalf("expected failure, got %q", result.Outcome)
	}
}

// TestMarkdownLinksControlResolvesAgainstSourceAndCandidate proves the control
// consults the independently reconstructed checkout as well as the candidate,
// and ignores external URLs and pure fragments.
func TestMarkdownLinksControlResolvesAgainstSourceAndCandidate(t *testing.T) {
	workspace := t.TempDir()
	if err := os.MkdirAll(filepath.Join(workspace, "docs"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(workspace, "docs", "existing.md"), []byte("# existing\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	content := "[source](existing.md) [sibling](sibling.md) [web](https://example.org) [anchor](#section)\n"
	bundle, err := workercontract.BindBundle(workercontract.CandidateBundle{
		SchemaVersion: workercontract.BundleVersion, SourceRevision: "commit",
		PlanDigest: fixedDigest('a'),
		Files: []workercontract.CandidateFile{
			{Path: "docs/notes.md", Content: []byte(content)},
			{Path: "docs/sibling.md", Content: []byte("# sibling\n")},
		},
	})
	if err != nil {
		t.Fatal(err)
	}
	plan := mustPlan(t, bundle, []verificationcontract.Control{
		{ID: ControlMarkdownLinks, Version: "1", Consequential: true},
	})
	engine, _ := NewEngine(time.Now().UTC, DefaultIdentity)
	result, err := engine.Verify(context.Background(), plan, bundle, workspace)
	if err != nil {
		t.Fatal(err)
	}
	if result.Outcome != "passed" {
		t.Fatalf("expected pass, got %q", result.Outcome)
	}
}

func mustBundle(t *testing.T, path, content string) workercontract.CandidateBundle {
	t.Helper()
	bundle, err := workercontract.BindBundle(workercontract.CandidateBundle{
		SchemaVersion: workercontract.BundleVersion, SourceRevision: "commit",
		PlanDigest: fixedDigest('a'),
		Files:      []workercontract.CandidateFile{{Path: path, Content: []byte(content)}},
	})
	if err != nil {
		t.Fatal(err)
	}
	return bundle
}

func mustPlan(t *testing.T, bundle workercontract.CandidateBundle, controls []verificationcontract.Control) verificationcontract.AssurancePlan {
	t.Helper()
	paths := make([]string, 0, len(bundle.Files))
	for _, file := range bundle.Files {
		paths = append(paths, file.Path)
	}
	plan, err := verificationcontract.BindPlan(verificationcontract.AssurancePlan{
		SchemaVersion:   verificationcontract.AssurancePlanVersion,
		CandidateDigest: bundle.Digest, PlanDigest: bundle.PlanDigest, RiskClass: "low",
		Classification: "internal", AffectedCapabilities: []string{"forge.capability.delivery"},
		ChangedPaths: paths, Controls: controls, ProducedBy: "forge.worker.subject.test",
		Environment: "developer-trial",
		ExpiresAt:   time.Now().UTC().Add(time.Hour), EvidenceCorrelationID: "corr",
	})
	if err != nil {
		t.Fatal(err)
	}
	return plan
}

func fixedDigest(value byte) string {
	result := make([]byte, 64)
	for index := range result {
		result[index] = value
	}
	return "sha256:" + string(result)
}
