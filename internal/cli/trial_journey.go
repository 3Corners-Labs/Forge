package cli

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/3corners/forge/internal/approval"
	forgecontext "github.com/3corners/forge/internal/context"
	changecontract "github.com/3corners/forge/internal/contracts/change"
	"github.com/3corners/forge/internal/contracts/discovery"
	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
	workercontract "github.com/3corners/forge/internal/contracts/worker"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/evidence"
	"github.com/3corners/forge/internal/execution"
	"github.com/3corners/forge/internal/experience"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/intelligence"
	"github.com/3corners/forge/internal/intent"
	"github.com/3corners/forge/internal/policy"
	"github.com/3corners/forge/internal/sourcecontrol"
	"github.com/3corners/forge/internal/tenant"
	"github.com/3corners/forge/internal/verification"
)

const trialCorrelation = "forge.m1.cli.developer-trial"

type trialIdentityProvider struct {
	assertion identity.Assertion
}

func (provider trialIdentityProvider) Verify(context.Context, string) (identity.Assertion, error) {
	return provider.assertion, nil
}

type trialBindingSource struct {
	calls int
}

func (source *trialBindingSource) ResolveRepositoryBinding(_ context.Context, tenantID, workspaceID, projectID, selector string) (tenant.RepositoryBindingDefinition, bool, error) {
	source.calls++
	if tenantID != "tenant-developer-trial" || workspaceID != "workspace-forge" ||
		projectID != "project-forge" || selector != "forge" {
		return tenant.RepositoryBindingDefinition{}, false, nil
	}
	return tenant.RepositoryBindingDefinition{
		TenantID: tenantID, WorkspaceID: workspaceID, ProjectID: projectID,
		RepositoryID: "repository-forge", Selector: selector,
		CanonicalLocator: "contract-double://3Corners-Labs/Forge", Classification: "Internal",
	}, true, nil
}

type trialContextCatalog struct{}

func (trialContextCatalog) Discover(_ context.Context, request forgecontext.SourceDiscoveryRequest) ([]forgecontext.SourceRecord, error) {
	if request.TenantID != "tenant-developer-trial" || request.WorkspaceID != "workspace-forge" ||
		request.ProjectID != "project-forge" || request.RepositoryID != "repository-forge" {
		return nil, errors.New("context scope is not authority-derived")
	}
	freshness := time.Now().UTC().Add(-time.Minute)
	return []forgecontext.SourceRecord{
		{
			TenantID: request.TenantID, WorkspaceID: request.WorkspaceID, ProjectID: request.ProjectID,
			RepositoryID: request.RepositoryID, SourceID: "accepted-architecture", SourceRevision: "m1-c",
			Locator: "docs/architecture/vertical-slice.md", Classification: "Internal",
			Authority: "accepted-knowledge", Purpose: delivery.M1ChangePurpose,
			FreshnessAt: &freshness, ContentDigest: "sha256:" + strings.Repeat("1", 64),
		},
		{
			TenantID: request.TenantID, WorkspaceID: request.WorkspaceID, ProjectID: request.ProjectID,
			RepositoryID: request.RepositoryID, SourceID: "promoted-intent", SourceRevision: "v1",
			Locator: "forge.intent", Classification: "Internal", Authority: "task-instruction",
			Purpose: delivery.M1ChangePurpose, FreshnessAt: &freshness,
			ContentDigest: "sha256:" + strings.Repeat("2", 64),
		},
	}, nil
}

type policyResolvers struct {
	values []*policy.Evaluator
}

func (resolver *policyResolvers) add(value *policy.Evaluator) {
	resolver.values = append(resolver.values, value)
}

func (resolver *policyResolvers) ResolvePolicyDecision(ctx context.Context, id string) (trust.PolicyDecision, bool, error) {
	for index := len(resolver.values) - 1; index >= 0; index-- {
		value, found, err := resolver.values[index].ResolvePolicyDecision(ctx, id)
		if err != nil || found {
			return value, found, err
		}
	}
	return trust.PolicyDecision{}, false, nil
}

type trialVerificationStore struct {
	plans   map[string]verificationcontract.AssurancePlan
	results map[string]verificationcontract.Result
}

func newTrialVerificationStore() *trialVerificationStore {
	return &trialVerificationStore{plans: make(map[string]verificationcontract.AssurancePlan), results: make(map[string]verificationcontract.Result)}
}

func (store *trialVerificationStore) SavePlan(_ context.Context, value verificationcontract.AssurancePlan) error {
	if existing, found := store.plans[value.ID]; found && existing.Digest != value.Digest {
		return errors.New("assurance plan conflict")
	}
	store.plans[value.ID] = value
	return nil
}

func (store *trialVerificationStore) SaveResult(_ context.Context, value verificationcontract.Result) error {
	if existing, found := store.results[value.ID]; found && existing.Digest != value.Digest {
		return errors.New("verification result conflict")
	}
	store.results[value.ID] = value
	return nil
}

type trialWorkerRunner struct {
	now func() time.Time
}

func (runner trialWorkerRunner) Run(_ context.Context, request workercontract.Request, _ string) (workercontract.CandidateBundle, workercontract.Result, error) {
	if err := request.Validate(runner.now()); err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	files := make([]workercontract.CandidateFile, len(request.Replacements))
	var outputBytes int64
	for index, replacement := range request.Replacements {
		files[index] = workercontract.CandidateFile{Path: replacement.Path, SourceDigest: replacement.ExpectedSourceDigest, Content: slices.Clone(replacement.Content)}
		outputBytes += int64(len(replacement.Content))
	}
	bundle, err := workercontract.BindBundle(workercontract.CandidateBundle{
		SchemaVersion: workercontract.BundleVersion, SourceRevision: request.SourceRevision,
		PlanDigest: request.PlanDigest, Files: files,
	})
	if err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	paths := make([]string, len(files))
	for index := range files {
		paths[index] = files[index].Path
	}
	completed := runner.now()
	return bundle, workercontract.Result{
		SchemaVersion: workercontract.ResultVersion, RequestID: request.ID, RequestDigest: request.Digest,
		CandidateDigest: bundle.Digest, TreeDigest: workercontract.CandidateTreeDigest(bundle),
		ChangedPaths: paths, BundleDigest: bundle.Digest,
		ResourceUse:    workercontract.ResourceUse{DurationMillis: 1, OutputBytes: outputBytes, ChangedFiles: int64(len(files))},
		CleanupOutcome: "completed", EvidenceReferences: []string{"candidate:" + bundle.Digest},
		CompletedAt: completed,
	}, nil
}

type trialVerifierRunner struct {
	now func() time.Time
}

func (runner trialVerifierRunner) Verify(_ context.Context, plan verificationcontract.AssurancePlan, bundle workercontract.CandidateBundle, _ string) (verificationcontract.Result, error) {
	if err := plan.Validate(runner.now()); err != nil {
		return verificationcontract.Result{}, err
	}
	findings := make([]verificationcontract.Finding, len(plan.Controls))
	for index, control := range plan.Controls {
		findings[index] = verificationcontract.Finding{
			ControlID: control.ID, Outcome: "passed", Code: "passed",
			EvidenceDigest: "sha256:" + strings.Repeat(fmt.Sprintf("%x", (index%15)+1), 64),
			CandidateOwned: control.CandidateOwned, Consequential: control.Consequential,
		}
	}
	return verificationcontract.BindResult(verificationcontract.Result{
		SchemaVersion: verificationcontract.ResultVersion, AssurancePlanID: plan.ID,
		AssurancePlanDigest: plan.Digest, CandidateDigest: bundle.Digest,
		Environment: plan.Environment, Findings: findings, Outcome: "passed",
		Independent: true, CompletedAt: runner.now(), EvidenceCorrelationID: plan.EvidenceCorrelationID,
	})
}

type trialSession struct {
	proposal              experience.RequirementsProposal
	contextPackage        trust.ContextPackage
	run                   delivery.DeliveryRun
	planRequest           trust.ApprovalRequest
	planRecord            trust.ApprovalRecordV2
	bundle                workercontract.CandidateBundle
	publicationRequest    trust.ApprovalRequest
	publicationRecord     trust.ApprovalRecordV2
	publicationEnvelope   trust.ExecutionEnvelope
	publicationDecisionID string
}

