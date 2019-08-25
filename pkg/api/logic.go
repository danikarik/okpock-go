package api

import "context"

// Logic implements method for business logic.
type Logic interface {
	// IsProjectExists ...
	// TODO: description
	IsProjectExists(ctx context.Context, title, name, desc string, passType PassType) (bool, error)

	// SaveNewProject ...
	// TODO: description
	SaveNewProject(ctx context.Context, user *User, project *Project) error

	// LoadProject ...
	// TODO: description
	LoadProject(ctx context.Context, user *User, id string) (*Project, error)

	// LoadProjects ...
	// TODO: description
	LoadProjects(ctx context.Context, user *User) ([]*Project, error)

	// UpdateProject ...
	// TODO: description
	UpdateProject(ctx context.Context, project *Project) error

	// SetBackgroundImage ...
	// TODO: description
	SetBackgroundImage(ctx context.Context, key string, project *Project) error

	// SetFooterImage ...
	// TODO: description
	SetFooterImage(ctx context.Context, key string, project *Project) error

	// SetIconImage ...
	// TODO: description
	SetIconImage(ctx context.Context, key string, project *Project) error

	// SetStripImage ...
	// TODO: description
	SetStripImage(ctx context.Context, key string, project *Project) error
}
