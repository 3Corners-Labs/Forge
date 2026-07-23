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
	"strings"
	"time"
)

type journeyFactory func() (guidedJourney, error)

type startInput struct {
	Workspace  string
	Project    string
	Repository string
	Intent     intentInput
}

type intentInput struct {
	Outcome            string   `json:"outcome"`
	AffectedUsers      []string `json:"affectedUsers"`
	Problem            string   `json:"problem"`
	FunctionalBehavior []string `json:"functionalBehavior"`
	UXExpectations     []string `json:"uxExpectations"`
	Accessibility      []string `json:"accessibility"`
	Privacy            []string `json:"privacy"`
	Security           []string `json:"security"`
	Operations         []string `json:"operations"`
	Constraints        []string `json:"constraints"`
	Exclusions         []string `json:"exclusions"`
	AcceptanceCriteria []string `json:"acceptanceCriteria"`
	Assumptions        []string `json:"assumptions"`
	Experiments        []string `json:"experiments"`
}

type proposalStage struct {
	SessionID      string   `json:"sessionId"`
	ProposalID     string   `json:"proposalId"`
	ProposalDigest string   `json:"proposalDigest"`
	Outcome        string   `json:"outcome"`
	Assumptions    []string `json:"assumptions"`
	Experiments    []string `json:"experiments"`
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
	Start(context.Context, startInput) (proposalStage, error)
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
		return errors.New("usage: forge delivery start --workspace <workspace> --project <project> --repository <repository> --intent-file <path> --output json")
	}
	fs := flag.NewFlagSet("delivery start", flag.ContinueOnError)
	fs.SetOutput(errorOutput)
	workspace := fs.String("workspace", "", "workspace selector")
	project := fs.String("project", "", "project selector")
	repository := fs.String("repository", "", "repository selector")
	intentFile := fs.String("intent-file", "", "requirements input JSON")
	outputFormat := fs.String("output", "json", "output format")
	if err := fs.Parse(args[1:]); err != nil {
		return err
	}
	if *workspace == "" || *project == "" || *repository == "" || *intentFile == "" || *outputFormat != "json" {
		return errors.New("delivery start selectors, intent file, and JSON output are required")
	}
	intent, err := readIntent(*intentFile)
	if err != nil {
		return err
	}
	journey, err := journeys()
	if err != nil {
		return err
	}
	ctx := context.Background()
	writeEvent(output, "progress", "proposal-preparing", "Structuring untrusted input as a Requirements Proposal.", nil)
	proposal, err := journey.Start(ctx, startInput{Workspace: *workspace, Project: *project, Repository: *repository, Intent: intent})
	if err != nil {
		return err
	}
	writeEvent(output, "requirements-proposal", "awaiting-promotion", "Review the proposal. It has no delivery authority.", proposal)
	if err := requireExact(input, output, "promote "+proposal.ProposalDigest, "Promote the exact proposal"); err != nil {
		return err
	}
	writeEvent(output, "progress", "planning", "Proposal promoted. Resolving bounded context and constructing the Delivery Plan.", nil)
	plan, err := journey.PromoteAndPlan(ctx, proposal.SessionID, proposal.ProposalDigest)
	if err != nil {
		return err
	}
	writeEvent(output, "delivery-plan", "awaiting-plan-approval", "Review risk, checks, budgets, limitations, and rollback.", plan)
	if err := requireExact(input, output, "approve-plan "+plan.PlanDigest, "Approve execution of the exact plan"); err != nil {
		return err
	}
	writeEvent(output, "progress", "executing", "Plan approved. Preparing the bounded candidate and running independent verification.", nil)
	candidate, err := journey.ApprovePlanAndVerify(ctx, plan.SessionID, plan.PlanDigest)
	if err != nil {
		return err
	}
	writeEvent(output, "verified-candidate", "awaiting-publication-approval", "Verification passed. Publication still has no authority.", candidate)
	if err := requireExact(input, output, "approve-publication "+candidate.CandidateDigest, "Approve contract-double publication of the exact candidate"); err != nil {
		return err
	}
	writeEvent(output, "progress", "publishing", "Publication approved. Reconciling the idempotent contract-double operation.", nil)
	publication, err := journey.ApprovePublicationAndPublish(ctx, candidate.SessionID, candidate.CandidateDigest)
	if err != nil {
		return err
	}
	writeEvent(output, "publication-receipt", "reviewable", "Pull request published. Human review and merge remain required.", publication)
	return nil
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

func writeEvent(output io.Writer, eventType, state, message string, data any) {
	_ = json.NewEncoder(output).Encode(event{
		SchemaVersion: "forge.cli-event.v1", Type: eventType, State: state,
		Message: message, At: time.Now().UTC().Format(time.RFC3339Nano), Data: data,
	})
}
