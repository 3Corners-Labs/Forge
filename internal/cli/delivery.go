package cli

import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/3corners/forge/internal/contracts/discovery"
	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/intent"
)

type journeyOptions struct {
	ProviderActivationPath   string
	ProviderActivationDigest string
}

type journeyFactory func(journeyOptions) (guidedJourney, error)

type startInput struct {
	Workspace  string
	Project    string
	Repository string
	RawInput   string
	Provenance []string
	Intent     intentInput
}

type intentInput struct {
	Outcome            string            `json:"outcome"`
	AffectedUsers      []string          `json:"affectedUsers"`
	Problem            string            `json:"problem"`
	FunctionalBehavior []string          `json:"functionalBehavior"`
	UXExpectations     []string          `json:"uxExpectations"`
	Accessibility      []string          `json:"accessibility"`
	Privacy            []string          `json:"privacy"`
	Security           []string          `json:"security"`
	Operations         []string          `json:"operations"`
	Constraints        []string          `json:"constraints"`
	Exclusions         []string          `json:"exclusions"`
	AcceptanceCriteria []string          `json:"acceptanceCriteria"`
	Inferences         []string          `json:"inferences"`
	Assumptions        []string          `json:"assumptions"`
	Hypotheses         []string          `json:"hypotheses"`
	Questions          []intent.Question `json:"questions"`
	Experiments        []string          `json:"experiments"`
	ProductAreaRefs    []string          `json:"productAreaRefs"`
}

type proposalStage struct {
	SessionID          string            `json:"sessionId"`
	ProposalID         string            `json:"proposalId"`
	ProposalVersion    uint64            `json:"proposalVersion"`
	ProposalDigest     string            `json:"proposalDigest"`
	PreviousDigest     string            `json:"previousDigest,omitempty"`
	ChangeExplanation  string            `json:"changeExplanation"`
	Outcome            string            `json:"outcome"`
	AffectedUsers      []string          `json:"affectedUsers"`
	Problem            string            `json:"problem"`
	FunctionalBehavior []string          `json:"functionalBehavior"`
	UXExpectations     []string          `json:"uxExpectations"`
	Accessibility      []string          `json:"accessibility"`
	Privacy            []string          `json:"privacy"`
	Security           []string          `json:"security"`
	Operations         []string          `json:"operations"`
	Constraints        []string          `json:"constraints"`
	Exclusions         []string          `json:"exclusions"`
	AcceptanceCriteria []string          `json:"acceptanceCriteria"`
	Inferences         []string          `json:"inferences"`
	Assumptions        []string          `json:"assumptions"`
	Hypotheses         []string          `json:"hypotheses"`
	Questions          []intent.Question `json:"questions"`
	Experiments        []string          `json:"experiments"`
	ProductAreaRefs    []string          `json:"productAreaRefs"`
}

type proposalEditInput struct {
	Field             string
	Item              int
	Replacement       string
	ChangeExplanation string
}

type discoveryInput struct {
	Workspace         string
	Project           string
	Repository        string
	SessionID         string
	UserStatements    []string
	CurrentDraft      discovery.Draft
	PreviousQuestions []discovery.Question
	ProvenanceRefs    []string
}

type discoveryStage struct {
	SessionID        string               `json:"sessionId"`
	Response         string               `json:"response"`
	Draft            discovery.Draft      `json:"draft"`
	Questions        []discovery.Question `json:"questions"`
	Uncertainty      []string             `json:"uncertainty"`
	Limitations      []string             `json:"limitations"`
	Refusal          bool                 `json:"refusal"`
	Incomplete       bool                 `json:"incomplete"`
	Usage            discovery.Usage      `json:"usage"`
	ActivationID     string               `json:"activationId"`
	ActivationDigest string               `json:"activationDigest"`
	Provider         string               `json:"provider"`
	ModelID          string               `json:"modelId"`
	RoleProfileID    string               `json:"roleProfileId"`
	PromptVersion    string               `json:"promptVersion"`
	ResponseSchema   string               `json:"responseSchema"`
}

type planStage struct {
	SessionID  string   `json:"sessionId"`
	RunID      string   `json:"runId"`
	PlanID     string   `json:"planId"`
	PlanDigest string   `json:"planDigest"`
	RiskClass  string   `json:"riskClass"`
	Checks     []string `json:"checks"`
}

