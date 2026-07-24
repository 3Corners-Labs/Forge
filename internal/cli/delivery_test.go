package cli

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
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
	discoveryInputs []discoveryInput
	discoveryErrors map[int]error
	startInputs     []startInput
	proposalEdits   []proposalEditInput
}

func (stub *journeyStub) Discover(_ context.Context, input discoveryInput) (discoveryStage, error) {
	stub.calls = append(stub.calls, "discover")
	attempt := len(stub.discoveryInputs)
	stub.discoveryInputs = append(stub.discoveryInputs, input)
	if err, found := stub.discoveryErrors[attempt]; found {
		return discoveryStage{}, err
	}
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
func (stub *journeyStub) ReviseProposal(_ context.Context, sessionID, digest string, edit proposalEditInput) (proposalStage, error) {
	stub.calls = append(stub.calls, "revise")
	stub.proposalEdits = append(stub.proposalEdits, edit)
	if sessionID != "session-1" || digest == "" {
		return proposalStage{}, errors.New("current proposal not found")
	}
	return proposalStage{
		SessionID: sessionID, ProposalID: "proposal-1", ProposalVersion: 2,
		ProposalDigest: "sha256:revised", PreviousDigest: digest,
		ChangeExplanation: edit.ChangeExplanation, Privacy: []string{edit.Replacement},
	}, nil
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
	activationPath := filepath.Join(t.TempDir(), "approved-activation.json")
	activationDigest := "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	var receivedOptions journeyOptions
	err := execute([]string{
		"delivery", "start", "--workspace", "fixture-workspace", "--project", "forge",
		"--repository", "forge", "--message", "Help employees choose a desk",
		"--provider-activation", activationPath,
		"--provider-activation-digest", activationDigest, "--output", "json",
	}, strings.NewReader("Employees who come to the office\npromote\nyes\nyes\n"),
		&output, &bytes.Buffer{}, func(options journeyOptions) (guidedJourney, error) {
			receivedOptions = options
			return journey, nil
		})
	if err != nil {
		t.Fatalf("conversational delivery: %v\n%s", err, output.String())
	}
	if strings.Join(journey.calls, ",") != "discover,discover,start,promote,execute,publish" {
		t.Fatalf("unexpected conversational sequence: %v", journey.calls)
	}
	if receivedOptions.ProviderActivationPath != activationPath ||
		receivedOptions.ProviderActivationDigest != activationDigest {
		t.Fatalf("provider activation was not passed to the journey: %#v", receivedOptions)
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
		"Review Requirements Proposal [promote/edit/reject]",
		"Approve execution of this exact plan? [yes/no]",
		"Approve contract-double publication of this exact candidate? [yes/no]",
	} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("conversational output missing %q:\n%s", expected, output.String())
		}
	}
}

func TestRequirementsProposalReviewCanEditThenPromoteExactRevision(t *testing.T) {
	journey := &journeyStub{}
	input := bufio.NewReader(strings.NewReader(
		"edit\nprivacy\n1\nDocument which personal data is used and why\nReplace the absolute prohibition\npromote\n",
	))
	var output bytes.Buffer
	proposal, promote, err := reviewRequirementsProposal(
		context.Background(), journey, input, &output,
		proposalStage{
			SessionID: "session-1", ProposalID: "proposal-1", ProposalVersion: 1,
			ProposalDigest: "sha256:proposal", Privacy: []string{"No personal data"},
		},
	)
	if err != nil {
		t.Fatal(err)
	}
	if !promote || proposal.ProposalVersion != 2 || proposal.ProposalDigest != "sha256:revised" ||
		proposal.PreviousDigest != "sha256:proposal" || len(journey.proposalEdits) != 1 ||
		journey.proposalEdits[0].Field != "privacy" || journey.proposalEdits[0].Item != 1 {
		t.Fatalf("review result=%#v promote=%t edits=%#v", proposal, promote, journey.proposalEdits)
	}
	for _, expected := range []string{
		"Review Requirements Proposal [promote/edit/reject]",
		`"proposalVersion":2`,
		`"previousDigest":"sha256:proposal"`,
		`"state":"awaiting-promotion"`,
	} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("review output missing %q:\n%s", expected, output.String())
		}
	}
}

func TestRequirementsProposalReviewCanRejectWithoutPromotion(t *testing.T) {
	journey := &journeyStub{}
	var output bytes.Buffer
	_, promote, err := reviewRequirementsProposal(
		context.Background(), journey, bufio.NewReader(strings.NewReader("reject\n")), &output,
		proposalStage{SessionID: "session-1", ProposalDigest: "sha256:proposal"},
	)
	if err != nil || promote || len(journey.calls) != 0 ||
		!strings.Contains(output.String(), `"state":"rejected"`) {
		t.Fatalf("reject promote=%t calls=%v err=%v output=%s", promote, journey.calls, err, output.String())
	}
}