type developerTrialJourney struct {
	now               func() time.Time
	human             identity.Principal
	taskAgent         identity.Principal
	worker            identity.Principal
	directory         *tenant.Directory
	bindingSource     *trialBindingSource
	scope             tenant.Scope
	binding           tenant.RepositoryBinding
	experience        *experience.Service
	discovery         *experience.DiscoveryService
	changePreparation *intelligence.ChangeService
	changeActivation  intelligence.ProviderActivation
	intent            *intent.Service
	delivery          *delivery.Service
	policies          *policyResolvers
	approvals         *approval.Service
	evidence          *evidence.Service
	sourceControl     *sourcecontrol.Service
	provider          sourcecontrol.Provider
	sessions          map[string]*trialSession
	discoveryLimits   discovery.Limits
	catalog           *forgecontext.RepositoryCatalog
	workerRunner      execution.Runner
	verifierRunner    verification.Runner
	runtimeRoot       string
	operational       bool
	persistence       journeyPersistence
	nextID            int
}

func newDeveloperTrialJourney(options journeyOptions) (guidedJourney, error) {
	anchor := time.Now().UTC()
	now := func() time.Time { return anchor }
	runtime, err := newOperationalRuntime(options, now)
	if err != nil {
		return nil, err
	}
	profile := os.Getenv("FORGE_IDENTITY_PROFILE")
	humanID := "synthetic-developer-allowed"
	if profile == "observer-denied" {
		humanID = "synthetic-observer-denied"
	}
	human, err := trialPrincipal(humanID, identity.HumanPrincipal, now)
	if err != nil {
		return nil, err
	}
	taskAgent, err := trialPrincipal("synthetic-task-agent", identity.AgentPrincipal, now)
	if err != nil {
		return nil, err
	}
	worker, err := trialPrincipal("synthetic-worker", identity.WorkloadPrincipal, now)
	if err != nil {
		return nil, err
	}
	bindings := &trialBindingSource{}
	memberships := []tenant.MembershipDefinition{}
	if humanID == "synthetic-developer-allowed" {
		memberships = append(memberships, tenant.MembershipDefinition{PrincipalID: humanID, WorkspaceID: "workspace-forge"})
	}
	directory, err := tenant.NewDirectory(
		"developer-trial-authority",
		[]tenant.WorkspaceDefinition{{TenantID: "tenant-developer-trial", ID: "workspace-forge", Selector: "fixture-workspace"}},
		[]tenant.ProjectDefinition{{WorkspaceID: "workspace-forge", ID: "project-forge", Selector: "forge"}},
		memberships, bindings, now,
	)
	if err != nil {
		return nil, err
	}
	persistence, err := newJourneyPersistence(runtime.operational)
	if err != nil {
		return nil, err
	}
	deliveryService, err := delivery.NewService(persistence.delivery, now)
	if err != nil {
		_ = persistence.Close()
		return nil, err
	}
	policies := &policyResolvers{}
	approvals, err := approval.NewServiceWithPublication(
		deliveryService, deliveryService, policies, persistence.approval,
		approval.Profile{
			Issuer: "forge.approval.developer-trial", MaximumDecisionLifetime: 30 * time.Minute,
			SupportedConditions:     []string{"require-independent-verification"},
			EnforcedObligationTypes: []string{"human-interaction"},
		}, now,
	)
	if err != nil {
		_ = persistence.Close()
		return nil, err
	}
	evidenceService, err := evidence.NewService(persistence.evidence, now)
	if err != nil {
		_ = persistence.Close()
		return nil, err
	}
	value := &developerTrialJourney{
		now: now, human: human, taskAgent: taskAgent, worker: worker,
		directory: directory, bindingSource: bindings, delivery: deliveryService,
		policies: policies, approvals: approvals, evidence: evidenceService,
		provider: persistence.publication,
		sessions: make(map[string]*trialSession), catalog: runtime.catalog,
		workerRunner: runtime.worker, verifierRunner: runtime.verifier,
		runtimeRoot: runtime.runtimeRoot, operational: runtime.operational, persistence: persistence,
	}
	value.experience, err = experience.NewService(persistence.experience, now, value.newIdentifier)
	if err != nil {
		_ = persistence.Close()
		return nil, err
	}
	value.intent, err = intent.NewService(experienceStoreResolver{value.experience}, deliveryService, persistence.intent, now, value.newIdentifier)
	if err != nil {
		_ = persistence.Close()
		return nil, err
	}
	value.sourceControl, err = sourcecontrol.NewService(
		deliveryService, approvals, policies, evidenceService, value.provider,
		persistence.sourceControl, sourcecontrol.PublicationProfile{
			Provider: "contract-double", BaseBranch: "main", BaseRevision: "developer-trial-base",
			Title: "Forge M1 Developer Trial candidate",
			Body:  "Verified synthetic candidate. Human review and merge remain required.",
		}, now,
	)
	if err != nil {
		_ = persistence.Close()
		return nil, err
	}
	if options.ProviderActivationPath != "" {
		activation, loadErr := intelligence.LoadProviderActivation(options.ProviderActivationPath, anchor)
		if loadErr != nil {
			return nil, loadErr
		}
		if activation.Digest != options.ProviderActivationDigest {
			return nil, errors.New("provider activation does not match its separately approved digest")
		}
		if activation.Provider != "openai" {
			return nil, errors.New("only the approved OpenAI Developer Trial adapter is runnable")
		}
		liveClock := func() time.Time { return time.Now().UTC() }
		advisor, serviceErr := intelligence.NewOpenAIDiscoveryService(activation, liveClock)
		if serviceErr != nil {
			return nil, serviceErr
		}
		if configureErr := value.configureDiscoveryWithClock(
			advisor, liveClock, activation.Budgets,
		); configureErr != nil {
			return nil, configureErr
		}
	}
	if options.ChangeActivationPath != "" {
		activation, loadErr := intelligence.LoadProviderActivation(options.ChangeActivationPath, anchor)
		if loadErr != nil {
			return nil, loadErr
		}
		if activation.Digest != options.ChangeActivationDigest {
			return nil, errors.New("change activation does not match its separately approved digest")
		}
		if activation.Provider != "openai" ||
			!slices.Equal(activation.AllowedPurposes, []string{changecontract.PurposeV1}) {
			return nil, errors.New("only the approved OpenAI M1 change-preparation adapter is runnable")
		}
		liveClock := func() time.Time { return time.Now().UTC() }
		preparer, serviceErr := intelligence.NewOpenAIChangeService(activation, liveClock)
		if serviceErr != nil {
			return nil, serviceErr
		}
		value.changePreparation = preparer
		value.changeActivation = activation
	}
	return value, nil
}

func (journey *developerTrialJourney) Close() error {
	if journey == nil {
		return nil
	}
	return journey.persistence.Close()
}

type experienceStoreResolver struct {
	service *experience.Service
}

func (resolver experienceStoreResolver) ResolvePromotionSubject(ctx context.Context, tenantID, proposalID string) (intent.PromotionSubject, bool, error) {
	return resolver.service.ResolvePromotionSubject(ctx, tenantID, proposalID)
}

func (journey *developerTrialJourney) newIdentifier() string {
	journey.nextID++
	return fmt.Sprintf("forge-m1-%d-%d", journey.now().UnixNano(), journey.nextID)
}

func (journey *developerTrialJourney) configureDiscovery(advisor experience.DiscoveryAdvisor) error {
	return journey.configureDiscoveryWithClock(advisor, journey.now, defaultDiscoveryLimits())
}

func (journey *developerTrialJourney) configureDiscoveryWithClock(
	advisor experience.DiscoveryAdvisor,
	clock func() time.Time,
	limits discovery.Limits,
) error {
	service, err := experience.NewDiscoveryService(
		advisor, clock, journey.newIdentifier, limits,
	)
	if err != nil {
		return err
	}
	journey.discovery = service
	journey.discoveryLimits = limits
	return nil
}

func defaultDiscoveryLimits() discovery.Limits {
	return discovery.Limits{
		MaxInputBytes: 32 << 10, MaxOutputBytes: 32 << 10, MaxTokens: 4096,
		MaxLatency: 20 * time.Second, MaxCostMicros: 250_000,
	}
}