type candidateStage struct {
	SessionID          string   `json:"sessionId"`
	RunID              string   `json:"runId"`
	CandidateDigest    string   `json:"candidateDigest"`
	VerificationDigest string   `json:"verificationDigest"`
	ChangedPaths       []string `json:"changedPaths"`
	EvidenceReady      bool     `json:"evidenceReady"`
}

type publicationStage struct {
	SessionID      string `json:"sessionId"`
	RunID          string `json:"runId"`
	OperationID    string `json:"operationId"`
	Branch         string `json:"branch"`
	PullRequestID  string `json:"pullRequestId"`
	PullRequestURL string `json:"pullRequestUrl"`
	Provider       string `json:"provider"`
	MergePermitted bool   `json:"mergePermitted"`
}

type guidedJourney interface {
	Discover(context.Context, discoveryInput) (discoveryStage, error)
	Start(context.Context, startInput) (proposalStage, error)
	ReviseProposal(context.Context, string, string, proposalEditInput) (proposalStage, error)
	PromoteAndPlan(context.Context, string, string) (planStage, error)
	ApprovePlanAndVerify(context.Context, string, string) (candidateStage, error)
	ApprovePublicationAndPublish(context.Context, string, string) (publicationStage, error)
}

type event struct {
	SchemaVersion string `json:"schemaVersion"`
	Type          string `json:"type"`
	State         string `json:"state"`
	Message       string `json:"message"`
	At            string `json:"at"`
	Data          any    `json:"data,omitempty"`
}

func runDelivery(args []string, input *bufio.Reader, output, errorOutput io.Writer, journeys journeyFactory) error {
	if len(args) == 0 || args[0] != "start" {
		return errors.New("usage: forge delivery start --workspace <workspace> --project <project> --repository <repository> (--intent-file <path> | --message <request> --provider-activation <path> --provider-activation-digest <sha256>) --output json")
	}
	fs := flag.NewFlagSet("delivery start", flag.ContinueOnError)
	fs.SetOutput(errorOutput)
	workspace := fs.String("workspace", "", "workspace selector")
	project := fs.String("project", "", "project selector")
	repository := fs.String("repository", "", "repository selector")
	intentFile := fs.String("intent-file", "", "requirements input JSON")
	message := fs.String("message", "", "natural-language product request")
	providerActivation := fs.String("provider-activation", "", "exact approved model-provider activation JSON")
	providerActivationDigest := fs.String("provider-activation-digest", "", "separately approved exact activation digest")
	outputFormat := fs.String("output", "json", "output format")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if *workspace == "" || *project == "" || *repository == "" || *outputFormat != "json" ||
		((*intentFile == "") == (*message == "")) {
		return errors.New("delivery start requires selectors, JSON output, and exactly one of intent file or message")
	}
	if *intentFile != "" && (*providerActivation != "" || *providerActivationDigest != "") {
		return errors.New("provider activation is only valid with a natural-language message")
	}
	if (*providerActivation == "") != (*providerActivationDigest == "") {
		return errors.New("provider activation file and approved digest must be supplied together")
	}
	journey, err := journeys(journeyOptions{
		ProviderActivationPath:   *providerActivation,
		ProviderActivationDigest: *providerActivationDigest,
	})
	if err != nil {
		return err
	}
	ctx := context.Background()
	var intent intentInput
	var rawInput string
	var provenance []string
	conversational := *message != ""
	if conversational {
		intent, rawInput, provenance, err = discoverIntent(ctx, journey, input, output, discoveryInput{
			Workspace: *workspace, Project: *project, Repository: *repository,
			UserStatements: []string{*message}, ProvenanceRefs: []string{"interaction:user:1"},
		})
		if err != nil {
			return err
		}
	} else {
		intent, err = readIntent(*intentFile)
		if err != nil {
			return err
		}
		rawInput = marshalIntent(intent)
		provenance = []string{"cli:intent-file"}
	}
	writeEvent(output, "progress", "proposal-preparing", "Structuring untrusted input as a Requirements Proposal.", nil)
	proposal, err := journey.Start(ctx, startInput{
		Workspace: *workspace, Project: *project, Repository: *repository,
		RawInput: rawInput, Provenance: provenance, Intent: intent,
	})
	if err != nil {
		return err
	}
	writeEvent(output, "requirements-proposal", "awaiting-promotion", "Review the proposal. It has no delivery authority.", proposal)
	if conversational {
		var promote bool
		proposal, promote, err = reviewRequirementsProposal(ctx, journey, input, output, proposal)
		if err != nil {
			return err
		}
		if !promote {
			return nil
		}
	} else {
		if err := requireExact(input, output, "promote "+proposal.ProposalDigest, "Promote the exact proposal"); err != nil {
			return err
		}
	}
	writeEvent(output, "progress", "planning", "Proposal promoted. Resolving bounded context and constructing the Delivery Plan.", nil)
	plan, err := journey.PromoteAndPlan(ctx, proposal.SessionID, proposal.ProposalDigest)
	if err != nil {
		return err
	}
	writeEvent(output, "delivery-plan", "awaiting-plan-approval", "Review risk, checks, budgets, limitations, and rollback.", plan)
	if conversational {
		if err := requireConfirmation(input, output, "Approve execution of this exact plan?"); err != nil {
			return err
		}
	} else {
		if err := requireExact(input, output, "approve-plan "+plan.PlanDigest, "Approve execution of the exact plan"); err != nil {
			return err
		}
	}
	writeEvent(output, "progress", "executing", "Plan approved. Preparing the bounded candidate and running independent verification.", nil)
	candidate, err := journey.ApprovePlanAndVerify(ctx, plan.SessionID, plan.PlanDigest)
	if err != nil {
		return err
	}
	writeEvent(output, "verified-candidate", "awaiting-publication-approval", "Verification passed. Publication still has no authority.", candidate)
	if conversational {
		if err := requireConfirmation(input, output, "Approve contract-double publication of this exact candidate?"); err != nil {
			return err
		}
	} else {
		if err := requireExact(input, output, "approve-publication "+candidate.CandidateDigest, "Approve contract-double publication of the exact candidate"); err != nil {
			return err
		}
	}
	writeEvent(output, "progress", "publishing", "Publication approved. Reconciling the idempotent contract-double operation.", nil)
	publication, err := journey.ApprovePublicationAndPublish(ctx, candidate.SessionID, candidate.CandidateDigest)
	if err != nil {
		return err
	}
	writeEvent(output, "publication-receipt", "reviewable", "Pull request published. Human review and merge remain required.", publication)
	return nil
}

