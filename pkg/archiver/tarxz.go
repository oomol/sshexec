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

	format, input, identifyErr := archives.Identify(context.Background(), tarball, archiveFile)
	if identifyErr != nil {
		return fmt.Errorf("identify format: %w", identifyErr)
	}

	extractor, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("unsupported format for extraction")
	}

	if dirErr := createDirWithPermissions(dst, dirPermissions); dirErr != nil {
		return fmt.Errorf("creating destination directory: %w", dirErr)
	}

	handler := func(ctx context.Context, f archives.FileInfo) error {
		return handleFile(f, dst)
	}

	if extractErr := extractor.Extract(context.Background(), input, handler); extractErr != nil {
		return fmt.Errorf("extracting files: %w", extractErr)
	}

	return nil
}

func createDirWithPermissions(path string, mode os.FileMode) error {
	if err := os.MkdirAll(path, mode); err != nil {
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
	if dirErr := createDirWithPermissions(parentDir, dirPermissions); dirErr != nil {
		return dirErr
	}

	// Handle directories
	if f.IsDir() {
		// Create the directory with permissions from the archive
		if dirErr := createDirWithPermissions(dstPath, f.Mode()); dirErr != nil {
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
		if linkErr := os.Symlink(f.LinkTarget, dstPath); linkErr != nil {
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
		if chmodErr := os.Chmod(parentDir, originalMode.Mode()|0o200); chmodErr != nil {
			return fmt.Errorf("chmod parent directory: %w", chmodErr)
		}
		defer func() {
			// Restore the original permissions after writing
			if chmodErr := os.Chmod(parentDir, originalMode.Mode()); chmodErr != nil {
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
