package sequel

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/danikarik/okpock/pkg/store"
	sqlx "github.com/jmoiron/sqlx"
)

const (
	// TimeFormat - MySQL Timestamp Layout Format
	TimeFormat = "2006-01-02 15:04:05"
)

// New returns MySQL store implementation.
func New(db *sqlx.DB) *MySQL {
	return &MySQL{
		db:      db,
		builder: sq.StatementBuilder.PlaceholderFormat(sq.Question),
		cacher:  sq.NewStmtCacher(db),
	}
}

// MySQL implements store specs using `MySQL` database.
type MySQL struct {
	db      *sqlx.DB
	builder sq.StatementBuilderType
	cacher  sq.DBProxyContext
}

func (m *MySQL) countQuery(ctx context.Context, query sq.SelectBuilder) (int64, error) {
	var cnt int64

	err := query.RunWith(m.cacher).QueryRowContext(ctx).Scan(&cnt)
	if err != nil {
		return -1, err
	}

	return cnt, nil
}

func (m *MySQL) selectQuery(ctx context.Context, query sq.SelectBuilder) (*sqlx.Rows, error) {
	rawsql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	rows, err := m.db.QueryxContext(ctx, rawsql, args...)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		return nil, err
	}

	return rows, nil
}

func (m *MySQL) selectRowQuery(ctx context.Context, query sq.SelectBuilder) (*sqlx.Row, error) {
	rawsql, args, err := query.ToSql()
	if err != nil {
		return nil, err
	}

	return m.db.QueryRowxContext(ctx, rawsql, args...), nil
}

func (m *MySQL) scanQuery(ctx context.Context, query sq.SelectBuilder, v interface{}) error {
	err := query.RunWith(m.cacher).QueryRowContext(ctx).Scan(&v)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) insertQuery(ctx context.Context, query sq.InsertBuilder) (int64, error) {
	res, err := query.RunWith(m.cacher).ExecContext(ctx)
	if err != nil {
		return -1, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return -1, err
	}

	if id == 0 {
		return -1, store.ErrZeroID
	}

	return id, nil
}

func (m *MySQL) updateQuery(ctx context.Context, query sq.UpdateBuilder) (int64, error) {
	res, err := query.RunWith(m.cacher).ExecContext(ctx)
	if err != nil {
		return -1, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	if rows == 0 {
		return -1, store.ErrZeroRowsAffected
	}

	return rows, nil
}

func (m *MySQL) deleteQuery(ctx context.Context, query sq.DeleteBuilder) (int64, error) {
	res, err := query.RunWith(m.cacher).ExecContext(ctx)
	if err != nil {
		return -1, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}

	if rows == 0 {
		return -1, store.ErrZeroRowsAffected
	}

	return rows, nil
}
