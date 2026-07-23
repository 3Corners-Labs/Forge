package approval

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/3corners/forge/internal/contracts/trust"
)

type PostgresStore struct{ db *sql.DB }

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	if db == nil {
		return nil, errors.New("approval PostgreSQL connection is required")
	}
	return &PostgresStore{db: db}, nil
}

func (store *PostgresStore) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_approval;
CREATE TABLE IF NOT EXISTS forge_approval.schema_version(version integer PRIMARY KEY CHECK(version=1));
INSERT INTO forge_approval.schema_version(version) VALUES(1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_approval.requests(
 tenant_id text NOT NULL, request_id text NOT NULL, payload jsonb NOT NULL,
 PRIMARY KEY(tenant_id,request_id));
CREATE TABLE IF NOT EXISTS forge_approval.request_keys(
 tenant_id text NOT NULL, idempotency_key text NOT NULL, fingerprint text NOT NULL, request_id text NOT NULL,
 PRIMARY KEY(tenant_id,idempotency_key));
CREATE TABLE IF NOT EXISTS forge_approval.records(
 tenant_id text NOT NULL, record_id text NOT NULL, request_id text NOT NULL, previous_record_id text, payload jsonb NOT NULL,
 PRIMARY KEY(tenant_id,record_id));
CREATE TABLE IF NOT EXISTS forge_approval.current_records(
 tenant_id text NOT NULL, request_id text NOT NULL, record_id text NOT NULL,
 PRIMARY KEY(tenant_id,request_id));
CREATE TABLE IF NOT EXISTS forge_approval.decision_keys(
 tenant_id text NOT NULL, idempotency_key text NOT NULL, fingerprint text NOT NULL, record_id text NOT NULL,
 PRIMARY KEY(tenant_id,idempotency_key));`
	_, err := store.db.ExecContext(ctx, migration)
	return err
}

func (store *PostgresStore) CreateRequest(ctx context.Context, tenantID, key string, value trust.ApprovalRequest) (trust.ApprovalRequest, bool, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return trust.ApprovalRequest{}, false, ErrStoreUnavailable
	}
	defer tx.Rollback()
	var fingerprint, requestID string
	err = tx.QueryRowContext(ctx, `SELECT fingerprint,request_id FROM forge_approval.request_keys WHERE tenant_id=$1 AND idempotency_key=$2 FOR UPDATE`, tenantID, key).Scan(&fingerprint, &requestID)
	if err == nil {
		if fingerprint != value.RequestDigest {
			return trust.ApprovalRequest{}, false, ErrIdempotencyConflict
		}
		existing, found, resolveErr := resolveApprovalRequest(ctx, tx, tenantID, requestID)
		if resolveErr != nil || !found {
			return trust.ApprovalRequest{}, false, ErrStoreUnavailable
		}
		return existing, true, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return trust.ApprovalRequest{}, false, ErrStoreUnavailable
	}
	payload, _ := json.Marshal(value)
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_approval.requests VALUES($1,$2,$3)`, tenantID, value.ID, payload); err != nil {
		return trust.ApprovalRequest{}, false, ErrDecisionConflict
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_approval.request_keys VALUES($1,$2,$3,$4)`, tenantID, key, value.RequestDigest, value.ID); err != nil {
		return trust.ApprovalRequest{}, false, ErrIdempotencyConflict
	}
	if err = tx.Commit(); err != nil {
		return trust.ApprovalRequest{}, false, ErrStoreUnavailable
	}
	return cloneRequest(value), false, nil
}

func (store *PostgresStore) ResolveRequest(ctx context.Context, tenantID, requestID string) (trust.ApprovalRequest, bool, error) {
	return resolveApprovalRequest(ctx, store.db, tenantID, requestID)
}

type queryRower interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func resolveApprovalRequest(ctx context.Context, query queryRower, tenantID, requestID string) (trust.ApprovalRequest, bool, error) {
	var payload []byte
	err := query.QueryRowContext(ctx, `SELECT payload FROM forge_approval.requests WHERE tenant_id=$1 AND request_id=$2`, tenantID, requestID).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return trust.ApprovalRequest{}, false, nil
	}
	if err != nil {
		return trust.ApprovalRequest{}, false, ErrStoreUnavailable
	}
	var value trust.ApprovalRequest
	if json.Unmarshal(payload, &value) != nil {
		return trust.ApprovalRequest{}, false, ErrStoreUnavailable
	}
	return cloneRequest(value), true, nil
}

func (store *PostgresStore) CreateDecision(ctx context.Context, tenantID, key, requestID string, value trust.ApprovalRecordV2) (trust.ApprovalRecordV2, bool, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return trust.ApprovalRecordV2{}, false, ErrStoreUnavailable
	}
	defer tx.Rollback()
	var fingerprint, recordID string
	err = tx.QueryRowContext(ctx, `SELECT fingerprint,record_id FROM forge_approval.decision_keys WHERE tenant_id=$1 AND idempotency_key=$2 FOR UPDATE`, tenantID, key).Scan(&fingerprint, &recordID)
	if err == nil {
		if fingerprint != value.Integrity.Digest {
			return trust.ApprovalRecordV2{}, false, ErrIdempotencyConflict
		}
		existing, found, resolveErr := resolveCurrentRecord(ctx, tx, tenantID, recordID)
		if resolveErr != nil || !found {
			return trust.ApprovalRecordV2{}, false, ErrStoreUnavailable
		}
		return existing, true, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return trust.ApprovalRecordV2{}, false, ErrStoreUnavailable
	}
	var currentID string
	currentErr := tx.QueryRowContext(ctx, `SELECT record_id FROM forge_approval.current_records WHERE tenant_id=$1 AND request_id=$2 FOR UPDATE`, tenantID, requestID).Scan(&currentID)
	if currentErr == nil {
		current, found, resolveErr := resolveRecord(ctx, tx, tenantID, currentID)
		if resolveErr != nil || !found || value.Decision != "revoked" ||
			value.PreviousRecordID != current.ID || current.Decision != "approved" {
			return trust.ApprovalRecordV2{}, false, ErrDecisionConflict
		}
	} else if !errors.Is(currentErr, sql.ErrNoRows) {
		return trust.ApprovalRecordV2{}, false, ErrStoreUnavailable
	}
	payload, _ := json.Marshal(value)
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_approval.records VALUES($1,$2,$3,$4,$5)`, tenantID, value.ID, requestID, nullable(value.PreviousRecordID), payload); err != nil {
		return trust.ApprovalRecordV2{}, false, ErrDecisionConflict
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_approval.current_records VALUES($1,$2,$3) ON CONFLICT(tenant_id,request_id) DO UPDATE SET record_id=EXCLUDED.record_id`, tenantID, requestID, value.ID); err != nil {
		return trust.ApprovalRecordV2{}, false, ErrDecisionConflict
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_approval.decision_keys VALUES($1,$2,$3,$4)`, tenantID, key, value.Integrity.Digest, value.ID); err != nil {
		return trust.ApprovalRecordV2{}, false, ErrIdempotencyConflict
	}
	if err = tx.Commit(); err != nil {
		return trust.ApprovalRecordV2{}, false, ErrStoreUnavailable
	}
	return cloneRecord(value), false, nil
}

