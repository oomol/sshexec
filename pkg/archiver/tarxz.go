package archiver

import (
	"context"
	"fmt"

	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mholt/archives"
	"github.com/sirupsen/logrus"
)

const (
	dirPermissions  = 0o700 // Default directory permissions
	filePermissions = 0o600 // Default file permissions
)

func Unarchive(tarball, dst string) error {
	archiveFile, openErr := os.Open(tarball)
	if openErr != nil {
		return fmt.Errorf("open tarball %s: %w", tarball, openErr)
	}
	defer archiveFile.Close()

	format, input, err := archives.Identify(context.Background(), tarball, archiveFile)
	if err != nil {
		return fmt.Errorf("identify format: %w", err)
	}

	extractor, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("unsupported format for extraction")
	}

	if err := createDirWithPermissions(dst, dirPermissions); err != nil {
		return fmt.Errorf("creating destination directory: %w", err)
	}

	handler := func(ctx context.Context, f archives.FileInfo) error {
		return handleFile(f, dst)
	}

	if err := extractor.Extract(context.Background(), input, handler); err != nil {
		return fmt.Errorf("extracting files: %w", err)
	}

	return nil
}

func createDirWithPermissions(path string, mode os.FileMode) error {
	err := os.MkdirAll(path, mode)
	if err != nil {
		return fmt.Errorf("mkdir: %w", err)
	}

	return nil
}

// securePath ensures the path is safely relative to the target directory.
func securePath(basePath, relativePath string) (string, error) {
	relativePath = filepath.Clean("/" + relativePath)                         // Normalize path with a leading slash
	relativePath = strings.TrimPrefix(relativePath, string(os.PathSeparator)) // Remove leading separator

	dstPath := filepath.Join(basePath, relativePath)

	if !strings.HasPrefix(filepath.Clean(dstPath)+string(os.PathSeparator), filepath.Clean(basePath)+string(os.PathSeparator)) {
		return "", fmt.Errorf("illegal file path: %s", dstPath)
	}

	return dstPath, nil
}

func handleFile(f archives.FileInfo, dst string) error {
	// Validate and construct the destination path
	dstPath, pathErr := securePath(dst, f.NameInArchive)
	if pathErr != nil {
		return pathErr
	}

	// Ensure the parent directory exists
	parentDir := filepath.Dir(dstPath)

	dirErr := createDirWithPermissions(parentDir, dirPermissions)
	if dirErr != nil {
		return dirErr
	}

	// Handle directories
	if f.IsDir() {
		// Create the directory with permissions from the archive
		dirErr := createDirWithPermissions(dstPath, f.Mode())
		if dirErr != nil {
			return fmt.Errorf("creating directory: %w", dirErr)
		}

		return nil
	}

	// Ignore symlinks (or hardlinks)
	if f.LinkTarget != "" {
		if runtime.GOOS == "windows" {
			logrus.Warnf("Ignoring symlink: %s -> %s", f.NameInArchive, f.LinkTarget)
			return nil
		}

		linkErr := os.Symlink(f.LinkTarget, dstPath)
		if linkErr != nil {
			return fmt.Errorf("creating symlink: %w", linkErr)
		}
	}

	// Check and handle parent directory permissions
	originalMode, statErr := os.Stat(parentDir)
	if statErr != nil {
		return fmt.Errorf("stat parent directory: %w", statErr)
	}

	// If parent directory is read-only, temporarily make it writable
	if originalMode.Mode().Perm()&0o200 == 0 {
		chmodErr := os.Chmod(parentDir, originalMode.Mode()|0o200)
		if chmodErr != nil {
			return fmt.Errorf("chmod parent directory: %w", chmodErr)
		}

		defer func() {
			// Restore the original permissions after writing
			chmodErr := os.Chmod(parentDir, originalMode.Mode())
			if chmodErr != nil {
				logrus.Warnf("Failed to restore original permissions for %s: %v", parentDir, chmodErr)
			}
		}()
	}

	// Handle regular files
	reader, openErr := f.Open()
	if openErr != nil {
		return fmt.Errorf("open file: %w", openErr)
	}
	defer reader.Close()

	dstFile, createErr := os.OpenFile(dstPath, os.O_CREATE|os.O_WRONLY, f.Mode())
	if createErr != nil {
		return fmt.Errorf("create file: %w", createErr)
	}
	defer dstFile.Close()

	if _, copyErr := io.Copy(dstFile, reader); copyErr != nil {
		return fmt.Errorf("copy: %w", copyErr)
	}

	return nil
}
