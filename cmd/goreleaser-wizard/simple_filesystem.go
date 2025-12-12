package main

import (
	"context"
	"io"
	"os"
	"path/filepath"
)

// SimpleFileSystemRepository is a basic implementation for demonstration
type SimpleFileSystemRepository struct{}

func (r *SimpleFileSystemRepository) ReadFile(ctx context.Context, path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (r *SimpleFileSystemRepository) WriteFile(ctx context.Context, path string, data []byte, perm os.FileMode) error {
	return os.WriteFile(path, data, perm)
}

func (r *SimpleFileSystemRepository) CreateFile(ctx context.Context, path string) (io.WriteCloser, error) {
	return os.Create(path)
}

func (r *SimpleFileSystemRepository) DeleteFile(ctx context.Context, path string) error {
	return os.Remove(path)
}

func (r *SimpleFileSystemRepository) FileExists(ctx context.Context, path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (r *SimpleFileSystemRepository) CreateDir(ctx context.Context, path string, perm os.FileMode) error {
	return os.Mkdir(path, perm)
}

func (r *SimpleFileSystemRepository) CreateDirAll(ctx context.Context, path string, perm os.FileMode) error {
	return os.MkdirAll(path, perm)
}

func (r *SimpleFileSystemRepository) DirExists(ctx context.Context, path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func (r *SimpleFileSystemRepository) ReadDir(ctx context.Context, path string) ([]os.DirEntry, error) {
	return os.ReadDir(path)
}

func (r *SimpleFileSystemRepository) GetFileInfo(ctx context.Context, path string) (os.FileInfo, error) {
	return os.Stat(path)
}

func (r *SimpleFileSystemRepository) CheckPermissions(ctx context.Context, path string) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, err
	}
	return info.Mode().Perm()&0o777 != 0, nil
}

func (r *SimpleFileSystemRepository) AbsPath(path string) (string, error) {
	return filepath.Abs(path)
}

func (r *SimpleFileSystemRepository) RelPath(base, target string) (string, error) {
	return filepath.Rel(base, target)
}

func (r *SimpleFileSystemRepository) CleanPath(path string) string {
	return filepath.Clean(path)
}

func (r *SimpleFileSystemRepository) JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

func (r *SimpleFileSystemRepository) TempDir(dir, pattern string) (string, error) {
	return os.MkdirTemp(dir, pattern)
}
