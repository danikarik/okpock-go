package sequel

import (
	"context"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"
	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

func checkPassCard(p *api.PassCardInfo, opts byte) error {
	if (opts & checkNilStruct) != 0 {
		if p == nil {
			return store.ErrNilStruct
		}
		if p.Data == nil {
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

// SaveNewPassCard ...
func (m *MySQL) SaveNewPassCard(ctx context.Context, project *api.Project, passcard *api.PassCardInfo) error {
	err := checkPassCard(passcard, checkNilStruct)
	if err != nil {
		return err
	}

	query := m.builder.Insert("pass_cards").
		Columns(
			"raw_data",
			"created_at",
			"updated_at",
		).
		Values(
			passcard.Data,
			passcard.CreatedAt,
			passcard.UpdatedAt,
		)

	id, err := m.insertQuery(ctx, query)
	if err != nil {
		return err
	}
	passcard.ID = id

	query = m.builder.Insert("project_pass_cards").
		Columns("project_id", "pass_card_id").
		Values(project.ID, passcard.ID)

	_, err = m.insertQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) loadPassCard(ctx context.Context, query sq.SelectBuilder) (*api.PassCardInfo, error) {
	row, err := m.selectRowQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	var p = &api.PassCardInfo{}

	err = row.StructScan(p)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return p, nil
}

// LoadPassCard ...
func (m *MySQL) LoadPassCard(ctx context.Context, project *api.Project, id int64) (*api.PassCardInfo, error) {
	if id == 0 {
		return nil, store.ErrZeroID
	}

	query := m.builder.Select("pc.*").
		From("pass_cards pc").
		LeftJoin("project_pass_cards ppc on ppc.pass_card_id = pc.id").
		Where(sq.Eq{
			"ppc.project_id": project.ID,
			"pc.id":          id,
		})

	return m.loadPassCard(ctx, query)
}

// LoadPassCardBySerialNumber ...
func (m *MySQL) LoadPassCardBySerialNumber(ctx context.Context, project *api.Project, serialNumber string) (*api.PassCardInfo, error) {
	if serialNumber == "" {
		return nil, store.ErrEmptyQueryParam
	}

	query := m.builder.Select("pc.*").
		From("pass_cards pc").
		LeftJoin("project_pass_cards ppc on ppc.pass_card_id = pc.id").
		Where(sq.Eq{
			"ppc.project_id":                 project.ID,
			"pc.raw_data->>'$.serialNumber'": serialNumber,
		})

	return m.loadPassCard(ctx, query)
}

// LoadPassCards ...
func (m *MySQL) LoadPassCards(ctx context.Context, project *api.Project) ([]*api.PassCardInfo, error) {
	err := checkProject(project, checkNilStruct|checkZeroID)
	if err != nil {
		return nil, err
	}

	var passcards = []*api.PassCardInfo{}

	query := m.builder.Select("pc.*").
		From("pass_cards pc").
		LeftJoin("project_pass_cards ppc on ppc.pass_card_id = pc.id").
		Where(sq.Eq{"ppc.project_id": project.ID}).
		OrderBy("created_at desc")

	rows, err := m.selectQuery(ctx, query)
	if err == store.ErrNotFound {
		return passcards, nil
	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var passcard = &api.PassCardInfo{}

		err = rows.StructScan(passcard)
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		if err != nil {
			return nil, err
		}

		passcards = append(passcards, passcard)
	}

	return passcards, nil
}

// LoadPassCardsByBarcodeMessage ...
func (m *MySQL) LoadPassCardsByBarcodeMessage(ctx context.Context, project *api.Project, message string) ([]*api.PassCardInfo, error) {
	err := checkProject(project, checkNilStruct|checkZeroID)
	if err != nil {
		return nil, err
	}

	var passcards = []*api.PassCardInfo{}

	query := m.builder.Select("pc.*").
		From("pass_cards pc").
		LeftJoin("project_pass_cards ppc on ppc.pass_card_id = pc.id").
		Where(sq.Eq{"ppc.project_id": project.ID}).
		Where("JSON_CONTAINS(pc.raw_data->>'$.barcodes[*].message', JSON_ARRAY('" + message + "'))").
		OrderBy("created_at desc")

	rows, err := m.selectQuery(ctx, query)
	if err == store.ErrNotFound {
		return passcards, nil
	}
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var passcard = &api.PassCardInfo{}

		err = rows.StructScan(passcard)
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		if err != nil {
			return nil, err
		}

		passcards = append(passcards, passcard)
	}

	return passcards, nil
}

// UpdatePassCard ...
func (m *MySQL) UpdatePassCard(ctx context.Context, data *api.PassCard, passcard *api.PassCardInfo) error {
	err := checkPassCard(passcard, checkNilStruct|checkZeroID)
	if err != nil {
		return err
	}

	data.CopyFrom(passcard.Data)
	passcard.Data = data
	passcard.UpdatedAt = time.Now()

	query := m.builder.Update("pass_cards").
		Set("raw_data", passcard.Data).
		Set("updated_at", passcard.UpdatedAt).
		Where(sq.Eq{"id": passcard.ID})

	_, err = m.updateQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
