// Package evidence owns append-oriented, integrity-linked M1 evidence. It is
// separate from telemetry and stores only minimized references and digests.
package evidence

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"slices"
	"sync"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/contracts/trust"
)

var (
	ErrUnavailable = errors.New("evidence store unavailable")
	ErrConflict    = errors.New("evidence chain conflict")
)

type AppendCommand struct {
	StreamKey        string
	IdempotencyKey   string
	Fingerprint      string
	ExpectedPrevious string
	Record           trust.EvidenceRecord
}

type Store interface {
	Head(context.Context, string) (string, error)
	Append(context.Context, AppendCommand) (trust.EvidenceRecord, bool, error)
	Stream(context.Context, string) ([]trust.EvidenceRecord, error)
}

type Service struct {
	store Store
	now   func() time.Time
}

func NewService(store Store, now func() time.Time) (*Service, error) {
	if store == nil || now == nil {
		return nil, errors.New("evidence service configuration is invalid")
	}
	return &Service{store: store, now: now}, nil
}

type AppendRequest struct {
	CorrelationID  string
	IdempotencyKey string
	Type           string
	ActorID        string
	Scope          trust.EvidenceScope
	Subject        trust.EvidenceSubject
	References     []string
	Outcome        string
	Classification string
	RetentionClass string
}

