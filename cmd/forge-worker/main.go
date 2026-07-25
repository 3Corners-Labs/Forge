package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"time"

	workercontract "github.com/3corners/forge/internal/contracts/worker"
	workerruntime "github.com/3corners/forge/internal/worker"
)

func main() {
	requestPath := flag.String("request", "", "closed worker request")
	workspace := flag.String("workspace", "", "authorized ephemeral checkout")
	bundlePath := flag.String("bundle", "", "candidate bundle output")
	resultPath := flag.String("result", "", "structured result output")
	flag.Parse()
	if *requestPath == "" || *workspace == "" || *bundlePath == "" || *resultPath == "" {
		exit("worker arguments are invalid")
	}
	data, err := os.ReadFile(*requestPath)
	if err != nil {
		exit("worker request is unavailable")
	}
	request, err := workercontract.DecodeRequest(data)
	if err != nil {
		exit("worker request is invalid")
	}
	engine, _ := workerruntime.NewEngine(func() time.Time { return time.Now().UTC() })
	bundle, result, err := engine.Prepare(context.Background(), request, *workspace)
	if err != nil {
		exit("worker preparation failed: " + preparationFailureCode(err))
	}
	bundleBytes, _ := workercontract.EncodeCanonical(bundle)
	resultBytes, _ := workercontract.EncodeCanonical(result)
	if err := os.WriteFile(*bundlePath, bundleBytes, 0o644); err != nil {
		exit("worker bundle could not be written")
	}
	if err := os.WriteFile(*resultPath, resultBytes, 0o644); err != nil {
		exit("worker result could not be written")
	}
}

func exit(message string) {
	_, _ = fmt.Fprintln(os.Stderr, message)
	os.Exit(1)
}

// preparationFailureCode keeps diagnostics useful without exposing workspace
// paths, request content, or provider data from the isolated runtime.
func preparationFailureCode(err error) string {
	switch {
	case errors.Is(err, workerruntime.ErrWorkspaceBoundary):
		return "workspace-boundary"
	case errors.Is(err, workerruntime.ErrSourceChanged):
		return "source-changed"
	case errors.Is(err, workerruntime.ErrOutputLimit):
		return "output-limit"
	case errors.Is(err, fs.ErrPermission):
		return "permission-denied"
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return "cancelled"
	default:
		return "invalid-request-or-io"
	}
}
