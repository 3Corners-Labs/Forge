package context

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
		return nil, errors.New("context PostgreSQL connection is required")
	}
	return &PostgresStore{db: db}, nil
}

func (store *PostgresStore) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_context;
CREATE TABLE IF NOT EXISTS forge_context.schema_version(version integer PRIMARY KEY CHECK (version = 1));
INSERT INTO forge_context.schema_version(version) VALUES (1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_context.packages(
 tenant_id text NOT NULL, package_id text NOT NULL, payload jsonb NOT NULL,
 PRIMARY KEY(tenant_id, package_id));
CREATE TABLE IF NOT EXISTS forge_context.idempotency(
 tenant_id text NOT NULL, idempotency_key text NOT NULL, fingerprint text NOT NULL,
 package_id text NOT NULL, PRIMARY KEY(tenant_id, idempotency_key));`
	_, err := store.db.ExecContext(ctx, migration)
	return err
}

func (store *PostgresStore) Get(ctx context.Context, tenantID, key, fingerprint string) (trust.ContextPackage, bool, error) {
	var storedFingerprint string
	var payload []byte
	err := store.db.QueryRowContext(ctx, `
SELECT i.fingerprint, p.payload
FROM forge_context.idempotency i
JOIN forge_context.packages p ON p.tenant_id=i.tenant_id AND p.package_id=i.package_id
WHERE i.tenant_id=$1 AND i.idempotency_key=$2`, tenantID, key).Scan(&storedFingerprint, &payload)
	if errors.Is(err, sql.ErrNoRows) {
		return trust.ContextPackage{}, false, nil
	}
	if err != nil {
		return trust.ContextPackage{}, false, ErrContextStoreUnavailable
	}
	if storedFingerprint != fingerprint {
		return trust.ContextPackage{}, false, ErrContextConflict
	}
	var value trust.ContextPackage
	if json.Unmarshal(payload, &value) != nil {
		return trust.ContextPackage{}, false, ErrContextStoreUnavailable
	}
	return clonePackage(value), true, nil
}

func (store *PostgresStore) Create(ctx context.Context, value trust.ContextPackage, tenantID, key, fingerprint string) (trust.ContextPackage, bool, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return trust.ContextPackage{}, false, ErrContextStoreUnavailable
	}
	defer tx.Rollback()
	var existingFingerprint, packageID string
	err = tx.QueryRowContext(ctx, `SELECT fingerprint, package_id FROM forge_context.idempotency WHERE tenant_id=$1 AND idempotency_key=$2 FOR UPDATE`, tenantID, key).Scan(&existingFingerprint, &packageID)
	if err == nil {
		if existingFingerprint != fingerprint {
			return trust.ContextPackage{}, false, ErrContextConflict
		}
		var payload []byte
		if tx.QueryRowContext(ctx, `SELECT payload FROM forge_context.packages WHERE tenant_id=$1 AND package_id=$2`, tenantID, packageID).Scan(&payload) != nil {
			return trust.ContextPackage{}, false, ErrContextStoreUnavailable
		}
		var existing trust.ContextPackage
		if json.Unmarshal(payload, &existing) != nil {
			return trust.ContextPackage{}, false, ErrContextStoreUnavailable
		}
		return clonePackage(existing), true, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return trust.ContextPackage{}, false, ErrContextStoreUnavailable
	}
	payload, _ := json.Marshal(value)
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_context.packages(tenant_id,package_id,payload) VALUES($1,$2,$3)`, tenantID, value.ID, payload); err != nil {
		return trust.ContextPackage{}, false, ErrContextConflict
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_context.idempotency(tenant_id,idempotency_key,fingerprint,package_id) VALUES($1,$2,$3,$4)`, tenantID, key, fingerprint, value.ID); err != nil {
		return trust.ContextPackage{}, false, ErrContextConflict
	}
	if err = tx.Commit(); err != nil {
		return trust.ContextPackage{}, false, ErrContextStoreUnavailable
	}
	return clonePackage(value), false, nil
}

func (store *PostgresStore) Resolve(ctx context.Context, tenantID, packageID string) (trust.ContextPackage, bool, error) {
	var payload []byte
	err := store.db.QueryRowContext(ctx, `SELECT payload FROM forge_context.packages WHERE tenant_id=$1 AND package_id=$2`, tenantID, packageID).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return trust.ContextPackage{}, false, nil
	}
	if err != nil {
		return trust.ContextPackage{}, false, ErrContextStoreUnavailable
	}
	var value trust.ContextPackage
	if json.Unmarshal(payload, &value) != nil {
		return trust.ContextPackage{}, false, ErrContextStoreUnavailable
	}
	return clonePackage(value), true, nil
}
