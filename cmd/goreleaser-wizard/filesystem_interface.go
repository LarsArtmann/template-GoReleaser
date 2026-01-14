package main

import (
	"context"
	"io"
	"os"
)

// FileSystemRepository interface for file system operations.
type FileSystemRepository interface {
	ReadFile(ctx context.Context, path string) ([]byte, error)
	WriteFile(ctx context.Context, path string, data []byte, perm os.FileMode) error
	CreateFile(ctx context.Context, path string) (io.WriteCloser, error)
	DeleteFile(ctx context.Context, path string) error
	FileExists(ctx context.Context, path string) (bool, error)
	CreateDir(ctx context.Context, path string, perm os.FileMode) error
	CreateDirAll(ctx context.Context, path string, perm os.FileMode) error
	DirExists(ctx context.Context, path string) (bool, error)
	ReadDir(ctx context.Context, path string) ([]os.DirEntry, error)
	GetFileInfo(ctx context.Context, path string) (os.FileInfo, error)
	CheckPermissions(ctx context.Context, path string) (bool, error)
	AbsPath(path string) (string, error)
	RelPath(base, target string) (string, error)
	CleanPath(path string) string
	JoinPath(elem ...string) string
	TempDir(dir, pattern string) (string, error)
}