func (service *Service) Append(ctx context.Context, request AppendRequest) (trust.EvidenceRecord, bool, error) {
	if ctx == nil || request.CorrelationID == "" || request.IdempotencyKey == "" {
		return trust.EvidenceRecord{}, false, failure.New(failure.InvalidRequest, "evidence append request is invalid", request.CorrelationID, failure.RetryNone)
	}
	at := service.now()
	if at.IsZero() {
		return trust.EvidenceRecord{}, false, failure.New(failure.DependencyUnavailable, "evidence clock is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	streamKey := request.Scope.TenantID + "\x00" + request.Scope.ProjectID + "\x00" + request.Scope.DeliveryRunID
	previous, err := service.store.Head(ctx, streamKey)
	if err != nil {
		return trust.EvidenceRecord{}, false, failure.New(failure.DependencyUnavailable, "evidence chain head is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	record, err := trust.BindEvidenceRecordIntegrity(trust.EvidenceRecord{
		SchemaVersion: trust.EvidenceRecordSchemaVersion, Type: request.Type, RecordedAt: at,
		ActorID: request.ActorID, Scope: request.Scope, Subject: request.Subject,
		References: slices.Clone(request.References), Outcome: request.Outcome,
		Classification: request.Classification, RetentionClass: request.RetentionClass,
	}, previous, request.CorrelationID)
	if err != nil {
		return trust.EvidenceRecord{}, false, err
	}
	stored, replayed, err := service.store.Append(ctx, AppendCommand{
		StreamKey: streamKey, IdempotencyKey: request.IdempotencyKey,
		Fingerprint: evidenceRequestFingerprint(request), ExpectedPrevious: previous, Record: record,
	})
	if errors.Is(err, ErrConflict) {
		return trust.EvidenceRecord{}, false, failure.New(failure.Conflict, "evidence chain changed", request.CorrelationID, failure.RetryReconcile)
	}
	if err != nil {
		return trust.EvidenceRecord{}, false, failure.New(failure.DependencyUnavailable, "evidence append is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if err := stored.VerifyIntegrity(request.CorrelationID); err != nil {
		return trust.EvidenceRecord{}, false, err
	}
	return stored, replayed, nil
}

func evidenceRequestFingerprint(request AppendRequest) string {
	value := request
	value.CorrelationID = ""
	value.IdempotencyKey = ""
	encoded, _ := json.Marshal(value)
	sum := sha256.Sum256(encoded)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func (service *Service) VerifyStream(ctx context.Context, scope trust.EvidenceScope, correlationID string) error {
	if ctx == nil || scope.TenantID == "" || scope.ProjectID == "" || correlationID == "" {
		return failure.New(failure.InvalidRequest, "evidence verification request is invalid", correlationID, failure.RetryNone)
	}
	streamKey := scope.TenantID + "\x00" + scope.ProjectID + "\x00" + scope.DeliveryRunID
	records, err := service.store.Stream(ctx, streamKey)
	if err != nil {
		return failure.New(failure.DependencyUnavailable, "evidence stream is unavailable", correlationID, failure.RetryBounded)
	}
	previous := ""
	for _, record := range records {
		if err := record.VerifyIntegrity(correlationID); err != nil ||
			record.Integrity.PreviousRecordDigest != previous ||
			record.Scope.TenantID != scope.TenantID || record.Scope.ProjectID != scope.ProjectID ||
			record.Scope.DeliveryRunID != scope.DeliveryRunID {
			return failure.New(failure.VerificationFailed, "evidence chain verification failed", correlationID, failure.RetryNone)
		}
		previous = record.Integrity.Digest
	}
	return nil
}

func (service *Service) VerifySubject(ctx context.Context, scope trust.EvidenceScope, subject trust.EvidenceSubject, requiredReferences []string, correlationID string) error {
	if ctx == nil || subject.Type == "" || subject.ID == "" || subject.Digest == "" ||
		len(requiredReferences) == 0 || correlationID == "" {
		return failure.New(failure.InvalidRequest, "evidence subject verification request is invalid", correlationID, failure.RetryNone)
	}
	if err := service.VerifyStream(ctx, scope, correlationID); err != nil {
		return err
	}
	streamKey := scope.TenantID + "\x00" + scope.ProjectID + "\x00" + scope.DeliveryRunID
	records, err := service.store.Stream(ctx, streamKey)
	if err != nil {
		return failure.New(failure.DependencyUnavailable, "evidence stream is unavailable", correlationID, failure.RetryBounded)
	}
	for _, record := range records {
		if record.Subject == subject && containsAll(record.References, requiredReferences) {
			return nil
		}
	}
	return failure.New(failure.EvidenceUnavailable, "required evidence subject is unavailable", correlationID, failure.RetryReconcile)
}

func containsAll(values, required []string) bool {
	for _, expected := range required {
		if !slices.Contains(values, expected) {
			return false
		}
	}
	return true
}

type memoryEntry struct {
	record      trust.EvidenceRecord
	fingerprint string
}

type MemoryStore struct {
	mu          sync.RWMutex
	streams     map[string][]trust.EvidenceRecord
	idempotency map[string]memoryEntry
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{streams: make(map[string][]trust.EvidenceRecord), idempotency: make(map[string]memoryEntry)}
}

func (store *MemoryStore) Head(_ context.Context, stream string) (string, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	values := store.streams[stream]
	if len(values) == 0 {
		return "", nil
	}
	return values[len(values)-1].Integrity.Digest, nil
}

func (store *MemoryStore) Append(_ context.Context, command AppendCommand) (trust.EvidenceRecord, bool, error) {
	store.mu.Lock()
	defer store.mu.Unlock()
	key := command.StreamKey + "\x00" + command.IdempotencyKey
	if existing, found := store.idempotency[key]; found {
		if existing.fingerprint != command.Fingerprint {
			return trust.EvidenceRecord{}, false, ErrConflict
		}
		return cloneRecord(existing.record), true, nil
	}
	values := store.streams[command.StreamKey]
	current := ""
	if len(values) > 0 {
		current = values[len(values)-1].Integrity.Digest
	}
	if current != command.ExpectedPrevious || command.Record.Integrity.PreviousRecordDigest != current {
		return trust.EvidenceRecord{}, false, ErrConflict
	}
	record := cloneRecord(command.Record)
	store.streams[command.StreamKey] = append(values, record)
	store.idempotency[key] = memoryEntry{record: record, fingerprint: command.Fingerprint}
	return cloneRecord(record), false, nil
}

func (store *MemoryStore) Stream(_ context.Context, stream string) ([]trust.EvidenceRecord, error) {
	store.mu.RLock()
	defer store.mu.RUnlock()
	values := store.streams[stream]
	result := make([]trust.EvidenceRecord, len(values))
	for index := range values {
		result[index] = cloneRecord(values[index])
	}
	return result, nil
}

func cloneRecord(value trust.EvidenceRecord) trust.EvidenceRecord {
	value.References = slices.Clone(value.References)
	return value
}