func (journey *developerTrialJourney) Discover(ctx context.Context, input discoveryInput) (discoveryStage, error) {
	scope, err := journey.directory.ResolveWorkspaceScope(journey.human, tenant.ScopeRequest{
		ContractVersion: tenant.ScopeRequestContractVersion, CorrelationID: trialCorrelation,
		WorkspaceSelector: input.Workspace, ProjectSelector: input.Project, RepositorySelector: input.Repository,
	})
	if err != nil {
		return discoveryStage{}, err
	}
	binding, err := journey.directory.GetRepositoryBinding(ctx, tenant.BindingRequest{
		ContractVersion: tenant.BindingRequestContractVersion, CorrelationID: trialCorrelation, Scope: scope,
	})
	if err != nil {
		return discoveryStage{}, err
	}
	session, err := experience.NewSession(journey.human, scope, binding, journey.now(), trialCorrelation)
	if err != nil {
		return discoveryStage{}, err
	}
	if journey.discovery == nil {
		return discoveryStage{}, failure.New(
			failure.DependencyUnavailable, "model-assisted discovery is not activated",
			trialCorrelation, failure.RetryNewAuthority,
		)
	}
	proposal, err := journey.discovery.Turn(ctx, session, journey.taskAgent, experience.DiscoveryTurn{
		SessionID: input.SessionID, UserStatements: input.UserStatements,
		CurrentDraft: input.CurrentDraft, PreviousQuestions: input.PreviousQuestions,
		ProvenanceRefs: input.ProvenanceRefs,
	})
	if err != nil {
		return discoveryStage{}, err
	}
	journey.scope, journey.binding = scope, binding
	return discoveryStage{
		SessionID: proposal.SessionID, Response: proposal.Response, Draft: proposal.Draft,
		Questions: proposal.Questions, Uncertainty: proposal.Uncertainty,
		Limitations: proposal.Limitations, Refusal: proposal.Refusal, Incomplete: proposal.Incomplete,
		Usage: proposal.Usage, Provider: proposal.Receipt.Provider, ModelID: proposal.Receipt.ModelID,
		ActivationID:     proposal.Receipt.ActivationID,
		ActivationDigest: proposal.Receipt.ActivationDigest,
		RoleProfileID:    proposal.Receipt.RoleProfileID, PromptVersion: proposal.Receipt.PromptVersion,
		ResponseSchema: proposal.Receipt.ResponseSchema,
	}, nil
}

func (journey *developerTrialJourney) Start(ctx context.Context, input startInput) (proposalStage, error) {
	scope, err := journey.directory.ResolveWorkspaceScope(journey.human, tenant.ScopeRequest{
		ContractVersion: tenant.ScopeRequestContractVersion, CorrelationID: trialCorrelation,
		WorkspaceSelector: input.Workspace, ProjectSelector: input.Project, RepositorySelector: input.Repository,
	})
	if err != nil {
		return proposalStage{}, err
	}
	binding, err := journey.directory.GetRepositoryBinding(ctx, tenant.BindingRequest{
		ContractVersion: tenant.BindingRequestContractVersion, CorrelationID: trialCorrelation, Scope: scope,
	})
	if err != nil {
		return proposalStage{}, err
	}
	session, err := experience.NewSession(journey.human, scope, binding, journey.now(), trialCorrelation)
	if err != nil {
		return proposalStage{}, err
	}
	rawInput := input.RawInput
	if rawInput == "" {
		raw, _ := json.Marshal(input.Intent)
		rawInput = string(raw)
	}
	provenance := input.Provenance
	if len(provenance) == 0 {
		provenance = []string{"cli:intent-file"}
	}
	proposal, err := journey.experience.Create(ctx, session, experience.CreateProposal{
		RawInput: rawInput, Outcome: input.Intent.Outcome, AffectedUsers: input.Intent.AffectedUsers,
		Problem: input.Intent.Problem, FunctionalBehavior: input.Intent.FunctionalBehavior,
		UXExpectations: input.Intent.UXExpectations, Accessibility: input.Intent.Accessibility,
		Privacy: input.Intent.Privacy, Security: input.Intent.Security, Operations: input.Intent.Operations,
		Constraints: input.Intent.Constraints, Exclusions: input.Intent.Exclusions,
		AcceptanceCriteria: input.Intent.AcceptanceCriteria, Inferences: input.Intent.Inferences,
		Assumptions: input.Intent.Assumptions, Hypotheses: input.Intent.Hypotheses,
		Questions: input.Intent.Questions, Experiments: input.Intent.Experiments,
		ProductAreaRefs: input.Intent.ProductAreaRefs, Provenance: provenance,
		ChangeExplanation: "Initial Developer Trial proposal.",
	})
	if err != nil {
		return proposalStage{}, err
	}
	journey.scope, journey.binding = scope, binding
	journey.sessions[proposal.ID] = &trialSession{proposal: proposal}
	return proposalToStage(proposal), nil
}

func (journey *developerTrialJourney) ReviseProposal(ctx context.Context, sessionID, expectedDigest string, edit proposalEditInput) (proposalStage, error) {
	current, ok := journey.sessions[sessionID]
	if !ok || current.proposal.Digest != expectedDigest {
		return proposalStage{}, errors.New("current requirements proposal was not found")
	}
	session, err := experience.NewSession(journey.human, journey.scope, journey.binding, journey.now(), trialCorrelation)
	if err != nil {
		return proposalStage{}, err
	}
	revised, err := journey.experience.Revise(ctx, session, current.proposal.ID, expectedDigest, experience.ProposalEdit{
		Field: edit.Field, Item: edit.Item, Replacement: edit.Replacement,
		ChangeExplanation: edit.ChangeExplanation,
	})
	if err != nil {
		return proposalStage{}, err
	}
	current.proposal = revised
	return proposalToStage(revised), nil
}

func proposalToStage(proposal experience.RequirementsProposal) proposalStage {
	return proposalStage{
		SessionID: proposal.ID, ProposalID: proposal.ID, ProposalVersion: proposal.Version,
		ProposalDigest: proposal.Digest, PreviousDigest: proposal.PreviousDigest,
		ChangeExplanation: proposal.ChangeExplanation,
		Outcome:           proposal.Outcome, AffectedUsers: proposal.AffectedUsers, Problem: proposal.Problem,
		FunctionalBehavior: proposal.FunctionalBehavior, UXExpectations: proposal.UXExpectations,
		Accessibility: proposal.Accessibility, Privacy: proposal.Privacy, Security: proposal.Security,
		Operations: proposal.Operations, Constraints: proposal.Constraints, Exclusions: proposal.Exclusions,
		AcceptanceCriteria: proposal.AcceptanceCriteria, Inferences: proposal.Inferences,
		Assumptions: proposal.Assumptions, Hypotheses: proposal.Hypotheses, Questions: proposal.Questions,
		Experiments: proposal.Experiments, ProductAreaRefs: proposal.ProductAreaRefs,
	}
}

