package delivery

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"path/filepath"
	"slices"
	"sort"
	"strings"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
)

const DeliveryPlanContractVersion = "forge.delivery-plan.v1"

type PlanStep struct {
	ID           string   `json:"id"`
	Dependencies []string `json:"dependencies"`
	OutputKinds  []string `json:"outputKinds"`
}

type PlanBudgets struct {
	MaxDurationSeconds int64 `json:"maxDurationSeconds"`
	MaxConcurrency     int64 `json:"maxConcurrency"`
	MaxAttempts        int64 `json:"maxAttempts"`
	MaxComputeUnits    int64 `json:"maxComputeUnits"`
	MaxInputTokens     int64 `json:"maxInputTokens"`
	MaxOutputTokens    int64 `json:"maxOutputTokens"`
	MaxCostMinorUnits  int64 `json:"maxCostMinorUnits"`
}

type DeliveryPlan struct {
	ContractVersion      string      `json:"contractVersion"`
	ID                   string      `json:"id"`
	Digest               string      `json:"digest"`
	RunID                string      `json:"runId"`
	PlanVersion          uint64      `json:"planVersion"`
	TenantID             string      `json:"tenantId"`
	WorkspaceID          string      `json:"workspaceId"`
	ProjectID            string      `json:"projectId"`
	RepositoryID         string      `json:"repositoryId"`
	SourceRevisions      []string    `json:"sourceRevisions"`
	IntentVersion        uint64      `json:"intentVersion"`
	ContextPackageID     string      `json:"contextPackageId"`
	ContextManifest      string      `json:"contextManifestDigest"`
	Objective            string      `json:"objective"`
	Constraints          []string    `json:"constraints"`
	AcceptanceCriteria   []string    `json:"acceptanceCriteria"`
	Steps                []PlanStep  `json:"steps"`
	AffectedCapabilities []string    `json:"affectedCapabilities"`
	ChangePaths          []string    `json:"changePaths"`
	Assumptions          []string    `json:"assumptions"`
	Omissions            []string    `json:"omissions"`
	Conflicts            []string    `json:"conflicts"`
	Risks                []string    `json:"risks"`
	Limitations          []string    `json:"limitations"`
	StopConditions       []string    `json:"stopConditions"`
	ModelIDs             []string    `json:"modelIds"`
	ToolIDs              []string    `json:"toolIds"`
	SecretReferenceIDs   []string    `json:"secretReferenceIds"`
	Destinations         []string    `json:"destinations"`
	NetworkRequirements  []string    `json:"networkRequirements"`
	FilesystemAccess     []string    `json:"filesystemAccess"`
	Budgets              PlanBudgets `json:"budgets"`
	VerificationIntent   []string    `json:"verificationIntent"`
	GuardrailReferences  []string    `json:"guardrailReferences"`
	Rollback             string      `json:"rollback"`
	ManualRecovery       string      `json:"manualRecovery"`
	Classification       string      `json:"classification"`
	RiskClass            string      `json:"riskClass"`
	EvidenceCorrelation  string      `json:"evidenceCorrelationId"`
	CreatedAt            time.Time   `json:"createdAt"`
	ExpiresAt            time.Time   `json:"expiresAt"`
}

type PlanProfile struct {
	Version             string
	Steps               []PlanStep
	AllowedCapabilities []string
	AllowedPathPrefixes []string
	ChangePaths         []string
	Assumptions         []string
	Risks               []string
	Limitations         []string
	StopConditions      []string
	Budgets             PlanBudgets
	VerificationIntent  []string
	GuardrailReferences []string
	Rollback            string
	ManualRecovery      string
	RiskClass           string
	MaximumLifetime     time.Duration
}

type PlanProfileResolver interface {
	ResolvePlanProfile(context.Context, string) (PlanProfile, bool, error)
}

type StaticPlanProfiles struct {
	profiles map[string]PlanProfile
}

func NewStaticPlanProfiles(profiles []PlanProfile) (*StaticPlanProfiles, error) {
	if len(profiles) == 0 {
		return nil, errors.New("at least one trusted delivery plan profile is required")
	}
	result := &StaticPlanProfiles{profiles: make(map[string]PlanProfile, len(profiles))}
	for _, candidate := range profiles {
		profile, err := NewPlanProfile(candidate)
		if err != nil {
			return nil, err
		}
		if _, duplicate := result.profiles[profile.Version]; duplicate {
			return nil, errors.New("delivery plan profile versions must be unique")
		}
		result.profiles[profile.Version] = clonePlanProfile(profile)
	}
	return result, nil
}

