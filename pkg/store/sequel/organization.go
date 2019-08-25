package sequel

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

func checkOrganization(o *api.Organization, opts byte) error {
	if (opts & checkNilStruct) != 0 {
		if o == nil {
			return store.ErrNilStruct
		}
	}

	if (opts & checkForeignID) != 0 {
		if o.UserID == "" {
			return store.ErrZeroID
		}
	}

	if (opts & checkZeroID) != 0 {
		if o.ID == "" {
			return store.ErrZeroID
		}
	}

	err := o.IsValid()
	if err != nil {
		return err
	}

	return nil
}

// IsOrganizationExists ...
func (m *MySQL) IsOrganizationExists(ctx context.Context, userID, title string) (bool, error) {
	query := m.builder.Select("count(1)").
		From("organizations").
		Where(sq.Eq{
			"user_id": userID,
			"title":   title,
		})

	cnt, err := m.countQuery(ctx, query)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// SaveNewOrganization ...
func (m *MySQL) SaveNewOrganization(ctx context.Context, org *api.Organization) error {
	err := checkOrganization(org, checkNilStruct|checkForeignID)
	if err != nil {
		return err
	}

	org.CreatedAt = time.Now()
	org.UpdatedAt = time.Now()

	query := m.builder.Insert("organizations").
		Columns(
			"id",
			"user_id",
			"title",
			"description",
			"raw_metadata",
			"created_at",
			"updated_at",
		).
		Values(
			org.ID,
			org.UserID,
			org.Title,
			org.Description,
			org.MetaData,
			org.CreatedAt,
			org.UpdatedAt,
		)

	err = m.insertQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// LoadOrganization ...
func (m *MySQL) LoadOrganization(ctx context.Context, id string) (*api.Organization, error) {
	if id == "" {
		return nil, store.ErrZeroID
	}

	query := m.builder.Select("*").
		From("organizations").
		Where(sq.Eq{"id": id})

	row, err := m.selectRowQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	var org = &api.Organization{}

	err = row.StructScan(org)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return org, nil
}

// LoadOrganizations ...
func (m *MySQL) LoadOrganizations(ctx context.Context, userID string) ([]*api.Organization, error) {
	if userID == "" {
		return nil, store.ErrZeroID
	}

	var orgs = []*api.Organization{}

	query := m.builder.Select("*").
		From("organizations").
		Where(sq.Eq{"user_id": userID}).
		OrderBy("created_at desc")

	rows, err := m.selectQuery(ctx, query)
	if err == store.ErrNotFound {
		return orgs, nil
	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var org = &api.Organization{}

		err = rows.StructScan(org)
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		if err != nil {
			return nil, err
		}

		orgs = append(orgs, org)
	}

	return orgs, nil
}

// UpdateOrganizationDescription ...
func (m *MySQL) UpdateOrganizationDescription(ctx context.Context, desc string, org *api.Organization) error {
	err := checkOrganization(org, checkNilStruct|checkZeroID|checkForeignID)
	if err != nil {
		return err
	}

	org.Description = desc
	org.UpdatedAt = time.Now()

	query := m.builder.Update("organizations").
		Set("description", org.Description).
		Set("updated_at", org.UpdatedAt).
		Where(sq.Eq{"id": org.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

// UpdateOrganizationMetaData ...
func (m *MySQL) UpdateOrganizationMetaData(ctx context.Context, data map[string]interface{}, org *api.Organization) error {
	err := checkOrganization(org, checkNilStruct|checkZeroID|checkForeignID)
	if err != nil {
		return err
	}

	org.MetaData = data
	org.UpdatedAt = time.Now()

	query := m.builder.Update("organizations").
		Set("raw_metadata", org.MetaData).
		Set("updated_at", org.UpdatedAt).
		Where(sq.Eq{"id": org.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
