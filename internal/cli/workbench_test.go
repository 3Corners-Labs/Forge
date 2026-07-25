package cli

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestWorkbenchAPIProjectsCompleteM1JourneyWithoutMergeAuthority(t *testing.T) {
	journey := &journeyStub{}
	api := newWorkbenchAPI(func(journeyOptions) (guidedJourney, error) { return journey, nil })

	run := workbenchRoundTrip(t, api, http.MethodPost, "/api/v1/runs",
		`{"message":"Build an inspectable agent workbench"}`, http.StatusCreated)
	if run.State != "awaiting-requirements-promotion" || len(run.Gates) != 1 ||
		run.Gates[0].Digest != "sha256:proposal" {
		t.Fatalf("proposal projection = %#v", run)
	}
	for _, action := range []string{"promote", "approve-plan", "approve-publication"} {
		run = workbenchRoundTrip(t, api, http.MethodPost, "/api/v1/runs/"+run.ID+"/actions",
			`{"action":"`+action+`"}`, http.StatusOK)
	}
	if run.State != "reviewable" || run.Progress != 100 || run.Publication == nil ||
		run.Publication.MergePermitted || len(run.Gates) != 0 {
		t.Fatalf("final projection = %#v", run)
	}
	if strings.Join(journey.calls, ",") != "start,promote,execute,publish" {
		t.Fatalf("capability sequence = %v", journey.calls)
	}
	for _, agent := range run.Agents {
		if agent.State != "completed" {
			t.Fatalf("agent %s state = %s", agent.ID, agent.State)
		}
	}
}

func TestWorkbenchAPIFailsClosedForInvalidAndOutOfOrderActions(t *testing.T) {
	api := newWorkbenchAPI(func(journeyOptions) (guidedJourney, error) { return &journeyStub{}, nil })
	workbenchStatus(t, api, http.MethodPost, "/api/v1/runs", `{"message":""}`, http.StatusBadRequest)
	run := workbenchRoundTrip(t, api, http.MethodPost, "/api/v1/runs",
		`{"message":"Prepare a bounded change"}`, http.StatusCreated)
	workbenchStatus(t, api, http.MethodPost, "/api/v1/runs/"+run.ID+"/actions",
		`{"action":"approve-plan"}`, http.StatusConflict)
	workbenchStatus(t, api, http.MethodPost, "/api/v1/runs/"+run.ID+"/actions",
		`{"action":"merge"}`, http.StatusConflict)
	workbenchStatus(t, api, http.MethodGet, "/api/v1/runs/missing", "", http.StatusNotFound)
}

func TestWorkbenchServesEmbeddedCockpitWithBrowserSecurityHeaders(t *testing.T) {
	api := newWorkbenchAPI(func(journeyOptions) (guidedJourney, error) { return &journeyStub{}, nil })
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response := httptest.NewRecorder()
	api.ServeHTTP(response, request)
	if response.Code != http.StatusOK || !strings.Contains(response.Body.String(), "Forge Workbench") {
		t.Fatalf("cockpit response status=%d body=%s", response.Code, response.Body.String())
	}
	if response.Header().Get("X-Content-Type-Options") != "nosniff" ||
		!strings.Contains(response.Header().Get("Content-Security-Policy"), "default-src 'self'") {
		t.Fatalf("security headers = %#v", response.Header())
	}
}

func TestWorkbenchChatUsesSameAPIAndRequiresThreeExplicitApprovals(t *testing.T) {
	journey := &journeyStub{}
	api := newWorkbenchAPI(func(journeyOptions) (guidedJourney, error) { return journey, nil })
	requester := func(method, url string, body any) (*workbenchRun, error) {
		encoded, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		path := strings.TrimPrefix(url, "http://workbench.test")
		return workbenchRoundTrip(t, api, method, path, string(encoded), map[string]int{
			"/api/v1/runs": http.StatusCreated,
		}[path]), nil
	}
	actionStatus := func(method, url string, body any) (*workbenchRun, error) {
		if strings.HasSuffix(url, "/actions") {
			encoded, _ := json.Marshal(body)
			return workbenchRoundTrip(t, api, method, strings.TrimPrefix(url, "http://workbench.test"), string(encoded), http.StatusOK), nil
		}
		return requester(method, url, body)
	}
	var output bytes.Buffer
	err := runWorkbenchChatWithRequester([]string{"--api", "http://workbench.test"}, strings.NewReader(
		"Build the Forge Workbench\napprove\napprove\napprove\n",
	), &output, &bytes.Buffer{}, func(journeyOptions) (guidedJourney, error) { return journey, nil }, actionStatus)
	if err != nil {
		t.Fatal(err)
	}
	for _, expected := range []string{
		"Free text becomes a proposal", "Promote requirements", "Approve bounded plan",
		"Publish reviewable change", "Reviewable result ready",
	} {
		if !strings.Contains(strings.ToUpper(output.String()), strings.ToUpper(expected)) {
			t.Fatalf("chat output missing %q:\n%s", expected, output.String())
		}
	}
	if strings.Join(journey.calls, ",") != "start,promote,execute,publish" {
		t.Fatalf("capability sequence = %v", journey.calls)
	}
}

func workbenchRoundTrip(t *testing.T, handler http.Handler, method, path, body string, status int) *workbenchRun {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != status {
		t.Fatalf("%s %s status=%d want=%d body=%s", method, path, response.Code, status, response.Body.String())
	}
	var run workbenchRun
	if err := json.NewDecoder(response.Body).Decode(&run); err != nil {
		t.Fatal(err)
	}
	return &run
}

func workbenchStatus(t *testing.T, handler http.Handler, method, path, body string, status int) {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	if response.Code != status {
		t.Fatalf("%s %s status=%d want=%d body=%s", method, path, response.Code, status, response.Body.String())
	}
}
