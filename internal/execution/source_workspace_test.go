package execution

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	changecontract "github.com/3corners/forge/internal/contracts/change"
)

func TestPrepareSourceWorkspaceMaterializesOnlyAdmittedBytes(t *testing.T) {
	runtimeRoot := t.TempDir()
	content := "package example\n"
	workspace, err := PrepareSourceWorkspace(runtimeRoot, []changecontract.SourceFile{{
		Path: "internal/example/example.go", Digest: sourceDigest([]byte(content)), Content: content,
	}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(workspace) })
	value, err := os.ReadFile(filepath.Join(workspace, "internal/example/example.go"))
	if err != nil || string(value) != content {
		t.Fatalf("workspace source = %q, %v", value, err)
	}
	if !strings.HasPrefix(workspace, runtimeRoot+string(filepath.Separator)) {
		t.Fatalf("workspace %q is outside runtime root %q", workspace, runtimeRoot)
	}
	info, err := os.Stat(filepath.Join(workspace, "internal/example/example.go"))
	if err != nil || info.Mode().Perm() != 0o666 {
		t.Fatalf("worker source mode = %v, %v", info.Mode().Perm(), err)
	}
}

func TestPrepareSourceWorkspaceFailsClosedOnTraversalOrDigestMismatch(t *testing.T) {
	for _, source := range []changecontract.SourceFile{
		{Path: "../escape", Digest: sourceDigest([]byte("x")), Content: "x"},
		{Path: "safe.txt", Digest: sourceDigest([]byte("different")), Content: "x"},
	} {
		if _, err := PrepareSourceWorkspace(t.TempDir(), []changecontract.SourceFile{source}); err == nil {
			t.Fatalf("source %#v was admitted", source)
		}
	}
}
