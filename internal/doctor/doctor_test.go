package doctor

import (
	"bytes"
	"strings"
	"testing"
)

func TestPrintToReportsEveryCheckWithoutRevealingConfiguration(t *testing.T) {
	var output bytes.Buffer
	ok := PrintTo(&output, []Check{
		{Name: "Docker daemon", OK: true},
		{Name: "FORGE_M1_POSTGRES_DSN", OK: false},
	})
	if ok {
		t.Fatal("expected missing configuration to fail the doctor")
	}
	text := output.String()
	if !strings.Contains(text, "Docker daemon") || !strings.Contains(text, "OK") {
		t.Fatalf("expected successful daemon check, got %q", text)
	}
	if !strings.Contains(text, "FORGE_M1_POSTGRES_DSN") || !strings.Contains(text, "MISSING") {
		t.Fatalf("expected missing configuration check, got %q", text)
	}
}
