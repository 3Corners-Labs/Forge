package sourcecontrol

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const githubAPIBase = "https://api.github.com"

// InstallationTokenSource returns a short-lived GitHub App installation
// token for one configured, non-secret credential reference.
type InstallationTokenSource interface {
	Token(context.Context, string) (string, error)
}

type GitHubAppConfig struct {
	Owner         string
	Repository    string
	RepositoryID  string
	CredentialRef string
}

// GitHubAppProvider implements only the ADR-0019 metadata, contents and
// pull-request endpoints needed for publication. It deliberately exposes no
// merge, release, workflow, secret, deployment or administration method.
type GitHubAppProvider struct {
	config GitHubAppConfig
	client *http.Client
	tokens InstallationTokenSource
}

func NewGitHubAppProvider(config GitHubAppConfig, client *http.Client, tokens InstallationTokenSource) (*GitHubAppProvider, error) {
	if config.Owner == "" || config.Repository == "" || config.RepositoryID == "" ||
		config.CredentialRef == "" || !safeGitHubSegment(config.Owner) ||
		!safeGitHubSegment(config.Repository) || client == nil || tokens == nil {
		return nil, errors.New("GitHub App publication configuration is invalid")
	}
	return &GitHubAppProvider{config: config, client: client, tokens: tokens}, nil
}

func (provider *GitHubAppProvider) Reconcile(ctx context.Context, operation Operation) (RemoteState, error) {
	if !provider.matches(operation) {
		return RemoteState{}, ErrProviderConflict
	}
	var ref githubRef
	status, err := provider.doJSON(ctx, http.MethodGet,
		"/repos/"+provider.config.Owner+"/"+provider.config.Repository+"/git/ref/heads/"+url.PathEscape(operation.Branch),
		nil, &ref, false)
	if err != nil {
		return RemoteState{}, err
	}
	if status == http.StatusNotFound {
		return RemoteState{}, nil
	}
	if status != http.StatusOK || ref.Object.SHA == "" {
		return RemoteState{}, ErrProviderUnavailable
	}
	var commit githubCommit
	status, err = provider.doJSON(ctx, http.MethodGet,
		"/repos/"+provider.config.Owner+"/"+provider.config.Repository+"/git/commits/"+url.PathEscape(ref.Object.SHA),
		nil, &commit, false)
	if err != nil || status != http.StatusOK {
		return RemoteState{}, ErrProviderUnavailable
	}
	marker, candidate, tree, base, ok := parseCommitMessage(commit.Message)
	if !ok || len(commit.Parents) != 1 || commit.Parents[0].SHA != base {
		return RemoteState{}, ErrProviderConflict
	}
	state := RemoteState{
		BranchExists: true, Branch: operation.Branch, CommitID: ref.Object.SHA,
		CandidateDigest: candidate, TreeDigest: tree, BaseBranch: operation.BaseBranch,
		BaseRevision: base, Marker: marker, RepositoryID: operation.RepositoryID,
	}
	var pulls []githubPull
	query := url.Values{
		"state": {"open"}, "base": {operation.BaseBranch},
		"head": {provider.config.Owner + ":" + operation.Branch},
	}.Encode()
	status, err = provider.doJSON(ctx, http.MethodGet,
		"/repos/"+provider.config.Owner+"/"+provider.config.Repository+"/pulls?"+query,
		nil, &pulls, false)
	if err != nil || status != http.StatusOK {
		return RemoteState{}, ErrProviderUnavailable
	}
	for _, pull := range pulls {
		if strings.Contains(pull.Body, operation.Marker) {
			if state.PullRequestID != "" {
				return RemoteState{}, ErrProviderConflict
			}
			state.PullRequestID = strconv.FormatInt(pull.Number, 10)
			state.PullRequestURL = pull.HTMLURL
			state.PullRequestOpen = pull.State == "open"
		}
	}
	return state, nil
}

