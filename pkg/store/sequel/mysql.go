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

const (
	_ = 1 << iota
	checkNilStruct
	checkZeroID
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

func (m *MySQL) finishTx(tx *sqlx.Tx, err error) error {
	if err != nil {
		tx.Rollback()
		return err
	}

	if commitErr := tx.Commit(); commitErr != nil {
		err = commitErr
	}

	return err
}

func (m *MySQL) insertQuery(ctx context.Context, query sq.InsertBuilder) (id int64, err error) {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return -1, err
	}

	defer func() { err = m.finishTx(tx, err) }()

	rawsql, args, err := query.ToSql()
	if err != nil {
		return -1, err
	}

	res, err := tx.ExecContext(ctx, rawsql, args...)
	if err != nil {
		return -1, err
	}

	id, err = res.LastInsertId()
	if err != nil {
		return -1, err
	}

	if id == 0 {
		return -1, store.ErrZeroID
	}

	return id, nil
}

func (m *MySQL) updateQuery(ctx context.Context, query sq.UpdateBuilder) (rows int64, err error) {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return -1, err
	}

	defer func() { err = m.finishTx(tx, err) }()

	rawsql, args, err := query.ToSql()
	if err != nil {
		return -1, err
	}

	res, err := tx.ExecContext(ctx, rawsql, args...)
	if err != nil {
		return -1, err
	}

	rows, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}

	if rows == 0 {
		return -1, store.ErrZeroRowsAffected
	}

	return rows, nil
}

func (m *MySQL) deleteQuery(ctx context.Context, query sq.DeleteBuilder) (rows int64, err error) {
	tx, err := m.db.BeginTxx(ctx, nil)
	if err != nil {
		return -1, err
	}

	defer func() { err = m.finishTx(tx, err) }()

	rawsql, args, err := query.ToSql()
	if err != nil {
		return -1, err
	}

	res, err := tx.ExecContext(ctx, rawsql, args...)
	if err != nil {
		return -1, err
	}

	rows, err = res.RowsAffected()
	if err != nil {
		return -1, err
	}

	if rows == 0 {
		return -1, store.ErrZeroRowsAffected
	}

	return rows, nil
}
