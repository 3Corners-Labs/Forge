package verification

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type SourceFile struct {
	Path    string
	Digest  string
	Content string
}

// PrepareFreshCheckout independently reconstructs the verifier's read-only
// source view from admitted bytes. It deliberately does not reuse the worker
// workspace or import Execution implementation.
func PrepareFreshCheckout(runtimeRoot string, sources []SourceFile) (string, error) {
	if runtimeRoot == "" || !filepath.IsAbs(runtimeRoot) {
		return "", errors.New("absolute project-local runtime root is required")
	}
	if err := os.MkdirAll(runtimeRoot, 0o700); err != nil {
		return "", err
	}
	root, err := os.MkdirTemp(runtimeRoot, "forge-verifier-source-")
	if err != nil {
		return "", err
	}
	ok := false
	defer func() {
		if !ok {
			_ = os.RemoveAll(root)
		}
	}()
	for _, source := range sources {
		if !safeCheckoutPath(source.Path) || source.Digest != checkoutDigest([]byte(source.Content)) {
			return "", errors.New("source projection is invalid")
		}
		target := filepath.Join(root, filepath.FromSlash(source.Path))
		relative, err := filepath.Rel(root, target)
		if err != nil || relative == ".." || strings.HasPrefix(relative, ".."+string(filepath.Separator)) {
			return "", errors.New("source projection exceeds checkout")
		}
		if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
			return "", err
		}
		if err := os.Chmod(filepath.Dir(target), 0o755); err != nil {
			return "", err
		}
		if err := os.WriteFile(target, []byte(source.Content), 0o644); err != nil {
			return "", err
		}
		if err := os.Chmod(target, 0o644); err != nil {
			return "", err
		}
	}
	if err := os.Chmod(root, 0o755); err != nil {
		return "", err
	}
	ok = true
	return root, nil
}

func safeCheckoutPath(value string) bool {
	return value != "" && value != "." && !filepath.IsAbs(value) &&
		filepath.Clean(value) == filepath.FromSlash(value) &&
		!strings.HasPrefix(value, "../") && !strings.Contains(value, "/../") &&
		!strings.Contains(value, "\\")
}

func checkoutDigest(value []byte) string {
	sum := sha256.Sum256(value)
	return "sha256:" + hex.EncodeToString(sum[:])
}
