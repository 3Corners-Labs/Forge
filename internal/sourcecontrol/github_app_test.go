package sourcecontrol

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"sync"
	"testing"
)

type githubTokenStub struct {
	reference string
	calls     int
}

func (source *githubTokenStub) Token(_ context.Context, reference string) (string, error) {
	source.reference = reference
	source.calls++
	return "short-lived-installation-token", nil
}

type githubAPIDouble struct {
	mu            sync.Mutex
	baseRevision  string
	branchCreated bool
	pullCreated   bool
	commitMessage string
	paths         []string
	authorization []string
}

func (double *githubAPIDouble) RoundTrip(request *http.Request) (*http.Response, error) {
	double.mu.Lock()
	defer double.mu.Unlock()
	double.paths = append(double.paths, request.Method+" "+request.URL.Path)
	double.authorization = append(double.authorization, request.Header.Get("Authorization"))
	status := http.StatusOK
	var output any
	switch {
	case request.Method == http.MethodGet && strings.Contains(request.URL.Path, "/git/ref/heads/main"):
		output = map[string]any{"object": map[string]any{"sha": double.baseRevision}}
	case request.Method == http.MethodGet && strings.Contains(request.URL.Path, "/git/ref/heads/forge/"):
		if !double.branchCreated {
			status = http.StatusNotFound
			output = map[string]any{"message": "Not Found"}
		} else {
			output = map[string]any{"object": map[string]any{"sha": "candidate-commit"}}
		}
	case request.Method == http.MethodGet && strings.HasSuffix(request.URL.Path, "/git/commits/"+double.baseRevision):
		output = map[string]any{"tree": map[string]any{"sha": "base-tree"}, "parents": []any{}}
	case request.Method == http.MethodGet && strings.HasSuffix(request.URL.Path, "/git/commits/candidate-commit"):
		output = map[string]any{
			"message": double.commitMessage, "tree": map[string]any{"sha": "candidate-tree"},
			"parents": []any{map[string]any{"sha": double.baseRevision}},
		}
	case request.Method == http.MethodPost && strings.HasSuffix(request.URL.Path, "/git/blobs"):
		status, output = http.StatusCreated, map[string]any{"sha": "blob-1"}
	case request.Method == http.MethodPost && strings.HasSuffix(request.URL.Path, "/git/trees"):
		status, output = http.StatusCreated, map[string]any{"sha": "candidate-tree"}
	case request.Method == http.MethodPost && strings.HasSuffix(request.URL.Path, "/git/commits"):
		var body struct {
			Message string `json:"message"`
		}
		_ = json.NewDecoder(request.Body).Decode(&body)
		double.commitMessage = body.Message
		status, output = http.StatusCreated, map[string]any{"sha": "candidate-commit"}
	case request.Method == http.MethodPost && strings.HasSuffix(request.URL.Path, "/git/refs"):
		double.branchCreated = true
		status, output = http.StatusCreated, map[string]any{"ref": "created"}
	case request.Method == http.MethodGet && strings.HasSuffix(request.URL.Path, "/pulls"):
		if double.pullCreated {
			output = []any{map[string]any{
				"number": 42, "html_url": "https://github.com/3Corners-Labs/Forge/pull/42",
				"state": "open", "body": "review\n\n" + testGitHubOperation().Marker,
			}}
		} else {
			output = []any{}
		}
	case request.Method == http.MethodPost && strings.HasSuffix(request.URL.Path, "/pulls"):
		double.pullCreated = true
		status, output = http.StatusCreated, map[string]any{
			"number": 42, "html_url": "https://github.com/3Corners-Labs/Forge/pull/42",
			"state": "open",
		}
	default:
		status, output = http.StatusNotFound, map[string]any{"message": "unexpected"}
	}
	encoded, _ := json.Marshal(output)
	return &http.Response{
		StatusCode: status, Header: make(http.Header),
		Body: io.NopCloser(bytes.NewReader(encoded)), Request: request,
	}, nil
}

