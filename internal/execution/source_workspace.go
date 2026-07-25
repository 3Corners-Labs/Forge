package execution

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"

	changecontract "github.com/3corners/forge/internal/contracts/change"
)

// PrepareSourceWorkspace creates a purpose-limited worker workspace from the
// exact source bytes admitted by Context and Change Preparation. It does not
// copy a repository, follow symlinks, or execute repository content.
func PrepareSourceWorkspace(runtimeRoot string, sources []changecontract.SourceFile) (string, error) {
	return prepareSourceWorkspace(runtimeRoot, "forge-worker-source-", sources)
}

func prepareSourceWorkspace(runtimeRoot, prefix string, sources []changecontract.SourceFile) (string, error) {
	if runtimeRoot == "" || !filepath.IsAbs(runtimeRoot) {
		return "", errors.New("absolute project-local runtime root is required")
	}
	if err := os.MkdirAll(runtimeRoot, 0o700); err != nil {
		return "", err
	}
	root, err := os.MkdirTemp(runtimeRoot, prefix)
	if err != nil {
		return "", err
	}
	// The isolated worker runs as a fixed non-root UID. Only this purpose-built,
	// ephemeral mount is made writable; the repository itself is never mounted.
	if err := os.Chmod(root, 0o777); err != nil {
		return "", err
	}
	ok := false
	defer func() {
		if !ok {
			_ = os.RemoveAll(root)
		}
	}()
	for _, source := range sources {
		if !safeSourcePath(source.Path) || source.Digest != sourceDigest([]byte(source.Content)) {
			return "", errors.New("source projection is invalid")
		}
		target := filepath.Join(root, filepath.FromSlash(source.Path))
		relative, err := filepath.Rel(root, target)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return "", errors.New("source projection exceeds workspace")
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o777); err != nil {
			return "", err
		}
		if err := os.Chmod(filepath.Dir(target), 0o777); err != nil {
			return "", err
		}
		if err := os.WriteFile(target, []byte(source.Content), 0o666); err != nil {
			return "", err
		}
		if err := os.Chmod(target, 0o666); err != nil {
			return "", err
		}
	}
	ok = true
	return root, nil
}

func safeSourcePath(value string) bool {
	return value != "" && value != "." && !filepath.IsAbs(value) &&
		filepath.Clean(value) == filepath.FromSlash(value) &&
		!strings.HasPrefix(value, "../") && !strings.Contains(value, "/../") &&
		!strings.Contains(value, "\\")
}

func sourceDigest(value []byte) string {
	sum := sha256.Sum256(value)
	return "sha256:" + hex.EncodeToString(sum[:])
}
