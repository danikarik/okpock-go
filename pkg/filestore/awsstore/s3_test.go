package awsstore_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"testing"

	"github.com/danikarik/okpock/pkg/env"
	"github.com/danikarik/okpock/pkg/filestore"
	"github.com/danikarik/okpock/pkg/filestore/awsstore"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
)

var requiredVars = []string{
	"AWS_ACCESS_KEY_ID",
	"AWS_SECRET_ACCESS_KEY",
	"AWS_REGION",
	"TEST_FILE",
	"TEST_FILE_IN_FOLDER",
	"TEST_TEMPLATES_BUCKET",
	"TEST_PASSES_BUCKET",
	"TEST_PROJECT",
}

func skipTest(t *testing.T) {
	if v, ok := os.LookupEnv("SKIP_S3_TEST"); ok {
		skip, err := strconv.ParseBool(v)
		if err == nil && skip {
			t.Skip(`skip test: SKIP_S3_TEST is present`)
		}
	}
}

func readFile(path string) (*filestore.Object, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fi, err := file.Stat()
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return &filestore.Object{
		Key:         fi.Name(),
		Body:        body,
		ContentType: http.DetectContentType(body),
	}, nil
}

func TestSingleFile(t *testing.T) {
	skipTest(t)

	env, err := env.NewLookup(requiredVars...)
	if err != nil {
		t.Skip(err)
	}

	testCases := []struct {
		Name string
		Path string
	}{
		{
			Name: "Pkpass",
			Path: "testdata/9973af9d-9cfa-4d9f-8c6c-32255de8d96b.pkpass",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)
			bucket := env.Get("TEST_PASSES_BUCKET")

			store, err := awsstore.New()
			if !assert.NoError(err) {
				return
			}

			obj, err := readFile(tc.Path)
			if !assert.NoError(err) {
				return
			}

			err = store.UploadFile(ctx, bucket, obj)
			if !assert.NoError(err) {
				return
			}

			loaded, err := store.GetFile(ctx, bucket, obj.Key)
			if !assert.NoError(err) {
				return
			}

			assert.Equal(obj.Body, loaded.Body)
		})
	}
}

func TestFolderFile(t *testing.T) {
	skipTest(t)

	env, err := env.NewLookup(requiredVars...)
	if err != nil {
		t.Skip(err)
	}

	testCases := []struct {
		Name   string
		Folder string
		Path   string
	}{
		{
			Name:   "Text",
			Folder: "demo",
			Path:   "testdata/hello.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)
			bucket := env.Get("TEST_PASSES_BUCKET")

			store, err := awsstore.New()
			if !assert.NoError(err) {
				return
			}

			obj, err := readFile(tc.Path)
			if !assert.NoError(err) {
				return
			}
			obj.Prefix = tc.Folder

			err = store.UploadFile(ctx, bucket, obj)
			if !assert.NoError(err) {
				return
			}

			loaded, err := store.GetFile(ctx, bucket, obj.Path())
			if !assert.NoError(err) {
				return
			}

			assert.Equal(obj.Body, loaded.Body)
		})
	}
}

func TestBucket(t *testing.T) {
	skipTest(t)

	env, err := env.NewLookup(requiredVars...)
	if err != nil {
		t.Skip(err)
	}

	testCases := []struct {
		Name   string
		Folder string
		Path   string
	}{
		{
			Name:   "Text",
			Folder: "demo",
			Path:   "testdata/hello.txt",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			ctx := context.Background()
			assert := assert.New(t)
			bucket := env.Get("TEST_TEMPLATES_BUCKET")

			store, err := awsstore.New()
			if !assert.NoError(err) {
				return
			}

			obj, err := readFile(tc.Path)
			if !assert.NoError(err) {
				return
			}
			obj.Prefix = tc.Folder

			err = store.UploadFile(ctx, bucket, obj)
			if !assert.NoError(err) {
				return
			}

			contents, err := store.GetBucketFiles(ctx, bucket, tc.Folder)
			if !assert.NoError(err) {
				return
			}

			assert.Len(contents, 1)
		})
	}
}

func TestUploadSingleFile(t *testing.T) {
	skipTest(t)

	env, err := env.NewLookup(requiredVars...)
	if err != nil {
		t.Skip(err)
	}

	ctx := context.Background()
	assert := assert.New(t)

	store, err := awsstore.New()
	if !assert.NoError(err) {
		return
	}

	obj := &filestore.Object{
		Key:         uuid.NewV4().String() + ".txt",
		Body:        []byte("Hello World\n"),
		ContentType: "text/plain",
	}

	err = store.UploadFile(ctx, env.Get("TEST_PASSES_BUCKET"), obj)
	assert.NoError(err)
}

func TestUploadExistingSingleFile(t *testing.T) {
	skipTest(t)

	env, err := env.NewLookup(requiredVars...)
	if err != nil {
		t.Skip(err)
	}

	ctx := context.Background()
	assert := assert.New(t)

	store, err := awsstore.New()
	if !assert.NoError(err) {
		return
	}

	obj := &filestore.Object{
		Key:         uuid.NewV4().String() + ".txt",
		Body:        []byte("Hello World\n"),
		ContentType: "text/plain",
	}

	err = store.UploadFile(ctx, env.Get("TEST_PASSES_BUCKET"), obj)
	if !assert.NoError(err) {
		return
	}

	obj.Body = []byte("Hello Okpock\n")

	err = store.UploadFile(ctx, env.Get("TEST_PASSES_BUCKET"), obj)
	if !assert.NoError(err) {
		return
	}

	loaded, err := store.GetFile(ctx, env.Get("TEST_PASSES_BUCKET"), obj.Key)
	assert.NoError(err)
	assert.Equal(obj.Body, loaded.Body)
}

func TestUploadFolderFile(t *testing.T) {
	skipTest(t)

	env, err := env.NewLookup(requiredVars...)
	if err != nil {
		t.Skip(err)
	}

	ctx := context.Background()
	assert := assert.New(t)

	store, err := awsstore.New()
	if !assert.NoError(err) {
		return
	}

	obj := &filestore.Object{
		Prefix:      uuid.NewV1().String(),
		Key:         uuid.NewV4().String() + ".txt",
		Body:        []byte("Hello World\n"),
		ContentType: "text/plain",
	}
	err = store.UploadFile(ctx, env.Get("TEST_PASSES_BUCKET"), obj)
	assert.NoError(err)
}
