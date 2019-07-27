package filestore

import (
	"context"
	"fmt"
	"net/http"
)

const (
	// ApplePkpass refers to `pkpass` extentsion's content-type.
	ApplePkpass = "application/vnd.apple.pkpass"
)

// Object represents file object in the bucket.
// It holds file's key and content.
type Object struct {
	Prefix      string
	Key         string
	Body        []byte
	ContentType string
}

// Path returns object key whether with prefix or not.
func (o *Object) Path() string {
	if o.Prefix != "" {
		return fmt.Sprintf("%s/%s", o.Prefix, o.Key)
	}
	return o.Key
}

// Serve serves object's content over http.
func (o *Object) Serve(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", o.ContentType)
	_, err := w.Write(o.Body)
	return err
}

// Storage holds method for accessing remote buckets.
type Storage interface {
	File(ctx context.Context, bucket, key string) (*Object, error)
	Bucket(ctx context.Context, bucket, prefix string) ([]*Object, error)
	Upload(ctx context.Context, bucket string, obj *Object) error
}
