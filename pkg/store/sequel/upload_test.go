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
			UploadNumber: 10,
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

			uploads, err := db.LoadUploads(ctx, user)
			if !assert.NoError(err) {
				return
			}
			if !assert.Len(uploads, tc.UploadNumber+1) {
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
