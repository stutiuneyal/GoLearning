package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var _ Store = (*LocalStore)(nil)

type LocalStore struct {
	rootDir       string
	publicURLbase string
}

func NewLocalStore(rootDir, publicURLBase string) (*LocalStore, error) {
	if strings.TrimSpace(rootDir) == "" {
		return nil, fmt.Errorf("root directory is required")
	}

	absoluteRoot, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, fmt.Errorf(
			"resolve storage root directory: %w",
			err,
		)
	}

	if err := os.MkdirAll(absoluteRoot, 0o755); err != nil {
		return nil, fmt.Errorf(
			"create storage root directory: %w",
			err,
		)
	}

	publicURLBase = strings.TrimSpace(publicURLBase)
	if publicURLBase == "" {
		return nil, fmt.Errorf("public URL base is required")
	}

	publicURLBase = "/" + strings.Trim(publicURLBase, "/")

	return &LocalStore{
		rootDir:       absoluteRoot,
		publicURLbase: publicURLBase,
	}, nil
}

func (s *LocalStore) RootDir() string {
	return s.rootDir
}

func (s *LocalStore) Put(ctx context.Context, key string, source io.Reader) (Object, error) {
	if source == nil {
		return Object{}, fmt.Errorf("source reader is required")
	}

	cleanKey, destination, err := s.resolveKey(key)
	if err != nil {
		return Object{}, err
	}

	if err := ctx.Err(); err != nil {
		return Object{}, err
	}

	if err := os.MkdirAll(filepath.Dir(destination), 0o755); err != nil {
		return Object{}, fmt.Errorf(
			"create object directory: %w",
			err,
		)
	}

	tempFile, err := os.CreateTemp(filepath.Dir(destination), ".upload-*")
	if err != nil {
		return Object{}, fmt.Errorf(
			"create temporary upload file: %w",
			err,
		)
	}

	temporaryPath := tempFile.Name()
	committed := false

	defer func() {
		tempFile.Close()

		if !committed {
			os.Remove(temporaryPath)
		}
	}()

	if err := copyWithContext(ctx, tempFile, source); err != nil {
		return Object{}, fmt.Errorf(
			"flush uploaded object: %w",
			err,
		)
	}

	if err := tempFile.Close(); err != nil {
		return Object{}, fmt.Errorf(
			"close uploaded object: %w",
			err,
		)
	}

	if err := os.Chmod(temporaryPath, 0o644); err != nil {
		return Object{}, fmt.Errorf(
			"set uploaded object permissions: %w",
			err,
		)
	}

	if err := os.Rename(temporaryPath, destination); err != nil {
		return Object{}, fmt.Errorf(
			"commit uploaded object: %w",
			err,
		)
	}

	committed = true

	objectURL, err := s.URL(cleanKey)
	if err != nil {
		return Object{}, err
	}

	return Object{
		Key: cleanKey,
		URL: objectURL,
	}, nil
}

func (s *LocalStore) Delete(ctx context.Context, key string) error {

	_, destination, err := s.resolveKey(key)

	if err != nil {
		return err
	}

	if err := ctx.Err(); err != nil {
		return err
	}

	if err := os.Remove(destination); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil
		}

		return fmt.Errorf(
			"delete object: %w",
			err,
		)
	}

	return nil
}

func (s *LocalStore) URL(key string) (string, error) {

	cleanKey, _, err := s.resolveKey(key)
	if err != nil {
		return "", err
	}

	urlKey := strings.ReplaceAll(cleanKey, string(filepath.Separator), "/")

	return s.publicURLbase + "/" + urlKey, nil
}

func (s *LocalStore) resolveKey(key string) (string, string, error) {
	key = strings.TrimSpace(key)
	if key == "" {
		return "", "", ErrInvalidObjectKey
	}

	// storage keys always use forward slashes
	key = strings.ReplaceAll(key, "\\", "/")

	cleanKey := filepath.ToSlash(filepath.Clean(filepath.FromSlash(key)))

	if cleanKey == "." ||
		cleanKey == ".." ||
		strings.HasPrefix(cleanKey, "../") {
		return "", "", ErrInvalidObjectKey
	}

	if filepath.IsAbs(filepath.FromSlash(cleanKey)) || strings.HasPrefix(cleanKey, "/") {
		return "", "", ErrInvalidObjectKey
	}

	destination := filepath.Join(s.rootDir, filepath.FromSlash(cleanKey))

	relative, err := filepath.Rel(s.rootDir, destination)
	if err != nil {
		return "", "", fmt.Errorf(
			"resolve object key: %w",
			err,
		)
	}

	if relative == ".." ||
		strings.HasPrefix(
			relative,
			".."+string(filepath.Separator),
		) {
		return "", "", ErrInvalidObjectKey
	}

	return cleanKey, destination, nil

}

func copyWithContext(ctx context.Context, destination io.Writer, source io.Reader) error {

	buffer := make([]byte, 32*1024)

	for {
		if err := ctx.Err(); err != nil {
			return err
		}

		readCount, readErr := source.Read(buffer)

		if readCount > 0 {
			_, writeErr := destination.Write(
				buffer[:readCount],
			)
			if writeErr != nil {
				return writeErr
			}
		}

		if errors.Is(readErr, io.EOF) {
			return nil
		}

		if readErr != nil {
			return readErr
		}

	}

}
