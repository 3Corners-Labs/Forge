package trust

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"slices"
	"sort"
	"strings"

	"github.com/3corners/forge/internal/contracts/failure"
)

const sha256IntegrityAlgorithm = "sha256"

// BindExecutionEnvelopeIntegrity returns a copy with a deterministic identity
// and digest. It performs no authorization and cannot turn caller data into
// authority; only the Policy capability may decide whether to issue it.
func BindExecutionEnvelopeIntegrity(value ExecutionEnvelope, correlationID string) (ExecutionEnvelope, error) {
	if value.ID != "" || value.Integrity != nil {
		return ExecutionEnvelope{}, invalid("execution envelope is already integrity-bound", correlationID)
	}
	value.ID = "forge.execution-envelope.pending"
	if err := value.Validate(correlationID); err != nil {
		return ExecutionEnvelope{}, err
	}
	value.ID = ""
	digest, err := executionEnvelopeDigest(value)
	if err != nil {
		return ExecutionEnvelope{}, failure.New(failure.Indeterminate, "execution envelope integrity could not be calculated", correlationID, failure.RetryNone)
	}
	value.ID = "forge.execution-envelope." + digest
	value.Integrity = &Integrity{Algorithm: sha256IntegrityAlgorithm, Digest: "sha256:" + digest}
	if err := value.VerifyIntegrity(correlationID); err != nil {
		return ExecutionEnvelope{}, err
	}
	return value, nil
}

// VerifyIntegrity detects mutation of every serialized envelope field except
// its derived ID and integrity record.
func (value ExecutionEnvelope) VerifyIntegrity(correlationID string) error {
	if err := value.Validate(correlationID); err != nil {
		return err
	}
	if value.Integrity == nil || value.Integrity.Algorithm != sha256IntegrityAlgorithm {
		return failure.New(failure.VerificationFailed, "execution envelope integrity is not verifiable", correlationID, failure.RetryNewAuthority)
	}
	digest, err := executionEnvelopeDigest(value)
	if err != nil {
		return failure.New(failure.Indeterminate, "execution envelope integrity could not be verified", correlationID, failure.RetryNone)
	}
	if value.ID != "forge.execution-envelope."+digest || value.Integrity.Digest != "sha256:"+digest {
		return failure.New(failure.VerificationFailed, "execution envelope integrity verification failed", correlationID, failure.RetryNewAuthority)
	}
	return nil
}

// BindContextPackageManifest returns a copy whose ID and manifest digest bind
// the exact information manifest. It does not add or carry authority.
func BindContextPackageManifest(value ContextPackage, correlationID string) (ContextPackage, error) {
	if value.ID != "" || value.ManifestDigest != "" {
		return ContextPackage{}, invalid("context package is already manifest-bound", correlationID)
	}
	value.ID = "forge.context-package.pending"
	value.ManifestDigest = "sha256:pending"
	if err := value.Validate(correlationID); err != nil {
		return ContextPackage{}, err
	}
	value.ID = ""
	value.ManifestDigest = ""
	digest, err := contextPackageDigest(value)
	if err != nil {
		return ContextPackage{}, failure.New(failure.Indeterminate, "context manifest digest could not be calculated", correlationID, failure.RetryNone)
	}
	value.ID = "forge.context-package." + digest
	value.ManifestDigest = "sha256:" + digest
	if err := value.VerifyManifest(correlationID); err != nil {
		return ContextPackage{}, err
	}
	return value, nil
}

// VerifyManifest detects mutation of every serialized context-package field
// except its derived ID and manifest digest.
func (value ContextPackage) VerifyManifest(correlationID string) error {
	if err := value.Validate(correlationID); err != nil {
		return err
	}
	digest, err := contextPackageDigest(value)
	if err != nil {
		return failure.New(failure.Indeterminate, "context manifest digest could not be verified", correlationID, failure.RetryNone)
	}
	if value.ID != "forge.context-package."+digest || value.ManifestDigest != "sha256:"+digest {
		return failure.New(failure.VerificationFailed, "context manifest verification failed", correlationID, failure.RetryNone)
	}
	return nil
}

