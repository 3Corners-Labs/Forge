package intelligence_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/3corners/forge/internal/intelligence"
)

func TestLoadProviderActivationBindsExactFileWithoutResolvingCredential(t *testing.T) {
	now := time.Date(2026, time.July, 24, 8, 0, 0, 0, time.UTC)
	path := filepath.Join(t.TempDir(), "activation.json")
	raw := validActivationJSON(now)
	if err := os.WriteFile(path, []byte(raw), 0o600); err != nil {
		t.Fatal(err)
	}
	activation, err := intelligence.LoadProviderActivation(path, now)
	if err != nil {
		t.Fatal(err)
	}
	if activation.Provider != "openai" ||
		activation.AllowedModelIDs[0] != "gpt-5.6-terra" ||
		activation.CredentialRef != "env:FORGE_OPENAI_API_KEY" ||
		activation.ReasoningEffort != "medium" ||
		activation.Pricing.InputMicrosPerMillionTokens != 2_500_000 ||
		!strings.HasPrefix(activation.Digest, "sha256:") {
		t.Fatalf("unexpected activation: %#v", activation)
	}
	if err := os.WriteFile(path, []byte(raw+"\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	reformatted, err := intelligence.LoadProviderActivation(path, now)
	if err != nil {
		t.Fatal(err)
	}
	if reformatted.Digest == activation.Digest {
		t.Fatal("exact activation bytes were not bound into the digest")
	}
}

func TestLoadProviderActivationFailsClosedForUnknownExpiredAndNonRegularInput(t *testing.T) {
	now := time.Date(2026, time.July, 24, 8, 0, 0, 0, time.UTC)
	directory := t.TempDir()
	tests := []struct {
		name string
		raw  string
	}{
		{
			name: "unknown field",
			raw:  strings.Replace(validActivationJSON(now), `"provider": "openai",`, `"provider": "openai", "authority": true,`, 1),
		},
		{
			name: "expired",
			raw:  strings.Replace(validActivationJSON(now), now.Add(time.Hour).Format(time.RFC3339), now.Add(-time.Second).Format(time.RFC3339), 1),
		},
		{
			name: "multiple objects",
			raw:  validActivationJSON(now) + `{}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(directory, strings.ReplaceAll(test.name, " ", "-")+".json")
			if err := os.WriteFile(path, []byte(test.raw), 0o600); err != nil {
				t.Fatal(err)
			}
			if _, err := intelligence.LoadProviderActivation(path, now); err == nil {
				t.Fatal("ineligible activation was accepted")
			}
		})
	}
	target := filepath.Join(directory, "target.json")
	if err := os.WriteFile(target, []byte(validActivationJSON(now)), 0o600); err != nil {
		t.Fatal(err)
	}
	link := filepath.Join(directory, "activation-link.json")
	if err := os.Symlink(target, link); err != nil {
		t.Fatal(err)
	}
	if _, err := intelligence.LoadProviderActivation(link, now); err == nil {
		t.Fatal("symlinked activation was accepted")
	}
}

func validActivationJSON(now time.Time) string {
	return fmt.Sprintf(`{
  "schemaVersion": "forge.model-provider-activation.v1",
  "id": "forge.provider-activation.openai-terra-developer-trial-v1",
  "provider": "openai",
  "endpoint": "https://eu.api.openai.com/v1/responses",
  "adapterVersion": "forge.openai-discovery-adapter.v1",
  "allowedModelIds": ["gpt-5.6-terra"],
  "allowedPurposes": ["forge.product.discovery"],
  "allowedClassifications": ["Internal"],
  "credentialRef": "env:FORGE_OPENAI_API_KEY",
  "humanApprovalRef": "approval:provider-trial-1",
  "promptVersion": "forge.discovery.prompt.v1",
  "retention": "provider-default-up-to-30-days",
  "trainingUse": "disabled",
  "region": "eu",
  "storeProviderState": false,
  "reasoningEffort": "medium",
  "pricing": {
    "currency": "USD",
    "inputMicrosPerMillionTokens": 2500000,
    "outputMicrosPerMillionTokens": 15000000
  },
  "budgets": {
    "maxInputBytes": 32768,
    "maxOutputBytes": 32768,
    "maxTokens": 4096,
    "maxLatencyMs": 20000,
    "maxCostMicros": 250000,
    "maxTurns": 8,
    "maxSessionTokens": 32768,
    "maxSessionCostMicros": 1000000,
    "maxProjectTokensPerDay": 131072,
    "maxProjectCostMicrosPerDay": 4000000,
    "maxConcurrentRequests": 1
  },
  "activatedAt": %q,
  "expiresAt": %q
}`, now.Add(-time.Minute).Format(time.RFC3339), now.Add(time.Hour).Format(time.RFC3339))
}
