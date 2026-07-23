package cli

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/discovery"
	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/intelligence"
)

type journeyStub struct {
	calls           []string
	discoveryStages []discoveryStage
	discoveryIndex  int
	startInputs     []startInput
}

func (stub *journeyStub) Discover(context.Context, discoveryInput) (discoveryStage, error) {
	stub.calls = append(stub.calls, "discover")
	if stub.discoveryIndex >= len(stub.discoveryStages) {
		return discoveryStage{}, errors.New("no discovery stage configured")
	}
	value := stub.discoveryStages[stub.discoveryIndex]
	stub.discoveryIndex++
	return value, nil
}
func (stub *journeyStub) Start(_ context.Context, input startInput) (proposalStage, error) {
	stub.calls = append(stub.calls, "start")
	stub.startInputs = append(stub.startInputs, input)
	return proposalStage{SessionID: "session-1", ProposalID: "proposal-1", ProposalDigest: "sha256:proposal", Outcome: "Test Forge"}, nil
}

func TestConversationalDeliveryClarifiesBeforeProposalAndUsesMeaningfulConfirmations(t *testing.T) {
	journey := &journeyStub{discoveryStages: []discoveryStage{
		{
			SessionID: "discovery-1", Response: "Who is the primary user?",
			Questions: []discovery.Question{{
				ID: "q-user", Text: "Who is the primary user?",
				Consequence: "The main actor remains unclear", Material: true,
			}},
			Provider: "contract-double", ModelID: "discovery-double",
		},
		{
			SessionID: "discovery-1", Response: "I can now propose a bounded trial.",
			Draft: discovery.Draft{
				Outcome: "Find a suitable desk quickly", AffectedUsers: []string{"employee"},
				Problem: "Desk choice is slow", Constraints: []string{"advice only"},
				AcceptanceCriteria: []string{"recommendation within three seconds"},
				Hypotheses:         []string{"Advice reduces search time"},
				ProductAreaRefs:    []string{"workplace.discovery"},
			},
			Questions: []discovery.Question{{
				ID: "q-copy", Text: "Should recommendations be copyable?", Material: false,
			}},
			Provider: "contract-double", ModelID: "discovery-double",
		},
	}}
	var output bytes.Buffer
	err := execute([]string{
		"delivery", "start", "--workspace", "fixture-workspace", "--project", "forge",
		"--repository", "forge", "--message", "Help employees choose a desk", "--output", "json",
	}, strings.NewReader("Employees who come to the office\nyes\nyes\nyes\n"),
		&output, &bytes.Buffer{}, func() (guidedJourney, error) { return journey, nil })
	if err != nil {
		t.Fatalf("conversational delivery: %v\n%s", err, output.String())
	}
	if strings.Join(journey.calls, ",") != "discover,discover,start,promote,execute,publish" {
		t.Fatalf("unexpected conversational sequence: %v", journey.calls)
	}
	if len(journey.startInputs) != 1 ||
		journey.startInputs[0].Intent.Outcome != "Find a suitable desk quickly" ||
		journey.startInputs[0].Intent.Hypotheses[0] != "Advice reduces search time" ||
		journey.startInputs[0].Intent.ProductAreaRefs[0] != "workplace.discovery" ||
		journey.startInputs[0].Intent.Questions[0].ID != "q-copy" ||
		!strings.Contains(journey.startInputs[0].RawInput, "Employees who come to the office") {
		t.Fatalf("discovery was not converted into proposal input: %#v", journey.startInputs)
	}
	for _, expected := range []string{
		`"state":"awaiting-clarification"`, `"state":"proposal-ready"`,
		"Promote this exact Requirements Proposal? [yes/no]",
		"Approve execution of this exact plan? [yes/no]",
		"Approve contract-double publication of this exact candidate? [yes/no]",
	} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("conversational output missing %q:\n%s", expected, output.String())
		}
	}
}

type scriptedDiscoveryAdapter struct {
	calls int
}

func (*scriptedDiscoveryAdapter) Provider() string { return "contract-double" }

func (adapter *scriptedDiscoveryAdapter) Propose(_ context.Context, _ intelligence.ProviderActivation, request discovery.TurnRequest) (discovery.TurnProposal, error) {
	adapter.calls++
	proposal := discovery.TurnProposal{
		SchemaVersion: discovery.ProposalVersion, SessionID: request.SessionID,
		CorrelationID: request.CorrelationID, ProvenanceRefs: request.ProvenanceRefs,
		Usage: discovery.Usage{InputTokens: 100, OutputTokens: 50, LatencyMs: 20, CostMicros: 100},
		Receipt: discovery.Receipt{
			Provider: "contract-double", ModelID: "discovery-double", AdapterVersion: "1",
			RoleProfileID: request.RoleProfileID, PromptVersion: "forge.discovery.prompt.v1",
			ResponseSchema: discovery.ProposalVersion, ProviderRequest: "provider-request",
		},
	}
	if adapter.calls == 1 {
		proposal.Response = "Who is the primary user?"
		proposal.Incomplete = true
		proposal.Questions = []discovery.Question{{
			ID: "q-user", Text: "Who is the primary user?",
			Consequence: "The main actor remains unclear", Material: true,
		}}
		return proposal, nil
	}
	proposal.Response = "I can now propose a bounded desk-advice trial."
	proposal.Draft = discovery.Draft{
		Outcome: "Find a suitable desk quickly", AffectedUsers: []string{"employee"},
		Problem: "Desk choice is slow", FunctionalBehavior: []string{"Recommend a suitable desk"},
		UXExpectations: []string{"Immediate interaction feedback"},
		Accessibility:  []string{"Keyboard-operable flow"}, Privacy: []string{"Opt-in colleague visibility"},
		Security: []string{"No production data"}, Operations: []string{"Bounded retry"},
		Constraints: []string{"Advice only"}, Exclusions: []string{"Desk reservation"},
		AcceptanceCriteria: []string{"Recommendation within three seconds"},
		Inferences:         []string{"Office visitors choose their own desk"},
		Assumptions:        []string{"Synthetic availability data"}, Hypotheses: []string{"Advice reduces search time"},
		Experiments: []string{"Disposable desk-finder mockup"}, ProductAreaRefs: []string{"workplace.discovery"},
	}
	return proposal, nil
}

