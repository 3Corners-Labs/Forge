package change

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"testing"
	"time"
)

var contractNow = time.Date(2026, time.July, 25, 10, 0, 0, 0, time.UTC)

func TestChangeProposalBindsExactRequestAndSourceDigests(t *testing.T) {
	request := validRequest(t)
	proposal, err := BindProposal(Proposal{
		SchemaVersion: ProposalVersion, RequestID: request.ID, RequestDigest: request.Digest,
		SourceRevision: request.SourceRevision, PlanDigest: request.PlanDigest,
		ContextDigest: request.ContextManifestDigest,
		Replacements: []FileReplacement{{
			Path: "internal/example/example.go", ExpectedSourceDigest: request.Sources[0].Digest,
			Content:   "package example\n\nfunc Value() string { return \"new\" }\n",
			Rationale: "Satisfies AC-1.", CandidateTest: false,
		}},
		ProposedChecks: []string{"go test ./internal/example"}, Assumptions: []string{},
		Uncertainty: []string{}, Omissions: []string{}, Usage: Usage{InputTokens: 100, OutputTokens: 50, LatencyMs: 25},
		Receipt: Receipt{
			ActivationID: "activation-1", ActivationDigest: digest("activation"),
			Provider: "contract-double", ModelID: "model-1", AdapterVersion: "adapter-1",
			RoleProfileID: RoleProfileV1, PromptVersion: "prompt-1",
			ResponseSchema: ProposalVersion, ProviderRequest: "provider-request-1",
		},
	}, request)
	if err != nil {
		t.Fatal(err)
	}
	if err := proposal.Validate(request); err != nil {
		t.Fatal(err)
	}

	stale := proposal
	stale.Replacements[0].ExpectedSourceDigest = digest("stale")
	if _, err := BindProposal(stale, request); !errorsIs(err, ErrInvalidContract) {
		t.Fatalf("stale source digest err=%v", err)
	}
}

func TestChangeProposalRejectsUndeclaredPathsAuthorityAndIncompleteBytes(t *testing.T) {
	request := validRequest(t)
	base := Proposal{
		SchemaVersion: ProposalVersion, RequestID: request.ID, RequestDigest: request.Digest,
		SourceRevision: request.SourceRevision, PlanDigest: request.PlanDigest,
		ContextDigest: request.ContextManifestDigest,
		Replacements: []FileReplacement{{
			Path: "internal/example/example.go", ExpectedSourceDigest: request.Sources[0].Digest,
			Content: "package example\n", Rationale: "Satisfies AC-1.",
		}},
		Receipt: Receipt{
			ActivationID: "activation-1", ActivationDigest: digest("activation"),
			Provider: "contract-double", ModelID: "model-1", AdapterVersion: "adapter-1",
			RoleProfileID: RoleProfileV1, PromptVersion: "prompt-1",
			ResponseSchema: ProposalVersion, ProviderRequest: "provider-request-1",
		},
	}
	crossScope := base
	crossScope.Replacements = []FileReplacement{{Path: ".github/workflows/unsafe.yml", Content: "x", Rationale: "outside scope"}}
	if _, err := BindProposal(crossScope, request); !errorsIs(err, ErrInvalidContract) {
		t.Fatalf("cross-scope err=%v", err)
	}
	incomplete := base
	incomplete.Incomplete = true
	if _, err := BindProposal(incomplete, request); !errorsIs(err, ErrInvalidContract) {
		t.Fatalf("incomplete bytes err=%v", err)
	}
}

func TestChangeRequestRejectsInventedOrMismatchedSource(t *testing.T) {
	request := validRequest(t)
	request.ID, request.Digest = "", ""
	request.Sources[0].Digest = digest("other")
	if _, err := BindRequest(request); !errorsIs(err, ErrInvalidContract) {
		t.Fatalf("mismatched source err=%v", err)
	}
	request = validRequest(t)
	request.ID, request.Digest = "", ""
	request.Sources[0].Path = "../secret"
	if _, err := BindRequest(request); !errorsIs(err, ErrInvalidContract) {
		t.Fatalf("unsafe source err=%v", err)
	}
}

func validRequest(t *testing.T) Request {
	t.Helper()
	content := "package example\n\nfunc Value() string { return \"old\" }\n"
	value, err := BindRequest(Request{
		SchemaVersion: RequestVersion, CorrelationID: "correlation-1", TenantID: "tenant-1",
		WorkspaceID: "workspace-1", ProjectID: "project-1", RepositoryID: "repository-1",
		DeliveryRunID: "run-1", TaskAgentID: "agent-1", PlanID: "plan-1",
		PlanDigest: digest("plan"), ContextPackageID: "context-1",
		ContextManifestDigest: digest("context"), RequirementsRecordID: "requirements-1",
		SourceRevision: strings.Repeat("a", 40), RoleProfileID: RoleProfileV1,
		Purpose: PurposeV1, Classification: "Internal",
		AllowedReadPaths: []string{"internal/example"}, AllowedWritePaths: []string{"internal/example"},
		Sources:            []SourceFile{{Path: "internal/example/example.go", Digest: textDigest(content), Content: content}},
		AcceptanceCriteria: []string{"AC-1 returns the new value."},
		TestIntent:         []string{"Unit test the changed behavior."}, Exclusions: []string{"No deployment."},
		ControlReferences: []string{"ADR-0023"},
		Limits: Limits{
			MaxInputBytes: 1 << 20, MaxOutputBytes: 1 << 20, MaxTokens: 4096,
			MaxLatency: 20 * time.Second, MaxCostMicros: 100_000, MaxChangedFiles: 2,
		},
		ExpiresAt: contractNow.Add(time.Hour),
	})
	if err != nil {
		t.Fatal(err)
	}
	return value
}

func digest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func errorsIs(actual, target error) bool {
	return actual == target
}