func BindApprovalRequestDigest(value ApprovalRequest, correlationID string) (ApprovalRequest, error) {
	if value.ID != "" || value.RequestDigest != "" {
		return ApprovalRequest{}, invalid("approval request is already digest-bound", correlationID)
	}
	value.EvidenceReferences = canonicalApprovalValues(value.EvidenceReferences)
	value.PermittedConditions = canonicalApprovalValues(value.PermittedConditions)
	value.PolicyDecisionIDs = canonicalApprovalValues(value.PolicyDecisionIDs)
	value.ID = "forge.approval-request.pending"
	value.RequestDigest = "sha256:" + strings.Repeat("0", 64)
	if err := value.Validate(correlationID); err != nil {
		return ApprovalRequest{}, err
	}
	value.ID, value.RequestDigest = "", ""
	digest, err := approvalRequestDigest(value)
	if err != nil {
		return ApprovalRequest{}, failure.New(failure.Indeterminate, "approval request digest could not be calculated", correlationID, failure.RetryNone)
	}
	value.ID = "forge.approval-request." + digest
	value.RequestDigest = "sha256:" + digest
	if err := value.VerifyDigest(correlationID); err != nil {
		return ApprovalRequest{}, err
	}
	return value, nil
}

func (value ApprovalRequest) VerifyDigest(correlationID string) error {
	if err := value.Validate(correlationID); err != nil {
		return err
	}
	digest, err := approvalRequestDigest(value)
	if err != nil {
		return failure.New(failure.Indeterminate, "approval request digest could not be verified", correlationID, failure.RetryNone)
	}
	if value.ID != "forge.approval-request."+digest || value.RequestDigest != "sha256:"+digest {
		return failure.New(failure.VerificationFailed, "approval request digest verification failed", correlationID, failure.RetryNone)
	}
	return nil
}

func BindApprovalRecordV2Integrity(value ApprovalRecordV2, correlationID string) (ApprovalRecordV2, error) {
	if value.ID != "" || value.Integrity != (Integrity{}) {
		return ApprovalRecordV2{}, invalid("approval record v2 is already integrity-bound", correlationID)
	}
	value.Conditions = canonicalApprovalValues(value.Conditions)
	value.PolicyDecisionIDs = canonicalApprovalValues(value.PolicyDecisionIDs)
	value.EvidenceReferences = canonicalApprovalValues(value.EvidenceReferences)
	value.ID = "forge.approval-record.pending"
	value.Integrity = Integrity{Algorithm: sha256IntegrityAlgorithm, Digest: "sha256:" + strings.Repeat("0", 64)}
	if err := value.Validate(correlationID); err != nil {
		return ApprovalRecordV2{}, err
	}
	value.ID = ""
	value.Integrity = Integrity{}
	digest, err := approvalRecordV2Digest(value)
	if err != nil {
		return ApprovalRecordV2{}, failure.New(failure.Indeterminate, "approval record v2 integrity could not be calculated", correlationID, failure.RetryNone)
	}
	value.ID = "forge.approval-record." + digest
	value.Integrity = Integrity{Algorithm: sha256IntegrityAlgorithm, Digest: "sha256:" + digest}
	if err := value.VerifyIntegrity(correlationID); err != nil {
		return ApprovalRecordV2{}, err
	}
	return value, nil
}

func (value ApprovalRecordV2) VerifyIntegrity(correlationID string) error {
	if err := value.Validate(correlationID); err != nil {
		return err
	}
	if value.Integrity.Algorithm != sha256IntegrityAlgorithm {
		return failure.New(failure.VerificationFailed, "approval record v2 integrity is not verifiable", correlationID, failure.RetryNone)
	}
	digest, err := approvalRecordV2Digest(value)
	if err != nil {
		return failure.New(failure.Indeterminate, "approval record v2 integrity could not be verified", correlationID, failure.RetryNone)
	}
	if value.ID != "forge.approval-record."+digest || value.Integrity.Digest != "sha256:"+digest {
		return failure.New(failure.VerificationFailed, "approval record v2 integrity verification failed", correlationID, failure.RetryNone)
	}
	return nil
}

func executionEnvelopeDigest(value ExecutionEnvelope) (string, error) {
	value.ID = ""
	value.Integrity = nil
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func contextPackageDigest(value ContextPackage) (string, error) {
	value.ID = ""
	value.ManifestDigest = ""
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func approvalRequestDigest(value ApprovalRequest) (string, error) {
	value.ID = ""
	value.RequestDigest = ""
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func approvalRecordV2Digest(value ApprovalRecordV2) (string, error) {
	value.ID = ""
	value.Integrity = Integrity{}
	encoded, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	digest := sha256.Sum256(encoded)
	return hex.EncodeToString(digest[:]), nil
}

func canonicalApprovalValues(values []string) []string {
	result := slices.Clone(values)
	sort.Strings(result)
	return result
}