func reviewRequirementsProposal(ctx context.Context, journey guidedJourney, input *bufio.Reader, output io.Writer, proposal proposalStage) (proposalStage, bool, error) {
	const maximumEdits = 16
	for edits := 0; ; {
		fmt.Fprintln(output, "Review Requirements Proposal [promote/edit/reject]")
		action, err := readLine(input, "proposal review action")
		if err != nil {
			return proposalStage{}, false, err
		}
		switch strings.ToLower(action) {
		case "promote":
			return proposal, true, nil
		case "reject":
			writeEvent(output, "requirements-proposal", "rejected", "Proposal rejected. No canonical intent or delivery authority was created.", proposal)
			return proposal, false, nil
		case "edit":
			if edits >= maximumEdits {
				return proposalStage{}, false, errors.New("requirements proposal exceeded the maximum review edits")
			}
			edit, readErr := readProposalEdit(input, output)
			if readErr != nil {
				return proposalStage{}, false, readErr
			}
			revised, reviseErr := journey.ReviseProposal(ctx, proposal.SessionID, proposal.ProposalDigest, edit)
			if reviseErr != nil {
				writeEvent(output, "proposal-edit-rejected", "awaiting-review", "The edit was not applied. The current proposal remains unchanged.", nil)
				continue
			}
			proposal = revised
			edits++
			writeEvent(output, "requirements-proposal", "awaiting-promotion", "Review the revised proposal. Only its exact current digest can be promoted.", proposal)
		default:
			fmt.Fprintln(output, "Choose promote, edit, or reject.")
		}
	}
}