func (profiles *StaticPlanProfiles) ResolvePlanProfile(ctx context.Context, version string) (PlanProfile, bool, error) {
	if profiles == nil || profiles.profiles == nil || ctx == nil || version == "" {
		return PlanProfile{}, false, errors.New("delivery plan profile registry is unavailable")
	}
	if err := ctx.Err(); err != nil {
		return PlanProfile{}, false, err
	}
	value, found := profiles.profiles[version]
	return clonePlanProfile(value), found, nil
}

func NewPlanProfile(profile PlanProfile) (PlanProfile, error) {
	if profile.Version == "" || profile.Rollback == "" || profile.ManualRecovery == "" ||
		profile.MaximumLifetime <= 0 || profile.MaximumLifetime > 2*time.Hour ||
		!oneOfRisk(profile.RiskClass) || len(profile.Steps) == 0 ||
		len(profile.AllowedCapabilities) == 0 || len(profile.AllowedPathPrefixes) == 0 ||
		len(profile.ChangePaths) == 0 || len(profile.VerificationIntent) == 0 ||
		len(profile.GuardrailReferences) == 0 || len(profile.StopConditions) == 0 ||
		!validBudgets(profile.Budgets) {
		return PlanProfile{}, errors.New("delivery plan profile is outside the M1 validation envelope")
	}
	for _, values := range [][]string{
		profile.AllowedCapabilities, profile.AllowedPathPrefixes, profile.ChangePaths,
		profile.Assumptions, profile.Risks, profile.Limitations, profile.StopConditions,
		profile.VerificationIntent, profile.GuardrailReferences,
	} {
		if !uniqueNonEmptySorted(values) {
			return PlanProfile{}, errors.New("delivery plan profile collections must be sorted, unique, and non-empty")
		}
	}
	if err := validateGraph(profile.Steps); err != nil {
		return PlanProfile{}, err
	}
	for _, path := range profile.ChangePaths {
		if !allowedPath(path, profile.AllowedPathPrefixes) {
			return PlanProfile{}, errors.New("delivery plan path exceeds the trusted profile")
		}
	}
	profile.Steps = cloneSteps(profile.Steps)
	profile.AllowedCapabilities = slices.Clone(profile.AllowedCapabilities)
	profile.AllowedPathPrefixes = slices.Clone(profile.AllowedPathPrefixes)
	profile.ChangePaths = slices.Clone(profile.ChangePaths)
	profile.Assumptions = slices.Clone(profile.Assumptions)
	profile.Risks = slices.Clone(profile.Risks)
	profile.Limitations = slices.Clone(profile.Limitations)
	profile.StopConditions = slices.Clone(profile.StopConditions)
	profile.VerificationIntent = slices.Clone(profile.VerificationIntent)
	profile.GuardrailReferences = slices.Clone(profile.GuardrailReferences)
	return profile, nil
}

func buildPlan(run DeliveryRun, profile PlanProfile, at time.Time, correlationID string) (DeliveryPlan, error) {
	contextRecord := run.Context()
	if run.State() != ContextReadyState || contextRecord == nil || at.IsZero() ||
		!at.Before(contextRecord.ExpiresAt) || run.CorrelationID() != correlationID {
		return DeliveryPlan{}, failure.New(failure.InvalidRequest, "delivery plan input is not current context-ready state", correlationID, failure.RetryNone)
	}
	expiresAt := at.Add(profile.MaximumLifetime)
	if contextRecord.ExpiresAt.Before(expiresAt) {
		expiresAt = contextRecord.ExpiresAt
	}
	value := DeliveryPlan{
		ContractVersion: DeliveryPlanContractVersion, RunID: run.ID(), PlanVersion: 1,
		TenantID: run.TenantID(), WorkspaceID: run.WorkspaceID(), ProjectID: run.ProjectID(),
		RepositoryID: run.RepositoryID(), SourceRevisions: slices.Clone(contextRecord.SourceRevisions),
		IntentVersion: 1, ContextPackageID: contextRecord.PackageID, ContextManifest: contextRecord.ManifestDigest,
		Objective: run.Objective(), Constraints: canonicalStrings(run.Constraints()), AcceptanceCriteria: canonicalStrings(run.AcceptanceCriteria()),
		Steps: cloneSteps(profile.Steps), AffectedCapabilities: slices.Clone(profile.AllowedCapabilities),
		ChangePaths: slices.Clone(profile.ChangePaths), Assumptions: slices.Clone(profile.Assumptions),
		Omissions: slices.Clone(contextRecord.Omissions), Conflicts: slices.Clone(contextRecord.Conflicts),
		Risks: slices.Clone(profile.Risks), Limitations: slices.Clone(profile.Limitations),
		StopConditions: slices.Clone(profile.StopConditions), ModelIDs: []string{}, ToolIDs: []string{},
		SecretReferenceIDs: []string{}, Destinations: []string{}, NetworkRequirements: []string{},
		FilesystemAccess: slices.Clone(profile.ChangePaths), Budgets: profile.Budgets,
		VerificationIntent:  slices.Clone(profile.VerificationIntent),
		GuardrailReferences: slices.Clone(profile.GuardrailReferences),
		Rollback:            profile.Rollback, ManualRecovery: profile.ManualRecovery,
		Classification: run.Classification(), RiskClass: profile.RiskClass,
		EvidenceCorrelation: correlationID, CreatedAt: at, ExpiresAt: expiresAt,
	}
	digest, err := planDigest(value)
	if err != nil {
		return DeliveryPlan{}, failure.New(failure.Indeterminate, "delivery plan digest could not be calculated", correlationID, failure.RetryNone)
	}
	value.ID = "forge.delivery-plan." + digest
	value.Digest = "sha256:" + digest
	if err := value.Validate(correlationID); err != nil {
		return DeliveryPlan{}, err
	}
	return value, nil
}