func (store *PostgresStore) ResolveCurrent(ctx context.Context, tenantID, recordID string) (trust.ApprovalRecordV2, bool, error) {
	return resolveCurrentRecord(ctx, store.db, tenantID, recordID)
}

func resolveCurrentRecord(ctx context.Context, query queryRower, tenantID, recordID string) (trust.ApprovalRecordV2, bool, error) {
	value, found, err := resolveRecord(ctx, query, tenantID, recordID)
	if err != nil || !found {
		return trust.ApprovalRecordV2{}, found, err
	}
	var current string
	if err := query.QueryRowContext(ctx, `SELECT record_id FROM forge_approval.current_records WHERE tenant_id=$1 AND request_id=$2`, tenantID, value.RequestID).Scan(&current); err != nil {
		return trust.ApprovalRecordV2{}, false, ErrStoreUnavailable
	}
	if current != recordID {
		return trust.ApprovalRecordV2{}, false, nil
	}
	return value, true, nil
}

func resolveRecord(ctx context.Context, query queryRower, tenantID, recordID string) (trust.ApprovalRecordV2, bool, error) {
	var payload []byte
	err := query.QueryRowContext(ctx, `SELECT payload FROM forge_approval.records WHERE tenant_id=$1 AND record_id=$2`, tenantID, recordID).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return trust.ApprovalRecordV2{}, false, nil
	}
	if err != nil {
		return trust.ApprovalRecordV2{}, false, ErrStoreUnavailable
	}
	var value trust.ApprovalRecordV2
	if json.Unmarshal(payload, &value) != nil {
		return trust.ApprovalRecordV2{}, false, ErrStoreUnavailable
	}
	return cloneRecord(value), true, nil
}

func nullable(value string) any {
	if value == "" {
		return nil
	}
	return value
}
