package verification

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	verificationcontract "github.com/3corners/forge/internal/contracts/verification"
)

var ErrStoreUnavailable = errors.New("verification store is unavailable")

// PostgresStore durably owns Assurance Plans and independent verification
// results. Digest conflicts fail closed; identical writes are idempotent.
type PostgresStore struct{ db *sql.DB }

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	if db == nil {
		return nil, errors.New("verification PostgreSQL connection is required")
	}
	return &PostgresStore{db: db}, nil
}

func (store *PostgresStore) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_verification;
CREATE TABLE IF NOT EXISTS forge_verification.schema_version(version integer PRIMARY KEY CHECK(version=1));
INSERT INTO forge_verification.schema_version(version) VALUES(1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_verification.assurance_plans(
 plan_id text PRIMARY KEY,digest text NOT NULL UNIQUE,candidate_digest text NOT NULL,payload jsonb NOT NULL);
CREATE TABLE IF NOT EXISTS forge_verification.results(
 result_id text PRIMARY KEY,digest text NOT NULL UNIQUE,plan_id text NOT NULL,candidate_digest text NOT NULL,payload jsonb NOT NULL);`
	_, err := store.db.ExecContext(ctx, migration)
	return err
}

func (store *PostgresStore) SavePlan(ctx context.Context, value verificationcontract.AssurancePlan) error {
	if err := value.Validate(value.ExpiresAt.Add(-1)); err != nil {
		return err
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return ErrStoreUnavailable
	}
	result, err := store.db.ExecContext(ctx, `
INSERT INTO forge_verification.assurance_plans(plan_id,digest,candidate_digest,payload)
VALUES($1,$2,$3,$4) ON CONFLICT(plan_id) DO UPDATE SET payload=forge_verification.assurance_plans.payload
WHERE forge_verification.assurance_plans.digest=EXCLUDED.digest`,
		value.ID, value.Digest, value.CandidateDigest, payload)
	if err != nil {
		return ErrStoreUnavailable
	}
	affected, err := result.RowsAffected()
	if err != nil || affected != 1 {
		return ErrStoreUnavailable
	}
	return nil
}

func (store *PostgresStore) SaveResult(ctx context.Context, value verificationcontract.Result) error {
	bound, bindErr := verificationcontract.BindResult(value)
	payload, err := json.Marshal(value)
	if bindErr != nil || bound.ID != value.ID || bound.Digest != value.Digest ||
		err != nil || value.AssurancePlanID == "" {
		return verificationcontract.ErrInvalidContract
	}
	result, err := store.db.ExecContext(ctx, `
INSERT INTO forge_verification.results(result_id,digest,plan_id,candidate_digest,payload)
VALUES($1,$2,$3,$4,$5) ON CONFLICT(result_id) DO UPDATE SET payload=forge_verification.results.payload
WHERE forge_verification.results.digest=EXCLUDED.digest`,
		value.ID, value.Digest, value.AssurancePlanID, value.CandidateDigest, payload)
	if err != nil {
		return ErrStoreUnavailable
	}
	affected, err := result.RowsAffected()
	if err != nil || affected != 1 {
		return ErrStoreUnavailable
	}
	return nil
}

func (store *PostgresStore) ResolvePlan(ctx context.Context, id string) (verificationcontract.AssurancePlan, bool, error) {
	var payload []byte
	err := store.db.QueryRowContext(ctx, `SELECT payload FROM forge_verification.assurance_plans WHERE plan_id=$1`, id).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return verificationcontract.AssurancePlan{}, false, nil
	}
	if err != nil {
		return verificationcontract.AssurancePlan{}, false, ErrStoreUnavailable
	}
	value, err := verificationcontract.DecodePlan(payload)
	if err != nil {
		return verificationcontract.AssurancePlan{}, false, ErrStoreUnavailable
	}
	return value, true, nil
}

func (store *PostgresStore) ResolveResult(ctx context.Context, id string) (verificationcontract.Result, bool, error) {
	var payload []byte
	err := store.db.QueryRowContext(ctx, `SELECT payload FROM forge_verification.results WHERE result_id=$1`, id).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return verificationcontract.Result{}, false, nil
	}
	if err != nil {
		return verificationcontract.Result{}, false, ErrStoreUnavailable
	}
	value, err := verificationcontract.DecodeResult(payload)
	if err != nil {
		return verificationcontract.Result{}, false, ErrStoreUnavailable
	}
	return value, true, nil
}
