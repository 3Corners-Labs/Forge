package cli

import (
	"path/filepath"
	"testing"
)

func TestRunComposeWithConfigInvokesDeveloperActions(t *testing.T) {
	path := filepath.Join("testdata", "developer.yaml")
	for _, action := range []string{"up", "down"} {
		t.Run(action, func(t *testing.T) {
			calls := 0
			gotAction := ""
			err := runComposeWithConfig(path, action, func(action string) error {
				calls++
				gotAction = action
				return nil
			})
			if err != nil {
				t.Fatalf("expected developer action to pass: %v", err)
			}
			if calls != 1 || gotAction != action {
				t.Fatalf("expected one %q invocation, got calls=%d action=%q", action, calls, gotAction)
			}
		})
	}
}

func TestRunComposeWithConfigRejectsProductionShapedProfileBeforeSideEffect(t *testing.T) {
	path := filepath.Join("testdata", "production-shaped.yaml")
	for _, action := range []string{"up", "down"} {
		t.Run(action, func(t *testing.T) {
			called := false
			err := runComposeWithConfig(path, action, func(string) error {
				called = true
				return nil
			})
			if err == nil {
				t.Fatal("expected production-shaped profile to be rejected")
			}
			if called {
				t.Fatal("compose side effect must not run after profile rejection")
			}
		})
	}
}

func TestRunComposeWithConfigRejectsMissingConfigurationBeforeSideEffect(t *testing.T) {
	called := false
	err := runComposeWithConfig(filepath.Join(t.TempDir(), "missing.yaml"), "up", func(string) error {
		called = true
		return nil
	})
	if err == nil {
		t.Fatal("expected missing configuration to be rejected")
	}
	if called {
		t.Fatal("compose side effect must not run when configuration cannot be loaded")
	}
}
