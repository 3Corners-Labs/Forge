package tenant

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

type PostgresBindingSource struct{ db *sql.DB }

func NewPostgresBindingSource(db *sql.DB) (*PostgresBindingSource, error) {
	if db == nil {
		return nil, errors.New("tenant PostgreSQL connection is required")
	}
	return &PostgresBindingSource{db: db}, nil
}

func (source *PostgresBindingSource) Migrate(ctx context.Context) error {
	const migration = `
CREATE SCHEMA IF NOT EXISTS forge_tenant;
CREATE TABLE IF NOT EXISTS forge_tenant.schema_version(version integer PRIMARY KEY CHECK(version=1));
INSERT INTO forge_tenant.schema_version(version) VALUES(1) ON CONFLICT DO NOTHING;
CREATE TABLE IF NOT EXISTS forge_tenant.workspaces(tenant_id text NOT NULL,workspace_id text PRIMARY KEY,selector text UNIQUE NOT NULL);
CREATE TABLE IF NOT EXISTS forge_tenant.projects(workspace_id text NOT NULL,project_id text PRIMARY KEY,selector text NOT NULL,UNIQUE(workspace_id,selector));
CREATE TABLE IF NOT EXISTS forge_tenant.memberships(principal_id text NOT NULL,workspace_id text NOT NULL,PRIMARY KEY(principal_id,workspace_id));
CREATE TABLE IF NOT EXISTS forge_tenant.repository_bindings(
 tenant_id text NOT NULL,workspace_id text NOT NULL,project_id text NOT NULL,repository_id text NOT NULL,
 selector text NOT NULL,canonical_locator text NOT NULL,classification text NOT NULL,
 PRIMARY KEY(tenant_id,workspace_id,project_id,selector));`
	_, err := source.db.ExecContext(ctx, migration)
	return err
}

func (source *PostgresBindingSource) ResolveRepositoryBinding(ctx context.Context, tenantID, workspaceID, projectID, selector string) (RepositoryBindingDefinition, bool, error) {
	var value RepositoryBindingDefinition
	err := source.db.QueryRowContext(ctx, `SELECT tenant_id,workspace_id,project_id,repository_id,selector,canonical_locator,classification FROM forge_tenant.repository_bindings WHERE tenant_id=$1 AND workspace_id=$2 AND project_id=$3 AND selector=$4`, tenantID, workspaceID, projectID, selector).Scan(
		&value.TenantID, &value.WorkspaceID, &value.ProjectID, &value.RepositoryID,
		&value.Selector, &value.CanonicalLocator, &value.Classification,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return RepositoryBindingDefinition{}, false, nil
	}
	return value, err == nil, err
}

func LoadPostgresDirectory(ctx context.Context, authorityID string, db *sql.DB, now func() time.Time) (*Directory, error) {
	source, err := NewPostgresBindingSource(db)
	if err != nil {
		return nil, err
	}
	workspaceRows, err := db.QueryContext(ctx, `SELECT tenant_id,workspace_id,selector FROM forge_tenant.workspaces ORDER BY workspace_id`)
	if err != nil {
		return nil, err
	}
	var workspaces []WorkspaceDefinition
	for workspaceRows.Next() {
		var value WorkspaceDefinition
		if workspaceRows.Scan(&value.TenantID, &value.ID, &value.Selector) != nil {
			workspaceRows.Close()
			return nil, errors.New("tenant workspace state is invalid")
		}
		workspaces = append(workspaces, value)
	}
	workspaceRows.Close()
	projectRows, err := db.QueryContext(ctx, `SELECT workspace_id,project_id,selector FROM forge_tenant.projects ORDER BY project_id`)
	if err != nil {
		return nil, err
	}
	var projects []ProjectDefinition
	for projectRows.Next() {
		var value ProjectDefinition
		if projectRows.Scan(&value.WorkspaceID, &value.ID, &value.Selector) != nil {
			projectRows.Close()
			return nil, errors.New("tenant project state is invalid")
		}
		projects = append(projects, value)
	}
	projectRows.Close()
	membershipRows, err := db.QueryContext(ctx, `SELECT principal_id,workspace_id FROM forge_tenant.memberships ORDER BY principal_id,workspace_id`)
	if err != nil {
		return nil, err
	}
	var memberships []MembershipDefinition
	for membershipRows.Next() {
		var value MembershipDefinition
		if membershipRows.Scan(&value.PrincipalID, &value.WorkspaceID) != nil {
			membershipRows.Close()
			return nil, errors.New("tenant membership state is invalid")
		}
		memberships = append(memberships, value)
	}
	membershipRows.Close()
	return NewDirectory(authorityID, workspaces, projects, memberships, source, now)
}
