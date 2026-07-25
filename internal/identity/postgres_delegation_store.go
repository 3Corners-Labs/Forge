package identity

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"time"
)

type PostgresDelegationStore struct{ db *sql.DB }

func NewPostgresDelegationStore(db *sql.DB) (*PostgresDelegationStore, error) {
	if db == nil {
		return nil, errors.New("identity PostgreSQL connection is required")
	}
	return &PostgresDelegationStore{db: db}, nil
}

func (store *PostgresDelegationStore) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_identity;
CREATE TABLE IF NOT EXISTS forge_identity.schema_version(version integer PRIMARY KEY CHECK(version=1));
INSERT INTO forge_identity.schema_version(version) VALUES(1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_identity.delegations(
 delegation_id text PRIMARY KEY,tenant_id text NOT NULL,request_digest text NOT NULL,payload jsonb NOT NULL,
 revoked_at timestamptz,revocation_reason text);
CREATE TABLE IF NOT EXISTS forge_identity.delegation_keys(
 tenant_id text NOT NULL,idempotency_key text NOT NULL,request_digest text NOT NULL,delegation_id text NOT NULL,
 PRIMARY KEY(tenant_id,idempotency_key));`
	_, err := store.db.ExecContext(ctx, migration)
	return err
}

func (store *PostgresDelegationStore) Create(ctx context.Context, value Delegation, key, requestDigest string) (Delegation, bool, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return Delegation{}, false, ErrDelegationStoreUnavailable
	}
	defer tx.Rollback()
	var digest, id string
	err = tx.QueryRowContext(ctx, `
SELECT request_digest,delegation_id FROM forge_identity.delegation_keys
WHERE tenant_id=$1 AND idempotency_key=$2 FOR UPDATE`, value.TenantID(), key).Scan(&digest, &id)
	if err == nil {
		if digest != requestDigest {
			return Delegation{}, false, ErrDelegationConflict
		}
		stored, found, loadErr := loadDelegation(ctx, tx, id)
		if loadErr != nil || !found || stored.Revoked {
			return Delegation{}, false, ErrDelegationRevoked
		}
		return stored.Delegation, true, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return Delegation{}, false, ErrDelegationStoreUnavailable
	}
	payload, _ := json.Marshal(delegationToPersistence(value))
	if _, err = tx.ExecContext(ctx, `
INSERT INTO forge_identity.delegations(delegation_id,tenant_id,request_digest,payload)
VALUES($1,$2,$3,$4)`, value.ID(), value.TenantID(), requestDigest, payload); err != nil {
		return Delegation{}, false, ErrDelegationConflict
	}
	if _, err = tx.ExecContext(ctx, `
INSERT INTO forge_identity.delegation_keys(tenant_id,idempotency_key,request_digest,delegation_id)
VALUES($1,$2,$3,$4)`, value.TenantID(), key, requestDigest, value.ID()); err != nil {
		return Delegation{}, false, ErrDelegationConflict
	}
	if err = tx.Commit(); err != nil {
		return Delegation{}, false, ErrDelegationStoreUnavailable
	}
	return cloneDelegation(value), false, nil
}

func (store *PostgresDelegationStore) Get(ctx context.Context, id string) (StoredDelegation, bool, error) {
	return loadDelegation(ctx, store.db, id)
}

func (store *PostgresDelegationStore) Revoke(ctx context.Context, id, reason string, at time.Time) (StoredDelegation, bool, error) {
	result, err := store.db.ExecContext(ctx, `
UPDATE forge_identity.delegations SET revoked_at=COALESCE(revoked_at,$2),
 revocation_reason=COALESCE(revocation_reason,$3) WHERE delegation_id=$1`, id, at, reason)
	if err != nil {
		return StoredDelegation{}, false, ErrDelegationStoreUnavailable
	}
	affected, _ := result.RowsAffected()
	if affected == 0 {
		return StoredDelegation{}, false, nil
	}
	return loadDelegation(ctx, store.db, id)
}

type delegationQuery interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func loadDelegation(ctx context.Context, query delegationQuery, id string) (StoredDelegation, bool, error) {
	var payload []byte
	var revoked sql.NullTime
	err := query.QueryRowContext(ctx, `
SELECT payload,revoked_at FROM forge_identity.delegations WHERE delegation_id=$1`, id).Scan(&payload, &revoked)
	if errors.Is(err, sql.ErrNoRows) {
		return StoredDelegation{}, false, nil
	}
	if err != nil {
		return StoredDelegation{}, false, ErrDelegationStoreUnavailable
	}
	var value delegationPersistence
	if json.Unmarshal(payload, &value) != nil {
		return StoredDelegation{}, false, ErrDelegationStoreUnavailable
	}
	return StoredDelegation{
		Delegation: value.delegation(), Revoked: revoked.Valid, RevokedAt: revoked.Time,
	}, true, nil
}

type delegationPersistence struct {
	ContractVersion, ID, ParentAuthorityID, IssuerSubjectID, HumanSubjectID, AgentSubjectID string
	TenantID, WorkspaceID, ProjectID, RepositoryID, DeliveryRunID                           string
	Purpose, Environment, MaximumAutonomy, RevocationSource, EvidenceCorrelationID          string
	RequestDigest                                                                           string
	Actions, Resources, Classifications, Destinations                                       []string
	IssuedAt, ExpiresAt                                                                     time.Time
}

func delegationToPersistence(value Delegation) delegationPersistence {
	return delegationPersistence{
		ContractVersion: value.contractVersion, ID: value.id, ParentAuthorityID: value.parentAuthorityID,
		IssuerSubjectID: value.issuerSubjectID, HumanSubjectID: value.humanSubjectID,
		AgentSubjectID: value.agentSubjectID, TenantID: value.tenantID, WorkspaceID: value.workspaceID,
		ProjectID: value.projectID, RepositoryID: value.repositoryID, DeliveryRunID: value.deliveryRunID,
		Purpose: value.purpose, Environment: value.environment, MaximumAutonomy: value.maximumAutonomy,
		RevocationSource: value.revocationSource, EvidenceCorrelationID: value.evidenceCorrelationID,
		RequestDigest: value.requestDigest, Actions: value.Actions(), Resources: value.Resources(),
		Classifications: value.Classifications(), Destinations: value.Destinations(),
		IssuedAt: value.issuedAt, ExpiresAt: value.expiresAt,
	}
}

func (value delegationPersistence) delegation() Delegation {
	return Delegation{
		contractVersion: value.ContractVersion, id: value.ID, parentAuthorityID: value.ParentAuthorityID,
		issuerSubjectID: value.IssuerSubjectID, humanSubjectID: value.HumanSubjectID,
		agentSubjectID: value.AgentSubjectID, tenantID: value.TenantID, workspaceID: value.WorkspaceID,
		projectID: value.ProjectID, repositoryID: value.RepositoryID, deliveryRunID: value.DeliveryRunID,
		purpose: value.Purpose, environment: value.Environment, maximumAutonomy: value.MaximumAutonomy,
		revocationSource: value.RevocationSource, evidenceCorrelationID: value.EvidenceCorrelationID,
		requestDigest: value.RequestDigest, actions: value.Actions, resources: value.Resources,
		classifications: value.Classifications, destinations: value.Destinations,
		issuedAt: value.IssuedAt, expiresAt: value.ExpiresAt,
	}
}

var _ DelegationStore = (*PostgresDelegationStore)(nil)
