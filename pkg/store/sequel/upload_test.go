package sequel_test

import (
	"context"
	"testing"

	"github.com/danikarik/okpock/pkg/api"
	"github.com/danikarik/okpock/pkg/store/sequel"
	"github.com/stretchr/testify/assert"
)

func TestIsUploadExists(t *testing.T) {
	testCases := []struct {
		Name       string
		SaveBefore bool
		Upload     *api.Upload
		Expected   bool
	}{
		{
			Name:       "NotTaken",
			SaveBefore: false,
			Upload: api.NewUpload(
				fakeString(),
				fakeString(),
				fakeString(),
			),
			Expected: false,
		},
		{
			Name:       "Taken",
			SaveBefore: true,
			Upload: api.NewUpload(
				fakeString(),
				fakeString(),
				fakeString(),
			),
			Expected: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			conn, err := testConnection(ctx, t)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			user := api.NewUser(fakeUsername(), fakeString(), fakeString(), nil)
			err = db.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			if tc.SaveBefore {
				err = db.SaveNewUpload(ctx, user, tc.Upload)
				if !assert.NoError(err) {
					return
				}
			}

			exists, err := db.IsUploadExists(ctx, user, tc.Upload.Filename, tc.Upload.Hash)
			assert.NoError(err)
			assert.Equal(tc.Expected, exists)
		})
	}
}

func TestSaveNewUpload(t *testing.T) {
	testCases := []struct {
		Name         string
		Upload       *api.Upload
		UploadNumber int
		LoadByUUID   bool
	}{
		{
			Name:   "SingleUpload",
			Upload: api.NewUpload(fakeString(), fakeString(), fakeString()),
		},
		{
			Name:         "ManyUpload",
			Upload:       api.NewUpload(fakeString(), fakeString(), fakeString()),
			UploadNumber: 9,
			LoadByUUID:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			conn, err := testConnection(ctx, t)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			user := api.NewUser(fakeUsername(), fakeString(), fakeString(), nil)
			err = db.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			for i := 0; i < tc.UploadNumber; i++ {
				upload := api.NewUpload(fakeString(), fakeString(), fakeString())
				err := db.SaveNewUpload(ctx, user, upload)
				if !assert.NoError(err) {
					return
				}
			}

			upload := api.NewUpload(tc.Upload.UUID, tc.Upload.Filename, tc.Upload.Hash)
			err = db.SaveNewUpload(ctx, user, upload)
			if !assert.NoError(err) {
				return
			}

			uploads, err := db.LoadUploads(ctx, user, nil)
			if !assert.NoError(err) {
				return
			}
			if !assert.Len(uploads.Data, tc.UploadNumber+1) {
				return
			}

			var loaded *api.Upload
			{
				if tc.LoadByUUID {
					loaded, err = db.LoadUploadByUUID(ctx, user, upload.UUID)
				} else {
					loaded, err = db.LoadUpload(ctx, user, upload.ID)
				}
				if !assert.NoError(err) {
					return
				}
			}

			assert.Equal(tc.Upload.UUID, loaded.UUID)
			assert.Equal(tc.Upload.Filename, loaded.Filename)
			assert.Equal(tc.Upload.Hash, loaded.Hash)
		})
	}
}

func TestUploadsPagination(t *testing.T) {
	testCases := []struct {
		Name     string
		Iters    int
		Uploads  int
		Limit    uint64
		HasNext  bool
		Expected int
	}{
		{
			Name:     "QueryFirstPage",
			Iters:    1,
			Uploads:  20,
			Limit:    10,
			HasNext:  true,
			Expected: 10,
		},
		{
			Name:     "QueryFirstPageFull",
			Iters:    1,
			Uploads:  7,
			Limit:    10,
			HasNext:  false,
			Expected: 7,
		},
		{
			Name:     "QuerySecondPage",
			Iters:    2,
			Uploads:  20,
			Limit:    7,
			HasNext:  true,
			Expected: 7,
		},
		{
			Name:     "QuerySecondPageFull",
			Iters:    2,
			Uploads:  20,
			Limit:    10,
			HasNext:  false,
			Expected: 10,
		},
		{
			Name:     "QueryAll",
			Iters:    3,
			Uploads:  9,
			Limit:    3,
			HasNext:  false,
			Expected: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)

			conn, err := testConnection(ctx, t)
			if !assert.NoError(err) {
				return
			}
			defer conn.Close()

			db := sequel.New(conn)

			user := api.NewUser(fakeUsername(), fakeEmail(), fakeString(), nil)
			err = db.SaveNewUser(ctx, user)
			if !assert.NoError(err) {
				return
			}

			for i := 0; i < tc.Uploads; i++ {
				upload := api.NewUpload(fakeString(), fakeString(), fakeString())
				err = db.SaveNewUpload(ctx, user, upload)
				if !assert.NoError(err) {
					return
				}
			}

			var (
				opts        = api.NewPagingOptions(0, tc.Limit)
				output      *api.Uploads
				hasNext     bool
				lastFetched int
			)

			for i := 0; i < tc.Iters; i++ {
				output, err = db.LoadUploads(ctx, user, opts)
				if !assert.NoError(err) {
					return
				}

				hasNext = output.Opts.HasNext()
				if hasNext {
					opts.Cursor = opts.Next
					opts.Next = 0
				}

				lastFetched = len(output.Data)
			}

			assert.Equal(tc.Expected, lastFetched)
			assert.Equal(tc.HasNext, hasNext)
		})
	}
}
