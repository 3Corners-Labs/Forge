package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoad(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "forge.yaml")
	if err := os.WriteFile(p, []byte("organization:\n  name: 3Corners\nruntime:\n  mode: compose\n  profile: developer\nintelligence:\n  enabled: false\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	cfg, err := Load(p)
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Organization.Name != "3Corners" {
		t.Fatalf("unexpected organization: %s", cfg.Organization.Name)
	}
	if cfg.Runtime.Profile != "developer" {
		t.Fatalf("unexpected runtime profile: %s", cfg.Runtime.Profile)
	}
}

func TestLoadRejectsUnknownKey(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "forge.yaml")
	content := "organization:\n  name: 3Corners\nruntime:\n  mode: compose\n  profile: developer\n  surprise: unsafe\n"
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(p); err == nil {
		t.Fatal("expected unknown key to be rejected")
	}
}

func TestLoadRequiresRuntimeProfile(t *testing.T) {
	d := t.TempDir()
	p := filepath.Join(d, "forge.yaml")
	content := "organization:\n  name: 3Corners\nruntime:\n  mode: compose\n"
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(p); err == nil {
		t.Fatal("expected missing runtime.profile to be rejected")
	}
}
