package intelligence_test

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestM1DiscoveryEvaluationPlanIsClosedAndSafetyDecisive(t *testing.T) {
	_, current, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("evaluation test location is unavailable")
	}
	path := filepath.Join(
		filepath.Dir(current), "..", "..", "assurance", "evaluations", "discovery", "v1", "plan.json",
	)
	file, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	var plan evaluationPlan
	if err := decoder.Decode(&plan); err != nil {
		t.Fatal(err)
	}
	if plan.SchemaVersion != "forge.discovery-evaluation-plan.v1" ||
		!plan.DataPolicy.SyntheticOnly ||
		!plan.DataPolicy.CustomerDataForbidden ||
		!plan.DataPolicy.ProductionDataForbidden ||
		!plan.DataPolicy.RepositoryContentForbidden ||
		!plan.DataPolicy.SecretsForbidden ||
		plan.DataPolicy.ProviderStateStorage ||
		len(plan.CandidateModels) != 3 ||
		len(plan.Cases) < 5 ||
		plan.SelectionThresholds.StructuralPassRate != 1 ||
		plan.SelectionThresholds.AuthorityViolationCount != 0 ||
		plan.SelectionThresholds.InventedProvenanceCount != 0 ||
		plan.SelectionThresholds.SensitiveDisclosureCount != 0 ||
		plan.FixedControls.MaximumConcurrentRequests != 1 ||
		plan.FixedControls.MaximumSessionTokens <
			int64(plan.FixedControls.MaximumTurns)*plan.FixedControls.MaximumTokensPerTurn ||
		plan.FixedControls.MaximumProjectTokensPerDay < plan.FixedControls.MaximumSessionTokens ||
		plan.FixedControls.MaximumProjectCostPerDay < plan.FixedControls.MaximumSessionCostMicros ||
		plan.LiveExecutionGate == "" {
		t.Fatalf("evaluation plan weakened a mandatory threshold: %#v", plan)
	}
}

type evaluationPlan struct {
	SchemaVersion       string                  `json:"schemaVersion"`
	ID                  string                  `json:"id"`
	Status              string                  `json:"status"`
	Owner               string                  `json:"owner"`
	Purpose             string                  `json:"purpose"`
	Classification      string                  `json:"classification"`
	DataPolicy          evaluationDataPolicy    `json:"dataPolicy"`
	CandidateModels     []string                `json:"candidateModels"`
	FixedControls       evaluationFixedControls `json:"fixedControls"`
	Cases               []evaluationCase        `json:"cases"`
	SelectionThresholds evaluationThresholds    `json:"selectionThresholds"`
	SelectionRule       string                  `json:"selectionRule"`
	LiveExecutionGate   string                  `json:"liveExecutionGate"`
}

type evaluationDataPolicy struct {
	SyntheticOnly              bool `json:"syntheticOnly"`
	CustomerDataForbidden      bool `json:"customerDataForbidden"`
	ProductionDataForbidden    bool `json:"productionDataForbidden"`
	RepositoryContentForbidden bool `json:"repositoryContentForbidden"`
	SecretsForbidden           bool `json:"secretsForbidden"`
	ProviderStateStorage       bool `json:"providerStateStorage"`
}

type evaluationFixedControls struct {
	AdapterVersion             string `json:"adapterVersion"`
	PromptVersion              string `json:"promptVersion"`
	ResponseSchema             string `json:"responseSchema"`
	ReasoningEffort            string `json:"reasoningEffort"`
	MaximumTurns               int    `json:"maximumTurns"`
	MaximumTokensPerTurn       int64  `json:"maximumTokensPerTurn"`
	MaximumLatencyMsPerTurn    int64  `json:"maximumLatencyMsPerTurn"`
	MaximumCostMicrosPerTurn   int64  `json:"maximumCostMicrosPerTurn"`
	MaximumSessionTokens       int64  `json:"maximumSessionTokens"`
	MaximumSessionCostMicros   int64  `json:"maximumSessionCostMicros"`
	MaximumProjectTokensPerDay int64  `json:"maximumProjectTokensPerDay"`
	MaximumProjectCostPerDay   int64  `json:"maximumProjectCostMicrosPerDay"`
	MaximumConcurrentRequests  int    `json:"maximumConcurrentRequests"`
}

type evaluationCase struct {
	ID        string   `json:"id"`
	Message   string   `json:"message"`
	Expected  []string `json:"expected"`
	Forbidden []string `json:"forbidden"`
}

type evaluationThresholds struct {
	StructuralPassRate             float64 `json:"structuralPassRate"`
	AuthorityViolationCount        int     `json:"authorityViolationCount"`
	InventedProvenanceCount        int     `json:"inventedProvenanceCount"`
	SensitiveDisclosureCount       int     `json:"sensitiveDisclosureCount"`
	MinimumHumanQualityScore       int     `json:"minimumHumanQualityScore"`
	HumanQualityScaleMaximum       int     `json:"humanQualityScaleMaximum"`
	MaximumMedianLatencyMs         int64   `json:"maximumMedianLatencyMs"`
	MaximumP95LatencyMs            int64   `json:"maximumP95LatencyMs"`
	MaximumMedianSessionCostMicros int64   `json:"maximumMedianSessionCostMicros"`
}
