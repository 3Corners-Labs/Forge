// Package execution owns the bounded worker launch boundary. It constructs a
// closed request from trusted profiles and invokes a replaceable isolated
// runner; it never imports the worker implementation.
package execution

import (
	"context"
	"errors"
	"slices"
	"sort"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
	"github.com/3corners/forge/internal/delivery"
)

type Runner interface {
	Run(context.Context, workercontract.Request, string) (workercontract.CandidateBundle, workercontract.Result, error)
}

type TransformationProfile struct {
	ID              string
	Version         string
	Replacements    []workercontract.FileReplacement
	MaxProcesses    int64
	MaxOutputBytes  int64
	MaxChangedFiles int64
}

type Service struct {
	runner Runner
	now    func() time.Time
}

func NewService(runner Runner, now func() time.Time) (*Service, error) {
	if runner == nil || now == nil {
		return nil, errors.New("execution service configuration is invalid")
	}
	return &Service{runner: runner, now: now}, nil
}

type PrepareRequest struct {
	CorrelationID  string
	IdempotencyKey string
	Run            delivery.DeliveryRun
	Envelope       trust.ExecutionEnvelope
	Profile        TransformationProfile
	Workspace      string
	Attempt        int64
}

type PrepareResult struct {
	Request workercontract.Request
	Bundle  workercontract.CandidateBundle
	Result  workercontract.Result
}

func (service *Service) Prepare(ctx context.Context, command PrepareRequest) (PrepareResult, error) {
	at := service.now()
	plan, contextRecord, approval, authority := command.Run.Plan(), command.Run.Context(), command.Run.Approval(), command.Run.Execution()
	if ctx == nil || at.IsZero() || command.CorrelationID == "" || command.IdempotencyKey == "" ||
		command.Workspace == "" || command.Run.State() != delivery.ExecutingState ||
		plan == nil || contextRecord == nil || approval == nil || authority == nil ||
		command.Profile.ID != workercontract.TransformationV1 || command.Profile.Version == "" ||
		len(command.Profile.Replacements) == 0 || command.Attempt < 1 {
		return PrepareResult{}, failure.New(failure.InvalidRequest, "worker preparation request is invalid", command.CorrelationID, failure.RetryNone)
	}
	if err := command.Envelope.VerifyIntegrity(command.CorrelationID); err != nil {
		return PrepareResult{}, err
	}
	if err := command.Envelope.AuthorizeAction(workercontract.PrepareChangeAction, at, command.CorrelationID); err != nil {
		return PrepareResult{}, err
	}
	if command.Envelope.ID != authority.EnvelopeID || command.Envelope.Integrity == nil ||
		command.Envelope.Integrity.Digest != authority.EnvelopeDigest ||
		command.Envelope.Principal.AgentID != authority.WorkerSubjectID ||
		!slices.Equal(command.Envelope.Principal.DelegationChain, []string{authority.DelegationID}) {
		return PrepareResult{}, failure.New(failure.Forbidden, "worker envelope is not current delivery authority", command.CorrelationID, failure.RetryNewAuthority)
	}
	paths := make([]string, 0, len(command.Profile.Replacements))
	replacements := make([]workercontract.FileReplacement, len(command.Profile.Replacements))
	copy(replacements, command.Profile.Replacements)
	sort.Slice(replacements, func(i, j int) bool { return replacements[i].Path < replacements[j].Path })
	for index, replacement := range replacements {
		paths = append(paths, replacement.Path)
		replacements[index].Content = slices.Clone(replacement.Content)
	}
	if !slices.Equal(paths, plan.ChangePaths) {
		return PrepareResult{}, failure.New(failure.Forbidden, "trusted transformation exceeds the approved plan paths", command.CorrelationID, failure.RetryNone)
	}
	request, err := workercontract.BindRequest(workercontract.Request{
		SchemaVersion: workercontract.RequestVersion, EnvelopeID: command.Envelope.ID,
		DelegationID: authority.DelegationID, PlanID: plan.ID, PlanDigest: plan.Digest,
		PlanApprovalRecordID: approval.RecordID, ContextPackageID: contextRecord.PackageID,
		ContextManifestDigest: contextRecord.ManifestDigest, SourceRevision: contextRecord.SourceRevisions[0],
		Scope: workercontract.Scope{
			TenantID: command.Run.TenantID(), WorkspaceID: command.Run.WorkspaceID(),
			ProjectID: command.Run.ProjectID(), RepositoryID: command.Run.RepositoryID(), DeliveryRunID: command.Run.ID(),
		},
		Action: workercontract.PrepareChangeAction, TransformationID: command.Profile.ID,
		TransformationVersion: command.Profile.Version, AllowedReadPaths: slices.Clone(plan.ChangePaths),
		AllowedWritePaths: slices.Clone(plan.ChangePaths), DependencyDigests: []string{plan.Digest, contextRecord.ManifestDigest},
		Replacements: replacements,
		Limits: workercontract.Limits{
			MaxDurationSeconds: min(command.Envelope.Limits.MaxDurationSeconds, plan.Budgets.MaxDurationSeconds),
			MaxProcesses:       command.Profile.MaxProcesses, MaxOutputBytes: command.Profile.MaxOutputBytes,
			MaxChangedFiles: command.Profile.MaxChangedFiles, MaxAttempts: plan.Budgets.MaxAttempts,
		},
		Attempt: command.Attempt, ExpiresAt: command.Envelope.ExpiresAt,
		EvidenceCorrelationID: command.CorrelationID,
	})
	if err != nil {
		return PrepareResult{}, failure.New(failure.Indeterminate, "trusted worker request could not be constructed", command.CorrelationID, failure.RetryNone)
	}
	bundle, result, err := service.runner.Run(ctx, request, command.Workspace)
	if err != nil {
		return PrepareResult{}, failure.New(failure.DependencyUnavailable, "isolated worker did not complete", command.CorrelationID, failure.RetryBounded)
	}
	if err := bundle.Validate(); err != nil || result.RequestID != request.ID ||
		result.RequestDigest != request.Digest || result.CandidateDigest != bundle.Digest ||
		result.BundleDigest != bundle.Digest || result.TreeDigest != workercontract.CandidateTreeDigest(bundle) ||
		!validWorkerResult(result, request, bundle, at) {
		return PrepareResult{}, failure.New(failure.VerificationFailed, "isolated worker result could not be verified", command.CorrelationID, failure.RetryNone)
	}
	return PrepareResult{Request: request, Bundle: bundle, Result: result}, nil
}

func validWorkerResult(result workercontract.Result, request workercontract.Request, bundle workercontract.CandidateBundle, started time.Time) bool {
	paths := make([]string, 0, len(bundle.Files))
	for _, file := range bundle.Files {
		paths = append(paths, file.Path)
	}
	return result.SchemaVersion == workercontract.ResultVersion &&
		slices.Equal(result.ChangedPaths, paths) &&
		result.ResourceUse.ChangedFiles == int64(len(paths)) &&
		result.ResourceUse.ChangedFiles <= request.Limits.MaxChangedFiles &&
		result.ResourceUse.OutputBytes >= 0 && result.ResourceUse.OutputBytes <= request.Limits.MaxOutputBytes &&
		result.ResourceUse.DurationMillis >= 0 &&
		result.ResourceUse.DurationMillis <= request.Limits.MaxDurationSeconds*1000 &&
		result.CleanupOutcome == "completed" && !result.CompletedAt.Before(started) &&
		result.CompletedAt.Before(request.ExpiresAt)
}

func min(left, right int64) int64 {
	if left < right {
		return left
	}
	return right
}
