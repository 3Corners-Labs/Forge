package delivery

import (
	"context"
	"errors"
	"sync"
)

// MemoryStore is the concurrency-safe Developer Trial adapter. It proves the
// atomic idempotency contract but intentionally provides no durability.
type MemoryStore struct {
	mu      sync.Mutex
	newID   func() string
	records map[string]StoredRecord
}

func NewMemoryStore(newID func() string) (*MemoryStore, error) {
	if newID == nil {
		return nil, errors.New("delivery run ID generator is required")
	}
	return &MemoryStore{newID: newID, records: make(map[string]StoredRecord)}, nil
}

func (store *MemoryStore) CreateCaptured(ctx context.Context, record CaptureRecord) (StoredRecord, bool, error) {
	if store == nil || store.newID == nil || store.records == nil {
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
	return cloneStoredRecord(stored), false, nil
}

func cloneStoredRecord(value StoredRecord) StoredRecord {
	value.CaptureRecord = cloneCaptureRecord(value.CaptureRecord)
	return value
}
