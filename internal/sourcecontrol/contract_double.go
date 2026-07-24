package sourcecontrol

import (
	"context"
	"errors"
	"sync"
)

type FailurePoint string

const (
	FailNone              FailurePoint = ""
	FailBeforeBranch      FailurePoint = "before-branch"
	FailAfterBranch       FailurePoint = "after-branch"
	FailBeforePullRequest FailurePoint = "before-pull-request"
	FailAfterPullRequest  FailurePoint = "after-pull-request"
	FailReconcile         FailurePoint = "reconcile"
)

// ContractDouble is deterministic and credential-free. Failure injection
// models transport ambiguity without introducing a live provider side effect.
type ContractDouble struct {
	mu           sync.Mutex
	states       map[string]RemoteState
	failurePoint FailurePoint
	failed       bool
	apiCalls     int
}

func NewContractDouble(failurePoint FailurePoint) *ContractDouble {
	return &ContractDouble{states: make(map[string]RemoteState), failurePoint: failurePoint}
}

func (provider *ContractDouble) Reconcile(ctx context.Context, operation Operation) (RemoteState, error) {
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.apiCalls++
	if ctx == nil {
		return RemoteState{}, ErrProviderUnavailable
	}
	if provider.failurePoint == FailReconcile && !provider.failed {
		provider.failed = true
		return RemoteState{}, ErrProviderUnavailable
	}
	return provider.states[operation.ID], nil
}

func (provider *ContractDouble) CreateBranch(ctx context.Context, operation Operation) error {
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.apiCalls++
	if ctx == nil {
		return ErrProviderUnavailable
	}
	if provider.failurePoint == FailBeforeBranch && !provider.failed {
		provider.failed = true
		return ErrProviderAmbiguous
	}
	state := provider.states[operation.ID]
	if state != (RemoteState{}) {
		return ErrProviderConflict
	}
	provider.states[operation.ID] = RemoteState{
		BranchExists: true, Branch: operation.Branch, CommitID: "double-commit-" + operation.ID,
		CandidateDigest: operation.CandidateDigest, TreeDigest: operation.TreeDigest,
		BaseBranch: operation.BaseBranch, BaseRevision: operation.BaseRevision,
		Marker: operation.Marker, RepositoryID: operation.RepositoryID,
	}
	if provider.failurePoint == FailAfterBranch && !provider.failed {
		provider.failed = true
		return ErrProviderAmbiguous
	}
	return nil
}

func (provider *ContractDouble) OpenPullRequest(ctx context.Context, operation Operation, title, body string) error {
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.apiCalls++
	if ctx == nil || title == "" || body == "" {
		return ErrProviderUnavailable
	}
	if provider.failurePoint == FailBeforePullRequest && !provider.failed {
		provider.failed = true
		return ErrProviderAmbiguous
	}
	state, found := provider.states[operation.ID]
	if !found || !state.BranchExists || state.PullRequestID != "" {
		return ErrProviderConflict
	}
	state.PullRequestID = "double-pr-" + operation.ID
	state.PullRequestURL = "https://contract-double.invalid/pulls/" + operation.ID
	state.PullRequestOpen = true
	provider.states[operation.ID] = state
	if provider.failurePoint == FailAfterPullRequest && !provider.failed {
		provider.failed = true
		return ErrProviderAmbiguous
	}
	return nil
}

func (provider *ContractDouble) APIRequests() int {
	provider.mu.Lock()
	defer provider.mu.Unlock()
	return provider.apiCalls
}

func (provider *ContractDouble) ReplaceState(operationID string, state RemoteState) error {
	if operationID == "" {
		return errors.New("operation identifier is required")
	}
	provider.mu.Lock()
	defer provider.mu.Unlock()
	provider.states[operationID] = state
	return nil
}
