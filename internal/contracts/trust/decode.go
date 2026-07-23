package trust

import (
	"bytes"
	"encoding/json"
	"io"

	"github.com/3corners/forge/internal/contracts/failure"
)

func DecodePolicyDecision(data []byte, correlationID string) (PolicyDecision, error) {
	var value PolicyDecision
	if err := decodeStrict(data, &value, "policy decision", correlationID); err != nil {
		return PolicyDecision{}, err
	}
	if err := value.Validate(correlationID); err != nil {
		return PolicyDecision{}, err
	}
	return value, nil
}

func DecodeExecutionEnvelope(data []byte, correlationID string) (ExecutionEnvelope, error) {
	var value ExecutionEnvelope
	if err := decodeStrict(data, &value, "execution envelope", correlationID); err != nil {
		return ExecutionEnvelope{}, err
	}
	if err := value.Validate(correlationID); err != nil {
		return ExecutionEnvelope{}, err
	}
	return value, nil
}

func DecodeContextPackage(data []byte, correlationID string) (ContextPackage, error) {
	var value ContextPackage
	if err := decodeStrict(data, &value, "context package", correlationID); err != nil {
		return ContextPackage{}, err
	}
	if err := value.Validate(correlationID); err != nil {
		return ContextPackage{}, err
	}
	return value, nil
}

func DecodeApprovalRecord(data []byte, correlationID string) (ApprovalRecord, error) {
	var value ApprovalRecord
	if err := decodeStrict(data, &value, "approval record", correlationID); err != nil {
		return ApprovalRecord{}, err
	}
	if err := value.Validate(correlationID); err != nil {
		return ApprovalRecord{}, err
	}
	return value, nil
}

func DecodeApprovalRequest(data []byte, correlationID string) (ApprovalRequest, error) {
	var value ApprovalRequest
	if err := decodeStrict(data, &value, "approval request", correlationID); err != nil {
		return ApprovalRequest{}, err
	}
	if err := value.VerifyDigest(correlationID); err != nil {
		return ApprovalRequest{}, err
	}
	return value, nil
}

func DecodeApprovalRecordV2(data []byte, correlationID string) (ApprovalRecordV2, error) {
	var value ApprovalRecordV2
	if err := decodeStrict(data, &value, "approval record v2", correlationID); err != nil {
		return ApprovalRecordV2{}, err
	}
	if err := value.VerifyIntegrity(correlationID); err != nil {
		return ApprovalRecordV2{}, err
	}
	return value, nil
}

func DecodeEvidenceRecord(data []byte, correlationID string) (EvidenceRecord, error) {
	var value EvidenceRecord
	if err := decodeStrict(data, &value, "evidence record", correlationID); err != nil {
		return EvidenceRecord{}, err
	}
	if err := value.Validate(correlationID); err != nil {
		return EvidenceRecord{}, err
	}
	return value, nil
}

func decodeStrict(data []byte, destination any, label, correlationID string) error {
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(destination); err != nil {
		return failure.New(failure.InvalidRequest, label+" is invalid", correlationID, failure.RetryNone)
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		return failure.New(failure.InvalidRequest, label+" contains trailing data", correlationID, failure.RetryNone)
	}
	return nil
}
