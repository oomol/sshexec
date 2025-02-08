package decompress

import (
	"context"
	"fmt"
	"github.com/mholt/archives"
	"io"
	"os"
	"path/filepath"
	"time"
)

func Extract(ctx context.Context, src, targetDir string) error {
	stream, err2 := os.Open(src)
	if err2 != nil {
		return err2
	}
	format, input, err := archives.Identify(ctx, src, stream)
	if err != nil {
		return err
	}

	ex, ok := format.(archives.Extractor)
	if !ok {
		return fmt.Errorf("failed to detect proper archive for extraction from %s got: %v", src, ex)
	}

	err = ex.Extract(ctx, input, func(_ context.Context, f archives.FileInfo) error {
		target := filepath.Join(targetDir, f.NameInArchive)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
		if f.IsDir() {
			return os.MkdirAll(target, f.Mode())
		} else if f.LinkTarget != "" {
			return os.Symlink(f.LinkTarget, target)
		}
		targetFile, err := os.Create(target)
		if err != nil {
			return fmt.Errorf("create %s: %w", target, err)
		}
		arc, err := f.Open()
		if err != nil {
			return err
		}
		if _, err := io.Copy(targetFile, arc); err != nil {
			return err
		}
		if err := arc.Close(); err != nil {
			return err
		}
		if err := targetFile.Close(); err != nil {
			return err
		}
		if err := os.Chmod(target, f.Mode()); err != nil {
			return err
		}
		if err := os.Chtimes(target, time.Time{}, f.ModTime()); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
