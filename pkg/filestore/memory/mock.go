package memory

import (
	"context"
	"errors"
	"strings"
	"sync"

	"github.com/danikarik/okpock/pkg/filestore"
)

type mockFile struct {
	Bucket        string
	Prefix        string
	Key           string
	ContentType   string
	ContentLength int
	Body          []byte
}

// New returns mock storage handler.
func New() filestore.Storage {
	return &mockHandler{files: map[string]mockFile{}}
}

type mockHandler struct {
	mu    sync.Mutex
	files map[string]mockFile
}

func mockIndex(bucket, key string) string {
	return bucket + "^" + key
}

func (m *mockHandler) GetFile(ctx context.Context, bucket, key string) (*filestore.Object, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	index := mockIndex(bucket, key)
	f, ok := m.files[index]
	if !ok {
		return nil, errors.New("file not found")
	}

	obj := &filestore.Object{
		Prefix:      f.Prefix,
		Key:         key,
		Body:        f.Body,
		ContentType: f.ContentType,
	}

	return obj, nil
}

func (m *mockHandler) GetBucketFiles(ctx context.Context, bucket, prefix string) ([]*filestore.Object, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	objs := []*filestore.Object{}
	for index, file := range m.files {
		parts := strings.Split(index, "^")
		buc, key := parts[0], parts[1]
		if buc == bucket && file.Prefix == prefix {
			obj := &filestore.Object{
				Prefix:      prefix,
				Key:         key,
				Body:        file.Body,
				ContentType: file.ContentType,
			}
			objs = append(objs, obj)
		}
	}

	return objs, nil
}

func (m *mockHandler) UploadFile(ctx context.Context, bucket string, obj *filestore.Object) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	f := mockFile{
		Bucket:        bucket,
		Prefix:        obj.Prefix,
		Key:           obj.Key,
		ContentType:   obj.ContentType,
		ContentLength: len(obj.Body),
		Body:          obj.Body,
	}

	index := mockIndex(f.Bucket, obj.Path())
	m.files[index] = f

	return nil
}
