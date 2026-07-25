// Package change defines the closed, provider-neutral contract for
// non-authoritative model-assisted change preparation.
package change

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"path"
	"slices"
	"sort"
	"strings"
	"time"
	"unicode/utf8"
)

const (
	RequestVersion  = "forge.change-preparation-request.v1"
	ProposalVersion = "forge.change-proposal.v1"
	PurposeV1       = "forge.m1.prepare-change"
	RoleProfileV1   = "forge.role.software-engineer.change-preparation.v1"
)

var (
	ErrInvalidContract = errors.New("change-preparation contract is invalid")
	ErrDigestMismatch  = errors.New("change-preparation contract digest mismatch")
)

type Limits struct {
	MaxInputBytes   int64         `json:"maxInputBytes"`
	MaxOutputBytes  int64         `json:"maxOutputBytes"`
	MaxTokens       int64         `json:"maxTokens"`
	MaxLatency      time.Duration `json:"maxLatency"`
	MaxCostMicros   int64         `json:"maxCostMicros"`
	MaxChangedFiles int64         `json:"maxChangedFiles"`
}

type SourceFile struct {
	Path    string `json:"path"`
	Digest  string `json:"digest"`
	Content string `json:"content"`
}

type Request struct {
	SchemaVersion         string       `json:"schemaVersion"`
	ID                    string       `json:"id"`
	Digest                string       `json:"digest"`
	CorrelationID         string       `json:"correlationId"`
	TenantID              string       `json:"tenantId"`
	WorkspaceID           string       `json:"workspaceId"`
	ProjectID             string       `json:"projectId"`
	RepositoryID          string       `json:"repositoryId"`
	DeliveryRunID         string       `json:"deliveryRunId"`
	TaskAgentID           string       `json:"taskAgentId"`
	PlanID                string       `json:"planId"`
	PlanDigest            string       `json:"planDigest"`
	ContextPackageID      string       `json:"contextPackageId"`
	ContextManifestDigest string       `json:"contextManifestDigest"`
	RequirementsRecordID  string       `json:"requirementsRecordId"`
	SourceRevision        string       `json:"sourceRevision"`
	RoleProfileID         string       `json:"roleProfileId"`
	Purpose               string       `json:"purpose"`
	Classification        string       `json:"classification"`
	AllowedReadPaths      []string     `json:"allowedReadPaths"`
	AllowedWritePaths     []string     `json:"allowedWritePaths"`
	Sources               []SourceFile `json:"sources"`
	AcceptanceCriteria    []string     `json:"acceptanceCriteria"`
	TestIntent            []string     `json:"testIntent"`
	Exclusions            []string     `json:"exclusions"`
	ControlReferences     []string     `json:"controlReferences"`
	Limits                Limits       `json:"limits"`
	ExpiresAt             time.Time    `json:"expiresAt"`
}

type FileReplacement struct {
	Path                 string `json:"path"`
	ExpectedSourceDigest string `json:"expectedSourceDigest,omitempty"`
	Content              string `json:"content"`
	Rationale            string `json:"rationale"`
	CandidateTest        bool   `json:"candidateTest"`
}

type Usage struct {
	InputTokens  int64 `json:"inputTokens"`
	OutputTokens int64 `json:"outputTokens"`
	LatencyMs    int64 `json:"latencyMs"`
	CostMicros   int64 `json:"costMicros"`
}

type Receipt struct {
	ActivationID     string `json:"activationId"`
	ActivationDigest string `json:"activationDigest"`
	Provider         string `json:"provider"`
	ModelID          string `json:"modelId"`
	AdapterVersion   string `json:"adapterVersion"`
	RoleProfileID    string `json:"roleProfileId"`
	PromptVersion    string `json:"promptVersion"`
	ResponseSchema   string `json:"responseSchema"`
	ProviderRequest  string `json:"providerRequest"`
}

type Proposal struct {
	SchemaVersion  string            `json:"schemaVersion"`
	ID             string            `json:"id"`
	Digest         string            `json:"digest"`
	RequestID      string            `json:"requestId"`
	RequestDigest  string            `json:"requestDigest"`
	SourceRevision string            `json:"sourceRevision"`
	PlanDigest     string            `json:"planDigest"`
	ContextDigest  string            `json:"contextDigest"`
	Replacements   []FileReplacement `json:"replacements"`
	ProposedChecks []string          `json:"proposedChecks"`
	Assumptions    []string          `json:"assumptions"`
	Uncertainty    []string          `json:"uncertainty"`
	Omissions      []string          `json:"omissions"`
	Refusal        bool              `json:"refusal"`
	Incomplete     bool              `json:"incomplete"`
	Usage          Usage             `json:"usage"`
	Receipt        Receipt           `json:"receipt"`
}

