package verification

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPrepareFreshCheckoutReconstructsIndependentSource(t *testing.T) {
	content := "# source\n"
	root, err := PrepareFreshCheckout(t.TempDir(), []SourceFile{{
		Path: "docs/source.md", Digest: checkoutDigest([]byte(content)), Content: content,
	}})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(root) })
	value, err := os.ReadFile(filepath.Join(root, "docs/source.md"))
	if err != nil || string(value) != content {
		t.Fatalf("fresh checkout = %q, %v", value, err)
	}
	info, err := os.Stat(filepath.Join(root, "docs/source.md"))
	if err != nil || info.Mode().Perm() != 0o644 {
		t.Fatalf("verifier source mode = %v, %v", info.Mode().Perm(), err)
	}
}

func TestPrepareFreshCheckoutRejectsUntrustedProjection(t *testing.T) {
	if _, err := PrepareFreshCheckout(t.TempDir(), []SourceFile{{
		Path: "/absolute", Digest: checkoutDigest([]byte("x")), Content: "x",
	}}); err == nil {
		t.Fatal("absolute source path was admitted")
	}
}
