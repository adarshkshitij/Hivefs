package main

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/adarshkshitij/hivefs/metrics"
)

type PathTransformFunc func(string) string

type StoreOpts struct {
	// Root is the folder name of the root, containing all the files/folders of the system.
	Root              string
	PathTransformFunc PathTransformFunc
}

var DefaultPathTransformFunc = func(key string) string {
	return key
}

type Store struct {
	StoreOpts
}

func NewStore(opts StoreOpts) *Store {
	if opts.PathTransformFunc == nil {
		opts.PathTransformFunc = DefaultPathTransformFunc
	}
	if opts.Root == "" {
		opts.Root = "data"
	}
	return &Store{
		StoreOpts: opts,
	}
}

func (s *Store) writeStream(key string, r io.Reader) error {
	pathName := s.PathTransformFunc(key)
	fullPath := filepath.Join(s.Root, pathName)

	if err := os.MkdirAll(fullPath, os.ModePerm); err != nil {
		return err
	}

	// For now, we use the key as the filename within the transformed path
	filePath := filepath.Join(fullPath, key)

	f, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	n, err := io.Copy(f, r)
	if err != nil {
		return err
	}

	slog.Info("file written to disk", "path", filePath, "bytes", n)

	// Track bytes stored on disk
	metrics.BytesStoredTotal.Add(float64(n))

	return nil
}