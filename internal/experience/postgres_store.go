package experience

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

type PostgresStore struct{ db *sql.DB }

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	if db == nil {
		return nil, errors.New("experience PostgreSQL connection is required")
	}
	return &PostgresStore{db: db}, nil
}

func (store *PostgresStore) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_experience;
CREATE TABLE IF NOT EXISTS forge_experience.schema_version(version integer PRIMARY KEY CHECK(version=1));
INSERT INTO forge_experience.schema_version(version) VALUES(1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_experience.requirements_proposals(
 tenant_id text NOT NULL,proposal_id text NOT NULL,version bigint NOT NULL,digest text NOT NULL,payload jsonb NOT NULL,
 current boolean NOT NULL,PRIMARY KEY(tenant_id,proposal_id,version),UNIQUE(tenant_id,proposal_id,digest));
CREATE UNIQUE INDEX IF NOT EXISTS requirements_proposals_one_current
 ON forge_experience.requirements_proposals(tenant_id,proposal_id) WHERE current;`
	_, err := store.db.ExecContext(ctx, migration)
	return err
}

func (store *PostgresStore) Create(ctx context.Context, value RequirementsProposal) error {
	payload, err := json.Marshal(value)
	if err != nil || !value.Current || value.Version != 1 {
		return ErrProposalConflict
	}
	_, err = store.db.ExecContext(ctx, `
INSERT INTO forge_experience.requirements_proposals(tenant_id,proposal_id,version,digest,payload,current)
VALUES($1,$2,$3,$4,$5,true)`, value.TenantID, value.ID, value.Version, value.Digest, payload)
	if err != nil {
		return ErrProposalConflict
	}
	return nil
}

func (store *PostgresStore) Revise(ctx context.Context, expectedDigest string, value RequirementsProposal) error {
	payload, err := json.Marshal(value)
	if err != nil || !value.Current || value.Version < 2 || value.PreviousDigest != expectedDigest {
		return ErrProposalConflict
	}
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var version uint64
	var digest string
	err = tx.QueryRowContext(ctx, `
SELECT version,digest FROM forge_experience.requirements_proposals
WHERE tenant_id=$1 AND proposal_id=$2 AND current=true FOR UPDATE`,
		value.TenantID, value.ID,
	).Scan(&version, &digest)
	if err != nil || digest != expectedDigest || value.Version != version+1 {
		return ErrProposalConflict
	}
	if _, err = tx.ExecContext(ctx, `
UPDATE forge_experience.requirements_proposals SET current=false
WHERE tenant_id=$1 AND proposal_id=$2 AND current=true`, value.TenantID, value.ID); err != nil {
		return err
	}
	if _, err = tx.ExecContext(ctx, `
INSERT INTO forge_experience.requirements_proposals(tenant_id,proposal_id,version,digest,payload,current)
VALUES($1,$2,$3,$4,$5,true)`, value.TenantID, value.ID, value.Version, value.Digest, payload); err != nil {
		return ErrProposalConflict
	}
	return tx.Commit()
}

func (store *PostgresStore) Resolve(ctx context.Context, tenantID, proposalID string) (RequirementsProposal, bool, error) {
	var payload []byte
	err := store.db.QueryRowContext(ctx, `
SELECT payload FROM forge_experience.requirements_proposals
WHERE tenant_id=$1 AND proposal_id=$2 AND current=true`, tenantID, proposalID).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return RequirementsProposal{}, false, nil
	}
	if err != nil {
		return RequirementsProposal{}, false, err
	}
	var value RequirementsProposal
	if json.Unmarshal(payload, &value) != nil {
		return RequirementsProposal{}, false, errors.New("requirements proposal payload is invalid")
	}
	return cloneProposal(value), true, nil
}

var _ Store = (*PostgresStore)(nil)
