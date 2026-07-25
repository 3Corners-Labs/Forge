// Package verification owns Assurance Plans and independent verification
// results. The isolated verifier implementation is supplied through a port.
package verification

import (
	"context"
	"errors"
	"slices"
	"sort"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
	"github.com/3corners/forge/internal/delivery"
)

type Runner interface {
	Verify(context.Context, verificationcontract.AssurancePlan, workercontract.CandidateBundle, string) (verificationcontract.Result, error)
}

type Store interface {
	SavePlan(context.Context, verificationcontract.AssurancePlan) error
	SaveResult(context.Context, verificationcontract.Result) error
}

type Profile struct {
	Version         string
	Controls        []verificationcontract.Control
	MaximumLifetime time.Duration
}

type Service struct {
	runner  Runner
	store   Store
	profile Profile
	now     func() time.Time
}

func NewService(runner Runner, store Store, profile Profile, now func() time.Time) (*Service, error) {
	sort.Slice(profile.Controls, func(i, j int) bool { return profile.Controls[i].ID < profile.Controls[j].ID })
	if runner == nil || store == nil || now == nil || profile.Version == "" ||
		profile.MaximumLifetime <= 0 || len(profile.Controls) == 0 {
		return nil, errors.New("verification service configuration is invalid")
	}
	return &Service{runner: runner, store: store, profile: profile, now: now}, nil
}

type VerifyRequest struct {
	CorrelationID string
	Run           delivery.DeliveryRun
	Bundle        workercontract.CandidateBundle
	FreshCheckout string
}

type VerifyResult struct {
	Plan   verificationcontract.AssurancePlan
	Result verificationcontract.Result
}

func (service *Service) Verify(ctx context.Context, request VerifyRequest) (VerifyResult, error) {
	at := service.now()
	plan, candidate := request.Run.Plan(), request.Run.Candidate()
	// The producing principal is the workload identity Policy issued for the
	// worker phase. It is recorded so that verification and, later, publication
	// can be held to AC.8.1: an agent may not certify what it produced. The
	// candidate cannot assert this value.
	execution := request.Run.Execution()
	if ctx == nil || at.IsZero() || request.CorrelationID == "" || request.FreshCheckout == "" ||
		request.Run.State() != delivery.VerifyingState || plan == nil || candidate == nil ||
		execution == nil || execution.WorkerSubjectID == "" ||
		request.Bundle.Digest != candidate.Digest || request.Bundle.PlanDigest != plan.Digest {
		return VerifyResult{}, failure.New(failure.InvalidRequest, "verification input is invalid", request.CorrelationID, failure.RetryNone)
	}
	expiresAt := at.Add(service.profile.MaximumLifetime)
	if plan.ExpiresAt.Before(expiresAt) {
		expiresAt = plan.ExpiresAt
	}
	assurancePlan, err := verificationcontract.BindPlan(verificationcontract.AssurancePlan{
		SchemaVersion: verificationcontract.AssurancePlanVersion, CandidateDigest: candidate.Digest,
		PlanDigest: plan.Digest, RiskClass: plan.RiskClass, Classification: plan.Classification,
		AffectedCapabilities: slices.Clone(plan.AffectedCapabilities), ChangedPaths: slices.Clone(candidate.ChangedPaths),
		Controls: slices.Clone(service.profile.Controls), ProducedBy: execution.WorkerSubjectID,
		Environment: "developer-trial",
		ExpiresAt:   expiresAt, EvidenceCorrelationID: request.CorrelationID,
	})
	if err != nil {
		return VerifyResult{}, failure.New(failure.Indeterminate, "assurance plan could not be constructed", request.CorrelationID, failure.RetryNone)
	}
	if err := service.store.SavePlan(ctx, assurancePlan); err != nil {
		return VerifyResult{}, failure.New(failure.DependencyUnavailable, "assurance plan could not be persisted", request.CorrelationID, failure.RetryBounded)
	}
	result, err := service.runner.Verify(ctx, assurancePlan, request.Bundle, request.FreshCheckout)
	if err != nil {
		return VerifyResult{}, failure.New(failure.DependencyUnavailable, "independent verifier did not complete", request.CorrelationID, failure.RetryBounded)
	}
	if err := result.Validate(assurancePlan); err != nil {
		return VerifyResult{}, failure.New(failure.VerificationFailed, "verification result could not be verified", request.CorrelationID, failure.RetryNone)
	}
	if err := service.store.SaveResult(ctx, result); err != nil {
		return VerifyResult{}, failure.New(failure.DependencyUnavailable, "verification result could not be persisted", request.CorrelationID, failure.RetryBounded)
	}
	return VerifyResult{Plan: assurancePlan, Result: result}, nil
}
