package storage

import (
	"context"
	"errors"
	"io"
)

var (
	ErrInvalidObjectKey = errors.New("invalid object key")
	ErrObjectNotFound   = errors.New("object not found")
)

type Object struct {
	Key string
	URL string
}

type Store interface {
	Put(ctx context.Context, key string, source io.Reader) (Object, error)
	Delete(ctx context.Context, key string) error
	URL(key string) (string, error)
}
