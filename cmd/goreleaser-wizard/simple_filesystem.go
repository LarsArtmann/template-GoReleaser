package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// SimpleFileSystemRepository is a basic implementation for demonstration.
type SimpleFileSystemRepository struct{}

func (r *SimpleFileSystemRepository) ReadFile(ctx context.Context, path string) ([]byte, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", path, err)
	}

	return data, nil
}

func (r *SimpleFileSystemRepository) WriteFile(
	ctx context.Context,
	path string,
	data []byte,
	perm os.FileMode,
) error {
	if err := os.WriteFile(path, data, perm); err != nil {
		return fmt.Errorf("failed to write file %q: %w", path, err)
	}

	return nil
}

func (r *SimpleFileSystemRepository) CreateFile(
	ctx context.Context,
	path string,
) (io.WriteCloser, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("failed to create file %q: %w", path, err)
	}

	return f, nil
}

func (r *SimpleFileSystemRepository) DeleteFile(ctx context.Context, path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete file %q: %w", path, err)
	}

	return nil
}

func (r *SimpleFileSystemRepository) FileExists(ctx context.Context, path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("failed to check file existence %q: %w", path, err)
}

func (r *SimpleFileSystemRepository) CreateDir(
	ctx context.Context,
	path string,
	perm os.FileMode,
) error {
	if err := os.Mkdir(path, perm); err != nil {
		return fmt.Errorf("failed to create directory %q: %w", path, err)
	}

	return nil
}

func (r *SimpleFileSystemRepository) CreateDirAll(
	ctx context.Context,
	path string,
	perm os.FileMode,
) error {
	if err := os.MkdirAll(path, perm); err != nil {
		return fmt.Errorf("failed to create directories %q: %w", path, err)
	}

	return nil
}

func (r *SimpleFileSystemRepository) DirExists(ctx context.Context, path string) (bool, error) {
	info, err := os.Stat(path)
	if err == nil {
		return info.IsDir(), nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, fmt.Errorf("failed to check directory existence %q: %w", path, err)
}

func (r *SimpleFileSystemRepository) ReadDir(
	ctx context.Context,
	path string,
) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read directory %q: %w", path, err)
	}

	return entries, nil
}

func (r *SimpleFileSystemRepository) GetFileInfo(
	ctx context.Context,
	path string,
) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("failed to get file info %q: %w", path, err)
	}

	return info, nil
}

func (r *SimpleFileSystemRepository) CheckPermissions(
	ctx context.Context,
	path string,
) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, fmt.Errorf("failed to check permissions for %q: %w", path, err)
	}

	return info.Mode().Perm()&0o777 != 0, nil
}

func (r *SimpleFileSystemRepository) AbsPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to get absolute path for %q: %w", path, err)
	}

	return abs, nil
}

func (r *SimpleFileSystemRepository) RelPath(base, target string) (string, error) {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return "", fmt.Errorf("failed to get relative path from %q to %q: %w", base, target, err)
	}

	return rel, nil
}

func (r *SimpleFileSystemRepository) CleanPath(path string) string {
	return filepath.Clean(path)
}

func (r *SimpleFileSystemRepository) JoinPath(elem ...string) string {
	return filepath.Join(elem...)
}

func (r *SimpleFileSystemRepository) TempDir(dir, pattern string) (string, error) {
	tempDir, err := os.MkdirTemp(dir, pattern)
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}

	return tempDir, nil
}
