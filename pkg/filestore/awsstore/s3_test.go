package awsstore_test

import (
	"context"
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

func TestSingleFile(t *testing.T) {
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

	obj, err := store.File(ctx, env.Get("TEST_PASSES_BUCKET"), env.Get("TEST_FILE"))
	if !assert.NoError(err) {
		assert.FailNow("could not read file")
	}

	if !assert.True(len(obj.Body) > 0) {
		assert.FailNow("wrong content length")
	}

	t.Logf("key: %s, content-type: %s, content-length: %v\n", obj.Key, obj.ContentType, len(obj.Body))
}

func TestFolderFile(t *testing.T) {
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

	obj, err := store.File(ctx, env.Get("TEST_TEMPLATES_BUCKET"), env.Get("TEST_FILE_IN_FOLDER"))
	if !assert.NoError(err) {
		assert.FailNow("could not read file")
	}

	if !assert.True(len(obj.Body) > 0) {
		assert.FailNow("wrong content length")
	}

	t.Logf("key: %s, content-type: %s, content-length: %v\n", obj.Key, obj.ContentType, len(obj.Body))
}

func TestBucket(t *testing.T) {
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

	contents, err := store.Bucket(ctx, env.Get("TEST_TEMPLATES_BUCKET"), env.Get("TEST_PROJECT"))
	if !assert.NoError(err) {
		assert.FailNow("could not read file")
	}

	if !assert.Len(contents, 1) {
		assert.FailNow("bucket cannot be empty")
	}

	for _, c := range contents {
		t.Logf("key: %s, content-type: %s, content-length: %v\n", c.Key, c.ContentType, len(c.Body))
	}
}

func TestUploadSingleFile(t *testing.T) {
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

	obj := &filestore.Object{
		Key:         uuid.NewV4().String() + ".txt",
		Body:        []byte("Hello World\n"),
		ContentType: "text/plain",
	}

	err = store.Upload(ctx, env.Get("TEST_PASSES_BUCKET"), obj)
	if !assert.NoError(err) {
		assert.FailNow("could not upload file")
	}

	t.Logf(
		"bucket: %s, key: %s, content-type: %s, content-length: %v\n",
		env.Get("TEST_PASSES_BUCKET"),
		obj.Key,
		obj.ContentType,
		len(obj.Body),
	)
}

func TestUploadFolderFile(t *testing.T) {
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

	obj := &filestore.Object{
		Prefix:      uuid.NewV1().String(),
		Key:         uuid.NewV4().String() + ".txt",
		Body:        []byte("Hello World\n"),
		ContentType: "text/plain",
	}
	err = store.Upload(ctx, env.Get("TEST_PASSES_BUCKET"), obj)
	if !assert.NoError(err) {
		assert.FailNow("could not upload file")
	}

	t.Logf(
		"bucket: %s, prefix: %s, key: %s, content-type: %s, content-length: %v\n",
		env.Get("TEST_PASSES_BUCKET"),
		obj.Prefix,
		obj.Key,
		obj.ContentType,
		len(obj.Body),
	)
}