func (value DeliveryPlan) Validate(correlationID string) error {
	for _, item := range []string{
		value.ContractVersion, value.ID, value.Digest, value.RunID, value.TenantID, value.WorkspaceID,
		value.ProjectID, value.RepositoryID, value.ContextPackageID, value.ContextManifest, value.Objective,
		value.Rollback, value.ManualRecovery, value.Classification, value.RiskClass, value.EvidenceCorrelation,
	} {
		if item == "" {
			return failure.New(failure.InvalidRequest, "delivery plan has an empty required field", correlationID, failure.RetryNone)
		}
	}
	if value.ContractVersion != DeliveryPlanContractVersion || value.PlanVersion != 1 || value.IntentVersion != 1 ||
		value.CreatedAt.IsZero() || value.ExpiresAt.IsZero() || !value.ExpiresAt.After(value.CreatedAt) ||
		!oneOfRisk(value.RiskClass) || !validBudgets(value.Budgets) || value.EvidenceCorrelation != correlationID {
		return failure.New(failure.InvalidRequest, "delivery plan is invalid", correlationID, failure.RetryNone)
	}
	if !oneOfClassification(value.Classification) || len(value.SourceRevisions) == 0 ||
		len(value.AffectedCapabilities) == 0 || len(value.ChangePaths) == 0 ||
		len(value.VerificationIntent) == 0 || len(value.GuardrailReferences) == 0 ||
		len(value.StopConditions) == 0 || len(value.FilesystemAccess) == 0 ||
		len(value.ModelIDs) != 0 || len(value.ToolIDs) != 0 || len(value.SecretReferenceIDs) != 0 ||
		len(value.Destinations) != 0 || len(value.NetworkRequirements) != 0 {
		return failure.New(failure.InvalidRequest, "delivery plan exceeds the deterministic M1 profile", correlationID, failure.RetryNone)
	}
	for _, values := range [][]string{
		value.SourceRevisions, value.Constraints, value.AcceptanceCriteria, value.AffectedCapabilities,
		value.ChangePaths, value.Assumptions, value.Omissions, value.Conflicts, value.Risks,
		value.Limitations, value.StopConditions, value.FilesystemAccess, value.VerificationIntent,
		value.GuardrailReferences,
	} {
		if !uniqueNonEmptySorted(values) {
			return failure.New(failure.InvalidRequest, "delivery plan collections must be canonical", correlationID, failure.RetryNone)
		}
	}
	if err := validateGraph(value.Steps); err != nil {
		return failure.New(failure.InvalidRequest, "delivery plan graph is invalid", correlationID, failure.RetryNone)
	}
	digest, err := planDigest(value)
	if err != nil || value.ID != "forge.delivery-plan."+digest || value.Digest != "sha256:"+digest {
		return failure.New(failure.VerificationFailed, "delivery plan digest verification failed", correlationID, failure.RetryNone)
	}
	return nil
}

func planDigest(value DeliveryPlan) (string, error) {
	value.ID, value.Digest = "", ""
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return hex.EncodeToString(sum[:]), nil
}

