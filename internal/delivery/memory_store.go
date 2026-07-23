package delivery

import (
	"context"
	"errors"
	"slices"
	"sync"
)

// MemoryStore is the concurrency-safe Developer Trial adapter. It proves the
// atomic idempotency contract but intentionally provides no durability.
type MemoryStore struct {
	mu      sync.Mutex
	newID   func() string
	records map[string]StoredRecord
	byID    map[string]string
	changes map[string]transitionEntry
}

func NewMemoryStore(newID func() string) (*MemoryStore, error) {
	if newID == nil {
		return nil, errors.New("delivery run ID generator is required")
	}
	return &MemoryStore{
		newID: newID, records: make(map[string]StoredRecord),
		byID: make(map[string]string), changes: make(map[string]transitionEntry),
	}, nil
}

func (store *MemoryStore) CreateCaptured(ctx context.Context, record CaptureRecord) (StoredRecord, bool, error) {
	if store == nil || store.newID == nil || store.records == nil || store.byID == nil || store.changes == nil {
		return StoredRecord{}, false, errors.New("delivery memory store is unavailable")
	}
	if ctx == nil {
		return StoredRecord{}, false, errors.New("delivery store context is required")
	}
	if err := ctx.Err(); err != nil {
		return StoredRecord{}, false, err
	}
	key := record.ScopeKey + "\x00" + record.IdempotencyKey
	store.mu.Lock()
	defer store.mu.Unlock()
	if existing, found := store.records[key]; found {
		if existing.RequestFingerprint != record.RequestFingerprint {
			return StoredRecord{}, false, ErrIdempotencyConflict
		}
		return cloneStoredRecord(existing), true, nil
	}
	id := store.newID()
	if id == "" {
		return StoredRecord{}, false, errors.New("delivery run ID generator returned an empty value")
	}
	stored := StoredRecord{
		CaptureRecord:   cloneCaptureRecord(record),
		ContractVersion: DeliveryRunContractVersion,
		ID:              id,
		State:           CapturedState,
		Version:         1,
	}
	store.records[key] = cloneStoredRecord(stored)
	store.byID[id] = key
	return cloneStoredRecord(stored), false, nil
}

func cloneStoredRecord(value StoredRecord) StoredRecord {
	value.CaptureRecord = cloneCaptureRecord(value.CaptureRecord)
	value.Context = cloneContextRecord(value.Context)
	value.Plan = clonePlan(value.Plan)
	value.Approval = clonePlanApprovalRecord(value.Approval)
	value.Execution = cloneExecutionRecord(value.Execution)
	value.Candidate = cloneCandidateRecord(value.Candidate)
	value.Verification = cloneVerificationRecord(value.Verification)
	value.Outbox = slices.Clone(value.Outbox)
	return value
}