func (journey *developerTrialJourney) PromoteAndPlan(ctx context.Context, sessionID, proposalDigest string) (planStage, error) {
	session, ok := journey.sessions[sessionID]
	if !ok || session.proposal.Digest != proposalDigest {
		return planStage{}, errors.New("current requirements proposal was not found")
	}
	promoted, err := journey.intent.Promote(ctx, intent.PromoteRequest{
		ContractVersion: intent.PromotionRequestVersion, CorrelationID: trialCorrelation,
		IdempotencyKey: "promote:" + proposalDigest, ProposalID: session.proposal.ID,
		ProposalDigest: proposalDigest, Principal: journey.human, Scope: journey.scope,
		RepositoryBinding: journey.binding,
	})
	if err != nil {
		return planStage{}, err
	}
	run := promoted.Run
	delegationRequest, err := identity.NewDelegationRequest(identity.DelegationRequestInput{
		ContractVersion: identity.DelegationRequestContractVersion, IdempotencyKey: "context-delegation:" + run.ID(),
		CorrelationID: trialCorrelation, HumanSubjectID: journey.human.SubjectID(),
		AgentSubjectID: journey.taskAgent.SubjectID(), TenantID: run.TenantID(),
		WorkspaceID: run.WorkspaceID(), ProjectID: run.ProjectID(), RepositoryID: run.RepositoryID(),
		DeliveryRunID: run.ID(), Purpose: identity.M1ChangePurpose,
		Actions: []string{identity.DelegatedActionAssembleContext}, Resources: []string{run.RepositoryID()},
		Classifications: []string{run.Classification()}, Destinations: []string{},
		Environment: identity.DeveloperTrialEnvironment, MaximumAutonomy: identity.ContextPreparationOnly,
		RequestedExpiresAt: journey.now().Add(20 * time.Minute),
	})
	if err != nil {
		return planStage{}, err
	}
	evaluator, err := policy.NewEvaluator(policy.PolicySet{
		Version: "policy.m1-c.context.v1",
		Rules: []policy.Rule{
			trialRule("permit-context-delegation", policy.ActionCreateDelegation, journey.human.SubjectID(), delegationRequest.Digest(), run, journey.scope, nil),
			trialRule("permit-context", policy.ActionAssembleContext, journey.human.SubjectID(), run.RepositoryID(), run, journey.scope,
				[]policy.Obligation{{Type: "context-evidence", Parameters: map[string]any{"required": true}}}),
		},
	}, []string{"context-evidence"}, 10*time.Minute, journey.now)
	if err != nil {
		return planStage{}, err
	}
	journey.policies.add(evaluator)
	delegationDecision, err := evaluator.EvaluateDelegation(ctx, policy.EvaluateDelegationRequest{
		ContractVersion: policy.EvaluateDelegationRequestContractVersion, CorrelationID: trialCorrelation,
		Human: journey.human, Agent: journey.taskAgent, Scope: journey.scope, Run: run,
		Delegation: delegationRequest, Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	})
	if err != nil {
		return planStage{}, err
	}
	delegations, err := identity.NewDelegationService(journey.persistence.delegation, 20*time.Minute, journey.now)
	if err != nil {
		return planStage{}, err
	}
	delegationResult, err := delegations.CreateDelegation(ctx, identity.CreateDelegationRequest{
		ContractVersion: identity.CreateDelegationContractVersion, Request: delegationRequest,
		Human: journey.human, Agent: journey.taskAgent, ParentDecision: delegationDecision,
		EnforcedObligationTypes: []string{},
	})
	if err != nil {
		return planStage{}, err
	}
	contextDecision, err := evaluator.Evaluate(ctx, policy.EvaluateRequest{
		ContractVersion: policy.EvaluateRequestContractVersion, CorrelationID: trialCorrelation,
		Principal: journey.human, Scope: journey.scope, Run: run, Action: policy.ActionAssembleContext,
		Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	})
	if err != nil {
		return planStage{}, err
	}
	envelopeIssuer, err := policy.NewEnvelopeIssuer(
		evaluator, delegations, journey.persistence.envelope,
		policy.EnvelopeProfile{
			Issuer: "forge.policy.developer-trial", Environment: policy.DeveloperTrialEnvironment,
			Risk: policy.RiskLow, AllowedPolicyVersions: []string{"policy.m1-c.context.v1"},
			EnforcedObligationTypes: []string{"context-evidence"}, MaximumLifetime: 10 * time.Minute,
			Limits: trust.ExecutionLimits{MaxDurationSeconds: 600, MaxConcurrency: 1},
		}, journey.now,
	)
	if err != nil {
		return planStage{}, err
	}
	envelopeResult, err := envelopeIssuer.IssueExecutionEnvelope(ctx, policy.IssueExecutionEnvelopeRequest{
		ContractVersion: policy.IssueExecutionEnvelopeContractVersion, CorrelationID: trialCorrelation,
		IdempotencyKey: "context-envelope:" + run.ID(), Human: journey.human, Agent: journey.taskAgent,
		Scope: journey.scope, Run: run, DelegationID: delegationResult.Delegation.ID(),
		PolicyDecisionID: contextDecision.ID,
	})
	if err != nil {
		return planStage{}, err
	}
	contextService, err := forgecontext.NewService(
		delegations, evaluator, journey.catalog, journey.persistence.context,
		forgecontext.AssemblyProfile{
			MaximumItems: 10, MaximumSerializedBytes: 1 << 20, MaximumLifetime: 10 * time.Minute,
			EnforcedObligationTypes: []string{"context-evidence"},
		}, journey.now,
	)
	if err != nil {
		return planStage{}, err
	}
	assembled, err := contextService.AssembleContext(ctx, forgecontext.AssembleRequest{
		ContractVersion: forgecontext.AssembleRequestContractVersion, CorrelationID: trialCorrelation,
		IdempotencyKey: "context:" + run.ID(), Human: journey.human, Agent: journey.taskAgent,
		Scope: journey.scope, Run: run, DelegationID: delegationResult.Delegation.ID(),
		Envelope: envelopeResult.Envelope, PolicyDecisionID: contextDecision.ID,
		RequestedSourceIDs: []string{"accepted-architecture", "planned-target"},
	})
	if err != nil {
		return planStage{}, err
	}
	authorizing, _, err := journey.delivery.BeginContextAuthorization(ctx, trialTransition(run, "context-authorizing"))
	if err != nil {
		return planStage{}, err
	}
	ready, _, err := journey.delivery.RecordContext(ctx, delivery.RecordContextRequest{
		TransitionRequest: trialTransition(authorizing, "context-ready"),
		Package:           assembled.Package, Envelope: envelopeResult.Envelope,
	})
	if err != nil {
		return planStage{}, err
	}
	profiles, err := delivery.NewStaticPlanProfiles([]delivery.PlanProfile{trialPlanProfile(journey.operational)})
	if err != nil {
		return planStage{}, err
	}
	planned, _, err := journey.delivery.RecordPlan(ctx, profiles, delivery.RecordPlanRequest{
		TransitionRequest: trialTransition(ready, "planned"), ProfileVersion: "forge.m1-c-plan-profile.v1",
	})
	if err != nil {
		return planStage{}, err
	}
	approvalRequest, _, err := journey.approvals.RequestPlanReview(ctx, approval.RequestPlanReview{
		ContractVersion: approval.RequestPlanReviewContractVersion, CorrelationID: trialCorrelation,
		IdempotencyKey: "plan-review:" + run.ID(), TenantID: run.TenantID(), RunID: run.ID(),
		RequestingActorID: run.ActorID(), IntendedApproverID: run.ActorID(),
		RequestedExpiry: journey.now().Add(20 * time.Minute),
	})
	if err != nil {
		return planStage{}, err
	}
	awaiting, _, err := journey.delivery.AwaitPlanApproval(ctx, journey.approvals, delivery.AwaitPlanApprovalRequest{
		TransitionRequest: trialTransition(planned, "await-plan-approval"), ApprovalRequestID: approvalRequest.ID,
	})
	if err != nil {
		return planStage{}, err
	}
	session.run, session.planRequest, session.contextPackage = awaiting, approvalRequest, assembled.Package
	plan := awaiting.Plan()
	return planStage{
		SessionID: sessionID, RunID: awaiting.ID(), PlanID: plan.ID,
		PlanDigest: plan.Digest, RiskClass: plan.RiskClass, Checks: plan.VerificationIntent,
	}, nil
}