func readProposalEdit(input *bufio.Reader, output io.Writer) (proposalEditInput, error) {
	fmt.Fprintln(output, "Editable fields: outcome, problem, affected-users, functional-behavior, ux-expectations, accessibility, privacy, security, operations, constraints, exclusions, acceptance-criteria, inferences, assumptions, hypotheses, experiments, product-area-refs")
	fmt.Fprint(output, "Field: ")
	field, err := readLine(input, "proposal field")
	if err != nil {
		return proposalEditInput{}, err
	}
	field = strings.ToLower(field)
	item := 0
	if field != "outcome" && field != "problem" {
		fmt.Fprint(output, "Item number (1-based): ")
		rawItem, readErr := readLine(input, "proposal item")
		if readErr != nil {
			return proposalEditInput{}, readErr
		}
		item, err = strconv.Atoi(rawItem)
		if err != nil || item < 1 {
			return proposalEditInput{}, errors.New("proposal item must be a positive number")
		}
	}
	fmt.Fprint(output, "Replacement: ")
	replacement, err := readLine(input, "proposal replacement")
	if err != nil {
		return proposalEditInput{}, err
	}
	fmt.Fprint(output, "Why this changes the proposal: ")
	explanation, err := readLine(input, "proposal change explanation")
	if err != nil {
		return proposalEditInput{}, err
	}
	return proposalEditInput{
		Field: field, Item: item,
		Replacement: replacement, ChangeExplanation: explanation,
	}, nil
}

func readLine(input *bufio.Reader, label string) (string, error) {
	line, err := input.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", fmt.Errorf("read %s: %w", label, err)
	}
	value := strings.TrimSpace(line)
	if value == "" {
		return "", fmt.Errorf("%s is required", label)
	}
	return value, nil
}

func discoverIntent(ctx context.Context, journey guidedJourney, input *bufio.Reader, output io.Writer, request discoveryInput) (intentInput, string, []string, error) {
	const maximumTurns = 8
	statements := append([]string(nil), request.UserStatements...)
	provenance := append([]string(nil), request.ProvenanceRefs...)
	for turn := 0; turn < maximumTurns; turn++ {
		request.UserStatements = append([]string(nil), statements...)
		request.ProvenanceRefs = append([]string(nil), provenance...)
		writeEvent(output, "progress", "discovering", "Forge is clarifying the product outcome. Model output has no authority.", nil)
		stage, err := journey.Discover(ctx, request)
		if err != nil {
			if normalized, retryable := retryableDiscoveryFailure(err); retryable {
				writeEvent(
					output,
					"discovery-paused",
					"awaiting-retry",
					"Forge preserved the current discovery checkpoint after a retryable failure.",
					normalized,
				)
				if retryErr := requireConfirmation(
					input,
					output,
					"Retry the current discovery turn? Existing answers remain in this session.",
				); retryErr == nil {
					continue
				}
				return intentInput{}, "", nil, fmt.Errorf(
					"guided discovery remains paused: %w", err,
				)
			}
			return intentInput{}, "", nil, err
		}
		state := "proposal-ready"
		if hasMaterialQuestions(stage.Questions) {
			state = "awaiting-clarification"
		}
		writeEvent(output, "discovery-turn", state, stage.Response, stage)
		if stage.Refusal {
			return intentInput{}, "", nil, errors.New("model-assisted discovery refused the request")
		}
		if !hasMaterialQuestions(stage.Questions) {
			if stage.Incomplete || !stage.Draft.Complete() {
				return intentInput{}, "", nil, errors.New("model-assisted discovery did not produce a complete proposal")
			}
			return intentFromDraft(stage.Draft, stage.Questions), strings.Join(statements, "\n"), provenance, nil
		}
		answers, answerRefs, err := readDiscoveryAnswers(input, output, stage.Questions, len(statements))
		if err != nil {
			return intentInput{}, "", nil, err
		}
		statements = append(statements, answers...)
		provenance = append(provenance, answerRefs...)
		request.SessionID = stage.SessionID
		request.CurrentDraft = stage.Draft
		request.PreviousQuestions = append([]discovery.Question(nil), stage.Questions...)
	}
	return intentInput{}, "", nil, errors.New("model-assisted discovery exceeded the maximum clarification turns")
}

func retryableDiscoveryFailure(err error) (*failure.Error, bool) {
	var normalized *failure.Error
	if !errors.As(err, &normalized) || normalized.Validate() != nil ||
		normalized.Retry != failure.RetryBounded {
		return nil, false
	}
	return normalized, true
}