func validateGraph(steps []PlanStep) error {
	if len(steps) == 0 || len(steps) > 32 {
		return errors.New("delivery plan graph size is invalid")
	}
	seen := make(map[string]int, len(steps))
	for index, step := range steps {
		if step.ID == "" || len(step.OutputKinds) == 0 || !uniqueNonEmptySorted(step.Dependencies) ||
			!uniqueNonEmptySorted(step.OutputKinds) {
			return errors.New("delivery plan step is invalid")
		}
		if _, duplicate := seen[step.ID]; duplicate {
			return errors.New("delivery plan step IDs must be unique")
		}
		seen[step.ID] = index
		for _, dependency := range step.Dependencies {
			if prior, found := seen[dependency]; !found || prior >= index {
				return errors.New("delivery plan graph must be ordered and acyclic")
			}
		}
	}
	return nil
}

func validBudgets(value PlanBudgets) bool {
	return value.MaxDurationSeconds > 0 && value.MaxDurationSeconds <= 7200 &&
		value.MaxConcurrency > 0 && value.MaxConcurrency <= 2 &&
		value.MaxAttempts > 0 && value.MaxAttempts <= 5 &&
		value.MaxComputeUnits >= 0 && value.MaxInputTokens >= 0 && value.MaxOutputTokens >= 0 &&
		value.MaxCostMinorUnits >= 0
}

func uniqueNonEmptySorted(values []string) bool {
	for index, value := range values {
		if value == "" || (index > 0 && value <= values[index-1]) {
			return false
		}
	}
	return true
}

func allowedPath(path string, prefixes []string) bool {
	clean := filepath.ToSlash(filepath.Clean(path))
	if clean == "." || strings.HasPrefix(clean, "../") || strings.HasPrefix(clean, "/") {
		return false
	}
	for _, prefix := range prefixes {
		cleanPrefix := strings.TrimSuffix(filepath.ToSlash(filepath.Clean(prefix)), "/")
		if clean == cleanPrefix || strings.HasPrefix(clean, cleanPrefix+"/") {
			return true
		}
	}
	return false
}

func clonePlan(value *DeliveryPlan) *DeliveryPlan {
	if value == nil {
		return nil
	}
	result := *value
	result.SourceRevisions = slices.Clone(value.SourceRevisions)
	result.Constraints = slices.Clone(value.Constraints)
	result.AcceptanceCriteria = slices.Clone(value.AcceptanceCriteria)
	result.Steps = cloneSteps(value.Steps)
	result.AffectedCapabilities = slices.Clone(value.AffectedCapabilities)
	result.ChangePaths = slices.Clone(value.ChangePaths)
	result.Assumptions = slices.Clone(value.Assumptions)
	result.Omissions = slices.Clone(value.Omissions)
	result.Conflicts = slices.Clone(value.Conflicts)
	result.Risks = slices.Clone(value.Risks)
	result.Limitations = slices.Clone(value.Limitations)
	result.StopConditions = slices.Clone(value.StopConditions)
	result.ModelIDs = slices.Clone(value.ModelIDs)
	result.ToolIDs = slices.Clone(value.ToolIDs)
	result.SecretReferenceIDs = slices.Clone(value.SecretReferenceIDs)
	result.Destinations = slices.Clone(value.Destinations)
	result.NetworkRequirements = slices.Clone(value.NetworkRequirements)
	result.FilesystemAccess = slices.Clone(value.FilesystemAccess)
	result.VerificationIntent = slices.Clone(value.VerificationIntent)
	result.GuardrailReferences = slices.Clone(value.GuardrailReferences)
	return &result
}

func cloneSteps(values []PlanStep) []PlanStep {
	result := make([]PlanStep, len(values))
	for index, value := range values {
		result[index] = PlanStep{ID: value.ID, Dependencies: slices.Clone(value.Dependencies), OutputKinds: slices.Clone(value.OutputKinds)}
	}
	return result
}

func clonePlanProfile(value PlanProfile) PlanProfile {
	value.Steps = cloneSteps(value.Steps)
	value.AllowedCapabilities = slices.Clone(value.AllowedCapabilities)
	value.AllowedPathPrefixes = slices.Clone(value.AllowedPathPrefixes)
	value.ChangePaths = slices.Clone(value.ChangePaths)
	value.Assumptions = slices.Clone(value.Assumptions)
	value.Risks = slices.Clone(value.Risks)
	value.Limitations = slices.Clone(value.Limitations)
	value.StopConditions = slices.Clone(value.StopConditions)
	value.VerificationIntent = slices.Clone(value.VerificationIntent)
	value.GuardrailReferences = slices.Clone(value.GuardrailReferences)
	return value
}

func canonicalStrings(values []string) []string {
	result := slices.Clone(values)
	sort.Strings(result)
	return result
}

func oneOfRisk(value string) bool {
	return value == "low" || value == "moderate" || value == "high" || value == "critical"
}

func oneOfClassification(value string) bool {
	return value == "Public" || value == "Internal" || value == "Confidential" || value == "Restricted"
}
