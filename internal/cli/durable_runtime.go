package cli

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/3corners/forge/internal/approval"
	forgecontext "github.com/3corners/forge/internal/context"
	"github.com/3corners/forge/internal/delivery"
	"github.com/3corners/forge/internal/evidence"
	"github.com/3corners/forge/internal/experience"
	"github.com/3corners/forge/internal/identity"
	"github.com/3corners/forge/internal/intent"
	"github.com/3corners/forge/internal/policy"
	"github.com/3corners/forge/internal/sourcecontrol"
	"github.com/3corners/forge/internal/verification"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type journeyPersistence struct {
	db               *sql.DB
	delivery         delivery.Store
	approval         approval.Store
	context          forgecontext.Store
	evidence         evidence.Store
	experience       experience.Store
	intent           intent.Store
	delegation       identity.DelegationStore
	workerDelegation identity.WorkerDelegationStore
	envelope         policy.EnvelopeStore
	verification     verification.Store
	sourceControl    sourcecontrol.Store
	publication      sourcecontrol.Provider
	artifacts        *evidence.MinIOArtifactStore
}

func newJourneyPersistence(operational bool) (journeyPersistence, error) {
	if !operational {
		deliveryStore, err := delivery.NewMemoryStore(func() string {
			return fmt.Sprintf("forge-run-%d", time.Now().UTC().UnixNano())
		})
		if err != nil {
			return journeyPersistence{}, err
		}
		return journeyPersistence{
			delivery: deliveryStore, approval: approval.NewMemoryStore(),
			context: forgecontext.NewMemoryStore(), evidence: evidence.NewMemoryStore(),
			experience: experience.NewMemoryStore(), intent: intent.NewMemoryStore(),
			delegation:       identity.NewMemoryDelegationStore(),
			workerDelegation: identity.NewMemoryWorkerDelegationStore(),
			envelope:         policy.NewMemoryEnvelopeStore(), verification: newTrialVerificationStore(),
			sourceControl: sourcecontrol.NewMemoryStore(),
			publication:   sourcecontrol.NewContractDouble(sourcecontrol.FailNone),
		}, nil
	}
	dsn := os.Getenv("FORGE_M1_POSTGRES_DSN")
	if dsn == "" {
		return journeyPersistence{}, errors.New("Docker operational runtime requires FORGE_M1_POSTGRES_DSN")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return journeyPersistence{}, errors.New("Developer Trial PostgreSQL configuration is invalid")
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial PostgreSQL is unavailable")
	}
	deliveryStore, err := delivery.NewPostgresStore(db, func() string {
		return fmt.Sprintf("forge-run-%d", time.Now().UTC().UnixNano())
	})
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial Delivery store is unavailable")
	}
	approvalStore, err := approval.NewPostgresStore(db)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial Approval store is unavailable")
	}
	contextStore, err := forgecontext.NewPostgresStore(db)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial Context store is unavailable")
	}
	evidenceStore, err := evidence.NewPostgresStore(db)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial Evidence store is unavailable")
	}
	experienceStore, err := experience.NewPostgresStore(db)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial Experience store is unavailable")
	}
	intentStore, err := intent.NewPostgresStore(db)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial Intent store is unavailable")
	}
	delegationStore, err := identity.NewPostgresDelegationStore(db)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial delegation store is unavailable")
	}
	workerStore, err := identity.NewPostgresWorkerDelegationStore(db)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial worker-delegation store is unavailable")
	}
	envelopeStore, err := policy.NewPostgresEnvelopeStore(db)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial Policy store is unavailable")
	}
	verificationStore, err := verification.NewPostgresStore(db)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial Verification store is unavailable")
	}
	sourceControlStore, err := sourcecontrol.NewPostgresStore(db)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial Source Control store is unavailable")
	}
	publicationProvider, err := sourcecontrol.NewPostgresContractDouble(db)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial publication provider is unavailable")
	}
	migrations := []func(context.Context) error{
		deliveryStore.Migrate, approvalStore.Migrate, contextStore.Migrate, evidenceStore.Migrate,
		experienceStore.Migrate, intentStore.Migrate, delegationStore.Migrate, workerStore.Migrate,
		envelopeStore.Migrate, verificationStore.Migrate, sourceControlStore.Migrate,
	}
	for _, migrate := range migrations {
		if err := migrate(ctx); err != nil {
			_ = db.Close()
			return journeyPersistence{}, errors.New("Developer Trial capability migration failed closed")
		}
	}
	endpoint, accessKey, secretKey := os.Getenv("FORGE_M1_MINIO_ENDPOINT"),
		os.Getenv("FORGE_M1_MINIO_ACCESS_KEY"), os.Getenv("FORGE_M1_MINIO_SECRET_KEY")
	if endpoint == "" || accessKey == "" || secretKey == "" {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Docker operational runtime requires bounded MinIO configuration")
	}
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""), Secure: false,
	})
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial artifact configuration is invalid")
	}
	const bucket = "forge-m1-artifacts"
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, errors.New("Developer Trial artifact store is unavailable")
	}
	if !exists {
		if err := minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			_ = db.Close()
			return journeyPersistence{}, errors.New("Developer Trial artifact bucket is unavailable")
		}
	}
	artifacts, err := evidence.NewMinIOArtifactStore(
		minioClient, bucket, evidenceStore, 1<<20, func() time.Time { return time.Now().UTC() },
	)
	if err != nil {
		_ = db.Close()
		return journeyPersistence{}, err
	}
	return journeyPersistence{
		db: db, delivery: deliveryStore, approval: approvalStore, context: contextStore,
		evidence: evidenceStore, experience: experienceStore, intent: intentStore,
		delegation: delegationStore, workerDelegation: workerStore, envelope: envelopeStore,
		verification: verificationStore, sourceControl: sourceControlStore,
		publication: publicationProvider, artifacts: artifacts,
	}, nil
}

func (persistence journeyPersistence) Close() error {
	if persistence.db == nil {
		return nil
	}
	return persistence.db.Close()
}
