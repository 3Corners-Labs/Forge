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

func TestValidateDeveloperBootstrap(t *testing.T) {
	var cfg Config
	cfg.Runtime.Mode = DeveloperRuntimeMode
	cfg.Runtime.Profile = DeveloperRuntimeProfile
	if err := cfg.ValidateDeveloperBootstrap(); err != nil {
		t.Fatalf("expected developer configuration to pass: %v", err)
	}
}

func TestValidateDeveloperBootstrapRejectsNonDeveloperConfiguration(t *testing.T) {
	tests := []struct {
		name         string
		mode         string
		profile      string
		intelligence bool
	}{
		{name: "connected self-hosted", mode: "compose", profile: "connected-self-hosted"},
		{name: "restricted self-hosted", mode: "compose", profile: "restricted-self-hosted"},
		{name: "disconnected self-hosted", mode: "compose", profile: "disconnected-self-hosted"},
		{name: "Forge SaaS", mode: "compose", profile: "forge-saas"},
		{name: "unsupported runtime", mode: "kubernetes", profile: "developer"},
		{name: "ungoverned intelligence", mode: "compose", profile: "developer", intelligence: true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var cfg Config
			cfg.Runtime.Mode = test.mode
			cfg.Runtime.Profile = test.profile
			cfg.Intelligence.Enabled = test.intelligence
			if err := cfg.ValidateDeveloperBootstrap(); err == nil {
				t.Fatal("expected non-developer configuration to be rejected")
			}
		})
	}
}
