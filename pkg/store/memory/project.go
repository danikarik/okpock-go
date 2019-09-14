package memory

import (
	"context"
	"time"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

// IsProjectExists ...
func (m *Memory) IsProjectExists(ctx context.Context, title, organizationName, desc string, passType api.PassType) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, p := range m.projects {
		if p.Title == title &&
			p.OrganizationName == organizationName &&
			p.Description == desc &&
			p.PassType == passType {
			return true, nil
		}
	}

	return false, nil
}

// SaveNewProject ...
func (m *Memory) SaveNewProject(ctx context.Context, user *api.User, project *api.Project) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.projects[project.ID] = project
	m.userProjects[project.ID] = user.ID

	return nil
}

// LoadProject ...
func (m *Memory) LoadProject(ctx context.Context, user *api.User, id int64) (*api.Project, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	project, ok := m.projects[id]
	if !ok {
		return nil, store.ErrNotFound
	}

	return project, nil
}

// LoadProjects ...
func (m *Memory) LoadProjects(ctx context.Context, user *api.User) ([]*api.Project, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	projects := []*api.Project{}
	for projectID, userID := range m.userProjects {
		if userID == user.ID {
			projects = append(projects, m.projects[projectID])
		}
	}

	return projects, nil
}

// UpdateProject ...
func (m *Memory) UpdateProject(ctx context.Context, title, organizationName, desc string, project *api.Project) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	project.Title = title
	project.OrganizationName = organizationName
	project.Description = desc
	project.UpdatedAt = time.Now()

	m.projects[project.ID] = project

	return nil
}

// SetBackgroundImage ...
func (m *Memory) SetBackgroundImage(ctx context.Context, key string, project *api.Project) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	project.BackgroundImage = key
	project.UpdatedAt = time.Now()

	m.projects[project.ID] = project

	return nil
}

// SetFooterImage ...
func (m *Memory) SetFooterImage(ctx context.Context, key string, project *api.Project) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	project.FooterImage = key
	project.UpdatedAt = time.Now()

	m.projects[project.ID] = project

	return nil
}

// SetIconImage ...
func (m *Memory) SetIconImage(ctx context.Context, key string, project *api.Project) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	project.IconImage = key
	project.UpdatedAt = time.Now()

	m.projects[project.ID] = project

	return nil
}

// SetLogoImage ...
func (m *Memory) SetLogoImage(ctx context.Context, key string, project *api.Project) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	project.LogoImage = key
	project.UpdatedAt = time.Now()

	m.projects[project.ID] = project

	return nil
}

// SetStripImage ...
func (m *Memory) SetStripImage(ctx context.Context, key string, project *api.Project) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	project.StripImage = key
	project.UpdatedAt = time.Now()

	m.projects[project.ID] = project

	return nil
}
