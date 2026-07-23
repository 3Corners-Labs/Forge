// Package worker defines the closed contracts shared by the control plane and
// the isolated M1 worker. It contains no runtime or control-plane dependency.
package worker

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"path"
	"slices"
	"sort"
	"strings"
	"time"
)

const (
	RequestVersion = "forge.worker-request.v1"
	ResultVersion  = "forge.worker-result.v1"
	BundleVersion  = "forge.candidate-bundle.v1"

	PrepareChangeAction = "forge.worker.prepare-change"
	TransformationV1    = "forge.deterministic-file-replacement.v1"
)

var (
	ErrInvalidContract = errors.New("worker contract is invalid")
	ErrDigestMismatch  = errors.New("worker contract digest mismatch")
)

type Scope struct {
	TenantID      string `json:"tenantId"`
	WorkspaceID   string `json:"workspaceId"`
	ProjectID     string `json:"projectId"`
	RepositoryID  string `json:"repositoryId"`
	DeliveryRunID string `json:"deliveryRunId"`
}

type Limits struct {
	MaxDurationSeconds int64 `json:"maxDurationSeconds"`
	MaxProcesses       int64 `json:"maxProcesses"`
	MaxOutputBytes     int64 `json:"maxOutputBytes"`
	MaxChangedFiles    int64 `json:"maxChangedFiles"`
	MaxAttempts        int64 `json:"maxAttempts"`
}

// FileReplacement is data for one Forge-owned transformation. It cannot name
// a command, executable, model, tool, destination, or secret.
type FileReplacement struct {
	Path                 string `json:"path"`
	ExpectedSourceDigest string `json:"expectedSourceDigest,omitempty"`
	Content              []byte `json:"content"`
}

type Request struct {
	SchemaVersion         string            `json:"schemaVersion"`
	ID                    string            `json:"id"`
	Digest                string            `json:"digest"`
	EnvelopeID            string            `json:"envelopeId"`
	DelegationID          string            `json:"delegationId"`
	PlanID                string            `json:"planId"`
	PlanDigest            string            `json:"planDigest"`
	PlanApprovalRecordID  string            `json:"planApprovalRecordId"`
	ContextPackageID      string            `json:"contextPackageId"`
	ContextManifestDigest string            `json:"contextManifestDigest"`
	SourceRevision        string            `json:"sourceRevision"`
	Scope                 Scope             `json:"scope"`
	Action                string            `json:"action"`
	TransformationID      string            `json:"transformationId"`
	TransformationVersion string            `json:"transformationVersion"`
	AllowedReadPaths      []string          `json:"allowedReadPaths"`
	AllowedWritePaths     []string          `json:"allowedWritePaths"`
	DependencyDigests     []string          `json:"dependencyDigests"`
	Replacements          []FileReplacement `json:"replacements"`
	Limits                Limits            `json:"limits"`
	Attempt               int64             `json:"attempt"`
	ExpiresAt             time.Time         `json:"expiresAt"`
	EvidenceCorrelationID string            `json:"evidenceCorrelationId"`
}

type CandidateFile struct {
	Path          string `json:"path"`
	SourceDigest  string `json:"sourceDigest,omitempty"`
	Content       []byte `json:"content"`
	ContentDigest string `json:"contentDigest"`
}

type CandidateBundle struct {
	SchemaVersion  string          `json:"schemaVersion"`
	SourceRevision string          `json:"sourceRevision"`
	PlanDigest     string          `json:"planDigest"`
	Files          []CandidateFile `json:"files"`
	Digest         string          `json:"digest"`
}

type ResourceUse struct {
	DurationMillis int64 `json:"durationMillis"`
	OutputBytes    int64 `json:"outputBytes"`
	ChangedFiles   int64 `json:"changedFiles"`
}

type Result struct {
	SchemaVersion      string      `json:"schemaVersion"`
	RequestID          string      `json:"requestId"`
	RequestDigest      string      `json:"requestDigest"`
	CandidateDigest    string      `json:"candidateDigest"`
	TreeDigest         string      `json:"treeDigest"`
	ChangedPaths       []string    `json:"changedPaths"`
	BundleDigest       string      `json:"bundleDigest"`
	ResourceUse        ResourceUse `json:"resourceUse"`
	CleanupOutcome     string      `json:"cleanupOutcome"`
	EvidenceReferences []string    `json:"evidenceReferences"`
	CompletedAt        time.Time   `json:"completedAt"`
}

