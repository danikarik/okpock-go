package memory_test

import (
	"context"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	"github.com/danikarik/okpock/pkg/filestore"
	"github.com/danikarik/okpock/pkg/filestore/memory"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var mockRequiredVars = []string{
	"TEST_FILE",
	"TEST_FILE_IN_FOLDER",
	"TEST_TEMPLATES_BUCKET",
	"TEST_PASSES_BUCKET",
	"TEST_PROJECT",
}

func TestMockUploadFile(t *testing.T) {
	env, err := env.NewLookup(mockRequiredVars...)
	if err != nil {
		t.Skip(err)
	}

	testCases := []struct {
		Name   string
		Object *filestore.Object
	}{
		{
			Name: "TestFile",
			Object: &filestore.Object{
				Key:         uuid.NewV4().String() + ".txt",
				Body:        []byte("Hello World\n"),
				ContentType: "text/plain",
			},
		},
		{
			Name: "TestFileWithPrefix",
			Object: &filestore.Object{
				Prefix:      uuid.NewV1().String(),
				Key:         uuid.NewV4().String() + ".txt",
				Body:        []byte("Hello World\n"),
				ContentType: "text/plain",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			var (
				ctx    = context.Background()
				store  = memory.New()
				assert = assert.New(t)
				bucket = env.Get("TEST_PASSES_BUCKET")
			)

			err = store.UploadFile(ctx, bucket, tc.Object)
			if !assert.NoError(err) {
				return
			}

			obj, err := store.GetFile(ctx, bucket, tc.Object.Key)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(tc.Object.Prefix, obj.Prefix)
			assert.Equal(tc.Object.Key, obj.Key)
			assert.Equal(tc.Object.Body, obj.Body)
			assert.Equal(tc.Object.ContentType, obj.ContentType)

			objs, err := store.GetBucketFiles(ctx, bucket, tc.Object.Prefix)
			if !assert.NoError(err) {
				return
			}

			if assert.Len(objs, 1) {
				assert.Equal(tc.Object.Prefix, objs[0].Prefix)
				assert.Equal(tc.Object.Key, objs[0].Key)
				assert.Equal(tc.Object.Body, objs[0].Body)
				assert.Equal(tc.Object.ContentType, objs[0].ContentType)
			}
		})
	}
}