func (journey *developerTrialJourney) ApprovePlanAndVerify(ctx context.Context, sessionID, planDigest string) (candidateStage, error) {
	session, ok := journey.sessions[sessionID]
	if !ok || session.run.Plan() == nil || session.run.Plan().Digest != planDigest {
		return candidateStage{}, errors.New("current delivery plan was not found")
	}
	approvalEvaluator, err := trialApprovalEvaluator(session.run, journey.scope, journey.human.SubjectID(), session.planRequest.RequestDigest, journey.now)
	if err != nil {
		return candidateStage{}, err
	}
	decision, err := approvalEvaluator.EvaluateApproval(ctx, policy.EvaluateApprovalRequest{
		ContractVersion: policy.EvaluateApprovalRequestContractVersion, CorrelationID: trialCorrelation,
		Principal: journey.human, Scope: journey.scope, Run: session.run, Request: session.planRequest,
		Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	})
	if err != nil {
		return candidateStage{}, err
	}
	journey.policies.add(approvalEvaluator)
	record, _, err := journey.approvals.DecideApproval(ctx, approval.DecideApproval{
		ContractVersion: approval.DecideApprovalContractVersion, CorrelationID: trialCorrelation,
		IdempotencyKey: "approve-plan:" + planDigest, TenantID: session.run.TenantID(),
		RequestID: session.planRequest.ID, Approver: journey.human, PolicyDecisionID: decision.ID,
		Decision: "approved", Conditions: []string{"require-independent-verification"},
		EvidenceReferences: session.planRequest.EvidenceReferences,
	})
	if err != nil {
		return candidateStage{}, err
	}
	approved, _, err := journey.delivery.RecordPlanApproval(ctx, journey.approvals, delivery.RecordPlanApprovalRequest{
		TransitionRequest: trialTransition(session.run, "plan-approved"), ApprovalRecordID: record.ID,
	})
	if err != nil {
		return candidateStage{}, err
	}
	plan, contextRecord := approved.Plan(), approved.Context()
	changeLimits := changecontract.Limits{
		MaxInputBytes: 64 << 10, MaxOutputBytes: 1 << 20, MaxTokens: 8192,
		MaxLatency: 60 * time.Second, MaxCostMicros: 500_000,
		MaxChangedFiles: int64(len(plan.ChangePaths)),
	}
	if journey.changePreparation != nil {
		changeLimits = changecontract.Limits{
			MaxInputBytes:   journey.changeActivation.Budgets.MaxInputBytes,
			MaxOutputBytes:  journey.changeActivation.Budgets.MaxOutputBytes,
			MaxTokens:       journey.changeActivation.Budgets.MaxTokens,
			MaxLatency:      journey.changeActivation.Budgets.MaxLatency,
			MaxCostMicros:   journey.changeActivation.Budgets.MaxCostMicros,
			MaxChangedFiles: int64(len(plan.ChangePaths)),
		}
	}
	projected, err := journey.catalog.Project([]string{m1ArchitecturePath, m1TargetPath})
	if err != nil {
		return candidateStage{}, err
	}
	projectedSources := make([]changecontract.SourceFile, len(projected))
	for index, source := range projected {
		projectedSources[index] = changecontract.SourceFile{
			Path: source.Path, Digest: source.Digest, Content: source.Content,
		}
	}
	changeRequest, err := changecontract.BindRequest(changecontract.Request{
		SchemaVersion: changecontract.RequestVersion, CorrelationID: trialCorrelation,
		TenantID: approved.TenantID(), WorkspaceID: approved.WorkspaceID(),
		ProjectID: approved.ProjectID(), RepositoryID: approved.RepositoryID(),
		DeliveryRunID: approved.ID(), TaskAgentID: journey.taskAgent.SubjectID(),
		PlanID: plan.ID, PlanDigest: plan.Digest, ContextPackageID: contextRecord.PackageID,
		ContextManifestDigest: contextRecord.ManifestDigest,
		RequirementsRecordID:  session.proposal.ID, SourceRevision: contextRecord.SourceRevisions[0],
		RoleProfileID: changecontract.RoleProfileV1, Purpose: changecontract.PurposeV1,
		Classification:    approved.Classification(),
		AllowedReadPaths:  []string{m1ArchitecturePath, m1TargetPath},
		AllowedWritePaths: plan.ChangePaths, Sources: projectedSources,
		AcceptanceCriteria: session.proposal.AcceptanceCriteria,
		TestIntent:         plan.VerificationIntent, Exclusions: session.proposal.Exclusions,
		ControlReferences: plan.GuardrailReferences,
		Limits:            changeLimits,
		ExpiresAt:         journey.now().Add(15 * time.Minute),
	})
	if err != nil {
		return candidateStage{}, err
	}
	var changeProposal changecontract.Proposal
	if journey.changePreparation != nil {
		changeDelegationRequest, delegationErr := identity.NewDelegationRequest(identity.DelegationRequestInput{
			ContractVersion: identity.DelegationRequestContractVersion,
			IdempotencyKey:  "change-delegation:" + approved.ID(), CorrelationID: trialCorrelation,
			HumanSubjectID: journey.human.SubjectID(), AgentSubjectID: journey.taskAgent.SubjectID(),
			TenantID: approved.TenantID(), WorkspaceID: approved.WorkspaceID(),
			ProjectID: approved.ProjectID(), RepositoryID: approved.RepositoryID(),
			DeliveryRunID: approved.ID(), Purpose: identity.M1ChangePurpose,
			Actions:            []string{identity.DelegatedActionProposeChange},
			Resources:          []string{changeRequest.Digest},
			Classifications:    []string{approved.Classification()},
			Destinations:       []string{"model:" + journey.changeActivation.Provider},
			Environment:        identity.DeveloperTrialEnvironment,
			MaximumAutonomy:    identity.ChangeProposalOnly,
			RequestedExpiresAt: journey.now().Add(15 * time.Minute),
		})
		if delegationErr != nil {
			return candidateStage{}, delegationErr
		}
		changeEvaluator, evaluatorErr := policy.NewEvaluator(policy.PolicySet{
			Version: "policy.m1-exit.change-proposal.v1",
			Rules: []policy.Rule{
				trialRule(
					"permit-change-delegation", policy.ActionCreateDelegation, journey.human.SubjectID(),
					changeDelegationRequest.Digest(), approved, journey.scope, nil,
				),
				trialRule(
					"permit-change-proposal", policy.ActionProposeChange, journey.taskAgent.SubjectID(),
					changeRequest.Digest, approved, journey.scope,
					[]policy.Obligation{
						{Type: "model-output-non-authoritative", Parameters: map[string]any{"required": true}},
						{Type: "purpose-minimized-context", Parameters: map[string]any{"required": true}},
					},
				),
			},
		}, []string{"model-output-non-authoritative", "purpose-minimized-context"}, 10*time.Minute, journey.now)
		if evaluatorErr != nil {
			return candidateStage{}, evaluatorErr
		}
		journey.policies.add(changeEvaluator)
		changeDelegationDecision, delegationErr := changeEvaluator.EvaluateDelegation(ctx, policy.EvaluateDelegationRequest{
			ContractVersion: policy.EvaluateDelegationRequestContractVersion,
			CorrelationID:   trialCorrelation, Human: journey.human, Agent: journey.taskAgent,
			Scope: journey.scope, Run: approved, Delegation: changeDelegationRequest,
			Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
		})
		if delegationErr != nil {
			return candidateStage{}, delegationErr
		}
		changeDelegations, delegationErr := identity.NewDelegationService(
			journey.persistence.delegation, 15*time.Minute, journey.now,
		)
		if delegationErr != nil {
			return candidateStage{}, delegationErr
		}
		changeDelegation, delegationErr := changeDelegations.CreateDelegation(ctx, identity.CreateDelegationRequest{
			ContractVersion: identity.CreateDelegationContractVersion, Request: changeDelegationRequest,
			Human: journey.human, Agent: journey.taskAgent, ParentDecision: changeDelegationDecision,
		})
		if delegationErr != nil {
			return candidateStage{}, delegationErr
		}
		issuer, issuerErr := policy.NewChangeAuthorityIssuer(
			changeEvaluator, changeDelegations, journey.persistence.envelope,
			policy.ChangeAuthorityProfile{
				Issuer: "forge.policy.developer-trial", Provider: journey.changeActivation.Provider,
				ModelID:                 journey.changeActivation.AllowedModelIDs[0],
				CredentialRef:           journey.changeActivation.CredentialRef,
				MaximumLifetime:         10 * time.Minute,
				EnforcedObligationTypes: []string{"model-output-non-authoritative", "purpose-minimized-context"},
			}, journey.now,
		)
		if issuerErr != nil {
			return candidateStage{}, issuerErr
		}
		authority, authorityErr := issuer.Issue(ctx, policy.IssueChangeAuthorityRequest{
			SchemaVersion: policy.IssueChangeAuthorityVersion, CorrelationID: trialCorrelation,
			IdempotencyKey: "change-proposal:" + approved.ID(), Human: journey.human,
			TaskAgent: journey.taskAgent, Scope: journey.scope, Run: approved,
			Request: changeRequest, DelegationID: changeDelegation.Delegation.ID(),
			Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
		})
		if authorityErr != nil {
			return candidateStage{}, authorityErr
		}
		changeProposal, err = journey.changePreparation.Prepare(ctx, intelligence.ChangePreparationCommand{
			Request: changeRequest, Decision: authority.Decision, Envelope: authority.Envelope,
		})
		if err != nil {
			return candidateStage{}, err
		}
	} else {
		changeProposal, err = deterministicTrialChangeProposal(changeRequest, session.proposal)
		if err != nil {
			return candidateStage{}, err
		}
	}
	changeProposalReference := ""
	if journey.persistence.artifacts != nil {
		encoded, encodeErr := json.Marshal(changeProposal)
		if encodeErr != nil {
			return candidateStage{}, encodeErr
		}
		metadata, artifactErr := journey.persistence.artifacts.Put(ctx, evidence.PutArtifactRequest{
			TenantID: approved.TenantID(), ProjectID: approved.ProjectID(), DeliveryRunID: approved.ID(),
			MediaType:      "application/vnd.forge.change-proposal+json",
			Classification: approved.Classification(), RetentionClass: "developer-trial-m1",
			CreatorCapability: "forge.capability.intelligence", Bytes: encoded,
		})
		if artifactErr != nil {
			return candidateStage{}, artifactErr
		}
		changeProposalReference = "artifact:" + metadata.ObjectKey + "#" + metadata.Digest
	}
	workerRequest, err := identity.BindWorkerDelegationRequest(identity.WorkerDelegationRequest{
		SchemaVersion: identity.WorkerDelegationRequestVersion, IdempotencyKey: "worker-delegation:" + approved.ID(),
		CorrelationID: trialCorrelation, HumanSubjectID: journey.human.SubjectID(),
		TaskAgentSubjectID: journey.taskAgent.SubjectID(), WorkerSubjectID: journey.worker.SubjectID(),
		TenantID: approved.TenantID(), WorkspaceID: approved.WorkspaceID(), ProjectID: approved.ProjectID(),
		RepositoryID: approved.RepositoryID(), DeliveryRunID: approved.ID(), PlanID: plan.ID,
		PlanDigest: plan.Digest, ChangeProposalDigest: changeProposal.Digest,
		ContextPackageID:      contextRecord.PackageID,
		ContextManifestDigest: contextRecord.ManifestDigest, ApprovalRequestID: session.planRequest.ID,
		ApprovalRequestDigest: session.planRequest.RequestDigest, ApprovalRecordID: record.ID,
		ApprovalRecordDigest: record.Integrity.Digest, SourceRevision: contextRecord.SourceRevisions[0],
		ChangePaths: plan.ChangePaths, Classification: approved.Classification(), RiskClass: plan.RiskClass,
		Environment: identity.DeveloperTrialEnvironment, Action: identity.WorkerPrepareChangeAction,
		Purpose: identity.WorkerM1Purpose, MaximumAutonomy: identity.WorkerMaximumAutonomy,
		MaxDurationSeconds: 60, MaxProcesses: 8, MaxOutputBytes: 1 << 20,
		MaxChangedFiles: int64(len(plan.ChangePaths)), MaxAttempts: plan.Budgets.MaxAttempts,
		RequestedExpiresAt: journey.now().Add(15 * time.Minute),
	})
	if err != nil {
		return candidateStage{}, err
	}
	workerEvaluator, err := policy.NewEvaluator(policy.PolicySet{
		Version: "policy.m1-c.worker.v1",
		Rules: []policy.Rule{
			trialRule("permit-worker-delegation", policy.ActionCreateDelegation, journey.human.SubjectID(), workerRequest.Digest, approved, journey.scope, nil),
			trialRule(
				"permit-worker-execution", policy.ActionPrepareChange, journey.worker.SubjectID(),
				changeProposal.Digest, approved, journey.scope,
				[]policy.Obligation{{Type: "proposal-admitted", Parameters: map[string]any{"required": true}}},
			),
		},
	}, []string{"proposal-admitted"}, 10*time.Minute, journey.now)
	if err != nil {
		return candidateStage{}, err
	}
	journey.policies.add(workerEvaluator)
	delegationDecision, err := workerEvaluator.EvaluateWorkerDelegation(ctx, journey.approvals, policy.EvaluateWorkerDelegationRequest{
		SchemaVersion: policy.EvaluateWorkerDelegationVersion, CorrelationID: trialCorrelation,
		Human: journey.human, TaskAgent: journey.taskAgent, Worker: journey.worker,
		Scope: journey.scope, Run: approved, Request: workerRequest,
		Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	})
	if err != nil {
		return candidateStage{}, err
	}
	workerDelegations, err := identity.NewWorkerDelegationService(journey.persistence.workerDelegation, 15*time.Minute, journey.now)
	if err != nil {
		return candidateStage{}, err
	}
	workerDelegation, _, err := workerDelegations.Create(ctx, identity.CreateWorkerDelegationRequest{
		SchemaVersion: identity.CreateWorkerDelegationVersion, Request: workerRequest,
		Human: journey.human, TaskAgent: journey.taskAgent, Worker: journey.worker,
		ParentDecision: delegationDecision, EnforcedObligationTypes: []string{},
	})
	if err != nil {
		return candidateStage{}, err
	}
	workerIssuer, err := policy.NewWorkerEnvelopeIssuer(
		workerEvaluator, journey.approvals, workerDelegations, journey.persistence.envelope,
		policy.WorkerEnvelopeProfile{
			Issuer: "forge.policy.developer-trial", MaximumLifetime: 10 * time.Minute,
			EnforcedObligationTypes: []string{"proposal-admitted"},
		},
		journey.now,
	)
	if err != nil {
		return candidateStage{}, err
	}
	workerEnvelope, _, err := workerIssuer.Issue(ctx, policy.IssueWorkerEnvelopeRequest{
		SchemaVersion: policy.IssueWorkerEnvelopeVersion, CorrelationID: trialCorrelation,
		IdempotencyKey: "worker-envelope:" + approved.ID(), Human: journey.human,
		TaskAgent: journey.taskAgent, Worker: journey.worker, Scope: journey.scope,
		Run: approved, DelegationID: workerDelegation.ID,
		Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	})
	if err != nil {
		return candidateStage{}, err
	}
	var materializeDecision trust.PolicyDecision
	found := false
	for _, decisionID := range workerEnvelope.PolicyDecisionIDs {
		candidate, resolved, resolveErr := workerEvaluator.ResolvePolicyDecision(ctx, decisionID)
		if resolveErr != nil {
			return candidateStage{}, resolveErr
		}
		if resolved && candidate.Action == workercontract.PrepareChangeAction &&
			candidate.Resource == changeProposal.Digest {
			materializeDecision, found = candidate, true
			break
		}
	}
	if !found {
		return candidateStage{}, errors.New("materialization policy decision was not found")
	}
	authorizing, _, err := journey.delivery.BeginExecutionAuthorization(ctx, trialTransition(approved, "execution-authorizing"))
	if err != nil {
		return candidateStage{}, err
	}
	executing, _, err := journey.delivery.RecordExecutionAuthority(ctx, delivery.RecordExecutionAuthorityRequest{
		TransitionRequest: trialTransition(authorizing, "executing"), Envelope: workerEnvelope,
		DelegationID: workerDelegation.ID, WorkerID: journey.worker.SubjectID(),
	})
	if err != nil {
		return candidateStage{}, err
	}
	executionService, err := execution.NewService(journey.workerRunner, journey.now)
	if err != nil {
		return candidateStage{}, err
	}
	workspace, err := execution.PrepareSourceWorkspace(
		journey.runtimeRoot, workerSources(changeRequest.Sources, plan.ChangePaths),
	)
	if err != nil {
		return candidateStage{}, err
	}
	defer os.RemoveAll(workspace)
	prepared, err := executionService.PrepareProposal(ctx, execution.PrepareProposalRequest{
		CorrelationID: trialCorrelation, IdempotencyKey: "prepare:" + approved.ID(),
		Run: executing, Envelope: workerEnvelope, MaterializeDecision: materializeDecision,
		ChangeRequest: changeRequest, Proposal: changeProposal,
		Workspace: workspace, Attempt: 1, MaxProcesses: 8, MaxOutputBytes: 1 << 20,
	})
	if err != nil {
		return candidateStage{}, err
	}
	artefactReferences := []string{"candidate:" + prepared.Bundle.Digest}
	if changeProposalReference != "" {
		artefactReferences = append(artefactReferences, changeProposalReference)
	}
	if journey.persistence.artifacts != nil {
		encoded, encodeErr := workercontract.EncodeCanonical(prepared.Bundle)
		if encodeErr != nil {
			return candidateStage{}, encodeErr
		}
		metadata, artifactErr := journey.persistence.artifacts.Put(ctx, evidence.PutArtifactRequest{
			TenantID: approved.TenantID(), ProjectID: approved.ProjectID(), DeliveryRunID: approved.ID(),
			MediaType:      "application/vnd.forge.candidate-bundle+json",
			Classification: approved.Classification(), RetentionClass: "developer-trial-m1",
			CreatorCapability: "forge.capability.execution", Bytes: encoded,
		})
		if artifactErr != nil {
			return candidateStage{}, artifactErr
		}
		artefactReferences = append(artefactReferences, "artifact:"+metadata.ObjectKey+"#"+metadata.Digest)
	}
	verifying, _, err := journey.delivery.RecordCandidate(ctx, delivery.RecordCandidateRequest{
		TransitionRequest: trialTransition(executing, "verifying"), WorkerRequest: prepared.Request,
		Bundle: prepared.Bundle, Result: prepared.Result,
		ArtefactReferences: artefactReferences,
	})
	if err != nil {
		return candidateStage{}, err
	}
	verificationService, err := verification.NewService(
		journey.verifierRunner, journey.persistence.verification,
		verification.Profile{
			Version: "forge.m1-c.verification.v1", MaximumLifetime: 10 * time.Minute,
			// The first two controls are integrity preconditions that cannot fail
			// once the lifecycle has admitted the candidate. The remaining three
			// are the controls capable of failing; the verifier refuses any plan
			// whose independent controls are all inapplicable to the changed
			// paths.
			Controls: []verificationcontract.Control{
				{ID: verificationcontract.ControlCandidateIntegrity, Version: "1", Consequential: true},
				{ID: verificationcontract.ControlDeclaredPathScope, Version: "1", Consequential: true},
				{ID: verificationcontract.ControlGoFormat, Version: "1", Consequential: true},
				{ID: verificationcontract.ControlDocumentMetadata, Version: "1", Consequential: true},
				{ID: verificationcontract.ControlMarkdownLinks, Version: "1", Consequential: true},
			},
		}, journey.now,
	)
	if err != nil {
		return candidateStage{}, err
	}
	verificationSources := make([]verification.SourceFile, 0, len(plan.ChangePaths))
	for _, source := range workerSources(changeRequest.Sources, plan.ChangePaths) {
		verificationSources = append(verificationSources, verification.SourceFile{
			Path: source.Path, Digest: source.Digest, Content: source.Content,
		})
	}
	fresh, err := verification.PrepareFreshCheckout(journey.runtimeRoot, verificationSources)
	if err != nil {
		return candidateStage{}, err
	}
	defer os.RemoveAll(fresh)
	verified, err := verificationService.Verify(ctx, verification.VerifyRequest{
		CorrelationID: trialCorrelation, Run: verifying, Bundle: prepared.Bundle, FreshCheckout: fresh,
	})
	if err != nil {
		return candidateStage{}, err
	}
	awaitingPublication, _, err := journey.delivery.RecordVerification(ctx, delivery.RecordVerificationRequest{
		TransitionRequest: trialTransition(verifying, "await-publication-approval"),
		AssurancePlan:     verified.Plan, Result: verified.Result,
		EvidenceReferences: []string{"verification:" + verified.Result.Digest},
	})
	if err != nil {
		return candidateStage{}, err
	}
	_, _, err = journey.evidence.Append(ctx, evidence.AppendRequest{
		CorrelationID: trialCorrelation, IdempotencyKey: "verification-evidence:" + prepared.Bundle.Digest,
		Type: "verification", ActorID: "forge.verification",
		Scope: trust.EvidenceScope{
			TenantID: awaitingPublication.TenantID(), WorkspaceID: awaitingPublication.WorkspaceID(),
			ProjectID: awaitingPublication.ProjectID(), DeliveryRunID: awaitingPublication.ID(),
		},
		Subject: trust.EvidenceSubject{Type: "verified-candidate", ID: awaitingPublication.ID() + ".candidate", Digest: prepared.Bundle.Digest},
		References: []string{
			"candidate:" + prepared.Bundle.Digest,
			"verification:" + verified.Result.Digest,
		}, Outcome: "passed",
		Classification: awaitingPublication.Classification(), RetentionClass: "developer-trial-m1",
	})
	if err != nil {
		return candidateStage{}, err
	}
	session.run, session.planRecord, session.bundle = awaitingPublication, record, prepared.Bundle
	return candidateStage{
		SessionID: sessionID, RunID: awaitingPublication.ID(), CandidateDigest: prepared.Bundle.Digest,
		ChangeProposalDigest: changeProposal.Digest, ChangeProvider: changeProposal.Receipt.Provider,
		ChangeModelID: changeProposal.Receipt.ModelID, ChangeUsage: changeProposal.Usage,
		Assumptions: changeProposal.Assumptions, Uncertainty: changeProposal.Uncertainty,
		Omissions:          changeProposal.Omissions,
		VerificationDigest: verified.Result.Digest, ChangedPaths: awaitingPublication.Candidate().ChangedPaths,
		EvidenceReady: true,
	}, nil
}

