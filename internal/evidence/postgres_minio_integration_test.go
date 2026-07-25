package evidence

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/trust"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

func TestPostgresEvidenceSurvivesServiceRestart(t *testing.T) {
	dsn := os.Getenv("FORGE_TEST_POSTGRES_DSN")
	if dsn == "" {
		t.Skip("FORGE_TEST_POSTGRES_DSN is not configured")
	}
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	store, _ := NewPostgresStore(db)
	if err := store.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	scope := trust.EvidenceScope{
		TenantID: "tenant", WorkspaceID: "workspace", ProjectID: "project",
		DeliveryRunID: fmt.Sprintf("run-%d", now.UnixNano()),
	}
	firstService, _ := NewService(store, func() time.Time { return now })
	first, replayed, err := firstService.Append(ctx, request(scope, "first", digest('a')))
	if err != nil || replayed {
		t.Fatalf("first append replay=%t err=%v", replayed, err)
	}
	restartedStore, _ := NewPostgresStore(db)
	restartedService, _ := NewService(restartedStore, func() time.Time { return now.Add(time.Second) })
	second, replayed, err := restartedService.Append(ctx, request(scope, "second", digest('b')))
	if err != nil || replayed {
		t.Fatalf("second append replay=%t err=%v", replayed, err)
	}
	if second.Integrity.PreviousRecordDigest != first.Integrity.Digest {
		t.Fatal("restart did not preserve the evidence chain head")
	}
	if err := restartedService.VerifyStream(ctx, scope, "corr"); err != nil {
		t.Fatal(err)
	}
}

func TestMinIOArtifactIsDigestAddressedAndMetadataIsDurable(t *testing.T) {
	dsn, endpoint := os.Getenv("FORGE_TEST_POSTGRES_DSN"), os.Getenv("FORGE_TEST_MINIO_ENDPOINT")
	if dsn == "" || endpoint == "" {
		t.Skip("PostgreSQL and MinIO integration endpoints are not configured")
	}
	ctx := context.Background()
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	metadataStore, _ := NewPostgresStore(db)
	if err := metadataStore.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	client, err := minio.New(endpoint, &minio.Options{
		Creds: credentials.NewStaticV4("minio", "minio-dev-only", ""), Secure: false,
	})
	if err != nil {
		t.Fatal(err)
	}
	const bucket = "forge-evidence-test"
	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		t.Fatal(err)
	}
	if !exists {
		if err := client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			t.Fatal(err)
		}
	}
	store, err := NewMinIOArtifactStore(client, bucket, metadataStore, 1<<20, func() time.Time { return time.Now().UTC() })
	if err != nil {
		t.Fatal(err)
	}
	runID := fmt.Sprintf("run-%d", time.Now().UnixNano())
	value, err := store.Put(ctx, PutArtifactRequest{
		TenantID: "tenant", ProjectID: "project", DeliveryRunID: runID,
		MediaType: "application/json", Classification: "Internal",
		RetentionClass: "developer-trial-m1", CreatorCapability: "forge.capability.evidence",
		Bytes: []byte(`{"outcome":"passed"}`),
	})
	if err != nil {
		t.Fatal(err)
	}
	if value.ObjectKey == "" || value.Digest == "" || value.SizeBytes == 0 {
		t.Fatalf("incomplete artifact metadata: %#v", value)
	}
	if _, err := client.StatObject(ctx, bucket, value.ObjectKey, minio.StatObjectOptions{}); err != nil {
		t.Fatal(err)
	}
	loaded, err := store.Get(ctx, GetArtifactRequest{
		TenantID: "tenant", ProjectID: "project", DeliveryRunID: runID,
		ObjectKey: value.ObjectKey, Digest: value.Digest,
	})
	if err != nil || string(loaded) != `{"outcome":"passed"}` {
		t.Fatalf("artifact read=%q err=%v", loaded, err)
	}
	if _, err := store.Put(ctx, PutArtifactRequest{
		TenantID: "../escape", ProjectID: "project", DeliveryRunID: runID,
		MediaType: "application/json", Classification: "Internal",
		RetentionClass: "developer-trial-m1", CreatorCapability: "forge.capability.evidence",
		Bytes: []byte(`{}`),
	}); err == nil {
		t.Fatal("object prefix traversal must fail closed")
	}
}
