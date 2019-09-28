package sequel

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store"
)

func checkUpload(u *api.Upload, opts byte) error {
	if (opts & checkNilStruct) != 0 {
		if u == nil {
			return store.ErrNilStruct
		}
	}

	if (opts & checkZeroID) != 0 {
		if u.ID == 0 {
			return store.ErrZeroID
		}
	}

	err := u.IsValid()
	if err != nil {
		return err
	}

	return nil
}

// IsUploadExists ...
func (m *MySQL) IsUploadExists(ctx context.Context, user *api.User, filename, hash string) (bool, error) {
	query := m.builder.Select("count(1)").
		From("uploads u").
		LeftJoin("user_uploads uu on uu.upload_id = u.id").
		Where(sq.Eq{
			"uu.user_id": user.ID,
			"u.filename": filename,
			"u.hash":     hash,
		})

	cnt, err := m.countQuery(ctx, query)
	if err != nil {
		return false, err
	}

	return cnt > 0, nil
}

// SaveNewUpload ...
func (m *MySQL) SaveNewUpload(ctx context.Context, user *api.User, upload *api.Upload) error {
	err := checkUpload(upload, checkNilStruct)
	if err != nil {
		return err
	}

	query := m.builder.Insert("uploads").
		Columns(
			"uuid",
			"filename",
			"hash",
			"created_at",
		).
		Values(
			upload.UUID,
			upload.Filename,
			upload.Hash,
			upload.CreatedAt,
		)

	id, err := m.insertQuery(ctx, query)
	if err != nil {
		return err
	}
	upload.ID = id

	query = m.builder.Insert("user_uploads").
		Columns("user_id", "upload_id").
		Values(user.ID, upload.ID)

	_, err = m.insertQuery(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (m *MySQL) loadUpload(ctx context.Context, query sq.SelectBuilder) (*api.Upload, error) {
	row, err := m.selectRowQuery(ctx, query)
	if err != nil {
		return nil, err
	}

	var u = &api.Upload{}

	err = row.StructScan(u)
	if err == sql.ErrNoRows {
		return nil, store.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return u, nil
}

// LoadUpload ...
func (m *MySQL) LoadUpload(ctx context.Context, user *api.User, id int64) (*api.Upload, error) {
	if id == 0 {
		return nil, store.ErrZeroID
	}

	query := m.builder.Select("u.*").
		From("uploads u").
		LeftJoin("user_uploads uu on uu.upload_id = u.id").
		Where(sq.Eq{
			"uu.user_id": user.ID,
			"u.id":       id,
		})

	return m.loadUpload(ctx, query)
}

// LoadUploadByUUID ...
func (m *MySQL) LoadUploadByUUID(ctx context.Context, user *api.User, uuid string) (*api.Upload, error) {
	if uuid == "" {
		return nil, store.ErrEmptyQueryParam
	}

	query := m.builder.Select("u.*").
		From("uploads u").
		LeftJoin("user_uploads uu on uu.upload_id = u.id").
		Where(sq.Eq{
			"uu.user_id": user.ID,
			"u.uuid":     uuid,
		})

	return m.loadUpload(ctx, query)
}

// LoadUploads ...
func (m *MySQL) LoadUploads(ctx context.Context, user *api.User, opts *api.PagingOptions) (*api.Uploads, error) {
	err := checkUser(user, checkNilStruct|checkZeroID)
	if err != nil {
		return nil, err
	}

	if opts == nil {
		opts = api.NewPagingOptions(0, 0)
	}

	var uploads = &api.Uploads{
		Opts: opts,
		Data: []*api.Upload{},
	}

	query := m.builder.Select("u.*").
		From("uploads u").
		LeftJoin("user_uploads uu on uu.upload_id = u.id").
		Where(sq.Eq{"uu.user_id": user.ID}).
		OrderBy("u.created_at desc", "u.id desc").
		Limit(opts.Limit + 1)

	if opts.Cursor > 0 {
		query = query.Where(sq.LtOrEq{"u.id": opts.Cursor})
	}

	rows, err := m.selectQuery(ctx, query)
	if err == store.ErrNotFound {
		return uploads, nil
	}
	if err != nil {
		return nil, err
	}

	var cnt uint64
	for rows.Next() {
		var upload = &api.Upload{}

		err = rows.StructScan(upload)
		if err == sql.ErrNoRows {
			return nil, store.ErrNotFound
		}
		if err != nil {
			return nil, err
		}

		if cnt++; cnt > opts.Limit {
			opts.Next = upload.ID
		} else {
			uploads.Data = append(uploads.Data, upload)
		}
	}

	return uploads, nil
}
