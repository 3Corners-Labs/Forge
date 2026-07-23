package evidence

import (
	"context"
	"testing"
	"time"

	"github.com/3corners/forge/internal/contracts/trust"
)

func TestEvidenceChainDetectsMutationAndPreventsForks(t *testing.T) {
	now := time.Now().UTC()
	store := NewMemoryStore()
	service, _ := NewService(store, func() time.Time { return now })
	scope := trust.EvidenceScope{TenantID: "tenant", WorkspaceID: "workspace", ProjectID: "project", DeliveryRunID: "run"}
	first, _, err := service.Append(context.Background(), request(scope, "one", digest('a')))
	if err != nil {
		t.Fatal(err)
	}
	second, _, err := service.Append(context.Background(), request(scope, "two", digest('b')))
	if err != nil {
		t.Fatal(err)
	}
	if second.Integrity.PreviousRecordDigest != first.Integrity.Digest {
		t.Fatal("evidence chain was not linked")
	}
	if err := service.VerifyStream(context.Background(), scope, "corr"); err != nil {
		t.Fatal(err)
	}
	store.streams["tenant\x00project\x00run"][0].Subject.Digest = digest('c')
	if err := service.VerifyStream(context.Background(), scope, "corr"); err == nil {
		t.Fatal("mutated evidence must fail verification")
	}
}

func TestEvidenceIdempotencyCannotChangeImmutableInput(t *testing.T) {
	now := time.Now().UTC()
	service, _ := NewService(NewMemoryStore(), func() time.Time { return now })
	scope := trust.EvidenceScope{TenantID: "tenant", WorkspaceID: "workspace", ProjectID: "project", DeliveryRunID: "run"}
	value := request(scope, "same", digest('a'))
	first, replayed, err := service.Append(context.Background(), value)
	if err != nil || replayed {
		t.Fatal("first append failed")
	}
	replay, replayed, err := service.Append(context.Background(), value)
	if err != nil || !replayed || replay.ID != first.ID {
		t.Fatal("identical replay did not return the record")
	}
	value.Subject.Digest = digest('b')
	if _, _, err := service.Append(context.Background(), value); err == nil {
		t.Fatal("changed immutable input must conflict")
	}
}

func request(scope trust.EvidenceScope, key, subjectDigest string) AppendRequest {
	return AppendRequest{
		CorrelationID: "corr", IdempotencyKey: key, Type: "verification",
		ActorID: "forge.verification", Scope: scope,
		Subject:    trust.EvidenceSubject{Type: "candidate", ID: "candidate-1", Digest: subjectDigest},
		References: []string{"verification:result"}, Outcome: "passed",
		Classification: "Internal", RetentionClass: "developer-trial-m1",
	}
}

func digest(value byte) string {
	result := make([]byte, 64)
	for index := range result {
		result[index] = value
	}
	return "sha256:" + string(result)
}
