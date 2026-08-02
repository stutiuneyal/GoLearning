package storage

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/url"
	"path"
	"path/filepath"
	"strings"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var _ Store = (*S3Store)(nil)

type S3Store struct {
	client        *s3.Client
	bucket        string
	publicURLBase string
}

func NewS3Store(ctx context.Context, region string, bucket string) (*S3Store, error) {

	region = strings.TrimSpace(region)
	bucket = strings.TrimSpace(bucket)

	if region == "" {
		return nil, fmt.Errorf("AWS region is required")
	}

	if bucket == "" {
		return nil, fmt.Errorf("S3 bucket is required")
	}

	awsConfig, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(region))
	if err != nil {
		return nil, fmt.Errorf("load AWS configuration: %w", err)
	}

	return &S3Store{
		client: s3.NewFromConfig(awsConfig),
		bucket: bucket,
		publicURLBase: fmt.Sprintf(
			"https://%s.s3.%s.amazonaws.com",
			bucket,
			region,
		),
	}, nil

}

func (s *S3Store) Put(ctx context.Context, key string, source io.Reader) (Object, error) {
	if source == nil {
		return Object{}, fmt.Errorf("source reader is required")
	}

	cleanKey, err := validateS3ObjectKey(key)
	if err != nil {
		return Object{}, err
	}

	contentType := mime.TypeByExtension(filepath.Ext(cleanKey))
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	_, err = s.client.PutObject(
		ctx,
		&s3.PutObjectInput{
			Bucket:      &s.bucket,
			Key:         &cleanKey,
			Body:        source,
			ContentType: &contentType,
		},
	)
	if err != nil {
		return Object{}, fmt.Errorf(
			"upload object to s3: %w",
			err,
		)
	}

	objectURl, err := s.URL(cleanKey)
	if err != nil {
		return Object{}, err
	}

	return Object{
		Key: cleanKey,
		URL: objectURl,
	}, nil
}

func (s *S3Store) Delete(ctx context.Context, key string) error {

	cleanKey, err := validateS3ObjectKey(key)
	if err != nil {
		return err
	}

	_, err = s.client.DeleteObject(
		ctx,
		&s3.DeleteObjectInput{
			Bucket: &s.bucket,
			Key:    &cleanKey,
		},
	)
	if err != nil {
		return fmt.Errorf(
			"delete object from S3: %w",
			err,
		)
	}

	return nil
}

func (s *S3Store) URL(key string) (string, error) {

	cleanKey, err := validateS3ObjectKey(key)
	if err != nil {
		return "", err
	}

	return s.publicURLBase + "/" + encodeS3ObjectKey(cleanKey), nil
}

func validateS3ObjectKey(key string) (string, error) {
	key = strings.TrimSpace(key)
	key = strings.ReplaceAll(key, "\\", "/")
	key = strings.TrimPrefix(key, "/")

	if key == "" {
		return "", ErrInvalidObjectKey
	}

	cleanKey := path.Clean(key)

	if cleanKey == "." || cleanKey == ".." || strings.HasPrefix(cleanKey, "../") || strings.HasPrefix(cleanKey, "/") {
		return "", ErrInvalidObjectKey
	}

	segments := strings.Split(cleanKey, "/")

	for _, segment := range segments {
		if segment == "" || segment == "." || segment == ".." {
			return "", ErrInvalidObjectKey
		}
	}

	if len(segments) < 2 || segments[0] != "users" {
		return "", ErrInvalidObjectKey
	}

	return cleanKey, nil
}

func encodeS3ObjectKey(key string) string {
	segments := strings.Split(key, "/")

	for i, segment := range segments {
		segments[i] = url.PathEscape(segment)
	}

	return strings.Join(segments, "/")
}