func (journey *developerTrialJourney) ApprovePublicationAndPublish(ctx context.Context, sessionID, candidateDigest string) (publicationStage, error) {
	session, ok := journey.sessions[sessionID]
	if !ok || session.run.Candidate() == nil || session.run.Candidate().Digest != candidateDigest {
		return publicationStage{}, errors.New("current verified candidate was not found")
	}
	request, _, err := journey.approvals.RequestPublicationReview(ctx, approval.RequestPublicationReview{
		ContractVersion: approval.RequestPublicationReviewContractVersion, CorrelationID: trialCorrelation,
		IdempotencyKey: "publication-review:" + candidateDigest, TenantID: session.run.TenantID(),
		RunID: session.run.ID(), RequestingActorID: session.run.ActorID(),
		IntendedApproverID: session.run.ActorID(), RequestedExpiry: journey.now().Add(15 * time.Minute),
	})
	if err != nil {
		return publicationStage{}, err
	}
	approvalEvaluator, err := trialApprovalEvaluator(session.run, journey.scope, journey.human.SubjectID(), request.RequestDigest, journey.now)
	if err != nil {
		return publicationStage{}, err
	}
	decision, err := approvalEvaluator.EvaluateApproval(ctx, policy.EvaluateApprovalRequest{
		ContractVersion: policy.EvaluateApprovalRequestContractVersion, CorrelationID: trialCorrelation,
		Principal: journey.human, Scope: journey.scope, Run: session.run, Request: request,
		Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	})
	if err != nil {
		return publicationStage{}, err
	}
	journey.policies.add(approvalEvaluator)
	record, _, err := journey.approvals.DecideApproval(ctx, approval.DecideApproval{
		ContractVersion: approval.DecideApprovalContractVersion, CorrelationID: trialCorrelation,
		IdempotencyKey: "approve-publication:" + candidateDigest, TenantID: session.run.TenantID(),
		RequestID: request.ID, Approver: journey.human, PolicyDecisionID: decision.ID,
		Decision: "approved", Conditions: []string{}, EvidenceReferences: request.EvidenceReferences,
	})
	if err != nil {
		return publicationStage{}, err
	}
	publicationEvaluator, err := policy.NewEvaluator(policy.PolicySet{
		Version: "policy.m1-c.publication.v1",
		Rules: []policy.Rule{trialRule(
			"permit-publication", policy.ActionPublishPullRequest, journey.worker.SubjectID(),
			candidateDigest, session.run, journey.scope,
			[]policy.Obligation{
				{Type: "evidence-ready", Parameters: map[string]any{"required": true}},
				{Type: "human-publication-approval", Parameters: map[string]any{"required": true}},
			},
		)},
	}, []string{"evidence-ready", "human-publication-approval"}, 10*time.Minute, journey.now)
	if err != nil {
		return publicationStage{}, err
	}
	journey.policies.add(publicationEvaluator)
	issuer, err := policy.NewPublicationEnvelopeIssuer(
		publicationEvaluator, journey.delivery, journey.approvals, journey.persistence.envelope,
		policy.PublicationEnvelopeProfile{
			Issuer: "forge.policy.developer-trial", MaximumLifetime: 10 * time.Minute,
			EnforcedObligationTypes: []string{"evidence-ready", "human-publication-approval"},
		}, journey.now,
	)
	if err != nil {
		return publicationStage{}, err
	}
	envelope, _, err := issuer.Issue(ctx, policy.IssuePublicationEnvelopeRequest{
		SchemaVersion: policy.IssuePublicationEnvelopeVersion, CorrelationID: trialCorrelation,
		IdempotencyKey: "publication:" + candidateDigest, Human: journey.human,
		Workload: journey.worker, Scope: journey.scope, RepositoryBinding: journey.binding,
		Run: session.run, ApprovalRecordID: record.ID, Provider: "contract-double",
		BaseBranch: "main", BaseRevision: "developer-trial-base",
		Environment: policy.DeveloperTrialEnvironment, Risk: policy.RiskLow,
	})
	if err != nil {
		return publicationStage{}, err
	}
	policyDecisionID := envelope.PolicyDecisionIDs[len(envelope.PolicyDecisionIDs)-1]
	receipt, _, err := journey.sourceControl.PublishPullRequest(ctx, sourcecontrol.PublishRequest{
		ContractVersion: sourcecontrol.PublishRequestContractVersion, CorrelationID: trialCorrelation,
		IdempotencyKey: envelope.IdempotencyKey, TenantID: session.run.TenantID(), RunID: session.run.ID(),
		Envelope: envelope, PolicyDecisionID: policyDecisionID, ApprovalRecordID: record.ID,
		RepositoryBinding: journey.binding, Provider: "contract-double", BaseBranch: "main",
		BaseRevision: "developer-trial-base", TreeDigest: session.run.Candidate().TreeDigest,
		ChangedPaths: session.run.Candidate().ChangedPaths, Title: "Forge M1 Developer Trial candidate",
		Body: "Verified synthetic candidate. Human review and merge remain required.",
	})
	if err != nil {
		return publicationStage{}, err
	}
	session.publicationRequest, session.publicationRecord = request, record
	session.publicationEnvelope, session.publicationDecisionID = envelope, policyDecisionID
	return publicationStage{
		SessionID: sessionID, RunID: session.run.ID(), OperationID: receipt.OperationID,
		Branch: receipt.Branch, PullRequestID: receipt.PullRequestID, PullRequestURL: receipt.PullRequestURL,
		Provider: receipt.Provider, MergePermitted: false,
	}, nil
}