func TestConversationalDeliveryPreservesCheckpointAcrossBoundedRetry(t *testing.T) {
	journey := &journeyStub{
		discoveryErrors: map[int]error{
			1: failure.New(
				failure.DependencyUnavailable,
				"model-assisted discovery timed out",
				"correlation-1",
				failure.RetryBounded,
			),
		},
		discoveryStages: []discoveryStage{
			{
				SessionID: "discovery-1",
				Response:  "Who is the primary user?",
				Draft: discovery.Draft{
					Outcome:     "Find a suitable desk quickly",
					Constraints: []string{"advice only"},
				},
				Questions: []discovery.Question{{
					ID: "q-user", Text: "Who is the primary user?",
					Consequence: "The main actor remains unclear", Material: true,
				}},
			},
			{
				SessionID: "discovery-1",
				Response:  "I can now propose a bounded trial.",
				Draft: discovery.Draft{
					Outcome: "Find a suitable desk quickly", AffectedUsers: []string{"employee"},
					Problem: "Desk choice is slow", Constraints: []string{"advice only"},
					AcceptanceCriteria: []string{"recommendation within three seconds"},
				},
			},
		},
	}
	var output bytes.Buffer
	err := runDelivery([]string{
		"start", "--workspace", "fixture-workspace", "--project", "forge",
		"--repository", "forge", "--message", "Help employees choose a desk",
		"--provider-activation", "approved-activation.json",
		"--provider-activation-digest", "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		"--output", "json",
	}, bufio.NewReader(strings.NewReader(
		"Employees who visit the office\nyes\npromote\nyes\nyes\n",
	)), &output, &bytes.Buffer{}, func(journeyOptions) (guidedJourney, error) {
		return journey, nil
	})
	if err != nil {
		t.Fatalf("retryable conversational delivery: %v\n%s", err, output.String())
	}
	if len(journey.discoveryInputs) != 3 {
		t.Fatalf("discovery attempts = %d, want 3", len(journey.discoveryInputs))
	}
	if strings.Count(output.String(), "Forge asks [q-user]") != 1 {
		t.Fatalf("retry repeated an answered question:\n%s", output.String())
	}
	failed := journey.discoveryInputs[1]
	retried := journey.discoveryInputs[2]
	if failed.SessionID != "discovery-1" || retried.SessionID != failed.SessionID ||
		!slices.Equal(retried.UserStatements, failed.UserStatements) ||
		!slices.Equal(retried.ProvenanceRefs, failed.ProvenanceRefs) ||
		retried.CurrentDraft.Outcome != failed.CurrentDraft.Outcome ||
		len(retried.PreviousQuestions) != len(failed.PreviousQuestions) {
		t.Fatalf("retry did not preserve checkpoint:\nfailed=%#v\nretried=%#v", failed, retried)
	}
	for _, expected := range []string{
		`"type":"discovery-paused"`,
		`"state":"awaiting-retry"`,
		`"retry":"bounded"`,
		"Retry the current discovery turn? Existing answers remain in this session. [yes/no]",
		`"state":"awaiting-promotion"`,
	} {
		if !strings.Contains(output.String(), expected) {
			t.Fatalf("retry output missing %q:\n%s", expected, output.String())
		}
	}
}

func TestConversationalDeliveryDoesNotRetryWithoutBoundedRetryClassification(t *testing.T) {
	expected := failure.New(
		failure.BudgetExhausted,
		"model-assisted discovery budget is exhausted",
		"correlation-1",
		failure.RetryNewAuthority,
	)
	journey := &journeyStub{discoveryErrors: map[int]error{0: expected}}
	var output bytes.Buffer
	_, _, _, err := discoverIntent(
		context.Background(),
		journey,
		bufio.NewReader(strings.NewReader("yes\n")),
		&output,
		discoveryInput{
			Workspace: "fixture-workspace", Project: "forge", Repository: "forge",
			UserStatements: []string{"Help employees choose a desk"},
			ProvenanceRefs: []string{"interaction:user:1"},
		},
	)
	var normalized *failure.Error
	if !errors.As(err, &normalized) || normalized != expected {
		t.Fatalf("non-retryable failure changed: %v", err)
	}
	if len(journey.discoveryInputs) != 1 ||
		strings.Contains(output.String(), `"type":"discovery-paused"`) {
		t.Fatalf("non-retryable failure was retried or paused:\n%s", output.String())
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
			ActivationID:     "activation-contract-double",
			ActivationDigest: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			Provider:         "contract-double", ModelID: "discovery-double", AdapterVersion: "1",
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
	value, err := newDeveloperTrialJourney(journeyOptions{})
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
	}, bufio.NewReader(strings.NewReader("Employees who visit the office\npromote\nyes\nyes\n")),
		&output, &bytes.Buffer{}, func(journeyOptions) (guidedJourney, error) { return journey, nil })
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
	value, err := newDeveloperTrialJourney(journeyOptions{})
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

