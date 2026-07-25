package context

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRepositoryCatalogBindsMetadataAndExactProjection(t *testing.T) {
	root := t.TempDir()
	path := filepath.Join(root, "docs", "architecture.md")
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("# Accepted\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	catalog, err := NewRepositoryCatalog(root, "revision-1", []RepositorySource{
		{ID: "architecture", Path: "docs/architecture.md", Authority: "accepted-knowledge", Classification: "Internal"},
		{ID: "target", Path: "docs/new.md", Authority: "reference", Classification: "Internal", Optional: true},
	}, func() time.Time { return time.Unix(100, 0).UTC() })
	if err != nil {
		t.Fatal(err)
	}
	records, err := catalog.Discover(context.Background(), SourceDiscoveryRequest{
		TenantID: "tenant", WorkspaceID: "workspace", ProjectID: "project", RepositoryID: "repository",
		RequestedSourceIDs: []string{"architecture", "target"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(records) != 1 || records[0].SourceRevision != "revision-1" {
		t.Fatalf("unexpected records: %#v", records)
	}
	projection, err := catalog.Project([]string{"docs/architecture.md", "docs/new.md"})
	if err != nil {
		t.Fatal(err)
	}
	if len(projection) != 1 || projection[0].Content != "# Accepted\n" {
		t.Fatalf("unexpected projection: %#v", projection)
	}
}

func TestRepositoryCatalogRejectsTraversalAndUnregisteredProjection(t *testing.T) {
	if _, err := NewRepositoryCatalog(t.TempDir(), "revision", []RepositorySource{{
		ID: "escape", Path: "../escape", Authority: "source", Classification: "Internal",
	}}, time.Now); err == nil {
		t.Fatal("catalog admitted traversal")
	}
	catalog, err := NewRepositoryCatalog(t.TempDir(), "revision", []RepositorySource{{
		ID: "target", Path: "docs/new.md", Authority: "source", Classification: "Internal", Optional: true,
	}}, time.Now)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := catalog.Project([]string{"docs/other.md"}); err == nil {
		t.Fatal("catalog projected an unregistered path")
	}
}