func BindRequest(value Request) (Request, error) {
	value.ID, value.Digest = "", ""
	canonicalizeRequest(&value)
	if err := validateRequestShape(value); err != nil {
		return Request{}, err
	}
	digest, err := digestJSON(value)
	if err != nil {
		return Request{}, err
	}
	value.Digest = digest
	value.ID = "forge.worker-request." + strings.TrimPrefix(digest, "sha256:")
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

func BindBundle(value CandidateBundle) (CandidateBundle, error) {
	value.Digest = ""
	sort.Slice(value.Files, func(i, j int) bool { return value.Files[i].Path < value.Files[j].Path })
	for index := range value.Files {
		if !safePath(value.Files[index].Path) || len(value.Files[index].Content) == 0 {
			return CandidateBundle{}, ErrInvalidContract
		}
		value.Files[index].ContentDigest = sha256Digest(value.Files[index].Content)
		if value.Files[index].SourceDigest != "" && !validDigest(value.Files[index].SourceDigest) {
			return CandidateBundle{}, ErrInvalidContract
		}
		if index > 0 && value.Files[index-1].Path == value.Files[index].Path {
			return CandidateBundle{}, ErrInvalidContract
		}
	}
	if value.SchemaVersion != BundleVersion || value.SourceRevision == "" ||
		!validDigest(value.PlanDigest) || len(value.Files) == 0 {
		return CandidateBundle{}, ErrInvalidContract
	}
	digest, err := digestJSON(value)
	if err != nil {
		return CandidateBundle{}, err
	}
	value.Digest = digest
	return value, nil
}

func (value CandidateBundle) Validate() error {
	bound, err := BindBundle(value)
	if err != nil {
		return err
	}
	if bound.Digest != value.Digest {
		return ErrDigestMismatch
	}
	return nil
}

func DecodeRequest(data []byte) (Request, error) {
	var value Request
	if err := decodeStrict(data, &value); err != nil {
		return Request{}, err
	}
	return value, nil
}

func DecodeBundle(data []byte) (CandidateBundle, error) {
	var value CandidateBundle
	if err := decodeStrict(data, &value); err != nil {
		return CandidateBundle{}, err
	}
	return value, nil
}

func DecodeResult(data []byte) (Result, error) {
	var value Result
	if err := decodeStrict(data, &value); err != nil {
		return Result{}, err
	}
	return value, nil
}

func EncodeCanonical(value any) ([]byte, error) {
	return json.Marshal(value)
}

func CandidateTreeDigest(bundle CandidateBundle) string {
	type entry struct{ Path, Digest string }
	values := make([]entry, 0, len(bundle.Files))
	for _, file := range bundle.Files {
		values = append(values, entry{file.Path, file.ContentDigest})
	}
	encoded, _ := json.Marshal(values)
	return sha256Digest(encoded)
}

func canonicalizeRequest(value *Request) {
	sort.Strings(value.AllowedReadPaths)
	value.AllowedReadPaths = slices.Compact(value.AllowedReadPaths)
	sort.Strings(value.AllowedWritePaths)
	value.AllowedWritePaths = slices.Compact(value.AllowedWritePaths)
	sort.Strings(value.DependencyDigests)
	value.DependencyDigests = slices.Compact(value.DependencyDigests)
	sort.Slice(value.Replacements, func(i, j int) bool { return value.Replacements[i].Path < value.Replacements[j].Path })
}

func validateRequestShape(value Request) error {
	if value.SchemaVersion != RequestVersion || value.EnvelopeID == "" || value.DelegationID == "" ||
		value.PlanID == "" || !validDigest(value.PlanDigest) || value.PlanApprovalRecordID == "" ||
		value.ContextPackageID == "" || !validDigest(value.ContextManifestDigest) || value.SourceRevision == "" ||
		value.Action != PrepareChangeAction || value.TransformationID != TransformationV1 ||
		value.TransformationVersion == "" || value.EvidenceCorrelationID == "" ||
		value.ExpiresAt.IsZero() || value.Attempt < 1 || value.Attempt > value.Limits.MaxAttempts ||
		value.Limits.MaxDurationSeconds < 1 || value.Limits.MaxProcesses < 1 ||
		value.Limits.MaxOutputBytes < 1 || value.Limits.MaxChangedFiles < 1 ||
		value.Limits.MaxAttempts < 1 || len(value.AllowedReadPaths) == 0 ||
		len(value.AllowedWritePaths) == 0 || len(value.Replacements) == 0 ||
		int64(len(value.Replacements)) > value.Limits.MaxChangedFiles {
		return ErrInvalidContract
	}
	for _, field := range []string{value.Scope.TenantID, value.Scope.WorkspaceID, value.Scope.ProjectID, value.Scope.RepositoryID, value.Scope.DeliveryRunID} {
		if field == "" {
			return ErrInvalidContract
		}
	}
	for _, item := range value.DependencyDigests {
		if !validDigest(item) {
			return ErrInvalidContract
		}
	}
	for index, replacement := range value.Replacements {
		if !safePath(replacement.Path) || len(replacement.Content) == 0 ||
			(replacement.ExpectedSourceDigest != "" && !validDigest(replacement.ExpectedSourceDigest)) ||
			!withinAny(replacement.Path, value.AllowedWritePaths) ||
			(index > 0 && value.Replacements[index-1].Path == replacement.Path) {
			return ErrInvalidContract
		}
	}
	return nil
}

func withinAny(candidate string, roots []string) bool {
	for _, root := range roots {
		if safePath(root) && (candidate == root || strings.HasPrefix(candidate, strings.TrimSuffix(root, "/")+"/")) {
			return true
		}
	}
	return false
}

func safePath(value string) bool {
	return value != "" && value != "." && !strings.HasPrefix(value, "/") &&
		!strings.Contains(value, "\\") && path.Clean(value) == value &&
		!strings.HasPrefix(value, "../") && !strings.Contains(value, "/../")
}

func digestJSON(value any) (string, error) {
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return sha256Digest(encoded), nil
}

func sha256Digest(value []byte) string {
	sum := sha256.Sum256(value)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func validDigest(value string) bool {
	if !strings.HasPrefix(value, "sha256:") || len(value) != 71 {
		return false
	}
	_, err := hex.DecodeString(strings.TrimPrefix(value, "sha256:"))
	return err == nil
}

func decodeStrict(data []byte, destination any) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return ErrInvalidContract
	}
	var trailing any
	if err := decoder.Decode(&trailing); !errors.Is(err, io.EOF) {
		return ErrInvalidContract
	}
	return nil
}
