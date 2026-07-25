package sourcecontrol

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPostgresContractDoubleReconcilesAcrossControlPlaneRestart(t *testing.T) {
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
	operationID := fmt.Sprintf("operation-%d", time.Now().UnixNano())
	operation := Operation{
		ContractVersion: ProviderContractVersion, ID: operationID, Provider: "contract-double",
		RepositoryID: "repository-forge", BaseBranch: "main", BaseRevision: "revision",
		Branch: "forge/" + operationID, CandidateDigest: integrationDigest('a'),
		TreeDigest: integrationDigest('b'), Marker: "<!-- forge-publication:" + operationID + " -->",
	}
	first, _ := NewPostgresContractDouble(db)
	if err := first.CreateBranch(ctx, operation); err != nil {
		t.Fatal(err)
	}
	if err := first.OpenPullRequest(ctx, operation, "title", "body"); err != nil {
		t.Fatal(err)
	}
	restarted, _ := NewPostgresContractDouble(db)
	state, err := restarted.Reconcile(ctx, operation)
	if err != nil || state.PullRequestID == "" || !state.PullRequestOpen {
		t.Fatalf("reconciled state=%#v err=%v", state, err)
	}
	if err := restarted.OpenPullRequest(ctx, operation, "title", "body"); err != nil {
		t.Fatalf("idempotent pull-request replay failed: %v", err)
	}
	if err := restarted.CreateBranch(ctx, operation); err != nil {
		t.Fatalf("idempotent branch replay after pull request failed: %v", err)
	}
	replayed, err := restarted.Reconcile(ctx, operation)
	if err != nil || replayed.PullRequestID != state.PullRequestID {
		t.Fatalf("replay created different publication: %#v %#v %v", state, replayed, err)
	}
}

func integrationDigest(character byte) string {
	value := make([]byte, 64)
	for index := range value {
		value[index] = character
	}
	return "sha256:" + string(value)
}