func (provider *GitHubAppProvider) CreateBranch(ctx context.Context, operation Operation) error {
	if !provider.matches(operation) || len(operation.Files) == 0 {
		return ErrProviderConflict
	}
	var baseRef githubRef
	status, err := provider.doJSON(ctx, http.MethodGet,
		"/repos/"+provider.config.Owner+"/"+provider.config.Repository+"/git/ref/heads/"+url.PathEscape(operation.BaseBranch),
		nil, &baseRef, false)
	if err != nil || status != http.StatusOK || baseRef.Object.SHA != operation.BaseRevision {
		return ErrProviderConflict
	}
	var baseCommit githubCommit
	status, err = provider.doJSON(ctx, http.MethodGet,
		"/repos/"+provider.config.Owner+"/"+provider.config.Repository+"/git/commits/"+url.PathEscape(operation.BaseRevision),
		nil, &baseCommit, false)
	if err != nil || status != http.StatusOK || baseCommit.Tree.SHA == "" {
		return ErrProviderUnavailable
	}
	entries := make([]githubTreeEntry, len(operation.Files))
	for index, file := range operation.Files {
		var blob githubSHA
		status, err = provider.doJSON(ctx, http.MethodPost,
			"/repos/"+provider.config.Owner+"/"+provider.config.Repository+"/git/blobs",
			map[string]any{"content": base64.StdEncoding.EncodeToString(file.Content), "encoding": "base64"},
			&blob, true)
		if err != nil || status != http.StatusCreated || blob.SHA == "" {
			return mutationFailure(err)
		}
		entries[index] = githubTreeEntry{Path: file.Path, Mode: "100644", Type: "blob", SHA: blob.SHA}
	}
	var tree githubSHA
	status, err = provider.doJSON(ctx, http.MethodPost,
		"/repos/"+provider.config.Owner+"/"+provider.config.Repository+"/git/trees",
		map[string]any{"base_tree": baseCommit.Tree.SHA, "tree": entries}, &tree, true)
	if err != nil || status != http.StatusCreated || tree.SHA == "" {
		return mutationFailure(err)
	}
	message := operation.Marker + "\n" +
		"candidate=" + operation.CandidateDigest + "\n" +
		"tree=" + operation.TreeDigest + "\n" +
		"base=" + operation.BaseRevision
	var commit githubSHA
	status, err = provider.doJSON(ctx, http.MethodPost,
		"/repos/"+provider.config.Owner+"/"+provider.config.Repository+"/git/commits",
		map[string]any{"message": message, "tree": tree.SHA, "parents": []string{operation.BaseRevision}},
		&commit, true)
	if err != nil || status != http.StatusCreated || commit.SHA == "" {
		return mutationFailure(err)
	}
	status, err = provider.doJSON(ctx, http.MethodPost,
		"/repos/"+provider.config.Owner+"/"+provider.config.Repository+"/git/refs",
		map[string]any{"ref": "refs/heads/" + operation.Branch, "sha": commit.SHA}, nil, true)
	if err != nil || status != http.StatusCreated {
		return mutationFailure(err)
	}
	return nil
}

func (provider *GitHubAppProvider) OpenPullRequest(ctx context.Context, operation Operation, title, body string) error {
	if !provider.matches(operation) || title == "" || body == "" {
		return ErrProviderConflict
	}
	var pull githubPull
	status, err := provider.doJSON(ctx, http.MethodPost,
		"/repos/"+provider.config.Owner+"/"+provider.config.Repository+"/pulls",
		map[string]any{
			"title": title, "head": operation.Branch, "base": operation.BaseBranch,
			"body": body + "\n\n" + operation.Marker, "draft": true,
		}, &pull, true)
	if err != nil || status != http.StatusCreated || pull.Number < 1 || pull.HTMLURL == "" {
		return mutationFailure(err)
	}
	return nil
}

func (provider *GitHubAppProvider) matches(operation Operation) bool {
	return provider != nil && provider.client != nil && provider.tokens != nil &&
		operation.ContractVersion == ProviderContractVersion && operation.Provider == "github-app" &&
		operation.RepositoryID == provider.config.RepositoryID &&
		operation.CredentialRef == provider.config.CredentialRef &&
		operation.BaseBranch != "" && operation.BaseRevision != "" && operation.Branch != "" &&
		validDigest(operation.CandidateDigest) && validDigest(operation.TreeDigest) &&
		operation.Marker != ""
}

