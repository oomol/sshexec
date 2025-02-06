package ffmpeg

import (
	"context"
	"errors"
	"fmt"
	"github.com/Code-Hex/pget"
	"github.com/gliderlabs/ssh"
	"os"
	"path/filepath"
	"sshd/pkg/define"
	"sshd/pkg/sio"
)

// Installer is an interface that defines the methods to download, unpack and setup ffmpeg binaries
type Installer struct {
	URL     string
	PREFIX  string
	Session ssh.Session
}

func (i Installer) Download(ctx context.Context) error {
	sio.Printf(i.Session, "Downloading ffmpeg binaries from %q\n", i.URL)
	pGet := pget.New()
	if err := pGet.Run(ctx, "1.0", []string{
		"-p", "4",
		"-o", filepath.Join(os.TempDir(), "ffmpeg.tar.xz"),
		define.FFReleaseURL,
	}); err != nil {
		return fmt.Errorf("download ffmpeg failed: %v", err)
	}

	// TODO Checksum of the downloaded file

	return nil
}

func (i Installer) Unpack(ctx context.Context) error {
	return nil
}

func (i Installer) Setup(ctx context.Context) error {
	return nil
}

func GetStudioHomeDir() (string, error) {
	d, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("failed to get user home directory")
	}
	d = filepath.Join(d, define.StudioDir)
	return d, nil
}
