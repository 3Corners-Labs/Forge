package evidence

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"

	"github.com/3corners/forge/internal/contracts/trust"
)

type PostgresStore struct{ db *sql.DB }

func NewPostgresStore(db *sql.DB) (*PostgresStore, error) {
	if db == nil {
		return nil, errors.New("evidence PostgreSQL connection is required")
	}
	return &PostgresStore{db: db}, nil
}

func (store *PostgresStore) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_evidence;
CREATE TABLE IF NOT EXISTS forge_evidence.schema_version(version integer PRIMARY KEY CHECK(version=1));
INSERT INTO forge_evidence.schema_version(version) VALUES(1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_evidence.records(
 stream_key text NOT NULL,sequence bigint NOT NULL,record_id text NOT NULL,payload jsonb NOT NULL,
 PRIMARY KEY(stream_key,sequence),UNIQUE(record_id));
CREATE TABLE IF NOT EXISTS forge_evidence.chain_heads(
 stream_key text PRIMARY KEY,sequence bigint NOT NULL,digest text NOT NULL);
CREATE TABLE IF NOT EXISTS forge_evidence.idempotency(
 stream_key text NOT NULL,idempotency_key text NOT NULL,fingerprint text NOT NULL,record_id text NOT NULL,
 PRIMARY KEY(stream_key,idempotency_key));
CREATE TABLE IF NOT EXISTS forge_evidence.artifacts(
 object_key text PRIMARY KEY,tenant_id text NOT NULL,project_id text NOT NULL,delivery_run_id text NOT NULL,
 digest text NOT NULL,size_bytes bigint NOT NULL,media_type text NOT NULL,classification text NOT NULL,
 retention_class text NOT NULL,creator_capability text NOT NULL,created_at timestamptz NOT NULL,
 UNIQUE(tenant_id,project_id,delivery_run_id,digest));`
	_, err := store.db.ExecContext(ctx, migration)
	return err
}

func (store *PostgresStore) Head(ctx context.Context, stream string) (string, error) {
	var digest string
	err := store.db.QueryRowContext(ctx, `SELECT digest FROM forge_evidence.chain_heads WHERE stream_key=$1`, evidenceDatabaseKey(stream)).Scan(&digest)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", ErrUnavailable
	}
	return digest, nil
}

func (store *PostgresStore) Append(ctx context.Context, command AppendCommand) (trust.EvidenceRecord, bool, error) {
	tx, err := store.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return trust.EvidenceRecord{}, false, ErrUnavailable
	}
	defer tx.Rollback()
	var fingerprint, recordID string
	streamKey := evidenceDatabaseKey(command.StreamKey)
	err = tx.QueryRowContext(ctx, `SELECT fingerprint,record_id FROM forge_evidence.idempotency WHERE stream_key=$1 AND idempotency_key=$2 FOR UPDATE`, streamKey, command.IdempotencyKey).Scan(&fingerprint, &recordID)
	if err == nil {
		if fingerprint != command.Fingerprint {
			return trust.EvidenceRecord{}, false, ErrConflict
		}
		value, found, loadErr := loadEvidenceRecord(ctx, tx, recordID)
		if loadErr != nil || !found {
			return trust.EvidenceRecord{}, false, ErrUnavailable
		}
		return value, true, tx.Commit()
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return trust.EvidenceRecord{}, false, ErrUnavailable
	}
	var sequence int64
	var current string
	headErr := tx.QueryRowContext(ctx, `SELECT sequence,digest FROM forge_evidence.chain_heads WHERE stream_key=$1 FOR UPDATE`, streamKey).Scan(&sequence, &current)
	if errors.Is(headErr, sql.ErrNoRows) {
		sequence, current = 0, ""
	} else if headErr != nil {
		return trust.EvidenceRecord{}, false, ErrUnavailable
	}
	if current != command.ExpectedPrevious || command.Record.Integrity.PreviousRecordDigest != current {
		return trust.EvidenceRecord{}, false, ErrConflict
	}
	sequence++
	payload, _ := json.Marshal(command.Record)
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_evidence.records VALUES($1,$2,$3,$4)`, streamKey, sequence, command.Record.ID, payload); err != nil {
		return trust.EvidenceRecord{}, false, ErrConflict
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_evidence.chain_heads VALUES($1,$2,$3) ON CONFLICT(stream_key) DO UPDATE SET sequence=EXCLUDED.sequence,digest=EXCLUDED.digest`, streamKey, sequence, command.Record.Integrity.Digest); err != nil {
		return trust.EvidenceRecord{}, false, ErrConflict
	}
	if _, err = tx.ExecContext(ctx, `INSERT INTO forge_evidence.idempotency VALUES($1,$2,$3,$4)`, streamKey, command.IdempotencyKey, command.Fingerprint, command.Record.ID); err != nil {
		return trust.EvidenceRecord{}, false, ErrConflict
	}
	if err = tx.Commit(); err != nil {
		return trust.EvidenceRecord{}, false, ErrUnavailable
	}
	return cloneRecord(command.Record), false, nil
}

type evidenceQuery interface {
	QueryRowContext(context.Context, string, ...any) *sql.Row
}

func loadEvidenceRecord(ctx context.Context, query evidenceQuery, id string) (trust.EvidenceRecord, bool, error) {
	var payload []byte
	err := query.QueryRowContext(ctx, `SELECT payload FROM forge_evidence.records WHERE record_id=$1`, id).Scan(&payload)
	if errors.Is(err, sql.ErrNoRows) {
		return trust.EvidenceRecord{}, false, nil
	}
	if err != nil {
		return trust.EvidenceRecord{}, false, ErrUnavailable
	}
	var value trust.EvidenceRecord
	if json.Unmarshal(payload, &value) != nil {
		return trust.EvidenceRecord{}, false, ErrUnavailable
	}
	return cloneRecord(value), true, nil
}

func (store *PostgresStore) Stream(ctx context.Context, stream string) ([]trust.EvidenceRecord, error) {
	rows, err := store.db.QueryContext(ctx, `SELECT payload FROM forge_evidence.records WHERE stream_key=$1 ORDER BY sequence`, evidenceDatabaseKey(stream))
	if err != nil {
		return nil, ErrUnavailable
	}
	defer rows.Close()
	var values []trust.EvidenceRecord
	for rows.Next() {
		var payload []byte
		var value trust.EvidenceRecord
		if rows.Scan(&payload) != nil || json.Unmarshal(payload, &value) != nil {
			return nil, ErrUnavailable
		}
		values = append(values, cloneRecord(value))
	}
	return values, rows.Err()
}

func evidenceDatabaseKey(value string) string {
	return base64.RawURLEncoding.EncodeToString([]byte(value))
}

type ArtifactMetadata struct {
	ObjectKey         string
	TenantID          string
	ProjectID         string
	DeliveryRunID     string
	Digest            string
	SizeBytes         int64
	MediaType         string
	Classification    string
	RetentionClass    string
	CreatorCapability string
	CreatedAt         string
}

func (store *PostgresStore) CommitArtifact(ctx context.Context, value ArtifactMetadata) error {
	_, err := store.db.ExecContext(ctx, `
INSERT INTO forge_evidence.artifacts(object_key,tenant_id,project_id,delivery_run_id,digest,size_bytes,media_type,classification,retention_class,creator_capability,created_at)
VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
ON CONFLICT(tenant_id,project_id,delivery_run_id,digest) DO NOTHING`,
		value.ObjectKey, value.TenantID, value.ProjectID, value.DeliveryRunID, value.Digest,
		value.SizeBytes, value.MediaType, value.Classification, value.RetentionClass,
		value.CreatorCapability, value.CreatedAt,
	)
	return err
}