func (provider *GitHubAppProvider) doJSON(ctx context.Context, method, path string, body any, output any, mutation bool) (int, error) {
	if ctx == nil || !allowedGitHubPath(method, path) {
		return 0, ErrProviderConflict
	}
	var input io.Reader
	if body != nil {
		encoded, err := json.Marshal(body)
		if err != nil {
			return 0, ErrProviderConflict
		}
		input = bytes.NewReader(encoded)
	}
	token, err := provider.tokens.Token(ctx, provider.config.CredentialRef)
	if err != nil || strings.TrimSpace(token) == "" {
		return 0, ErrProviderUnavailable
	}
	request, err := http.NewRequestWithContext(ctx, method, githubAPIBase+path, input)
	if err != nil {
		return 0, ErrProviderUnavailable
	}
	request.Header.Set("Authorization", "Bearer "+token)
	request.Header.Set("Accept", "application/vnd.github+json")
	request.Header.Set("X-GitHub-Api-Version", "2022-11-28")
	if body != nil {
		request.Header.Set("Content-Type", "application/json")
	}
	client := *provider.client
	client.CheckRedirect = func(*http.Request, []*http.Request) error {
		return errors.New("GitHub App publication redirects are forbidden")
	}
	response, err := client.Do(request)
	if err != nil {
		if mutation {
			return 0, ErrProviderAmbiguous
		}
		return 0, ErrProviderUnavailable
	}
	defer response.Body.Close()
	if output == nil {
		_, _ = io.Copy(io.Discard, io.LimitReader(response.Body, 64<<10))
		return response.StatusCode, nil
	}
	encoded, err := io.ReadAll(io.LimitReader(response.Body, 1<<20+1))
	if err != nil || len(encoded) > 1<<20 {
		if mutation {
			return response.StatusCode, ErrProviderAmbiguous
		}
		return response.StatusCode, ErrProviderUnavailable
	}
	if response.StatusCode >= 200 && response.StatusCode < 300 {
		if err := json.Unmarshal(encoded, output); err != nil {
			if mutation {
				return response.StatusCode, ErrProviderAmbiguous
			}
			return response.StatusCode, ErrProviderUnavailable
		}
	}
	return response.StatusCode, nil
}

func allowedGitHubPath(method, value string) bool {
	if strings.Contains(value, "..") || strings.ContainsAny(value, "\r\n") {
		return false
	}
	allowedMethod := method == http.MethodGet || method == http.MethodPost
	if !allowedMethod || !strings.HasPrefix(value, "/repos/") {
		return false
	}
	for _, forbidden := range []string{"/merge", "/releases", "/deployments", "/actions", "/secrets", "/hooks", "/collaborators"} {
		if strings.Contains(value, forbidden) {
			return false
		}
	}
	return strings.Contains(value, "/git/ref") || strings.Contains(value, "/git/commits") ||
		strings.Contains(value, "/git/blobs") || strings.Contains(value, "/git/trees") ||
		strings.Contains(value, "/pulls")
}

func mutationFailure(err error) error {
	if errors.Is(err, ErrProviderConflict) {
		return ErrProviderConflict
	}
	return ErrProviderAmbiguous
}

func safeGitHubSegment(value string) bool {
	if value == "" || value == "." || value == ".." {
		return false
	}
	for _, character := range value {
		if (character < 'a' || character > 'z') && (character < 'A' || character > 'Z') &&
			(character < '0' || character > '9') && character != '-' && character != '_' && character != '.' {
			return false
		}
	}
	return true
}

func parseCommitMessage(message string) (marker, candidate, tree, base string, ok bool) {
	lines := strings.Split(message, "\n")
	if len(lines) != 4 || !strings.HasPrefix(lines[0], "<!-- forge-publication:") ||
		!strings.HasSuffix(lines[0], " -->") {
		return "", "", "", "", false
	}
	marker = lines[0]
	candidate = strings.TrimPrefix(lines[1], "candidate=")
	tree = strings.TrimPrefix(lines[2], "tree=")
	base = strings.TrimPrefix(lines[3], "base=")
	return marker, candidate, tree, base,
		validDigest(candidate) && validDigest(tree) && base != "" &&
			strings.HasPrefix(lines[1], "candidate=") &&
			strings.HasPrefix(lines[2], "tree=") &&
			strings.HasPrefix(lines[3], "base=")
}

type githubRef struct {
	Object struct {
		SHA string `json:"sha"`
	} `json:"object"`
}

type githubSHA struct {
	SHA string `json:"sha"`
}

type githubCommit struct {
	Message string      `json:"message"`
	Tree    githubSHA   `json:"tree"`
	Parents []githubSHA `json:"parents"`
}

type githubTreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	SHA  string `json:"sha"`
}

type githubPull struct {
	Number  int64  `json:"number"`
	HTMLURL string `json:"html_url"`
	State   string `json:"state"`
	Body    string `json:"body"`
}

var _ Provider = (*GitHubAppProvider)(nil)
