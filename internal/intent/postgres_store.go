package intent

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

type PostgresStore struct{ db *sql.DB }

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	if db == nil {
		return nil, errors.New("intent PostgreSQL connection is required")
	}
	return &PostgresStore{db: db}, nil
}

func (store *PostgresStore) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_intent;
CREATE TABLE IF NOT EXISTS forge_intent.schema_version(version integer PRIMARY KEY CHECK(version=1));
INSERT INTO forge_intent.schema_version(version) VALUES(1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_intent.records(
 tenant_id text NOT NULL,idempotency_key text NOT NULL,fingerprint text NOT NULL,record_id text NOT NULL,payload jsonb NOT NULL,
 PRIMARY KEY(tenant_id,idempotency_key),UNIQUE(record_id));`
	_, err := store.db.ExecContext(ctx, migration)
	return err
}

func (store *PostgresStore) Create(ctx context.Context, tenantID, key, fingerprint string, record Record) (Record, bool, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return Record{}, false, err
	}
	defer tx.Rollback()
	var existingFingerprint string
	var payload []byte
	err = tx.QueryRowContext(ctx, `
SELECT fingerprint,payload FROM forge_intent.records
WHERE tenant_id=$1 AND idempotency_key=$2 FOR UPDATE`, tenantID, key).Scan(&existingFingerprint, &payload)
	if err == nil {
		if existingFingerprint != fingerprint || json.Unmarshal(payload, &record) != nil {
			return Record{}, false, ErrConflict
		}
		return record, true, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return Record{}, false, err
	}
	payload, err = json.Marshal(record)
	if err != nil {
		return Record{}, false, err
	}
	if _, err = tx.ExecContext(ctx, `
INSERT INTO forge_intent.records(tenant_id,idempotency_key,fingerprint,record_id,payload)
VALUES($1,$2,$3,$4,$5)`, tenantID, key, fingerprint, record.ID, payload); err != nil {
		return Record{}, false, ErrConflict
	}
	return record, false, tx.Commit()
}

var _ Store = (*PostgresStore)(nil)
