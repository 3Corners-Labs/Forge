package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "forge.yaml")
	if err := os.WriteFile(p, []byte("organization:\n  name: 3Corners\nruntime:\n  mode: compose\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Organization.Name != "3Corners" {
		t.Fatalf("unexpected organization: %s", cfg.Organization.Name)
	}
}
