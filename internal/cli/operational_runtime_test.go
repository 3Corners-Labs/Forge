package cli

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestOperationalRuntimeUsesProjectLocalRootAndExactRevision(t *testing.T) {
	root, revision := operationalRepositoryFixture(t)
	runtime, err := newOperationalRuntime(journeyOptions{
		RuntimeMode: "docker", RepositoryRoot: root, RuntimeRoot: filepath.Join(root, ".forge", "runtime"),
		SourceRevision: revision, WorkerImage: "sha256:" + strings.Repeat("a", 64),
		VerifierImage: "sha256:" + strings.Repeat("b", 64),
	}, time.Now)
	if err != nil {
		t.Fatal(err)
	}
	if !runtime.operational || runtime.sourceRevision != revision || !withinRoot(runtime.root, runtime.runtimeRoot) {
		t.Fatalf("unexpected runtime: %#v", runtime)
	}
}

func TestOperationalRuntimeRejectsExternalRootAndMutableImages(t *testing.T) {
	root, revision := operationalRepositoryFixture(t)
	if _, err := newOperationalRuntime(journeyOptions{
		RepositoryRoot: root, RuntimeRoot: filepath.Dir(root),
	}, time.Now); err == nil {
		t.Fatal("external runtime root was admitted")
	}
	if _, err := newOperationalRuntime(journeyOptions{
		RuntimeMode: "docker", RepositoryRoot: root, SourceRevision: revision,
		WorkerImage: "forge-worker:latest", VerifierImage: "forge-verifier:latest",
	}, time.Now); err == nil {
		t.Fatal("mutable images were admitted")
	}
}

func TestOperationalRuntimeRejectsDirtyAuthorizedSource(t *testing.T) {
	root, revision := operationalRepositoryFixture(t)
	path := filepath.Join(root, filepath.FromSlash(m1ArchitecturePath))
	if err := os.WriteFile(path, []byte("# changed\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := newOperationalRuntime(journeyOptions{
		RuntimeMode: "docker", RepositoryRoot: root, SourceRevision: revision,
		WorkerImage:   "sha256:" + strings.Repeat("a", 64),
		VerifierImage: "sha256:" + strings.Repeat("b", 64),
	}, time.Now); err == nil {
		t.Fatal("dirty authorized source was admitted as an exact revision")
	}
}

func operationalRepositoryFixture(t *testing.T) (string, string) {
	t.Helper()
	root := t.TempDir()
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module fixture\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(root, filepath.FromSlash(m1ArchitecturePath))
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("# architecture\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	for _, arguments := range [][]string{
		{"init", "--quiet"},
		{"config", "user.email", "forge-test@example.invalid"},
		{"config", "user.name", "Forge Test"},
		{"add", "go.mod", m1ArchitecturePath},
		{"commit", "--quiet", "-m", "fixture"},
	} {
		command := exec.Command("git", arguments...)
		command.Dir = root
		if output, err := command.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v: %s", arguments, err, output)
		}
	}
	command := exec.Command("git", "rev-parse", "HEAD")
	command.Dir = root
	output, err := command.Output()
	if err != nil {
		t.Fatal(err)
	}
	return root, strings.TrimSpace(string(output))
}