func trialPrincipal(subject string, kind identity.PrincipalKind, now func() time.Time) (identity.Principal, error) {
	at := now()
	verifier, err := identity.NewVerifier(trialIdentityProvider{assertion: identity.Assertion{
		SubjectID: subject, Kind: kind, Issuer: "forge.synthetic", Audience: "forge-cli",
		Assurance: "synthetic-local", ProfileID: "developer-trial",
		IssuedAt: at.Add(-time.Minute), ExpiresAt: at.Add(time.Hour), Synthetic: true,
	}}, "forge-cli", now)
	if err != nil {
		return identity.Principal{}, err
	}
	return verifier.Verify(context.Background(), identity.VerifyRequest{
		ContractVersion: identity.VerifyRequestContractVersion, Credential: "opaque",
		CorrelationID: trialCorrelation,
	})
}

func trialRule(id string, action policy.Action, subject, resource string, run delivery.DeliveryRun, scope tenant.Scope, obligations []policy.Obligation) policy.Rule {
	return policy.Rule{
		ID: id, Effect: policy.EffectPermit, Actions: []policy.Action{action},
		Environments: []policy.Environment{policy.DeveloperTrialEnvironment}, Risks: []policy.Risk{policy.RiskLow},
		SubjectIDs: []string{subject}, ResourceIDs: []string{resource}, DeliveryRunIDs: []string{run.ID()},
		Purposes: []string{run.Purpose()}, Classifications: []string{run.Classification()},
		TenantIDs: []string{scope.TenantID()}, WorkspaceIDs: []string{scope.WorkspaceID()},
		ProjectIDs: []string{scope.ProjectID()}, Obligations: obligations,
	}
}