func BindRequest(value Request) (Request, error) {
	value.ID, value.Digest = "", ""
	canonicalizeRequest(&value)
	if !validRequestShape(value) {
		return Request{}, ErrInvalidContract
	}
	digest, err := digestJSON(value)
	if err != nil {
		return Request{}, ErrInvalidContract
	}
	value.Digest = digest
	value.ID = "forge.change-request." + strings.TrimPrefix(digest, "sha256:")
	return value, nil
}

func (value Request) Validate(at time.Time) error {
	bound, err := BindRequest(value)
	if err != nil {
		return err
	}
	if bound.ID != value.ID || bound.Digest != value.Digest {
		return ErrDigestMismatch
	}
	if at.IsZero() || !at.Before(value.ExpiresAt) {
		return ErrInvalidContract
	}
	return nil
}

func BindProposal(value Proposal, request Request) (Proposal, error) {
	value.ID, value.Digest = "", ""
	sort.Slice(value.Replacements, func(i, j int) bool { return value.Replacements[i].Path < value.Replacements[j].Path })
	value.ProposedChecks = canonicalStrings(value.ProposedChecks)
	value.Assumptions = canonicalStrings(value.Assumptions)
	value.Uncertainty = canonicalStrings(value.Uncertainty)
	value.Omissions = canonicalStrings(value.Omissions)
	if err := request.Validate(request.ExpiresAt.Add(-time.Nanosecond)); err != nil ||
		!validProposalShape(value, request) {
		return Proposal{}, ErrInvalidContract
	}
	digest, err := digestJSON(value)
	if err != nil {
		return Proposal{}, ErrInvalidContract
	}
	value.Digest = digest
	value.ID = "forge.change-proposal." + strings.TrimPrefix(digest, "sha256:")
	return value, nil
}

func (value Proposal) Validate(request Request) error {
	bound, err := BindProposal(value, request)
	if err != nil {
		return err
	}
	if bound.ID != value.ID || bound.Digest != value.Digest {
		return ErrDigestMismatch
	}
	return nil
}

func validRequestShape(value Request) bool {
	if value.SchemaVersion != RequestVersion ||
		blank(value.CorrelationID, value.TenantID, value.WorkspaceID, value.ProjectID, value.RepositoryID,
			value.DeliveryRunID, value.TaskAgentID, value.PlanID, value.PlanDigest, value.ContextPackageID,
			value.ContextManifestDigest, value.RequirementsRecordID, value.SourceRevision, value.RoleProfileID,
			value.Purpose, value.Classification) ||
		value.RoleProfileID != RoleProfileV1 || value.Purpose != PurposeV1 ||
		!validDigest(value.PlanDigest) || !validDigest(value.ContextManifestDigest) ||
		len(value.AllowedReadPaths) == 0 || len(value.AllowedWritePaths) == 0 ||
		len(value.AcceptanceCriteria) == 0 ||
		len(value.TestIntent) == 0 || len(value.ControlReferences) == 0 ||
		value.Limits.MaxInputBytes < 1 || value.Limits.MaxOutputBytes < 1 ||
		value.Limits.MaxTokens < 1 || value.Limits.MaxLatency <= 0 ||
		value.Limits.MaxCostMicros < 0 || value.Limits.MaxChangedFiles < 1 ||
		value.ExpiresAt.IsZero() {
		return false
	}
	if !validTexts(value.AcceptanceCriteria) || !validTexts(value.TestIntent) ||
		!validTexts(value.Exclusions) || !validTexts(value.ControlReferences) {
		return false
	}
	for _, item := range append(slices.Clone(value.AllowedReadPaths), value.AllowedWritePaths...) {
		if !safePath(item) {
			return false
		}
	}
	seen := map[string]struct{}{}
	for _, source := range value.Sources {
		if !safePath(source.Path) || !withinAny(source.Path, value.AllowedReadPaths) ||
			!validDigest(source.Digest) || source.Digest != textDigest(source.Content) ||
			!validText(source.Content) {
			return false
		}
		if _, exists := seen[source.Path]; exists {
			return false
		}
		seen[source.Path] = struct{}{}
	}
	encoded, err := json.Marshal(value)
	return err == nil && int64(len(encoded)) <= value.Limits.MaxInputBytes
}

