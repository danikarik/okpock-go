package awsstore

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/danikarik/okpock/pkg/filestore"
)

// New returns AWS S3 Storage Handler.
func New() (filestore.Storage, error) {
	cfg := &aws.Config{Credentials: credentials.NewEnvCredentials()}
	sess, err := session.NewSession(cfg)
	if err != nil {
		return nil, fmt.Errorf("could not create session: %v", err)
	}
	s3h := &s3handler{
		srv: s3.New(sess),
	}
	return s3h, nil
}

type s3handler struct {
	srv *s3.S3
}

func (s3h *s3handler) File(ctx context.Context, bucket, key string) (*filestore.Object, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	output, err := s3h.srv.GetObject(input)
	if err != nil {
		return nil, fmt.Errorf("could not get object %s: %v", key, err)
	}
	buf, err := ioutil.ReadAll(output.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read body %s: %v", key, err)
	}
	return &filestore.Object{
		Key:         key,
		Body:        buf,
		ContentType: *output.ContentType,
	}, nil
}

func (s3h *s3handler) Bucket(ctx context.Context, bucket, prefix string) ([]*filestore.Object, error) {
	contents := make([]*filestore.Object, 0)
	input := &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(prefix),
	}
	output, err := s3h.srv.ListObjects(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchBucket:
				return nil, fmt.Errorf("%s %s: %v", s3.ErrCodeNoSuchBucket, bucket, aerr.Error())
			default:
				return nil, fmt.Errorf("%s: %v", bucket, aerr.Error())
			}
		}
		return nil, fmt.Errorf("could not list objects %s: %v", bucket, err)
	}
	for _, obj := range output.Contents {
		content, err := s3h.File(ctx, bucket, *obj.Key)
		if err != nil {
			return nil, fmt.Errorf("could not read content %s: %v", bucket, err)
		}
		if content.ContentType != "application/x-directory" {
			contents = append(contents, content)
		}
	}
	return contents, nil
}

func (s3h *s3handler) Upload(ctx context.Context, bucket string, obj *filestore.Object) error {
	if obj == nil {
		return errors.New("object cannot be empty")
	}
	input := &s3.PutObjectInput{
		Bucket:        aws.String(bucket),
		Key:           aws.String(obj.Path()),
		Body:          bytes.NewReader(obj.Body),
		ContentLength: aws.Int64(int64(len(obj.Body))),
		ContentType:   aws.String(obj.ContentType),
	}
	_, err := s3h.srv.PutObject(input)
	if err != nil {
		return fmt.Errorf("could not put object %s: %v", obj.Key, err)
	}
	return nil
}
