package identity

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type PostgresWorkerDelegationStore struct{ db *sql.DB }

func NewPostgresWorkerDelegationStore(db *sql.DB) (*PostgresWorkerDelegationStore, error) {
	if db == nil {
		return nil, errors.New("identity PostgreSQL connection is required")
	}
	return &PostgresWorkerDelegationStore{db: db}, nil
}

func (store *PostgresWorkerDelegationStore) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_identity;
CREATE TABLE IF NOT EXISTS forge_identity.schema_version(version integer PRIMARY KEY CHECK(version=1));
INSERT INTO forge_identity.schema_version(version) VALUES(1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_identity.worker_delegations(
 delegation_id text PRIMARY KEY, tenant_id text NOT NULL, request_digest text NOT NULL,
 payload jsonb NOT NULL, revoked_at timestamptz);
CREATE TABLE IF NOT EXISTS forge_identity.worker_delegation_keys(
 tenant_id text NOT NULL,idempotency_key text NOT NULL,request_digest text NOT NULL,delegation_id text NOT NULL,
 PRIMARY KEY(tenant_id,idempotency_key));`
	_, err := store.db.ExecContext(ctx, migration)
	return err
}

func (store *PostgresWorkerDelegationStore) CreateWorker(ctx context.Context, value WorkerDelegation, key string) (WorkerDelegation, bool, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return WorkerDelegation{}, false, ErrDelegationStoreUnavailable
	}
	defer tx.Rollback()
	var digest, id string
	err = tx.QueryRowContext(ctx, `SELECT request_digest,delegation_id FROM forge_identity.worker_delegation_keys WHERE tenant_id=$1 AND idempotency_key=$2 FOR UPDATE`, value.TenantID, key).Scan(&digest, &id)
	if err == nil {
		if digest != value.RequestDigest {
			return WorkerDelegation{}, false, ErrDelegationConflict
		}
		stored, found, loadErr := getWorker(ctx, tx, id)
		if loadErr != nil || !found || stored.Revoked {
			return WorkerDelegation{}, false, ErrDelegationRevoked
		}
		return stored.Delegation, true, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return WorkerDelegation{}, false, ErrDelegationStoreUnavailable
	}
	payload, _ := json.Marshal(value)
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_identity.worker_delegations VALUES($1,$2,$3,$4,NULL)`, value.ID, value.TenantID, value.RequestDigest, payload); err != nil {
		return WorkerDelegation{}, false, ErrDelegationConflict
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_identity.worker_delegation_keys VALUES($1,$2,$3,$4)`, value.TenantID, key, value.RequestDigest, value.ID); err != nil {
		return WorkerDelegation{}, false, ErrDelegationConflict
	}
	if err = tx.Commit(); err != nil {
		return WorkerDelegation{}, false, ErrDelegationStoreUnavailable
	}
	return cloneWorkerDelegation(value), false, nil
}

type workerQuery interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func getWorker(ctx context.Context, query workerQuery, id string) (StoredWorkerDelegation, bool, error) {
	var payload []byte
	var revoked sql.NullTime
	err := query.QueryRowContext(ctx, `SELECT payload,revoked_at FROM forge_identity.worker_delegations WHERE delegation_id=$1`, id).Scan(&payload, &revoked)
	if errors.Is(err, sql.ErrNoRows) {
		return StoredWorkerDelegation{}, false, nil
	}
	if err != nil {
		return StoredWorkerDelegation{}, false, ErrDelegationStoreUnavailable
	}
	var value WorkerDelegation
	if json.Unmarshal(payload, &value) != nil {
		return StoredWorkerDelegation{}, false, ErrDelegationStoreUnavailable
	}
	return StoredWorkerDelegation{Delegation: cloneWorkerDelegation(value), Revoked: revoked.Valid, RevokedAt: revoked.Time}, true, nil
}

func (store *PostgresWorkerDelegationStore) GetWorker(ctx context.Context, id string) (StoredWorkerDelegation, bool, error) {
	return getWorker(ctx, store.db, id)
}

func (store *PostgresWorkerDelegationStore) RevokeWorker(ctx context.Context, id string, at time.Time) (StoredWorkerDelegation, bool, error) {
	result, err := store.db.ExecContext(ctx, `UPDATE forge_identity.worker_delegations SET revoked_at=COALESCE(revoked_at,$2) WHERE delegation_id=$1`, id, at)
	if err != nil {
		return StoredWorkerDelegation{}, false, ErrDelegationStoreUnavailable
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return StoredWorkerDelegation{}, false, nil
	}
	return getWorker(ctx, store.db, id)
}
