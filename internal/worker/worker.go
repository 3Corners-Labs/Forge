// Package worker implements the deterministic M1 change-preparation runtime.
// It deliberately has no control-plane dependency and never invokes a shell.
package worker

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	workercontract "github.com/3corners/forge/internal/contracts/worker"
)

var (
	ErrWorkspaceBoundary = errors.New("worker workspace boundary violation")
	ErrSourceChanged     = errors.New("worker source digest changed")
	ErrOutputLimit       = errors.New("worker output limit exceeded")
)

type Engine struct {
	now func() time.Time
}

func NewEngine(now func() time.Time) (*Engine, error) {
	if now == nil {
		return nil, errors.New("worker clock is required")
	}
	return &Engine{now: now}, nil
}

func (engine *Engine) Prepare(ctx context.Context, request workercontract.Request, workspace string) (workercontract.CandidateBundle, workercontract.Result, error) {
	started := engine.now()
	if ctx == nil || workspace == "" || started.IsZero() {
		return workercontract.CandidateBundle{}, workercontract.Result{}, workercontract.ErrInvalidContract
	}
	if err := request.Validate(started); err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	root, err := secureRoot(workspace)
	if err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	files := make([]workercontract.CandidateFile, 0, len(request.Replacements))
	var outputBytes int64
	for _, replacement := range request.Replacements {
		if err := ctx.Err(); err != nil {
			return workercontract.CandidateBundle{}, workercontract.Result{}, err
		}
		target, err := secureTarget(root, replacement.Path)
		if err != nil {
			return workercontract.CandidateBundle{}, workercontract.Result{}, err
		}
		sourceDigest := ""
		current, readErr := os.ReadFile(target)
		switch {
		case readErr == nil:
			sourceDigest = digest(current)
		case errors.Is(readErr, os.ErrNotExist):
		default:
			return workercontract.CandidateBundle{}, workercontract.Result{}, ErrWorkspaceBoundary
		}
		if replacement.ExpectedSourceDigest != "" && replacement.ExpectedSourceDigest != sourceDigest {
			return workercontract.CandidateBundle{}, workercontract.Result{}, ErrSourceChanged
		}
		outputBytes += int64(len(replacement.Content))
		if outputBytes > request.Limits.MaxOutputBytes {
			return workercontract.CandidateBundle{}, workercontract.Result{}, ErrOutputLimit
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o700); err != nil {
			return workercontract.CandidateBundle{}, workercontract.Result{}, err
		}
		if err := rejectSymlinkPath(root, filepath.Dir(target)); err != nil {
			return workercontract.CandidateBundle{}, workercontract.Result{}, err
		}
		if err := os.WriteFile(target, replacement.Content, 0o600); err != nil {
			return workercontract.CandidateBundle{}, workercontract.Result{}, err
		}
		files = append(files, workercontract.CandidateFile{
			Path: replacement.Path, SourceDigest: sourceDigest,
			Content: slices.Clone(replacement.Content),
		})
	}
	bundle, err := workercontract.BindBundle(workercontract.CandidateBundle{
		SchemaVersion: workercontract.BundleVersion, SourceRevision: request.SourceRevision,
		PlanDigest: request.PlanDigest, Files: files,
	})
	if err != nil {
		return workercontract.CandidateBundle{}, workercontract.Result{}, err
	}
	changed := make([]string, 0, len(bundle.Files))
	for _, file := range bundle.Files {
		changed = append(changed, file.Path)
	}
	completed := engine.now()
	result := workercontract.Result{
		SchemaVersion: workercontract.ResultVersion, RequestID: request.ID, RequestDigest: request.Digest,
		CandidateDigest: bundle.Digest, TreeDigest: workercontract.CandidateTreeDigest(bundle),
		ChangedPaths: changed, BundleDigest: bundle.Digest,
		ResourceUse: workercontract.ResourceUse{
			DurationMillis: completed.Sub(started).Milliseconds(), OutputBytes: outputBytes,
			ChangedFiles: int64(len(changed)),
		},
		CleanupOutcome:     "pending-launcher-cleanup",
		EvidenceReferences: []string{"worker-request:" + request.Digest, "candidate:" + bundle.Digest},
		CompletedAt:        completed,
	}
	return bundle, result, nil
}

func secureRoot(workspace string) (string, error) {
	root, err := filepath.Abs(workspace)
	if err != nil {
		return "", ErrWorkspaceBoundary
	}
	root, err = filepath.EvalSymlinks(root)
	if err != nil {
		return "", ErrWorkspaceBoundary
	}
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return "", ErrWorkspaceBoundary
	}
	return root, nil
}

func secureTarget(root, relative string) (string, error) {
	target := filepath.Join(root, filepath.FromSlash(relative))
	rel, err := filepath.Rel(root, target)
	if err != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) || filepath.IsAbs(rel) {
		return "", ErrWorkspaceBoundary
	}
	if info, err := os.Lstat(target); err == nil && info.Mode()&os.ModeSymlink != 0 {
		return "", ErrWorkspaceBoundary
	} else if err != nil && !errors.Is(err, os.ErrNotExist) {
		return "", ErrWorkspaceBoundary
	}
	if err := rejectSymlinkPath(root, filepath.Dir(target)); err != nil {
		return "", err
	}
	return target, nil
}

func rejectSymlinkPath(root, target string) error {
	relative, err := filepath.Rel(root, target)
	if err != nil || strings.HasPrefix(relative, "..") {
		return ErrWorkspaceBoundary
	}
	current := root
	for _, part := range strings.Split(relative, string(filepath.Separator)) {
		if part == "" || part == "." {
			continue
		}
		current = filepath.Join(current, part)
		info, err := os.Lstat(current)
		if errors.Is(err, os.ErrNotExist) {
			continue
		}
		if err != nil || info.Mode()&os.ModeSymlink != 0 {
			return ErrWorkspaceBoundary
		}
	}
	return nil
}

func digest(value []byte) string {
	sum := sha256.Sum256(value)
	return "sha256:" + hex.EncodeToString(sum[:])
}
