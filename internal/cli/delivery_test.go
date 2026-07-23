package cli

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/3corners/forge/internal/contracts/failure"
)

type journeyStub struct {
	calls []string
}

func (stub *journeyStub) Start(context.Context, startInput) (proposalStage, error) {
	stub.calls = append(stub.calls, "start")
	return proposalStage{SessionID: "session-1", ProposalID: "proposal-1", ProposalDigest: "sha256:proposal", Outcome: "Test Forge"}, nil
}
func (stub *journeyStub) PromoteAndPlan(context.Context, string, string) (planStage, error) {
	stub.calls = append(stub.calls, "promote")
	return planStage{SessionID: "session-1", RunID: "run-1", PlanID: "plan-1", PlanDigest: "sha256:plan", RiskClass: "low", Checks: []string{"go test ./..."}}, nil
}
func (stub *journeyStub) ApprovePlanAndVerify(context.Context, string, string) (candidateStage, error) {
	stub.calls = append(stub.calls, "execute")
	return candidateStage{SessionID: "session-1", RunID: "run-1", CandidateDigest: "sha256:candidate", VerificationDigest: "sha256:verification", ChangedPaths: []string{"docs/example.md"}, EvidenceReady: true}, nil
}
func (stub *journeyStub) ApprovePublicationAndPublish(context.Context, string, string) (publicationStage, error) {
	stub.calls = append(stub.calls, "publish")
	return publicationStage{SessionID: "session-1", RunID: "run-1", OperationID: "operation-1", Branch: "forge/run-1-candidate", PullRequestID: "pr-1", PullRequestURL: "https://contract-double.invalid/pr-1", Provider: "contract-double"}, nil
}

func TestDeliveryStartRequiresThreeExactHumanActions(t *testing.T) {
	intentFile := writeIntentFixture(t)
	journey := &journeyStub{}
	input := strings.NewReader("promote sha256:proposal\napprove-plan sha256:plan\napprove-publication sha256:candidate\n")
	var output, errorOutput bytes.Buffer
	err := execute([]string{
		"delivery", "start", "--workspace", "fixture-workspace", "--project", "forge",
		"--repository", "forge", "--intent-file", intentFile, "--output", "json",
	}, input, &output, &errorOutput, func() (guidedJourney, error) { return journey, nil })
	if err != nil {
		t.Fatalf("delivery start: %v\n%s", err, errorOutput.String())
	}
	if strings.Join(journey.calls, ",") != "start,promote,execute,publish" {
		t.Fatalf("unexpected capability sequence: %v", journey.calls)
	}
	for _, expected := range []string{
		`"state":"awaiting-promotion"`, `"state":"awaiting-plan-approval"`,
		`"state":"awaiting-publication-approval"`, `"state":"reviewable"`,
		`"mergePermitted":false`,
	} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("output missing %q:\n%s", expected, output.String())
		}
	}
}

func TestDeliveryStartStopsBeforePlanForWrongProposalDigest(t *testing.T) {
	intentFile := writeIntentFixture(t)
	journey := &journeyStub{}
	var output bytes.Buffer
	err := execute([]string{
		"delivery", "start", "--workspace", "fixture-workspace", "--project", "forge",
		"--repository", "forge", "--intent-file", intentFile, "--output", "json",
	}, strings.NewReader("promote sha256:wrong\n"), &output, &bytes.Buffer{}, func() (guidedJourney, error) { return journey, nil })
	if err == nil || strings.Join(journey.calls, ",") != "start" {
		t.Fatalf("wrong digest advanced journey: calls=%v err=%v", journey.calls, err)
	}
}

type promptReader struct {
	output *bytes.Buffer
	offset int
}

func (reader *promptReader) Read(target []byte) (int, error) {
	value := reader.output.String()
	if reader.offset >= len(value) {
		return 0, io.EOF
	}
	remaining := value[reader.offset:]
	marker := " by entering: "
	index := strings.Index(remaining, marker)
	if index < 0 {
		return 0, io.EOF
	}
	start := index + len(marker)
	end := strings.IndexByte(remaining[start:], '\n')
	if end < 0 {
		return 0, io.EOF
	}
	answer := remaining[start:start+end] + "\n"
	reader.offset += start + end + 1
	return copy(target, answer), nil
}

func TestDeveloperTrialPublicJourneyReachesReviewableContractDoublePullRequest(t *testing.T) {
	t.Setenv("FORGE_IDENTITY_PROFILE", "developer-allowed")
	intentFile := writeIntentFixture(t)
	var output bytes.Buffer
	input := &promptReader{output: &output}
	err := execute([]string{
		"delivery", "start", "--workspace", "fixture-workspace", "--project", "forge",
		"--repository", "forge", "--intent-file", intentFile, "--output", "json",
	}, input, &output, &bytes.Buffer{}, newDeveloperTrialJourney)
	if err != nil {
		t.Fatalf("full Developer Trial journey: %v\n%s", err, output.String())
	}
	for _, expected := range []string{
		`"state":"awaiting-promotion"`, `"state":"awaiting-plan-approval"`,
		`"state":"awaiting-publication-approval"`, `"state":"reviewable"`,
		`"provider":"contract-double"`, `"mergePermitted":false`,
	} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("journey output missing %q:\n%s", expected, output.String())
		}
	}
	if strings.Contains(output.String(), "FORGE_TEST_SECRET_DO_NOT_RECORD") {
		t.Fatal("sensitive fixture marker leaked into public output")
	}
}

func TestDeniedDeveloperTrialJourneyStopsBeforeRepositoryDiscovery(t *testing.T) {
	t.Setenv("FORGE_IDENTITY_PROFILE", "observer-denied")
	value, err := newDeveloperTrialJourney()
	if err != nil {
		t.Fatal(err)
	}
	journey := value.(*developerTrialJourney)
	_, err = journey.Start(context.Background(), startInput{
		Workspace: "fixture-workspace", Project: "forge", Repository: "forge",
		Intent: intentInput{
			Outcome: "unauthorized", Problem: "must fail", AffectedUsers: []string{"observer"},
			Constraints: []string{"no discovery"}, AcceptanceCriteria: []string{"forbidden"},
		},
	})
	var normalized *failure.Error
	if err == nil || !errors.As(err, &normalized) || normalized.Code != failure.Forbidden {
		t.Fatalf("denial=%v", err)
	}
	if journey.bindingSource.calls != 0 {
		t.Fatalf("repository discovery occurred before denial: calls=%d", journey.bindingSource.calls)
	}
}

func writeIntentFixture(t *testing.T) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "intent.json")
	content := `{
  "outcome": "Show the guided Forge journey",
  "affectedUsers": ["Forge developer"],
  "problem": "The user cannot yet experience the M1 trust path",
  "functionalBehavior": ["Create one reviewable contract-double pull request"],
  "uxExpectations": ["Show immediate progress and explicit gates"],
  "accessibility": ["Keyboard-operable prompts"],
  "privacy": ["Use synthetic data only"],
  "security": ["No live credentials or merge"],
  "operations": ["Emit bounded progress"],
  "constraints": ["Developer Trial only"],
  "exclusions": ["Release", "Deployment"],
  "acceptanceCriteria": ["Independent verification passes"],
  "assumptions": ["The contract double is available"],
  "experiments": ["Run the deterministic CLI journey"]
}`
	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}
