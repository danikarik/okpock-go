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

	if (opts & checkForeignID) != 0 {
		if p.OrganizationID == "" {
			return store.ErrZeroID
		}
	}

	if (opts & checkZeroID) != 0 {
		if p.ID == "" {
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
func (m *MySQL) IsProjectExists(ctx context.Context, orgID, desc string, passType api.PassType) (bool, error) {
	query := m.builder.Select("count(1)").
		From("projects").
		Where(sq.Eq{
			"organization_id": orgID,
			"description":     desc,
			"pass_type":       passType,
		})

	cnt, err := m.countQuery(ctx, query)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// SaveNewProject ...
func (m *MySQL) SaveNewProject(ctx context.Context, proj *api.Project) error {
	err := checkProject(proj, checkNilStruct|checkForeignID)
	if err != nil {
		return err
	}

	proj.CreatedAt = time.Now()
	proj.UpdatedAt = time.Now()

	query := m.builder.Insert("projects").
		Columns(
			"id",
			"organization_id",
			"description",
			"pass_type",
			"created_at",
			"updated_at",
		).
		Values(
			proj.ID,
			proj.OrganizationID,
			proj.Description,
			proj.PassType,
			proj.CreatedAt,
			proj.UpdatedAt,
		)

	err = m.insertQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// LoadProject ...
func (m *MySQL) LoadProject(ctx context.Context, id string) (*api.Project, error) {
	if id == "" {
		return nil, store.ErrZeroID
	}

	query := m.builder.Select("*").
		From("projects").
		Where(sq.Eq{"id": id})

	row, err := m.selectRowQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	var proj = &api.Project{}

	err = row.StructScan(proj)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return proj, nil
}

// LoadProjects ...
func (m *MySQL) LoadProjects(ctx context.Context, userID string) ([]*api.Project, error) {
	if userID == "" {
		return nil, store.ErrZeroID
	}

	var projects = []*api.Project{}

	query := m.builder.Select("p.*").
		From("projects p").
		LeftJoin("organizations o on o.id = p.organization_id").
		Where(sq.Eq{"o.user_id": userID})

	rows, err := m.selectQuery(ctx, query)
	if err == store.ErrNotFound {
		return projects, nil
	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var proj = &api.Project{}

		err = rows.StructScan(proj)
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		if err != nil {
			return nil, err
		}

		projects = append(projects, proj)
	}

	return projects, nil
}

// UpdateProjectDescription ...
func (m *MySQL) UpdateProjectDescription(ctx context.Context, desc string, proj *api.Project) error {
	err := checkProject(proj, checkNilStruct|checkZeroID|checkForeignID)
	if err != nil {
		return err
	}

	proj.Description = desc
	proj.UpdatedAt = time.Now()

	query := m.builder.Update("projects").
		Set("description", proj.Description).
		Set("updated_at", proj.UpdatedAt).
		Where(sq.Eq{"id": proj.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// SetBackgroundImage ...
func (m *MySQL) SetBackgroundImage(ctx context.Context, key string, proj *api.Project) error {
	err := checkProject(proj, checkNilStruct|checkZeroID|checkForeignID)
	if err != nil {
		return err
	}

	proj.BackgroundImage = key
	proj.UpdatedAt = time.Now()

	query := m.builder.Update("projects").
		Set("background_image", proj.BackgroundImage).
		Set("updated_at", proj.UpdatedAt).
		Where(sq.Eq{"id": proj.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// SetFooterImage ...
func (m *MySQL) SetFooterImage(ctx context.Context, key string, proj *api.Project) error {
	err := checkProject(proj, checkNilStruct|checkZeroID|checkForeignID)
	if err != nil {
		return err
	}

	proj.FooterImage = key
	proj.UpdatedAt = time.Now()

	query := m.builder.Update("projects").
		Set("footer_image", proj.FooterImage).
		Set("updated_at", proj.UpdatedAt).
		Where(sq.Eq{"id": proj.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// SetIconImage ...
func (m *MySQL) SetIconImage(ctx context.Context, key string, proj *api.Project) error {
	err := checkProject(proj, checkNilStruct|checkZeroID|checkForeignID)
	if err != nil {
		return err
	}

	proj.IconImage = key
	proj.UpdatedAt = time.Now()

	query := m.builder.Update("projects").
		Set("icon_image", proj.IconImage).
		Set("updated_at", proj.UpdatedAt).
		Where(sq.Eq{"id": proj.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// SetStripImage ...
func (m *MySQL) SetStripImage(ctx context.Context, key string, proj *api.Project) error {
	err := checkProject(proj, checkNilStruct|checkZeroID|checkForeignID)
	if err != nil {
		return err
	}

	proj.StripImage = key
	proj.UpdatedAt = time.Now()

	query := m.builder.Update("projects").
		Set("strip_image", proj.StripImage).
		Set("updated_at", proj.UpdatedAt).
		Where(sq.Eq{"id": proj.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
