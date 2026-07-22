// Package tenant owns Developer Trial tenant, workspace, project, membership,
// and repository-binding resolution. It consumes verified identity facts but
// does not authenticate callers or evaluate final action policy.
package tenant

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/3corners/forge/internal/contracts/failure"
	"github.com/3corners/forge/internal/identity"
)

const (
	ScopeRequestContractVersion      = "forge.workspace-scope-request.v1"
	ScopeContractVersion             = "forge.workspace-scope.v1"
	BindingRequestContractVersion    = "forge.repository-binding-request.v1"
	RepositoryBindingContractVersion = "forge.repository-binding.v1"
)

type WorkspaceDefinition struct {
	TenantID string
	ID       string
	Selector string
}

type ProjectDefinition struct {
	WorkspaceID string
	ID          string
	Selector    string
}

type MembershipDefinition struct {
	PrincipalID string
	WorkspaceID string
}

type ScopeRequest struct {
	ContractVersion    string
	CorrelationID      string
	WorkspaceSelector  string
	ProjectSelector    string
	RepositorySelector string
}

type BindingRequest struct {
	ContractVersion string
	CorrelationID   string
	Scope           Scope
}

// RepositoryBindingDefinition is returned by a Tenant-owned persistence or
// contract-test adapter. It contains no provider credential.
type RepositoryBindingDefinition struct {
	TenantID         string
	WorkspaceID      string
	ProjectID        string
	RepositoryID     string
	Selector         string
	CanonicalLocator string
	Classification   string
}

type BindingSource interface {
	ResolveRepositoryBinding(context.Context, string, string, string, string) (RepositoryBindingDefinition, bool, error)
}

// Scope is minted only after membership resolution. Private fields prevent
// selectors or prompt content from constructing authoritative scope values.
type Scope struct {
	contractVersion    string
	authorityID        string
	tenantID           string
	workspaceID        string
	projectID          string
	principalID        string
	repositorySelector string
	principalExpiresAt time.Time
}

func (value Scope) ContractVersion() string { return value.contractVersion }
func (value Scope) TenantID() string        { return value.tenantID }
func (value Scope) WorkspaceID() string     { return value.workspaceID }
func (value Scope) ProjectID() string       { return value.projectID }
func (value Scope) PrincipalID() string     { return value.principalID }

// ValidateFor lets receiving capabilities verify that a scope was minted for
// the active principal without gaining access to Tenant-owned selectors or
// membership records.
func (value Scope) ValidateFor(principal identity.Principal, at time.Time, correlationID string) error {
	if err := principal.ValidateAt(at, correlationID); err != nil {
		return err
	}
	if value.contractVersion != ScopeContractVersion || value.authorityID == "" || value.tenantID == "" || value.workspaceID == "" || value.projectID == "" || value.principalID == "" || value.repositorySelector == "" {
		return unavailableScope(correlationID)
	}
	if value.principalID != principal.SubjectID() {
		return unavailableScope(correlationID)
	}
	if value.principalExpiresAt.IsZero() || !at.Before(value.principalExpiresAt) {
		return failure.New(failure.Expired, "authenticated scope is not active", correlationID, failure.RetryReauthenticate)
	}
	return nil
}

type RepositoryBinding struct {
	contractVersion  string
	authorityID      string
	tenantID         string
	workspaceID      string
	projectID        string
	repositoryID     string
	selector         string
	canonicalLocator string
	classification   string
}

func (value RepositoryBinding) ContractVersion() string  { return value.contractVersion }
func (value RepositoryBinding) TenantID() string         { return value.tenantID }
func (value RepositoryBinding) WorkspaceID() string      { return value.workspaceID }
func (value RepositoryBinding) ProjectID() string        { return value.projectID }
func (value RepositoryBinding) RepositoryID() string     { return value.repositoryID }
func (value RepositoryBinding) CanonicalLocator() string { return value.canonicalLocator }
func (value RepositoryBinding) Classification() string   { return value.classification }

// ValidateFor verifies exact Tenant-owned scope binding before another
// capability records or acts on the repository identity.
func (value RepositoryBinding) ValidateFor(scope Scope, correlationID string) error {
	if value.contractVersion != RepositoryBindingContractVersion || scope.contractVersion != ScopeContractVersion || value.authorityID == "" || scope.authorityID == "" || value.tenantID == "" || value.workspaceID == "" || value.projectID == "" || value.repositoryID == "" || value.selector == "" || value.canonicalLocator == "" || !validClassification(value.classification) {
		return unavailableScope(correlationID)
	}
	if value.authorityID != scope.authorityID || value.tenantID != scope.tenantID || value.workspaceID != scope.workspaceID || value.projectID != scope.projectID || value.selector != scope.repositorySelector {
		return unavailableScope(correlationID)
	}
	return nil
}

