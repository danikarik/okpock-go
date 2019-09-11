package sequel

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

func checkProject(p *api.Project, opts byte) error {
	if (opts & checkNilStruct) != 0 {
		if p == nil {
			return store.ErrNilStruct
		}
	}

	if (opts & checkZeroID) != 0 {
		if p.ID == 0 {
			return store.ErrZeroID
		}
	}

	err := p.IsValid()
	if err != nil {
		return err
	}

	return nil
}

// IsProjectExists ...
func (m *MySQL) IsProjectExists(ctx context.Context, title, organizationName, desc string, passType api.PassType) (bool, error) {
	query := m.builder.Select("count(1)").
		From("projects").
		Where(sq.Eq{
			"title":             title,
			"organization_name": organizationName,
			"description":       desc,
			"pass_type":         passType,
		})

	cnt, err := m.countQuery(ctx, query)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// SaveNewProject ...
func (m *MySQL) SaveNewProject(ctx context.Context, user *api.User, project *api.Project) error {
	err := checkProject(project, checkNilStruct)
	if err != nil {
		return err
	}

	query := m.builder.Insert("projects").
		Columns(
			"title",
			"organization_name",
			"description",
			"pass_type",
			"created_at",
			"updated_at",
		).
		Values(
			project.Title,
			project.OrganizationName,
			project.Description,
			project.PassType,
			project.CreatedAt,
			project.UpdatedAt,
		)

	id, err := m.insertQuery(ctx, query)
	if err != nil {
		return err
	}
	project.ID = id

	query = m.builder.Insert("user_projects").
		Columns("user_id", "project_id").
		Values(user.ID, project.ID)

	_, err = m.insertQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// LoadProject ...
func (m *MySQL) LoadProject(ctx context.Context, user *api.User, id int64) (*api.Project, error) {
	if id == 0 {
		return nil, store.ErrZeroID
	}

	query := m.builder.Select("p.*").
		From("projects p").
		LeftJoin("user_projects up on up.project_id = p.id").
		Where(sq.Eq{
			"p.id":       id,
			"up.user_id": user.ID,
		})

	row, err := m.selectRowQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	var project = &api.Project{}

	err = row.StructScan(project)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return project, nil
}

// LoadProjects ...
func (m *MySQL) LoadProjects(ctx context.Context, user *api.User) ([]*api.Project, error) {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return nil, err
	}

	var projects = []*api.Project{}

	query := m.builder.Select("p.*").
		From("projects p").
		LeftJoin("user_projects up on up.project_id = p.id").
		Where(sq.Eq{"up.user_id": user.ID}).
		OrderBy("created_at desc")

	rows, err := m.selectQuery(ctx, query)
	if err == store.ErrNotFound {
		return projects, nil
	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var project = &api.Project{}

		err = rows.StructScan(project)
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		if err != nil {
			return nil, err
		}

		projects = append(projects, project)
	}

	return projects, nil
}

// UpdateProject ...
func (m *MySQL) UpdateProject(ctx context.Context, title, organizationName, desc string, project *api.Project) error {
	err := checkProject(project, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	project.Title = title
	project.OrganizationName = organizationName
	project.Description = desc
	project.UpdatedAt = time.Now()

	query := m.builder.Update("projects").
		Set("title", project.Title).
		Set("organization_name", project.OrganizationName).
		Set("description", project.Description).
		Set("updated_at", project.UpdatedAt).
		Where(sq.Eq{"id": project.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// SetBackgroundImage ...
func (m *MySQL) SetBackgroundImage(ctx context.Context, key string, project *api.Project) error {
	err := checkProject(project, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	project.BackgroundImage = key
	project.UpdatedAt = time.Now()

	query := m.builder.Update("projects").
		Set("background_image", project.BackgroundImage).
		Set("updated_at", project.UpdatedAt).
		Where(sq.Eq{"id": project.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// SetFooterImage ...
func (m *MySQL) SetFooterImage(ctx context.Context, key string, project *api.Project) error {
	err := checkProject(project, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	project.FooterImage = key
	project.UpdatedAt = time.Now()

	query := m.builder.Update("projects").
		Set("footer_image", project.FooterImage).
		Set("updated_at", project.UpdatedAt).
		Where(sq.Eq{"id": project.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// SetIconImage ...
func (m *MySQL) SetIconImage(ctx context.Context, key string, project *api.Project) error {
	err := checkProject(project, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	project.IconImage = key
	project.UpdatedAt = time.Now()

	query := m.builder.Update("projects").
		Set("icon_image", project.IconImage).
		Set("updated_at", project.UpdatedAt).
		Where(sq.Eq{"id": project.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// SetStripImage ...
func (m *MySQL) SetStripImage(ctx context.Context, key string, project *api.Project) error {
	err := checkProject(project, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	project.StripImage = key
	project.UpdatedAt = time.Now()

	query := m.builder.Update("projects").
		Set("strip_image", project.StripImage).
		Set("updated_at", project.UpdatedAt).
		Where(sq.Eq{"id": project.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
