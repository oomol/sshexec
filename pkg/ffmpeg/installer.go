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
	// URL: https://github.com/oomol/builded/releases/download
	URL string // URL to download ffmpeg binaries
	// PREFIX: ~/.oomol-studio/host-shared/<package_name>
	PREFIX  string
	Session ssh.Session
}

func (i Installer) Download() error {
	sio.Printf(i.Session, "Downloading ffmpeg binaries from %q\n", i.URL)
	pGet := pget.New()
	err := pGet.Run(context.Background(), "1.0", []string{
		"-p", "4",
		"-o", filepath.Join(os.TempDir(), "ffmpeg.tar.xz"),
		define.FFReleaseURL,
	})
	if err != nil {
		return fmt.Errorf("download ffmpeg failed: %v", err)
	}
	return nil

}

func (i Installer) Unpack() error {
	return nil
}

func (i Installer) Setup() error {
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
