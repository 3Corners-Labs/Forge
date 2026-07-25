package sourcecontrol

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
)

type PostgresStore struct{ db *sql.DB }

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	if db == nil {
		return nil, errors.New("source-control PostgreSQL connection is required")
	}
	return &PostgresStore{db: db}, nil
}

func (store *PostgresStore) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_source_control;
CREATE TABLE IF NOT EXISTS forge_source_control.schema_version(version integer PRIMARY KEY CHECK(version=1));
INSERT INTO forge_source_control.schema_version(version) VALUES(1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_source_control.operations(
 tenant_id text NOT NULL,idempotency_key text NOT NULL,fingerprint text NOT NULL,operation_id text NOT NULL,payload jsonb NOT NULL,
 PRIMARY KEY(tenant_id,idempotency_key),UNIQUE(operation_id));
CREATE TABLE IF NOT EXISTS forge_source_control.receipts(
 operation_id text PRIMARY KEY,payload jsonb NOT NULL);
CREATE TABLE IF NOT EXISTS forge_source_control.contract_double_remote(
 operation_id text PRIMARY KEY,payload jsonb NOT NULL);`
	_, err := store.db.ExecContext(ctx, migration)
	return err
}

func (store *PostgresStore) Create(ctx context.Context, tenantID, key string, operation Operation) (Operation, bool, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return Operation{}, false, ErrStoreUnavailable
	}
	defer tx.Rollback()
	fingerprint := operationFingerprint(operation)
	var existingFingerprint string
	var payload []byte
	err = tx.QueryRowContext(ctx, `
SELECT fingerprint,payload FROM forge_source_control.operations
WHERE tenant_id=$1 AND idempotency_key=$2 FOR UPDATE`, tenantID, key).Scan(&existingFingerprint, &payload)
	if err == nil {
		var stored Operation
		if existingFingerprint != fingerprint || json.Unmarshal(payload, &stored) != nil {
			return Operation{}, false, ErrStoreConflict
		}
		return cloneOperation(stored), true, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return Operation{}, false, ErrStoreUnavailable
	}
	payload, err = json.Marshal(operation)
	if err != nil {
		return Operation{}, false, ErrStoreUnavailable
	}
	if _, err = tx.ExecContext(ctx, `
INSERT INTO forge_source_control.operations(tenant_id,idempotency_key,fingerprint,operation_id,payload)
VALUES($1,$2,$3,$4,$5)`, tenantID, key, fingerprint, operation.ID, payload); err != nil {
		return Operation{}, false, ErrStoreConflict
	}
	return cloneOperation(operation), false, tx.Commit()
}

func (store *PostgresStore) SaveReceipt(ctx context.Context, operationID string, receipt Receipt) (Receipt, bool, error) {
	payload, err := json.Marshal(receipt)
	if err != nil {
		return Receipt{}, false, ErrStoreUnavailable
	}
	result, err := store.db.ExecContext(ctx, `
INSERT INTO forge_source_control.receipts(operation_id,payload) VALUES($1,$2)
ON CONFLICT(operation_id) DO NOTHING`, operationID, payload)
	if err != nil {
		return Receipt{}, false, ErrStoreUnavailable
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return Receipt{}, false, ErrStoreUnavailable
	}
	if affected == 1 {
		return receipt, false, nil
	}
	var existingPayload []byte
	if err = store.db.QueryRowContext(ctx, `
SELECT payload FROM forge_source_control.receipts WHERE operation_id=$1`, operationID).Scan(&existingPayload); err != nil {
		return Receipt{}, false, ErrStoreUnavailable
	}
	var existing Receipt
	if json.Unmarshal(existingPayload, &existing) != nil ||
		existing.RepositoryID != receipt.RepositoryID || existing.PullRequestID != receipt.PullRequestID ||
		existing.CandidateDigest != receipt.CandidateDigest || existing.Branch != receipt.Branch {
		return Receipt{}, false, ErrStoreConflict
	}
	return existing, true, nil
}

var _ Store = (*PostgresStore)(nil)

// PostgresContractDouble models an external idempotent provider while keeping
// its state durable across control-plane restarts. It remains credential-free
// and has no merge, release, deployment, or production operation.
type PostgresContractDouble struct{ db *sql.DB }

func NewPostgresContractDouble(db *sql.DB) (*PostgresContractDouble, error) {
	if db == nil {
		return nil, errors.New("contract-double PostgreSQL connection is required")
	}
	return &PostgresContractDouble{db: db}, nil
}

func (provider *PostgresContractDouble) Reconcile(ctx context.Context, operation Operation) (RemoteState, error) {
	var payload []byte
	err := provider.db.QueryRowContext(ctx, `
SELECT payload FROM forge_source_control.contract_double_remote WHERE operation_id=$1`, operation.ID).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return RemoteState{}, nil
	}
	if err != nil {
		return RemoteState{}, ErrProviderUnavailable
	}
	var state RemoteState
	if json.Unmarshal(payload, &state) != nil || !remoteMatchesOperation(state, operation) {
		return RemoteState{}, ErrProviderConflict
	}
	return state, nil
}

func (provider *PostgresContractDouble) CreateBranch(ctx context.Context, operation Operation) error {
	state := RemoteState{
		BranchExists: true, Branch: operation.Branch, CommitID: "double-commit-" + operation.ID,
		CandidateDigest: operation.CandidateDigest, TreeDigest: operation.TreeDigest,
		BaseBranch: operation.BaseBranch, BaseRevision: operation.BaseRevision,
		Marker: operation.Marker, RepositoryID: operation.RepositoryID,
	}
	payload, _ := json.Marshal(state)
	result, err := provider.db.ExecContext(ctx, `
INSERT INTO forge_source_control.contract_double_remote(operation_id,payload)
VALUES($1,$2) ON CONFLICT(operation_id) DO NOTHING`, operation.ID, payload)
	if err != nil {
		return ErrProviderUnavailable
	}
	affected, _ := result.RowsAffected()
	if affected == 1 {
		return nil
	}
	existing, reconcileErr := provider.Reconcile(ctx, operation)
	if reconcileErr != nil || !existing.BranchExists {
		return ErrProviderConflict
	}
	return nil
}

func (provider *PostgresContractDouble) OpenPullRequest(ctx context.Context, operation Operation, title, body string) error {
	if title == "" || body == "" {
		return ErrProviderConflict
	}
	tx, err := provider.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return ErrProviderUnavailable
	}
	defer tx.Rollback()
	var payload []byte
	err = tx.QueryRowContext(ctx, `
SELECT payload FROM forge_source_control.contract_double_remote
WHERE operation_id=$1 FOR UPDATE`, operation.ID).Scan(&payload)
	if err != nil {
		return ErrProviderUnavailable
	}
	var state RemoteState
	if json.Unmarshal(payload, &state) != nil || !remoteMatchesOperation(state, operation) {
		return ErrProviderConflict
	}
	if state.PullRequestID == "" {
		state.PullRequestID = "double-pr-" + operation.ID
		state.PullRequestURL = "https://contract-double.invalid/pulls/" + operation.ID
		state.PullRequestOpen = true
		payload, _ = json.Marshal(state)
		if _, err = tx.ExecContext(ctx, `
UPDATE forge_source_control.contract_double_remote SET payload=$2 WHERE operation_id=$1`,
			operation.ID, payload); err != nil {
			return ErrProviderUnavailable
		}
	}
	return tx.Commit()
}

func remoteMatchesOperation(state RemoteState, operation Operation) bool {
	return state.BranchExists && state.Branch == operation.Branch &&
		state.CandidateDigest == operation.CandidateDigest && state.TreeDigest == operation.TreeDigest &&
		state.BaseBranch == operation.BaseBranch && state.BaseRevision == operation.BaseRevision &&
		state.Marker == operation.Marker && state.RepositoryID == operation.RepositoryID
}

var _ Provider = (*PostgresContractDouble)(nil)