func TestDeveloperTrialConversationalJourneyUsesIntelligenceAndReachesReviewableState(t *testing.T) {
	t.Setenv("FORGE_IDENTITY_PROFILE", "developer-allowed")
	value, err := newDeveloperTrialJourney()
	if err != nil {
		t.Fatal(err)
	}
	journey := value.(*developerTrialJourney)
	adapter := &scriptedDiscoveryAdapter{}
	advisor, err := intelligence.NewService(testDiscoveryActivation(journey.now()), adapter, journey.now)
	if err != nil {
		t.Fatal(err)
	}
	if err := journey.configureDiscovery(advisor); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	err = runDelivery([]string{
		"start", "--workspace", "fixture-workspace", "--project", "forge",
		"--repository", "forge", "--message", "Help employees choose a desk", "--output", "json",
	}, bufio.NewReader(strings.NewReader("Employees who visit the office\nyes\nyes\nyes\n")),
		&output, &bytes.Buffer{}, func() (guidedJourney, error) { return journey, nil })
	if err != nil {
		t.Fatalf("conversational Developer Trial: %v\n%s", err, output.String())
	}
	if adapter.calls != 2 || !strings.Contains(output.String(), `"state":"reviewable"`) ||
		!strings.Contains(output.String(), `"provider":"contract-double"`) ||
		!strings.Contains(output.String(), `"hypotheses":["Advice reduces search time"]`) ||
		!strings.Contains(output.String(), `"productAreaRefs":["workplace.discovery"]`) {
		t.Fatalf("unexpected conversational journey: calls=%d\n%s", adapter.calls, output.String())
	}
}

func TestDeniedConversationalJourneyStopsBeforeRepositoryAndModelDiscovery(t *testing.T) {
	t.Setenv("FORGE_IDENTITY_PROFILE", "observer-denied")
	value, err := newDeveloperTrialJourney()
	if err != nil {
		t.Fatal(err)
	}
	journey := value.(*developerTrialJourney)
	adapter := &scriptedDiscoveryAdapter{}
	advisor, err := intelligence.NewService(testDiscoveryActivation(journey.now()), adapter, journey.now)
	if err != nil {
		t.Fatal(err)
	}
	if err := journey.configureDiscovery(advisor); err != nil {
		t.Fatal(err)
	}
	_, err = journey.Discover(context.Background(), discoveryInput{
		Workspace: "fixture-workspace", Project: "forge", Repository: "forge",
		UserStatements: []string{"show me the repository"},
		ProvenanceRefs: []string{"interaction:user:1"},
	})
	var normalized *failure.Error
	if !errors.As(err, &normalized) || normalized.Code != failure.Forbidden {
		t.Fatalf("unexpected denial: %v", err)
	}
	if journey.bindingSource.calls != 0 || adapter.calls != 0 {
		t.Fatalf("denied discovery crossed a boundary: repository=%d model=%d", journey.bindingSource.calls, adapter.calls)
	}
}

func testDiscoveryActivation(now time.Time) intelligence.ProviderActivation {
	return intelligence.ProviderActivation{
		SchemaVersion: intelligence.ProviderActivationVersion,
		ID:            "activation-contract-double", Provider: "contract-double",
		Endpoint: "contract-double://discovery", AdapterVersion: "1",
		AllowedModelIDs: []string{"discovery-double"},
		AllowedPurposes: []string{discovery.PurposeV1}, AllowedClassifications: []string{"Internal"},
		CredentialRef: "none:contract-double", HumanApprovalRef: "approval:test",
		PromptVersion: "forge.discovery.prompt.v1", Retention: "none", TrainingUse: "disabled",
		Budgets: discovery.Limits{
			MaxInputBytes: 32 << 10, MaxOutputBytes: 32 << 10, MaxTokens: 4096,
			MaxLatency: 20 * time.Second, MaxCostMicros: 250_000,
		},
		MaxTurns: 8,
		Region:   "local", ActivatedAt: now.Add(-time.Minute), ExpiresAt: now.Add(time.Hour),
	}
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
