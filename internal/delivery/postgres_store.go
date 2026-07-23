package delivery

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
)

// PostgresStore is the Developer Trial durable Delivery adapter. State,
// optimistic version, idempotency command result, and outbox event commit in
// one capability-local transaction.
type PostgresStore struct {
	db    *sql.DB
	newID func() string
}

func NewPostgresStore(db *sql.DB, newID func() string) (*PostgresStore, error) {
	if db == nil || newID == nil {
		return nil, errors.New("delivery PostgreSQL configuration is invalid")
	}
	return &PostgresStore{db: db, newID: newID}, nil
}

func (store *PostgresStore) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_delivery;
CREATE TABLE IF NOT EXISTS forge_delivery.schema_version(version integer PRIMARY KEY CHECK(version=1));
INSERT INTO forge_delivery.schema_version(version) VALUES(1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_delivery.runs(
 run_id text PRIMARY KEY, scope_key text NOT NULL, state text NOT NULL, version bigint NOT NULL, payload jsonb NOT NULL);
CREATE TABLE IF NOT EXISTS forge_delivery.capture_keys(
 scope_key text NOT NULL, idempotency_key text NOT NULL, fingerprint text NOT NULL, run_id text NOT NULL,
 PRIMARY KEY(scope_key,idempotency_key));
CREATE TABLE IF NOT EXISTS forge_delivery.transition_keys(
 scope_key text NOT NULL, idempotency_key text NOT NULL, fingerprint text NOT NULL, run_id text NOT NULL, result jsonb NOT NULL,
 PRIMARY KEY(scope_key,idempotency_key));
CREATE TABLE IF NOT EXISTS forge_delivery.outbox(
 event_id text PRIMARY KEY, run_id text NOT NULL, run_version bigint NOT NULL, payload jsonb NOT NULL,
 delivered_at timestamptz);`
	_, err := store.db.ExecContext(ctx, migration)
	return err
}

func (store *PostgresStore) CreateCaptured(ctx context.Context, record CaptureRecord) (StoredRecord, bool, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	defer tx.Rollback()
	var fingerprint, runID string
	scopeKey := deliveryDatabaseKey(record.ScopeKey)
	err = tx.QueryRowContext(ctx, `SELECT fingerprint,run_id FROM forge_delivery.capture_keys WHERE scope_key=$1 AND idempotency_key=$2 FOR UPDATE`, scopeKey, record.IdempotencyKey).Scan(&fingerprint, &runID)
	if err == nil {
		if fingerprint != record.RequestFingerprint {
			return StoredRecord{}, false, ErrIdempotencyConflict
		}
		value, found, loadErr := loadRun(ctx, tx, runID, false)
		if loadErr != nil || !found {
			return StoredRecord{}, false, ErrDeliveryUnavailable
		}
		return value, true, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	runID = store.newID()
	if runID == "" {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	value := StoredRecord{CaptureRecord: cloneCaptureRecord(record), ContractVersion: DeliveryRunContractVersion, ID: runID, State: CapturedState, Version: 1}
	payload, _ := json.Marshal(deliveryDatabaseRecord(value))
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_delivery.runs VALUES($1,$2,$3,$4,$5)`, runID, scopeKey, value.State, value.Version, payload); err != nil {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_delivery.capture_keys VALUES($1,$2,$3,$4)`, scopeKey, record.IdempotencyKey, record.RequestFingerprint, runID); err != nil {
		return StoredRecord{}, false, ErrIdempotencyConflict
	}
	if err = tx.Commit(); err != nil {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	return cloneStoredRecord(value), false, nil
}

func (store *PostgresStore) Load(ctx context.Context, runID string) (StoredRecord, bool, error) {
	return loadRun(ctx, store.db, runID, false)
}

type deliveryQuery interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func loadRun(ctx context.Context, query deliveryQuery, runID string, lock bool) (StoredRecord, bool, error) {
	statement := `SELECT payload FROM forge_delivery.runs WHERE run_id=$1`
	if lock {
		statement += ` FOR UPDATE`
	}
	var payload []byte
	err := query.QueryRowContext(ctx, statement, runID).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return StoredRecord{}, false, nil
	}
	if err != nil {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	var value StoredRecord
	if json.Unmarshal(payload, &value) != nil || restoreDeliveryScopeKey(&value) != nil {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	return cloneStoredRecord(value), true, nil
}

func (store *PostgresStore) Transition(ctx context.Context, command TransitionCommand) (StoredRecord, bool, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	defer tx.Rollback()
	var fingerprint string
	var result []byte
	scopeKey := deliveryDatabaseKey(command.ScopeKey)
	err = tx.QueryRowContext(ctx, `SELECT fingerprint,result FROM forge_delivery.transition_keys WHERE scope_key=$1 AND idempotency_key=$2 FOR UPDATE`, scopeKey, command.IdempotencyKey).Scan(&fingerprint, &result)
	if err == nil {
		if fingerprint != command.Fingerprint {
			return StoredRecord{}, false, ErrTransitionConflict
		}
		var value StoredRecord
		if json.Unmarshal(result, &value) != nil || restoreDeliveryScopeKey(&value) != nil {
			return StoredRecord{}, false, ErrDeliveryUnavailable
		}
		return cloneStoredRecord(value), true, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	current, found, err := loadRun(ctx, tx, command.RunID, true)
	if err != nil || !found || current.ScopeKey != command.ScopeKey ||
		current.Version != command.ExpectedVersion || current.State != command.FromState {
		return StoredRecord{}, false, ErrTransitionConflict
	}
	next := cloneStoredRecord(current)
	next.State, next.Version = command.ToState, current.Version+1
	if command.Context != nil {
		next.Context = cloneContextRecord(command.Context)
	}
	if command.Plan != nil {
		next.Plan = clonePlan(command.Plan)
	}
	if command.Approval != nil {
		next.Approval = clonePlanApprovalRecord(command.Approval)
	}
	if command.Execution != nil {
		next.Execution = cloneExecutionRecord(command.Execution)
	}
	if command.Candidate != nil {
		next.Candidate = cloneCandidateRecord(command.Candidate)
	}
	if command.Verification != nil {
		next.Verification = cloneVerificationRecord(command.Verification)
	}
	event := OutboxEvent{
		ID: command.Fingerprint, RunID: command.RunID, FromState: command.FromState,
		ToState: command.ToState, RunVersion: next.Version,
		CorrelationID: command.CorrelationID, OccurredAt: command.OccurredAt,
	}
	next.Outbox = append(next.Outbox, event)
	payload, _ := json.Marshal(deliveryDatabaseRecord(next))
	update, err := tx.ExecContext(ctx, `UPDATE forge_delivery.runs SET state=$1,version=$2,payload=$3 WHERE run_id=$4 AND version=$5`, next.State, next.Version, payload, next.ID, current.Version)
	if err != nil {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	if affected, _ := update.RowsAffected(); affected != 1 {
		return StoredRecord{}, false, ErrTransitionConflict
	}
	eventPayload, _ := json.Marshal(event)
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_delivery.outbox VALUES($1,$2,$3,$4,NULL)`, event.ID, event.RunID, event.RunVersion, eventPayload); err != nil {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_delivery.transition_keys VALUES($1,$2,$3,$4,$5)`, scopeKey, command.IdempotencyKey, command.Fingerprint, command.RunID, payload); err != nil {
		return StoredRecord{}, false, ErrTransitionConflict
	}
	if err = tx.Commit(); err != nil {
		return StoredRecord{}, false, ErrDeliveryUnavailable
	}
	return cloneStoredRecord(next), false, nil
}

func deliveryDatabaseKey(value string) string {
	return "b64:" + base64.RawURLEncoding.EncodeToString([]byte(value))
}

func deliveryDatabaseRecord(value StoredRecord) StoredRecord {
	value = cloneStoredRecord(value)
	value.ScopeKey = deliveryDatabaseKey(value.ScopeKey)
	return value
}

func restoreDeliveryScopeKey(value *StoredRecord) error {
	const prefix = "b64:"
	if value == nil || len(value.ScopeKey) <= len(prefix) || value.ScopeKey[:len(prefix)] != prefix {
		return ErrDeliveryUnavailable
	}
	decoded, err := base64.RawURLEncoding.DecodeString(value.ScopeKey[len(prefix):])
	if err != nil {
		return ErrDeliveryUnavailable
	}
	value.ScopeKey = string(decoded)
	return nil
}

func (store *PostgresStore) PendingOutbox(ctx context.Context, limit int) ([]OutboxEvent, error) {
	if limit < 1 || limit > 1000 {
		return nil, errors.New("outbox limit is invalid")
	}
	rows, err := store.db.QueryContext(ctx, `SELECT payload FROM forge_delivery.outbox WHERE delivered_at IS NULL ORDER BY run_version,event_id LIMIT $1`, limit)
	if err != nil {
		return nil, ErrDeliveryUnavailable
	}
	defer rows.Close()
	var values []OutboxEvent
	for rows.Next() {
		var payload []byte
		var value OutboxEvent
		if rows.Scan(&payload) != nil || json.Unmarshal(payload, &value) != nil {
			return nil, ErrDeliveryUnavailable
		}
		values = append(values, value)
	}
	return values, rows.Err()
}

func (store *PostgresStore) MarkOutboxDelivered(ctx context.Context, eventID string) error {
	result, err := store.db.ExecContext(ctx, `UPDATE forge_delivery.outbox SET delivered_at=clock_timestamp() WHERE event_id=$1 AND delivered_at IS NULL`, eventID)
	if err != nil {
		return ErrDeliveryUnavailable
	}
	affected, _ := result.RowsAffected()
	if affected > 1 {
		return ErrDeliveryUnavailable
	}
	return nil
}