func validProposalShape(value Proposal, request Request) bool {
	if value.SchemaVersion != ProposalVersion || value.RequestID != request.ID ||
		value.RequestDigest != request.Digest || value.SourceRevision != request.SourceRevision ||
		value.PlanDigest != request.PlanDigest || value.ContextDigest != request.ContextManifestDigest ||
		blank(value.Receipt.ActivationID, value.Receipt.ActivationDigest, value.Receipt.Provider,
			value.Receipt.ModelID, value.Receipt.AdapterVersion, value.Receipt.RoleProfileID,
			value.Receipt.PromptVersion, value.Receipt.ResponseSchema, value.Receipt.ProviderRequest) ||
		!validDigest(value.Receipt.ActivationDigest) ||
		value.Receipt.RoleProfileID != request.RoleProfileID ||
		value.Receipt.ResponseSchema != ProposalVersion ||
		value.Usage.InputTokens < 0 || value.Usage.OutputTokens < 0 ||
		value.Usage.LatencyMs < 0 || value.Usage.CostMicros < 0 ||
		value.Usage.InputTokens+value.Usage.OutputTokens > request.Limits.MaxTokens ||
		value.Usage.CostMicros > request.Limits.MaxCostMicros ||
		time.Duration(value.Usage.LatencyMs)*time.Millisecond > request.Limits.MaxLatency ||
		!validTexts(value.ProposedChecks) || !validTexts(value.Assumptions) ||
		!validTexts(value.Uncertainty) || !validTexts(value.Omissions) {
		return false
	}
	if value.Refusal || value.Incomplete {
		return len(value.Replacements) == 0
	}
	if len(value.Replacements) == 0 || int64(len(value.Replacements)) > request.Limits.MaxChangedFiles {
		return false
	}
	sourceDigests := make(map[string]string, len(request.Sources))
	for _, source := range request.Sources {
		sourceDigests[source.Path] = source.Digest
	}
	seen := map[string]struct{}{}
	for _, replacement := range value.Replacements {
		if !safePath(replacement.Path) || !withinAny(replacement.Path, request.AllowedWritePaths) ||
			!validText(replacement.Content) || strings.TrimSpace(replacement.Rationale) == "" ||
			len(replacement.Rationale) > 16<<10 {
			return false
		}
		if expected, exists := sourceDigests[replacement.Path]; exists {
			if replacement.ExpectedSourceDigest != expected {
				return false
			}
		} else if replacement.ExpectedSourceDigest != "" {
			return false
		}
		if _, exists := seen[replacement.Path]; exists {
			return false
		}
		seen[replacement.Path] = struct{}{}
	}
	encoded, err := json.Marshal(value)
	return err == nil && int64(len(encoded)) <= request.Limits.MaxOutputBytes
}

func canonicalizeRequest(value *Request) {
	value.AllowedReadPaths = canonicalStrings(value.AllowedReadPaths)
	value.AllowedWritePaths = canonicalStrings(value.AllowedWritePaths)
	value.AcceptanceCriteria = canonicalStrings(value.AcceptanceCriteria)
	value.TestIntent = canonicalStrings(value.TestIntent)
	value.Exclusions = canonicalStrings(value.Exclusions)
	value.ControlReferences = canonicalStrings(value.ControlReferences)
	sort.Slice(value.Sources, func(i, j int) bool { return value.Sources[i].Path < value.Sources[j].Path })
}

func canonicalStrings(values []string) []string {
	result := slices.Clone(values)
	sort.Strings(result)
	return slices.Compact(result)
}

func safePath(value string) bool {
	return value != "" && value != "." && !strings.HasPrefix(value, "/") &&
		!strings.Contains(value, "\\") && path.Clean(value) == value &&
		!strings.HasPrefix(value, "../") && !strings.Contains(value, "/../")
}

func withinAny(candidate string, roots []string) bool {
	for _, root := range roots {
		if safePath(root) && (candidate == root || strings.HasPrefix(candidate, strings.TrimSuffix(root, "/")+"/")) {
			return true
		}
	}
	return false
}

func validTexts(values []string) bool {
	if len(values) > 128 {
		return false
	}
	for _, value := range values {
		if !validText(value) {
			return false
		}
	}
	return true
}

func validText(value string) bool {
	return value != "" && len(value) <= 256<<10 && utf8.ValidString(value) && !strings.ContainsRune(value, '\x00')
}

func validDigest(value string) bool {
	if !strings.HasPrefix(value, "sha256:") || len(value) != 71 {
		return false
	}
	for _, character := range value[7:] {
		if (character < '0' || character > '9') && (character < 'a' || character > 'f') {
			return false
		}
	}
	return true
}

func blank(values ...string) bool {
	for _, value := range values {
		if strings.TrimSpace(value) == "" {
			return true
		}
	}
	return false
}

func textDigest(value string) string {
	sum := sha256.Sum256([]byte(value))
	return "sha256:" + hex.EncodeToString(sum[:])
}

func digestJSON(value any) (string, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
