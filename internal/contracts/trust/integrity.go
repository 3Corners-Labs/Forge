package trust

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"

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
