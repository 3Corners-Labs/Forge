package cli

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"
	"sync"
	"syscall"
	"time"

	workbenchassets "github.com/3corners/forge/internal/experience/workbench"
)

const workbenchAPIVersion = "forge.workbench-api.v1"

type workbenchAgent struct {
	ID        string `json:"id"`
	Role      string `json:"role"`
	Label     string `json:"label"`
	State     string `json:"state"`
	Summary   string `json:"summary"`
	Authority string `json:"authority"`
	WaitingOn string `json:"waitingOn,omitempty"`
}

type workbenchGate struct {
	Action  string `json:"action"`
	Label   string `json:"label"`
	Digest  string `json:"digest"`
	Summary string `json:"summary"`
}

type workbenchActivity struct {
	Sequence int64  `json:"sequence"`
	At       string `json:"at"`
	Actor    string `json:"actor"`
	Type     string `json:"type"`
	State    string `json:"state"`
	Message  string `json:"message"`
}

type workbenchRun struct {
	SchemaVersion string              `json:"schemaVersion"`
	ID            string              `json:"id"`
	Title         string              `json:"title"`
	Request       string              `json:"request"`
	State         string              `json:"state"`
	Progress      int                 `json:"progress"`
	StartedAt     string              `json:"startedAt"`
	UpdatedAt     string              `json:"updatedAt"`
	Agents        []workbenchAgent    `json:"agents"`
	Gates         []workbenchGate     `json:"gates"`
	Activity      []workbenchActivity `json:"activity"`
	Proposal      *proposalStage      `json:"proposal,omitempty"`
	Plan          *planStage          `json:"plan,omitempty"`
	Candidate     *candidateStage     `json:"candidate,omitempty"`
	Publication   *publicationStage   `json:"publication,omitempty"`

	journey     guidedJourney
	subscribers map[chan workbenchActivity]struct{}
}

type workbenchAPI struct {
	mu       sync.RWMutex
	journeys journeyFactory
	runs     map[string]*workbenchRun
	nextID   int64
	now      func() time.Time
	assets   fs.FS
}

func newWorkbenchAPI(journeys journeyFactory) *workbenchAPI {
	return &workbenchAPI{
		journeys: journeys, runs: map[string]*workbenchRun{},
		now: func() time.Time { return time.Now().UTC() }, assets: workbenchassets.Assets(),
	}
}