type Directory struct {
	authorityID string
	workspaces  map[string]WorkspaceDefinition
	projects    map[string]ProjectDefinition
	memberships map[string]map[string]struct{}
	bindings    BindingSource
	now         func() time.Time
}

func NewDirectory(authorityID string, workspaces []WorkspaceDefinition, projects []ProjectDefinition, memberships []MembershipDefinition, bindings BindingSource, now func() time.Time) (*Directory, error) {
	if authorityID == "" {
		return nil, errors.New("tenant authority ID is required")
	}
	if bindings == nil {
		return nil, errors.New("repository binding source is required")
	}
	if now == nil {
		return nil, errors.New("tenant clock is required")
	}
	directory := &Directory{
		authorityID: authorityID,
		workspaces:  make(map[string]WorkspaceDefinition, len(workspaces)),
		projects:    make(map[string]ProjectDefinition, len(projects)),
		memberships: make(map[string]map[string]struct{}),
		bindings:    bindings,
		now:         now,
	}
	workspaceIDs := make(map[string]struct{}, len(workspaces))
	tenantIDs := make(map[string]struct{}, 1)
	for _, workspace := range workspaces {
		if workspace.TenantID == "" || workspace.ID == "" || workspace.Selector == "" {
			return nil, errors.New("workspace definition has an empty required field")
		}
		if _, exists := directory.workspaces[workspace.Selector]; exists {
			return nil, fmt.Errorf("duplicate workspace selector %q", workspace.Selector)
		}
		if _, exists := workspaceIDs[workspace.ID]; exists {
			return nil, fmt.Errorf("duplicate workspace ID %q", workspace.ID)
		}
		directory.workspaces[workspace.Selector] = workspace
		workspaceIDs[workspace.ID] = struct{}{}
		tenantIDs[workspace.TenantID] = struct{}{}
		if len(tenantIDs) > 1 {
			return nil, errors.New("Developer Trial directory cannot contain multiple tenants")
		}
	}
	projectIDs := make(map[string]struct{}, len(projects))
	for _, project := range projects {
		if project.WorkspaceID == "" || project.ID == "" || project.Selector == "" {
			return nil, errors.New("project definition has an empty required field")
		}
		if _, exists := workspaceIDs[project.WorkspaceID]; !exists {
			return nil, fmt.Errorf("project %q references an unknown workspace", project.ID)
		}
		key := projectKey(project.WorkspaceID, project.Selector)
		if _, exists := directory.projects[key]; exists {
			return nil, fmt.Errorf("duplicate project selector %q", project.Selector)
		}
		if _, exists := projectIDs[project.ID]; exists {
			return nil, fmt.Errorf("duplicate project ID %q", project.ID)
		}
		directory.projects[key] = project
		projectIDs[project.ID] = struct{}{}
	}
	for _, membership := range memberships {
		if membership.PrincipalID == "" || membership.WorkspaceID == "" {
			return nil, errors.New("membership definition has an empty required field")
		}
		if _, exists := workspaceIDs[membership.WorkspaceID]; !exists {
			return nil, fmt.Errorf("membership references an unknown workspace %q", membership.WorkspaceID)
		}
		if directory.memberships[membership.PrincipalID] == nil {
			directory.memberships[membership.PrincipalID] = make(map[string]struct{})
		}
		directory.memberships[membership.PrincipalID][membership.WorkspaceID] = struct{}{}
	}
	return directory, nil
}

