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
	UpdateProject(ctx context.Context, title, organizationName, desc string, project *Project) error

	// SetBackgroundImage ...
	SetBackgroundImage(ctx context.Context, key string, project *Project) error

	// SetFooterImage ...
	SetFooterImage(ctx context.Context, key string, project *Project) error

	// SetIconImage ...
	SetIconImage(ctx context.Context, key string, project *Project) error

	// SetStripImage ...
	SetStripImage(ctx context.Context, key string, project *Project) error
}

// UploadStore implements user upload related methods.
type UploadStore interface {
	// IsUploadExists ...
	IsUploadExists(ctx context.Context, user *User, filename, hash string) (bool, error)
	// SaveNewUpload ...
	SaveNewUpload(ctx context.Context, user *User, upload *Upload) error
	// LoadUpload ...
	LoadUpload(ctx context.Context, user *User, id int64) (*Upload, error)
	// LoadUploadByUUID ...
	LoadUploadByUUID(ctx context.Context, user *User, uuid string) (*Upload, error)
	// LoadUploads ...
	LoadUploads(ctx context.Context, user *User) ([]*Upload, error)
}

// Logic implements method for business logic.
type Logic interface {
	ProjectStore
	UploadStore
}