func runWorkbench(args []string, input io.Reader, output, errorOutput io.Writer, journeys journeyFactory) error {
	if len(args) > 0 && args[0] == "chat" {
		return runWorkbenchChat(args[1:], input, output, errorOutput, journeys)
	}
	if len(args) > 0 && args[0] == "serve" {
		args = args[1:]
	}
	flags := flag.NewFlagSet("workbench", flag.ContinueOnError)
	flags.SetOutput(errorOutput)
	listen := flags.String("listen", "127.0.0.1:7331", "loopback Workbench API address")
	if err := flags.Parse(args); err != nil {
		return err
	}
	host, _, err := net.SplitHostPort(*listen)
	if err != nil || (host != "127.0.0.1" && host != "localhost" && host != "::1") {
		return errors.New("Workbench Trial must listen on loopback")
	}
	api := newWorkbenchAPI(journeys)
	server := &http.Server{Addr: *listen, Handler: api, ReadHeaderTimeout: 5 * time.Second}
	stopped := make(chan os.Signal, 1)
	signal.Notify(stopped, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(stopped)
	go func() {
		<-stopped
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = server.Shutdown(ctx)
	}()
	fmt.Fprintf(output, "Forge Workbench ready at http://%s\n", *listen)
	err = server.ListenAndServe()
	if errors.Is(err, http.ErrServerClosed) {
		return nil
	}
	return err
}

func (api *workbenchAPI) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("X-Content-Type-Options", "nosniff")
	response.Header().Set("Referrer-Policy", "no-referrer")
	response.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' 'unsafe-inline'; connect-src 'self'; img-src 'self' data:; script-src 'self'")
	if strings.HasPrefix(request.URL.Path, "/api/v1/") {
		api.serveAPI(response, request)
		return
	}
	if request.Method != http.MethodGet && request.Method != http.MethodHead {
		http.Error(response, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	path := strings.TrimPrefix(request.URL.Path, "/")
	if path == "" {
		path = "index.html"
	}
	if _, err := fs.Stat(api.assets, path); err != nil {
		path = "index.html"
	}
	http.ServeFileFS(response, request, api.assets, path)
}

func (api *workbenchAPI) serveAPI(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Cache-Control", "no-store")
	segments := strings.Split(strings.Trim(request.URL.Path, "/"), "/")
	switch {
	case request.Method == http.MethodGet && request.URL.Path == "/api/v1/overview":
		api.overview(response)
	case request.Method == http.MethodPost && request.URL.Path == "/api/v1/runs":
		api.createRun(response, request)
	case len(segments) == 4 && segments[2] == "runs" && request.Method == http.MethodGet:
		api.getRun(response, segments[3])
	case len(segments) == 5 && segments[2] == "runs" && segments[4] == "actions" && request.Method == http.MethodPost:
		api.act(response, request, segments[3])
	case len(segments) == 5 && segments[2] == "runs" && segments[4] == "events" && request.Method == http.MethodGet:
		api.events(response, request, segments[3])
	default:
		writeAPIError(response, http.StatusNotFound, "not_found", "Workbench resource was not found")
	}
}

func (api *workbenchAPI) overview(response http.ResponseWriter) {
	api.mu.RLock()
	runs := make([]*workbenchRun, 0, len(api.runs))
	for _, run := range api.runs {
		runs = append(runs, cloneWorkbenchRun(run))
	}
	api.mu.RUnlock()
	slices.SortFunc(runs, func(left, right *workbenchRun) int {
		return strings.Compare(right.UpdatedAt, left.UpdatedAt)
	})
	writeJSON(response, http.StatusOK, map[string]any{
		"schemaVersion": workbenchAPIVersion, "product": "Forge Workbench",
		"mode": "exploratory-trial", "runs": runs,
	})
}

func (api *workbenchAPI) createRun(response http.ResponseWriter, request *http.Request) {
	var input struct {
		Message string `json:"message"`
	}
	if err := decodeJSON(request, &input); err != nil || strings.TrimSpace(input.Message) == "" || len(input.Message) > 4000 {
		writeAPIError(response, http.StatusBadRequest, "invalid_request", "Describe one bounded product outcome")
		return
	}
	journey, err := api.journeys(journeyOptions{})
	if err != nil {
		writeAPIError(response, http.StatusServiceUnavailable, "journey_unavailable", "Forge could not start the Trial journey")
		return
	}
	message := strings.TrimSpace(input.Message)
	proposal, err := journey.Start(request.Context(), startInput{
		Workspace: "fixture-workspace", Project: "forge", Repository: "forge",
		RawInput: message, Provenance: []string{"workbench:message"},
		Intent: intentInput{
			Outcome: message, AffectedUsers: []string{"Forge product team"},
			Problem:            "The requested outcome is not yet implemented.",
			FunctionalBehavior: []string{message},
			UXExpectations:     []string{"Immediate feedback and an inspectable result"},
			Accessibility:      []string{"Keyboard operable and screen-reader understandable"},
			Privacy:            []string{"Use only data necessary for the declared purpose"},
			Security:           []string{"Preserve Forge authority and verification boundaries"},
			Operations:         []string{"Expose progress, failure and recovery state"},
			Constraints:        []string{"Exploratory Developer Trial"},
			Exclusions:         []string{"Production release", "Automatic merge"},
			AcceptanceCriteria: []string{"Independent verification passes", "Human can inspect the reviewable result"},
			Assumptions:        []string{"The M1 Developer Trial services are available"},
			Experiments:        []string{"Run the Workbench journey and collect user feedback"},
		},
	})
	if err != nil {
		if closer, ok := journey.(io.Closer); ok {
			_ = closer.Close()
		}
		writeAPIError(response, http.StatusUnprocessableEntity, "proposal_failed", "Forge could not prepare a Requirements Proposal")
		return
	}
	api.mu.Lock()
	api.nextID++
	id := fmt.Sprintf("workbench-run-%d", api.nextID)
	now := api.now().Format(time.RFC3339Nano)
	run := &workbenchRun{
		SchemaVersion: workbenchAPIVersion, ID: id, Title: truncate(message, 72), Request: message,
		State: "awaiting-requirements-promotion", Progress: 18, StartedAt: now, UpdatedAt: now,
		Proposal: &proposal, journey: journey, subscribers: map[chan workbenchActivity]struct{}{},
	}
	run.Agents = workbenchAgents("proposal")
	run.Gates = []workbenchGate{{
		Action: "promote", Label: "Promote requirements", Digest: proposal.ProposalDigest,
		Summary: "Accept exactly what Forge understood; this authorizes planning, not implementation.",
	}}
	api.appendLocked(run, "virtual-product-manager", "proposal", run.State, "Requirements Proposal ready for review")
	api.runs[id] = run
	result := cloneWorkbenchRun(run)
	api.mu.Unlock()
	writeJSON(response, http.StatusCreated, result)
}

func (api *workbenchAPI) getRun(response http.ResponseWriter, id string) {
	api.mu.RLock()
	run := cloneWorkbenchRun(api.runs[id])
	api.mu.RUnlock()
	if run == nil {
		writeAPIError(response, http.StatusNotFound, "run_not_found", "Delivery run was not found")
		return
	}
	writeJSON(response, http.StatusOK, run)
}

func (api *workbenchAPI) act(response http.ResponseWriter, request *http.Request, id string) {
	var input struct {
		Action string `json:"action"`
	}
	if err := decodeJSON(request, &input); err != nil {
		writeAPIError(response, http.StatusBadRequest, "invalid_action", "Action is invalid")
		return
	}
	api.mu.Lock()
	run := api.runs[id]
	if run == nil {
		api.mu.Unlock()
		writeAPIError(response, http.StatusNotFound, "run_not_found", "Delivery run was not found")
		return
	}
	ctx := request.Context()
	var err error
	switch input.Action {
	case "promote":
		if run.Proposal == nil || run.Plan != nil {
			err = errors.New("promotion is not available")
			break
		}
		var plan planStage
		plan, err = run.journey.PromoteAndPlan(ctx, run.Proposal.SessionID, run.Proposal.ProposalDigest)
		if err == nil {
			run.Plan, run.State, run.Progress = &plan, "awaiting-plan-approval", 38
			run.Agents = workbenchAgents("plan")
			run.Gates = []workbenchGate{{Action: "approve-plan", Label: "Approve bounded plan", Digest: plan.PlanDigest, Summary: "Authorize the exact bounded worker plan."}}
			api.appendLocked(run, "delivery-coordinator", "plan", run.State, "Bounded Delivery Plan ready")
		}
	case "approve-plan":
		if run.Plan == nil || run.Candidate != nil {
			err = errors.New("plan approval is not available")
			break
		}
		var candidate candidateStage
		candidate, err = run.journey.ApprovePlanAndVerify(ctx, run.Plan.SessionID, run.Plan.PlanDigest)
		if err == nil {
			run.Candidate, run.State, run.Progress = &candidate, "awaiting-publication-approval", 82
			run.Agents = workbenchAgents("candidate")
			run.Gates = []workbenchGate{{Action: "approve-publication", Label: "Publish reviewable change", Digest: candidate.CandidateDigest, Summary: "Publish the independently verified candidate to the contract-double PR boundary."}}
			api.appendLocked(run, "quality-verifier", "verification", run.State, "Candidate prepared and independently verified")
		}
	case "approve-publication":
		if run.Candidate == nil || run.Publication != nil {
			err = errors.New("publication approval is not available")
			break
		}
		var publication publicationStage
		publication, err = run.journey.ApprovePublicationAndPublish(ctx, run.Candidate.SessionID, run.Candidate.CandidateDigest)
		if err == nil {
			run.Publication, run.State, run.Progress = &publication, "reviewable", 100
			run.Agents, run.Gates = workbenchAgents("published"), []workbenchGate{}
			api.appendLocked(run, "source-control-publisher", "publication", run.State, "Reviewable contract-double pull request published")
		}
	default:
		err = errors.New("unsupported action")
	}
	if err != nil {
		api.mu.Unlock()
		writeAPIError(response, http.StatusConflict, "action_rejected", err.Error())
		return
	}
	result := cloneWorkbenchRun(run)
	api.mu.Unlock()
	writeJSON(response, http.StatusOK, result)
}

func (api *workbenchAPI) events(response http.ResponseWriter, request *http.Request, id string) {
	flusher, ok := response.(http.Flusher)
	if !ok {
		writeAPIError(response, http.StatusNotImplemented, "stream_unavailable", "Streaming is unavailable")
		return
	}
	channel := make(chan workbenchActivity, 16)
	api.mu.Lock()
	run := api.runs[id]
	if run == nil {
		api.mu.Unlock()
		writeAPIError(response, http.StatusNotFound, "run_not_found", "Delivery run was not found")
		return
	}
	run.subscribers[channel] = struct{}{}
	history := slices.Clone(run.Activity)
	api.mu.Unlock()
	defer func() {
		api.mu.Lock()
		delete(run.subscribers, channel)
		api.mu.Unlock()
	}()
	response.Header().Set("Content-Type", "text/event-stream")
	response.Header().Set("Connection", "keep-alive")
	for _, item := range history {
		writeSSE(response, item)
	}
	flusher.Flush()
	for {
		select {
		case item := <-channel:
			writeSSE(response, item)
			flusher.Flush()
		case <-request.Context().Done():
			return
		}
	}
}

func (api *workbenchAPI) appendLocked(run *workbenchRun, actor, kind, state, message string) {
	item := workbenchActivity{
		Sequence: int64(len(run.Activity) + 1), At: api.now().Format(time.RFC3339Nano),
		Actor: actor, Type: kind, State: state, Message: message,
	}
	run.Activity = append(run.Activity, item)
	run.UpdatedAt = item.At
	for subscriber := range run.subscribers {
		select {
		case subscriber <- item:
		default:
		}
	}
}

func workbenchAgents(stage string) []workbenchAgent {
	agents := []workbenchAgent{
		{ID: "virtual-product-manager", Role: "product", Label: "Virtual Product Manager", State: "queued", Summary: "Clarifies intent and requirements", Authority: "Advise"},
		{ID: "delivery-coordinator", Role: "delivery", Label: "Delivery Coordinator", State: "queued", Summary: "Builds the bounded delivery graph", Authority: "Prepare"},
		{ID: "engineering-worker", Role: "engineering", Label: "Engineering Worker", State: "queued", Summary: "Prepares candidate bytes in isolation", Authority: "Prepare"},
		{ID: "quality-verifier", Role: "quality", Label: "Independent Quality", State: "queued", Summary: "Runs consequential verification", Authority: "Verify"},
		{ID: "source-control-publisher", Role: "delivery", Label: "Source Control Publisher", State: "queued", Summary: "Creates the reviewable hand-off", Authority: "Publish contract double"},
	}
	set := func(id, state, waiting string) {
		for index := range agents {
			if agents[index].ID == id {
				agents[index].State, agents[index].WaitingOn = state, waiting
			}
		}
	}
	switch stage {
	case "proposal":
		set("virtual-product-manager", "completed", "")
		set("delivery-coordinator", "waiting", "Human requirements promotion")
	case "plan":
		set("virtual-product-manager", "completed", "")
		set("delivery-coordinator", "completed", "")
		set("engineering-worker", "waiting", "Human plan approval")
	case "candidate":
		for _, id := range []string{"virtual-product-manager", "delivery-coordinator", "engineering-worker", "quality-verifier"} {
			set(id, "completed", "")
		}
		set("source-control-publisher", "waiting", "Human publication approval")
	case "published":
		for index := range agents {
			agents[index].State = "completed"
		}
	}
	return agents
}

func cloneWorkbenchRun(source *workbenchRun) *workbenchRun {
	if source == nil {
		return nil
	}
	copy := *source
	copy.Agents = slices.Clone(source.Agents)
	copy.Gates = slices.Clone(source.Gates)
	copy.Activity = slices.Clone(source.Activity)
	copy.journey, copy.subscribers = nil, nil
	return &copy
}

func decodeJSON(request *http.Request, target any) error {
	defer request.Body.Close()
	decoder := json.NewDecoder(io.LimitReader(request.Body, 16<<10))
	decoder.DisallowUnknownFields()
	return decoder.Decode(target)
}

func writeJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(status)
	_ = json.NewEncoder(response).Encode(value)
}

func writeAPIError(response http.ResponseWriter, status int, code, message string) {
	writeJSON(response, status, map[string]any{"schemaVersion": workbenchAPIVersion, "error": map[string]string{"code": code, "message": message}})
}

func writeSSE(output io.Writer, item workbenchActivity) {
	encoded, _ := json.Marshal(item)
	_, _ = fmt.Fprintf(output, "id: %d\nevent: activity\ndata: %s\n\n", item.Sequence, encoded)
}

func truncate(value string, maximum int) string {
	value = strings.TrimSpace(value)
	if len(value) <= maximum {
		return value
	}
	return value[:maximum-1] + "…"
}