func readDiscoveryAnswers(input *bufio.Reader, output io.Writer, questions []discovery.Question, offset int) ([]string, []string, error) {
	answers := make([]string, 0, len(questions))
	references := make([]string, 0, len(questions))
	for _, question := range questions {
		if !question.Material {
			continue
		}
		fmt.Fprintf(output, "Forge asks [%s]: %s\nYour answer: ", question.ID, question.Text)
		line, err := input.ReadString('\n')
		if err != nil && !errors.Is(err, io.EOF) {
			return nil, nil, fmt.Errorf("read discovery answer: %w", err)
		}
		answer := strings.TrimSpace(line)
		if answer == "" {
			return nil, nil, errors.New("material discovery question requires an answer")
		}
		answers = append(answers, answer)
		references = append(references, fmt.Sprintf("interaction:user:%d", offset+len(answers)))
	}
	return answers, references, nil
}

func intentFromDraft(value discovery.Draft, questions []discovery.Question) intentInput {
	return intentInput{
		Outcome: value.Outcome, AffectedUsers: value.AffectedUsers, Problem: value.Problem,
		FunctionalBehavior: value.FunctionalBehavior, UXExpectations: value.UXExpectations,
		Accessibility: value.Accessibility, Privacy: value.Privacy, Security: value.Security,
		Operations: value.Operations, Constraints: value.Constraints, Exclusions: value.Exclusions,
		AcceptanceCriteria: value.AcceptanceCriteria, Inferences: value.Inferences,
		Assumptions: value.Assumptions, Hypotheses: value.Hypotheses,
		Questions: discoveryQuestionsToIntent(questions), Experiments: value.Experiments,
		ProductAreaRefs: value.ProductAreaRefs,
	}
}

func discoveryQuestionsToIntent(values []discovery.Question) []intent.Question {
	result := make([]intent.Question, len(values))
	for index, value := range values {
		result[index] = intent.Question{
			ID: value.ID, Text: value.Text, Consequence: value.Consequence,
			Material: value.Material, Resolved: false,
		}
	}
	return result
}

func hasMaterialQuestions(values []discovery.Question) bool {
	for _, value := range values {
		if value.Material {
			return true
		}
	}
	return false
}

func marshalIntent(value intentInput) string {
	encoded, _ := json.Marshal(value)
	return string(encoded)
}

func readIntent(path string) (intentInput, error) {
	file, err := os.Open(path)
	if err != nil {
		return intentInput{}, fmt.Errorf("read intent: %w", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(io.LimitReader(file, 1<<20))
	decoder.DisallowUnknownFields()
	var value intentInput
	if err := decoder.Decode(&value); err != nil {
		return intentInput{}, fmt.Errorf("decode intent: %w", err)
	}
	if value.Outcome == "" || value.Problem == "" || len(value.AffectedUsers) == 0 ||
		len(value.AcceptanceCriteria) == 0 || value.Constraints == nil {
		return intentInput{}, errors.New("intent outcome, problem, affected users, constraints, and acceptance criteria are required")
	}
	var trailing any
	if decoder.Decode(&trailing) != io.EOF {
		return intentInput{}, errors.New("intent file must contain exactly one JSON object")
	}
	return value, nil
}

func requireExact(input *bufio.Reader, output io.Writer, expected, label string) error {
	fmt.Fprintf(output, "%s by entering: %s\n", label, expected)
	line, err := input.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("read human decision: %w", err)
	}
	if strings.TrimSpace(line) != expected {
		return errors.New("explicit human decision did not match the current immutable digest")
	}
	return nil
}

func requireConfirmation(input *bufio.Reader, output io.Writer, label string) error {
	fmt.Fprintf(output, "%s [yes/no]\n", label)
	line, err := input.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return fmt.Errorf("read human decision: %w", err)
	}
	if !strings.EqualFold(strings.TrimSpace(line), "yes") {
		return errors.New("explicit human decision was not approved")
	}
	return nil
}

func writeEvent(output io.Writer, eventType, state, message string, data any) {
	_ = json.NewEncoder(output).Encode(event{
		SchemaVersion: "forge.cli-event.v1", Type: eventType, State: state,
		Message: message, At: time.Now().UTC().Format(time.RFC3339Nano), Data: data,
	})
}