func TestGitHubAppProviderCreatesAndReconcilesExactDraftPullRequest(t *testing.T) {
	operation := testGitHubOperation()
	api := &githubAPIDouble{baseRevision: operation.BaseRevision}
	tokens := &githubTokenStub{}
	provider, err := NewGitHubAppProvider(
		GitHubAppConfig{
			Owner: "3Corners-Labs", Repository: "Forge", RepositoryID: operation.RepositoryID,
			CredentialRef: operation.CredentialRef,
		},
		&http.Client{Transport: api}, tokens,
	)
	if err != nil {
		t.Fatal(err)
	}
	if state, err := provider.Reconcile(context.Background(), operation); err != nil || state != (RemoteState{}) {
		t.Fatalf("initial reconcile=%#v err=%v", state, err)
	}
	if err := provider.CreateBranch(context.Background(), operation); err != nil {
		t.Fatal(err)
	}
	state, err := provider.Reconcile(context.Background(), operation)
	if err != nil || !state.BranchExists || state.CommitID != "candidate-commit" ||
		state.CandidateDigest != operation.CandidateDigest || state.TreeDigest != operation.TreeDigest ||
		state.BaseRevision != operation.BaseRevision {
		t.Fatalf("branch state=%#v err=%v", state, err)
	}
	if err := provider.OpenPullRequest(context.Background(), operation, "Verified candidate", "Review required."); err != nil {
		t.Fatal(err)
	}
	state, err = provider.Reconcile(context.Background(), operation)
	if err != nil || state.PullRequestID != "42" || !state.PullRequestOpen {
		t.Fatalf("PR state=%#v err=%v", state, err)
	}
	if tokens.reference != operation.CredentialRef || tokens.calls == 0 {
		t.Fatalf("token source was not invoked for the exact reference: %#v", tokens)
	}
	for index, path := range api.paths {
		if strings.Contains(path, "/merge") || strings.Contains(path, "/releases") ||
			api.authorization[index] != "Bearer short-lived-installation-token" {
			t.Fatalf("forbidden or unauthenticated API request: %s", path)
		}
	}
}

func TestGitHubAppProviderDeniesForbiddenEndpointsAndForeignOperations(t *testing.T) {
	for _, path := range []string{
		"/repos/3Corners-Labs/Forge/pulls/42/merge",
		"/repos/3Corners-Labs/Forge/releases",
		"/repos/3Corners-Labs/Forge/actions/workflows",
		"/orgs/3Corners-Labs/members",
	} {
		if allowedGitHubPath(http.MethodPost, path) {
			t.Fatalf("forbidden endpoint allowed: %s", path)
		}
	}
	operation := testGitHubOperation()
	provider, err := NewGitHubAppProvider(
		GitHubAppConfig{
			Owner: "3Corners-Labs", Repository: "Forge", RepositoryID: operation.RepositoryID,
			CredentialRef: operation.CredentialRef,
		},
		&http.Client{Transport: &githubAPIDouble{baseRevision: operation.BaseRevision}},
		&githubTokenStub{},
	)
	if err != nil {
		t.Fatal(err)
	}
	foreign := operation
	foreign.RepositoryID = "other-repository"
	if err := provider.CreateBranch(context.Background(), foreign); err != ErrProviderConflict {
		t.Fatalf("foreign operation err=%v", err)
	}
}

func testGitHubOperation() Operation {
	candidate := "sha256:" + strings.Repeat("a", 64)
	tree := "sha256:" + strings.Repeat("b", 64)
	return Operation{
		ContractVersion: ProviderContractVersion, ID: "publication-1", TenantID: "tenant-1",
		RunID: "run-1", RepositoryID: "repository-forge", Provider: "github-app",
		BaseBranch: "main", BaseRevision: strings.Repeat("c", 40),
		Branch: "forge/run-1-aaaaaaaaaaaa", CandidateID: "candidate-1",
		CandidateDigest: candidate, TreeDigest: tree,
		ChangedPaths:  []string{"internal/example/example.go"},
		Marker:        "<!-- forge-publication:publication-1 -->",
		CredentialRef: "github-app-installation:forge-trial",
		Files: []PublicationFile{{
			Path: "internal/example/example.go", Content: []byte("package example\n"),
			ContentDigest: bytesDigest([]byte("package example\n")),
		}},
	}
}
