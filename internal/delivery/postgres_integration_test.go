package delivery

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPostgresDeliveryTransitionAndOutboxSurviveRestart(t *testing.T) {
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
	suffix := time.Now().UnixNano()
	store, _ := NewPostgresStore(db, func() string { return fmt.Sprintf("run-%d", suffix) })
	if err := store.Migrate(ctx); err != nil {
		t.Fatal(err)
	}
	record := CaptureRecord{
		ScopeKey:           "tenant\x00workspace\x00project\x00repository\x00human",
		RequestFingerprint: fmt.Sprintf("sha256:%064x", suffix), IdempotencyKey: fmt.Sprintf("capture-%d", suffix),
		CorrelationID: "correlation", TenantID: "tenant", WorkspaceID: "workspace", ProjectID: "project",
		RepositoryID: "repository", ActorID: "human", DelegationChain: []string{},
		Purpose: M1ChangePurpose, Classification: "Internal", Objective: "bounded change",
		Constraints: []string{}, AcceptanceCriteria: []string{"verified"}, CreatedAt: time.Now().UTC(),
	}
	captured, replayed, err := store.CreateCaptured(ctx, record)
	if err != nil || replayed {
		t.Fatalf("capture replay=%t err=%v", replayed, err)
	}
	restarted, _ := NewPostgresStore(db, func() string { return "must-not-be-used" })
	loaded, found, err := restarted.Load(ctx, captured.ID)
	if err != nil || !found || loaded.Version != 1 {
		t.Fatalf("restart load found=%t version=%d err=%v", found, loaded.Version, err)
	}
	command := TransitionCommand{
		RunID: captured.ID, ScopeKey: record.ScopeKey, ExpectedVersion: 1,
		FromState: CapturedState, ToState: ContextAuthorizingState,
		IdempotencyKey: fmt.Sprintf("transition-%d", suffix), Fingerprint: fmt.Sprintf("sha256:%064x", suffix+1),
		CorrelationID: record.CorrelationID, OccurredAt: time.Now().UTC(),
	}
	transitioned, replayed, err := restarted.Transition(ctx, command)
	if err != nil || replayed || transitioned.Version != 2 {
		t.Fatalf("transition replay=%t version=%d err=%v", replayed, transitioned.Version, err)
	}
	events, err := restarted.PendingOutbox(ctx, 10)
	if err != nil {
		t.Fatal(err)
	}
	foundEvent := false
	for _, event := range events {
		if event.RunID == captured.ID && event.RunVersion == 2 {
			foundEvent = true
		}
	}
	if !foundEvent {
		t.Fatal("committed transition is missing its outbox event")
	}
}