func trialApprovalEvaluator(run delivery.DeliveryRun, scope tenant.Scope, humanID, requestDigest string, now func() time.Time) (*policy.Evaluator, error) {
	return policy.NewEvaluator(policy.PolicySet{
		Version: "policy.m1-c.approval." + requestDigest[7:19],
		Rules: []policy.Rule{trialRule(
			"permit-approval-"+requestDigest[7:19], policy.ActionDecideApproval, humanID,
			requestDigest, run, scope,
			[]policy.Obligation{{Type: "human-interaction", Parameters: map[string]any{"required": true}}},
		)},
	}, []string{"human-interaction"}, 10*time.Minute, now)
}

func trialTransition(run delivery.DeliveryRun, key string) delivery.TransitionRequest {
	return delivery.TransitionRequest{
		CorrelationID: trialCorrelation, IdempotencyKey: key + ":" + run.ID(),
		Run: run, ExpectedVersion: run.Version(),
	}
}

func trialPlanProfile(operational bool) delivery.PlanProfile {
	risks := []string{"Contract runners prove trust contracts but not process isolation."}
	assumptions := []string{"The contract double represents the publication provider."}
	if operational {
		risks = []string{"Developer Trial Docker isolation is not a production isolation claim."}
		assumptions = []string{"Worker and verifier images are bound by immutable image digests."}
	}
	return delivery.PlanProfile{
		Version: "forge.m1-c-plan-profile.v1",
		Steps: []delivery.PlanStep{
			{ID: "prepare", Dependencies: []string{}, OutputKinds: []string{"candidate"}},
			{ID: "verify", Dependencies: []string{"prepare"}, OutputKinds: []string{"verification-result"}},
			{ID: "publish", Dependencies: []string{"verify"}, OutputKinds: []string{"reviewable-pull-request"}},
		},
		AllowedCapabilities: []string{
			"forge.capability.delivery", "forge.capability.source-control", "forge.capability.verification",
		},
		AllowedPathPrefixes: []string{"docs"}, ChangePaths: []string{m1TargetPath},
		Assumptions:         assumptions,
		Risks:               risks,
		Limitations:         []string{"No live GitHub credential, merge, release, deployment, or production effect."},
		StopConditions:      []string{"Architecture change", "Changed digest", "Failed verification", "Unavailable evidence"},
		Budgets:             delivery.PlanBudgets{MaxDurationSeconds: 600, MaxConcurrency: 1, MaxAttempts: 2},
		VerificationIntent:  []string{"candidate integrity", "declared path scope"},
		GuardrailReferences: []string{"ADR-0016", "ADR-0017", "ADR-0019", "ADR-0020", "ADR-0021"},
		Rollback:            "Discard the synthetic contract-double operation.",
		ManualRecovery:      "Resume from the last capability-owned record with fresh authority.",
		RiskClass:           "low", MaximumLifetime: 30 * time.Minute,
	}
}

func deterministicTrialChangeProposal(request changecontract.Request, proposal experience.RequirementsProposal) (changecontract.Proposal, error) {
	var content strings.Builder
	content.WriteString("# Forge M1 Developer Trial candidate\n\n")
	content.WriteString("This candidate was prepared from the promoted Requirements Proposal.\n\n")
	content.WriteString("## Outcome\n\n")
	content.WriteString(proposal.Outcome)
	content.WriteString("\n\n## Problem\n\n")
	content.WriteString(proposal.Problem)
	content.WriteString("\n\n## Acceptance criteria\n")
	for _, criterion := range proposal.AcceptanceCriteria {
		content.WriteString("\n- ")
		content.WriteString(criterion)
	}
	content.WriteString("\n\n## Verification intent\n")
	for _, check := range request.TestIntent {
		content.WriteString("\n- ")
		content.WriteString(check)
	}
	content.WriteString("\n")
	replacements := make([]changecontract.FileReplacement, len(request.AllowedWritePaths))
	sourceDigests := make(map[string]string, len(request.Sources))
	for _, source := range request.Sources {
		sourceDigests[source.Path] = source.Digest
	}
	for index, target := range request.AllowedWritePaths {
		replacements[index] = changecontract.FileReplacement{
			Path: target, ExpectedSourceDigest: sourceDigests[target], Content: content.String(),
			Rationale: "Materializes the promoted requirements and approved verification intent for M1 review.",
		}
	}
	return changecontract.BindProposal(changecontract.Proposal{
		SchemaVersion: changecontract.ProposalVersion,
		RequestID:     request.ID, RequestDigest: request.Digest, SourceRevision: request.SourceRevision,
		PlanDigest: request.PlanDigest, ContextDigest: request.ContextManifestDigest,
		Replacements: replacements, ProposedChecks: request.TestIntent,
		Assumptions: []string{"Contract-double preparation proves the closed change contract without external effects."},
		Uncertainty: []string{}, Omissions: []string{"No live provider or GitHub publication is exercised."},
		Usage: changecontract.Usage{},
		Receipt: changecontract.Receipt{
			ActivationID:     "forge.contract-double.change-preparation.v1",
			ActivationDigest: "sha256:" + strings.Repeat("d", 64),
			Provider:         "contract-double", ModelID: "deterministic-template-v1",
			AdapterVersion: "forge.contract-double-change-adapter.v1",
			RoleProfileID:  request.RoleProfileID, PromptVersion: "forge.change-preparation.contract-double.v1",
			ResponseSchema: changecontract.ProposalVersion, ProviderRequest: "local:" + request.ID,
		},
	}, request)
}

var _ execution.Runner = trialWorkerRunner{}
var _ verification.Runner = trialVerifierRunner{}
var _ forgecontext.SourceCatalog = (*forgecontext.RepositoryCatalog)(nil)
var _ approval.PolicyResolver = (*policyResolvers)(nil)
var _ sourcecontrol.PolicyResolver = (*policyResolvers)(nil)