func TestDeveloperTrialLoadsExactProviderActivationWithoutResolvingCredential(t *testing.T) {
	t.Setenv("FORGE_OPENAI_API_KEY", "")
	now := time.Now().UTC()
	path := filepath.Join(t.TempDir(), "activation.json")
	raw := fmt.Sprintf(`{
  "schemaVersion": "forge.model-provider-activation.v1",
  "id": "forge.provider-activation.openai-terra-developer-trial-v1",
  "provider": "openai",
  "endpoint": "https://eu.api.openai.com/v1/responses",
  "adapterVersion": "forge.openai-discovery-adapter.v1",
  "allowedModelIds": ["gpt-5.6-terra"],
  "allowedPurposes": ["forge.product.discovery"],
  "allowedClassifications": ["Internal"],
  "credentialRef": "env:FORGE_OPENAI_API_KEY",
  "humanApprovalRef": "approval:provider-trial-1",
  "promptVersion": "forge.discovery.prompt.v1",
  "retention": "provider-default-up-to-30-days",
  "trainingUse": "disabled",
  "region": "eu",
  "storeProviderState": false,
  "reasoningEffort": "medium",
  "pricing": {
    "currency": "USD",
    "inputMicrosPerMillionTokens": 2500000,
    "outputMicrosPerMillionTokens": 15000000
  },
  "budgets": {
    "maxInputBytes": 32768,
    "maxOutputBytes": 32768,
    "maxTokens": 4096,
    "maxLatencyMs": 45000,
    "maxCostMicros": 250000,
    "maxTurns": 8,
    "maxSessionTokens": 32768,
    "maxSessionCostMicros": 1000000,
    "maxProjectTokensPerDay": 131072,
    "maxProjectCostMicrosPerDay": 4000000,
    "maxConcurrentRequests": 1
  },
  "activatedAt": %q,
  "expiresAt": %q
}`, now.Add(-time.Minute).Format(time.RFC3339), now.Add(time.Hour).Format(time.RFC3339))
	if err := os.WriteFile(path, []byte(raw), 0o600); err != nil {
		t.Fatal(err)
	}
	loaded, err := intelligence.LoadProviderActivation(path, now)
	if err != nil {
		t.Fatal(err)
	}
	value, err := newDeveloperTrialJourney(journeyOptions{
		ProviderActivationPath:   path,
		ProviderActivationDigest: loaded.Digest,
	})
	if err != nil {
		t.Fatal(err)
	}
	journey := value.(*developerTrialJourney)
	if journey.discovery == nil {
		t.Fatal("approved provider activation was not wired into discovery")
	}
	if journey.discoveryLimits != loaded.Budgets ||
		journey.discoveryLimits.MaxLatency != 45*time.Second {
		t.Fatalf(
			"approved provider budgets were not wired into Experience: got=%#v want=%#v",
			journey.discoveryLimits,
			loaded.Budgets,
		)
	}
	if _, err := newDeveloperTrialJourney(journeyOptions{
		ProviderActivationPath:   path,
		ProviderActivationDigest: "sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
	}); err == nil || !strings.Contains(err.Error(), "separately approved digest") {
		t.Fatalf("activation digest substitution was accepted: %v", err)
	}
}

func testDiscoveryActivation(now time.Time) intelligence.ProviderActivation {
	return intelligence.ProviderActivation{
		SchemaVersion: intelligence.ProviderActivationVersion,
		ID:            "activation-contract-double",
		Digest:        "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		Provider:      "contract-double",
		Endpoint:      "contract-double://discovery", AdapterVersion: "1",
		AllowedModelIDs: []string{"discovery-double"},
		AllowedPurposes: []string{discovery.PurposeV1}, AllowedClassifications: []string{"Internal"},
		CredentialRef: "none:contract-double", HumanApprovalRef: "approval:test",
		PromptVersion: "forge.discovery.prompt.v1", Retention: "none", TrainingUse: "disabled",
		ReasoningEffort: "none", Pricing: intelligence.Pricing{Currency: "USD"},
		Budgets: discovery.Limits{
			MaxInputBytes: 32 << 10, MaxOutputBytes: 32 << 10, MaxTokens: 4096,
			MaxLatency: 20 * time.Second, MaxCostMicros: 250_000,
		},
		MaxTurns:         8,
		MaxSessionTokens: 32768, MaxSessionCostMicros: 1_000_000,
		MaxProjectTokensPerDay: 131072, MaxProjectCostPerDay: 4_000_000,
		MaxConcurrentRequests: 1,
		Region:                "local", ActivatedAt: now.Add(-time.Minute), ExpiresAt: now.Add(time.Hour),
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
	}, input, &output, &errorOutput, func(journeyOptions) (guidedJourney, error) { return journey, nil })
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
	}, strings.NewReader("promote sha256:wrong\n"), &output, &bytes.Buffer{}, func(journeyOptions) (guidedJourney, error) { return journey, nil })
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
	value, err := newDeveloperTrialJourney(journeyOptions{})
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
