package policy

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	"github.com/3corners/forge/internal/contracts/trust"
)

type PostgresEnvelopeStore struct{ db *sql.DB }

func NewPostgresEnvelopeStore(db *sql.DB) (*PostgresEnvelopeStore, error) {
	if db == nil {
		return nil, errors.New("policy PostgreSQL connection is required")
	}
	return &PostgresEnvelopeStore{db: db}, nil
}

func (store *PostgresEnvelopeStore) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_policy;
CREATE TABLE IF NOT EXISTS forge_policy.schema_version(version integer PRIMARY KEY CHECK(version=1));
INSERT INTO forge_policy.schema_version(version) VALUES(1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_policy.envelopes(
 envelope_id text PRIMARY KEY,tenant_id text NOT NULL,payload jsonb NOT NULL);
CREATE TABLE IF NOT EXISTS forge_policy.envelope_keys(
 tenant_id text NOT NULL,idempotency_key text NOT NULL,fingerprint text NOT NULL,envelope_id text NOT NULL,
 PRIMARY KEY(tenant_id,idempotency_key));`
	_, err := store.db.ExecContext(ctx, migration)
	return err
}

func (store *PostgresEnvelopeStore) Create(ctx context.Context, value trust.ExecutionEnvelope, tenantID, key, fingerprint string) (trust.ExecutionEnvelope, bool, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return trust.ExecutionEnvelope{}, false, ErrEnvelopeStoreUnavailable
	}
	defer tx.Rollback()
	var existingFingerprint, id string
	err = tx.QueryRowContext(ctx, `SELECT fingerprint,envelope_id FROM forge_policy.envelope_keys WHERE tenant_id=$1 AND idempotency_key=$2 FOR UPDATE`, tenantID, key).Scan(&existingFingerprint, &id)
	if err == nil {
		if existingFingerprint != fingerprint {
			return trust.ExecutionEnvelope{}, false, ErrEnvelopeConflict
		}
		existing, found, loadErr := loadEnvelope(ctx, tx, id)
		if loadErr != nil || !found {
			return trust.ExecutionEnvelope{}, false, ErrEnvelopeStoreUnavailable
		}
		return existing, true, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return trust.ExecutionEnvelope{}, false, ErrEnvelopeStoreUnavailable
	}
	payload, _ := json.Marshal(value)
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_policy.envelopes VALUES($1,$2,$3)`, value.ID, tenantID, payload); err != nil {
		return trust.ExecutionEnvelope{}, false, ErrEnvelopeConflict
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_policy.envelope_keys VALUES($1,$2,$3,$4)`, tenantID, key, fingerprint, value.ID); err != nil {
		return trust.ExecutionEnvelope{}, false, ErrEnvelopeConflict
	}
	if err = tx.Commit(); err != nil {
		return trust.ExecutionEnvelope{}, false, ErrEnvelopeStoreUnavailable
	}
	return cloneEnvelope(value), false, nil
}

type envelopeQuery interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func loadEnvelope(ctx context.Context, query envelopeQuery, id string) (trust.ExecutionEnvelope, bool, error) {
	var payload []byte
	err := query.QueryRowContext(ctx, `SELECT payload FROM forge_policy.envelopes WHERE envelope_id=$1`, id).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return trust.ExecutionEnvelope{}, false, nil
	}
	if err != nil {
		return trust.ExecutionEnvelope{}, false, ErrEnvelopeStoreUnavailable
	}
	var value trust.ExecutionEnvelope
	if json.Unmarshal(payload, &value) != nil {
		return trust.ExecutionEnvelope{}, false, ErrEnvelopeStoreUnavailable
	}
	return cloneEnvelope(value), true, nil
}

func (store *PostgresEnvelopeStore) ResolveEnvelope(ctx context.Context, id string) (trust.ExecutionEnvelope, bool, error) {
	return loadEnvelope(ctx, store.db, id)
}
