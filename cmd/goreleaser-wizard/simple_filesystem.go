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
		return nil, wrapFSError("read file", path, err)
	}

	return data, nil
}

func (r *SimpleFileSystemRepository) WriteFile(
	ctx context.Context,
	path string,
	data []byte,
	perm os.FileMode,
) error {
	err := os.WriteFile(path, data, perm)
	if err != nil {
		return wrapFSError("write file", path, err)
	}

	return nil
}

func (r *SimpleFileSystemRepository) CreateFile(
	ctx context.Context,
	path string,
) (io.WriteCloser, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, wrapFSError("create file", path, err)
	}

	return f, nil
}

func (r *SimpleFileSystemRepository) DeleteFile(ctx context.Context, path string) error {
	err := os.Remove(path)
	if err != nil {
		return wrapFSError("delete file", path, err)
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

	return false, wrapFSError("check file existence", path, err)
}

func (r *SimpleFileSystemRepository) CreateDir(
	ctx context.Context,
	path string,
	perm os.FileMode,
) error {
	err := os.Mkdir(path, perm)
	if err != nil {
		return wrapFSError("create directory", path, err)
	}

	return nil
}

func (r *SimpleFileSystemRepository) CreateDirAll(
	ctx context.Context,
	path string,
	perm os.FileMode,
) error {
	err := os.MkdirAll(path, perm)
	if err != nil {
		return wrapFSError("create directories", path, err)
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

	return false, wrapFSError("check directory existence", path, err)
}

func (r *SimpleFileSystemRepository) ReadDir(
	ctx context.Context,
	path string,
) ([]os.DirEntry, error) {
	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, wrapFSError("read directory", path, err)
	}

	return entries, nil
}

func (r *SimpleFileSystemRepository) GetFileInfo(
	ctx context.Context,
	path string,
) (os.FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, wrapFSError("get file info", path, err)
	}

	return info, nil
}

func (r *SimpleFileSystemRepository) CheckPermissions(
	ctx context.Context,
	path string,
) (bool, error) {
	info, err := os.Stat(path)
	if err != nil {
		return false, wrapFSError("check permissions", path, err)
	}

	return info.Mode().Perm()&0o777 != 0, nil
}

func (r *SimpleFileSystemRepository) AbsPath(path string) (string, error) {
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", wrapFSError("get absolute path", path, err)
	}

	return abs, nil
}

func (r *SimpleFileSystemRepository) RelPath(base, target string) (string, error) {
	rel, err := filepath.Rel(base, target)
	if err != nil {
		return "", wrapFSError(
			fmt.Sprintf("get relative path from %q to %q", base, target),
			"",
			err,
		)
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
		return "", fmt.Errorf(
			"failed to create temp directory (dir=%q, pattern=%q): %w",
			dir,
			pattern,
			err,
		)
	}

	return tempDir, nil
}

// wrapFSError wraps a file system operation error with the operation name and path.
func wrapFSError(op, path string, err error) error {
	if path == "" {
		return fmt.Errorf("failed to %s: %w", op, err)
	}

	return fmt.Errorf("failed to %s %q: %w", op, path, err)
}