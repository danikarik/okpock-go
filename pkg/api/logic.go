package api

import "context"

// ProjectStore implements project related methods.
type ProjectStore interface {
	// IsProjectExists ...
	IsProjectExists(ctx context.Context, title, organizationName, desc string, passType PassType) (bool, error)

	// SaveNewProject ...
	SaveNewProject(ctx context.Context, user *User, project *Project) error

	// LoadProject ...
	LoadProject(ctx context.Context, user *User, id int64) (*Project, error)

	// LoadProjects ...
	LoadProjects(ctx context.Context, user *User) ([]*Project, error)

	// UpdateProject ...
	UpdateProject(ctx context.Context, project *Project) error

	// SetBackgroundImage ...
	SetBackgroundImage(ctx context.Context, key string, project *Project) error

	// SetFooterImage ...
	SetFooterImage(ctx context.Context, key string, project *Project) error

	// SetIconImage ...
	SetIconImage(ctx context.Context, key string, project *Project) error

	// SetStripImage ...
	SetStripImage(ctx context.Context, key string, project *Project) error
}

// Logic implements method for business logic.
type Logic interface {
	ProjectStore
}
