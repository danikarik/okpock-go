package api

import "context"

// BusinessLogic implements app's main logic.
type BusinessLogic interface {
	// IsOrganizationExists ...
	// TODO: description
	IsOrganizationExists(ctx context.Context, title string, userID int64) (bool, error)
	// SaveNewOrganization ...
	// TODO: description
	SaveNewOrganization(ctx context.Context, org *Organization) error
	// LoadOrganization ...
	// TODO: description
	LoadOrganization(ctx context.Context, id int64) (*Organization, error)
	// LoadOrganizations ...
	// TODO: description
	LoadOrganizations(ctx context.Context, userID int64) (*[]Organization, error)
	// UpdateOrganizationDescription ...
	// TODO: description
	UpdateOrganizationDescription(ctx context.Context, desc string, org *Organization) error
	// UpdateOrganizationMetaData ...
	// TODO: description
	UpdateOrganizationMetaData(ctx context.Context, data map[string]interface{}, org *Organization) error

	// IsProjectExists ...
	// TODO: description
	IsProjectExists(ctx context.Context, desc string, orgID int64, passType PassType) (bool, error)
	// SaveNewProject ...
	// TODO: description
	SaveNewProject(ctx context.Context, proj *Project) error
	// LoadProject ...
	// TODO: description
	LoadProject(ctx context.Context, id int64) (*Project, error)
	// LoadProjects ...
	// TODO: description
	LoadProjects(ctx context.Context, userID int64) (*[]Project, error)
	// UpdateProjectDescription ...
	// TODO: description
	UpdateProjectDescription(ctx context.Context, desc string, proj *Project) error
	// SetBackgroundImage ...
	// TODO: description
	SetBackgroundImage(ctx context.Context, key string, proj *Project) error
	// SetFooterImage ...
	// TODO: description
	SetFooterImage(ctx context.Context, key string, proj *Project) error
	// SetIconImage ...
	// TODO: description
	SetIconImage(ctx context.Context, key string, proj *Project) error
	// SetStripImage ...
	// TODO: description
	SetStripImage(ctx context.Context, key string, proj *Project) error
}
