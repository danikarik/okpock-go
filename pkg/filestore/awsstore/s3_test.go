package awsstore_test

import (
	"context"
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

func TestSingleFile(t *testing.T) {
	skipTest(t)

	env, err := env.NewLookup(requiredVars...)
	if err != nil {
		t.Skip(err)
	}

	ctx := context.Background()
	assert := assert.New(t)

	store, err := awsstore.New()
	if !assert.NoError(err) {
		assert.FailNow("could not init handler")
	}

	obj, err := store.GetFile(ctx, env.Get("TEST_PASSES_BUCKET"), env.Get("TEST_FILE"))
	assert.NoError(err)
	assert.True(len(obj.Body) > 0)
}

func TestFolderFile(t *testing.T) {
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

	obj, err := store.GetFile(ctx, env.Get("TEST_TEMPLATES_BUCKET"), env.Get("TEST_FILE_IN_FOLDER"))
	assert.NoError(err)
	assert.True(len(obj.Body) > 0)
}

func TestBucket(t *testing.T) {
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

	contents, err := store.GetBucketFiles(ctx, env.Get("TEST_TEMPLATES_BUCKET"), env.Get("TEST_PROJECT"))
	assert.NoError(err)
	assert.Len(contents, 1)
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
