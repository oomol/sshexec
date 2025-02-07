package decompress

import (
	"context"
	"fmt"
	"github.com/mholt/archives"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

func Decompress(ctx context.Context, src, dest string) error {
	format, _, err := archives.Identify(ctx, src, nil)
	if err != nil {
		return fmt.Errorf("failed to identify %s: %s", src, err)
	}

	// Decompress tar.xz to tar file
	if decomp, ok := format.(archives.Decompressor); ok {
		f, err := os.Open(src)
		if err != nil {
			return fmt.Errorf("failed to open %s: %s", src, err)
		}
		rc, err := decomp.OpenReader(f)
		if err != nil {
			return err
		}
		defer rc.Close()
		destFile, err := os.Create(dest)
		if err != nil {
			return fmt.Errorf("failed to create destination file %s: %s", dest, err)
		}
		defer destFile.Close()

		if _, err := io.Copy(destFile, rc); err != nil {
			return fmt.Errorf("failed to copy decompressed data to %s: %s", dest, err)
		}
	}

	return nil
}

func Extract(ctx context.Context, src, dest string) error {
	defer os.Remove(src)
	format, _, err := archives.Identify(ctx, src, nil)
	if err != nil {
		return fmt.Errorf("failed to identify %s: %s", src, err)
	}

	if ex, ok := format.(archives.Extractor); ok {
		tarFile, err := os.Open(src)
		if err != nil {
			return fmt.Errorf("failed to open %s: %s", src, err)
		}
		if err = ex.Extract(ctx, tarFile, func(ctx context.Context, f archives.FileInfo) error {
			if f.FileInfo.IsDir() {
				logrus.Infof("Creating directory %s", filepath.Join(dest, f.NameInArchive))
				if err := os.MkdirAll(filepath.Join(dest, f.NameInArchive), f.Mode()); err != nil {
					return fmt.Errorf("failed to create directory %s: %s", filepath.Join(dest, f.NameInArchive), err)
				}
			} else {
				logrus.Infof("Extracting file %s", filepath.Join(dest, f.NameInArchive))
				file_, err := f.Open()
				if err != nil {
					return fmt.Errorf("failed to open file %s: %s", f.NameInArchive, err)
				}
				defer file_.Close()

				myFile_, err := os.OpenFile(filepath.Join(dest, f.NameInArchive), os.O_CREATE|os.O_WRONLY, f.Mode())
				if err != nil {
					return fmt.Errorf("failed to open file %s: %s", filepath.Join(dest, f.NameInArchive), err)
				}
				_, _ = io.Copy(myFile_, file_)
				defer myFile_.Close()
			}

			return nil
		}); err != nil {
			return fmt.Errorf("failed to extract %s: %s", src, err)
		}

	}

	return nil
}
