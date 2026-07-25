package main

import (
	"context"
	"errors"
	"io/fs"
	"testing"

	workerruntime "github.com/3corners/forge/internal/worker"
)

func TestPreparationFailureCodeIsBounded(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{"workspace", workerruntime.ErrWorkspaceBoundary, "workspace-boundary"},
		{"source", workerruntime.ErrSourceChanged, "source-changed"},
		{"output", workerruntime.ErrOutputLimit, "output-limit"},
		{"permission", fs.ErrPermission, "permission-denied"},
		{"cancelled", context.Canceled, "cancelled"},
		{"unknown", errors.New("contains a sensitive local path"), "invalid-request-or-io"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := preparationFailureCode(test.err); got != test.want {
				t.Fatalf("preparationFailureCode() = %q, want %q", got, test.want)
			}
		})
	}
}