// ResolveWorkspaceScope derives authoritative IDs from a verified principal
// and untrusted selectors. Every unavailable scope returns the same safe
// denial and repository binding discovery has not occurred at this point.
func (directory *Directory) ResolveWorkspaceScope(principal identity.Principal, request ScopeRequest) (Scope, error) {
	if directory == nil || directory.now == nil || directory.bindings == nil {
		return Scope{}, failure.New(failure.DependencyUnavailable, "workspace resolution is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if request.ContractVersion != ScopeRequestContractVersion {
		return Scope{}, failure.New(failure.InvalidRequest, "workspace resolution request version is invalid", request.CorrelationID, failure.RetryNone)
	}
	if request.WorkspaceSelector == "" || request.ProjectSelector == "" || request.RepositorySelector == "" {
		return Scope{}, failure.New(failure.InvalidRequest, "workspace, project, and repository selectors are required", request.CorrelationID, failure.RetryNone)
	}
	if err := principal.ValidateAt(directory.now(), request.CorrelationID); err != nil {
		return Scope{}, err
	}
	workspace, exists := directory.workspaces[request.WorkspaceSelector]
	if !exists || !directory.hasMembership(principal.SubjectID(), workspace.ID) {
		return Scope{}, unavailableScope(request.CorrelationID)
	}
	project, exists := directory.projects[projectKey(workspace.ID, request.ProjectSelector)]
	if !exists {
		return Scope{}, unavailableScope(request.CorrelationID)
	}
	return Scope{
		contractVersion:    ScopeContractVersion,
		authorityID:        directory.authorityID,
		tenantID:           workspace.TenantID,
		workspaceID:        workspace.ID,
		projectID:          project.ID,
		principalID:        principal.SubjectID(),
		repositorySelector: request.RepositorySelector,
		principalExpiresAt: principal.ExpiresAt(),
	}, nil
}

// GetRepositoryBinding performs repository discovery only for a scope minted
// by this Directory and still covered by active principal authority.
func (directory *Directory) GetRepositoryBinding(ctx context.Context, request BindingRequest) (RepositoryBinding, error) {
	if directory == nil || directory.now == nil || directory.bindings == nil {
		return RepositoryBinding{}, failure.New(failure.DependencyUnavailable, "repository binding resolution is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if request.ContractVersion != BindingRequestContractVersion {
		return RepositoryBinding{}, failure.New(failure.InvalidRequest, "repository binding request version is invalid", request.CorrelationID, failure.RetryNone)
	}
	if ctx == nil {
		return RepositoryBinding{}, failure.New(failure.InvalidRequest, "repository binding context is required", request.CorrelationID, failure.RetryNone)
	}
	scope := request.Scope
	if scope.contractVersion != ScopeContractVersion || scope.authorityID != directory.authorityID || scope.tenantID == "" || scope.workspaceID == "" || scope.projectID == "" || scope.principalID == "" || scope.repositorySelector == "" {
		return RepositoryBinding{}, unavailableScope(request.CorrelationID)
	}
	if at := directory.now(); at.IsZero() || !at.Before(scope.principalExpiresAt) {
		return RepositoryBinding{}, failure.New(failure.Expired, "authenticated scope is not active", request.CorrelationID, failure.RetryReauthenticate)
	}
	definition, found, err := directory.bindings.ResolveRepositoryBinding(ctx, scope.tenantID, scope.workspaceID, scope.projectID, scope.repositorySelector)
	if err != nil {
		return RepositoryBinding{}, failure.New(failure.DependencyUnavailable, "repository binding resolution is unavailable", request.CorrelationID, failure.RetryBounded)
	}
	if !found {
		return RepositoryBinding{}, unavailableScope(request.CorrelationID)
	}
	if definition.TenantID != scope.tenantID || definition.WorkspaceID != scope.workspaceID || definition.ProjectID != scope.projectID || definition.Selector != scope.repositorySelector || definition.RepositoryID == "" || definition.CanonicalLocator == "" || !validClassification(definition.Classification) {
		return RepositoryBinding{}, failure.New(failure.Indeterminate, "repository binding could not be verified", request.CorrelationID, failure.RetryNone)
	}
	return RepositoryBinding{
		contractVersion:  RepositoryBindingContractVersion,
		authorityID:      directory.authorityID,
		tenantID:         definition.TenantID,
		workspaceID:      definition.WorkspaceID,
		projectID:        definition.ProjectID,
		repositoryID:     definition.RepositoryID,
		selector:         definition.Selector,
		canonicalLocator: definition.CanonicalLocator,
		classification:   definition.Classification,
	}, nil
}

func (directory *Directory) hasMembership(principalID, workspaceID string) bool {
	_, exists := directory.memberships[principalID][workspaceID]
	return exists
}

func projectKey(workspaceID, selector string) string {
	return workspaceID + "\x00" + selector
}

func validClassification(value string) bool {
	switch value {
	case "Public", "Internal", "Confidential", "Restricted":
		return true
	default:
		return false
	}
}

func unavailableScope(correlationID string) error {
	return failure.New(failure.Forbidden, "requested scope is not available to the authenticated principal", correlationID, failure.RetryNone)
}
