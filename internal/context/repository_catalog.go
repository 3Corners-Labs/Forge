package context

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

	"github.com/3corners/forge/internal/delivery"
)

type RepositorySource struct {
	ID             string
	Path           string
	Authority      string
	Classification string
	Optional       bool
}

type SourceProjection struct {
	Path    string
	Digest  string
	Content string
}

// RepositoryCatalog is the Developer Trial adapter for an exact local Forge
// checkout. It exposes metadata through Context and separately returns only
// explicitly selected source bytes for Change Preparation.
type RepositoryCatalog struct {
	root     string
	revision string
	sources  map[string]RepositorySource
	now      func() time.Time
}

func NewRepositoryCatalog(root, revision string, sources []RepositorySource, now func() time.Time) (*RepositoryCatalog, error) {
	absolute, err := filepath.Abs(root)
	if err != nil || revision == "" || len(sources) == 0 || now == nil {
		return nil, errors.New("repository catalog configuration is invalid")
	}
	realRoot, err := filepath.EvalSymlinks(absolute)
	if err != nil {
		return nil, errors.New("repository catalog root is unavailable")
	}
	catalog := &RepositoryCatalog{
		root: realRoot, revision: revision, sources: make(map[string]RepositorySource, len(sources)), now: now,
	}
	for _, source := range sources {
		if source.ID == "" || !safeRepositoryPath(source.Path) || source.Authority == "" ||
			source.Classification == "" {
			return nil, errors.New("repository catalog source is invalid")
		}
		if _, duplicate := catalog.sources[source.ID]; duplicate {
			return nil, errors.New("repository catalog source identifiers must be unique")
		}
		catalog.sources[source.ID] = source
	}
	return catalog, nil
}

func (catalog *RepositoryCatalog) Revision() string {
	if catalog == nil {
		return ""
	}
	return catalog.revision
}

func (catalog *RepositoryCatalog) Discover(_ context.Context, request SourceDiscoveryRequest) ([]SourceRecord, error) {
	if catalog == nil || request.TenantID == "" || request.WorkspaceID == "" ||
		request.ProjectID == "" || request.RepositoryID == "" || len(request.RequestedSourceIDs) == 0 {
		return nil, errors.New("repository source discovery request is invalid")
	}
	freshness := catalog.now()
	records := make([]SourceRecord, 0, len(request.RequestedSourceIDs))
	for _, id := range request.RequestedSourceIDs {
		source, found := catalog.sources[id]
		if !found {
			return nil, errors.New("repository source is not registered")
		}
		content, exists, err := catalog.read(source.Path)
		if err != nil {
			return nil, err
		}
		if !exists && !source.Optional {
			return nil, errors.New("required repository source is unavailable")
		}
		if !exists {
			continue
		}
		records = append(records, SourceRecord{
			TenantID: request.TenantID, WorkspaceID: request.WorkspaceID, ProjectID: request.ProjectID,
			RepositoryID: request.RepositoryID, SourceID: id, SourceRevision: catalog.revision,
			Locator: source.Path, Classification: source.Classification, Authority: source.Authority,
			Purpose: delivery.M1ChangePurpose, FreshnessAt: &freshness, ContentDigest: repositoryDigest(content),
		})
	}
	return records, nil
}

// Project returns exact UTF-8 source bytes for registered paths. Optional
// missing targets remain absent so a proposal must bind them as new files.
func (catalog *RepositoryCatalog) Project(paths []string) ([]SourceProjection, error) {
	if catalog == nil || len(paths) == 0 {
		return nil, errors.New("repository source projection is invalid")
	}
	requested := slices.Clone(paths)
	slices.Sort(requested)
	requested = slices.Compact(requested)
	result := make([]SourceProjection, 0, len(requested))
	for _, sourcePath := range requested {
		registered := false
		for _, source := range catalog.sources {
			if source.Path != sourcePath {
				continue
			}
			registered = true
			content, exists, err := catalog.read(source.Path)
			if err != nil {
				return nil, err
			}
			if exists {
				result = append(result, SourceProjection{
					Path: source.Path, Digest: repositoryDigest(content), Content: string(content),
				})
			} else if !source.Optional {
				return nil, errors.New("required repository source is unavailable")
			}
			break
		}
		if !registered {
			return nil, errors.New("repository source path is not registered")
		}
	}
	return result, nil
}

func (catalog *RepositoryCatalog) read(relative string) ([]byte, bool, error) {
	target := filepath.Join(catalog.root, filepath.FromSlash(relative))
	within, err := filepath.Rel(catalog.root, target)
	if err != nil || within == ".." || strings.HasPrefix(within, ".."+string(filepath.Separator)) {
		return nil, false, errors.New("repository source exceeds catalog root")
	}
	info, err := os.Lstat(target)
	if errors.Is(err, os.ErrNotExist) {
		return nil, false, nil
	}
	if err != nil || !info.Mode().IsRegular() || info.Mode()&os.ModeSymlink != 0 || info.Size() > 256<<10 {
		return nil, false, errors.New("repository source is not an eligible bounded text file")
	}
	content, err := os.ReadFile(target)
	if err != nil || strings.ContainsRune(string(content), '\x00') {
		return nil, false, errors.New("repository source cannot be read as bounded text")
	}
	return content, true, nil
}

func safeRepositoryPath(value string) bool {
	return value != "" && value != "." && !filepath.IsAbs(value) &&
		filepath.Clean(value) == filepath.FromSlash(value) &&
		!strings.HasPrefix(value, "../") && !strings.Contains(value, "/../") &&
		!strings.Contains(value, "\\")
}

func repositoryDigest(value []byte) string {
	sum := sha256.Sum256(value)
	return "sha256:" + hex.EncodeToString(sum[:])
}

var _ SourceCatalog = (*RepositoryCatalog)(nil)
