package ffmpeg

import (
	"context"
	"errors"
	"fmt"
	"github.com/Code-Hex/pget"
	"github.com/gliderlabs/ssh"
	"os"
	"path/filepath"
	"sshd/pkg/decompress"
	"sshd/pkg/define"
	"sshd/pkg/hash"
	"sshd/pkg/sio"
)

// Installer is an interface that defines the methods to download, unpack and setup ffmpeg binaries
type Installer struct {
	URL         string
	PREFIX      string
	FFMPEGTarXZ string
	Session     ssh.Session
}

func (i *Installer) Download(ctx context.Context) error {
	sio.Printf(i.Session, "Downloading ffmpeg binaries from %q\n", i.URL)
	ffmpegTarXZFile := filepath.Join(os.TempDir() + "ffmpeg.tar.xz")
	pGet := pget.New()
	if err := pGet.Run(ctx, "1.0", []string{
		"-p", "4",
		"-o", ffmpegTarXZFile,
		define.FFReleaseURL,
	}); err != nil {
		return fmt.Errorf("download ffmpeg failed: %v", err)
	}

	err := hash.CmpFileChecksum(ffmpegTarXZFile, define.FFMSha256)
	if err != nil {
		return fmt.Errorf("checksum failed: %v", err)
	}
	sio.Println(i.Session, "Download successful")
	i.FFMPEGTarXZ = ffmpegTarXZFile
	return nil
}

func (i *Installer) Unpack(ctx context.Context) error {
	if i.FFMPEGTarXZ == "" {
		return errors.New("ffmpeg tar.xz file not found")
	}
	sio.Println(i.Session, "Unpacking ffmpeg binaries")
	err := decompress.Unpack(ctx, i.FFMPEGTarXZ, i.PREFIX)
	if err != nil {
		return fmt.Errorf("unpack ffmpeg failed: %v", err)
	}

	return nil
}

func (i *Installer) Setup(ctx context.Context) error {
	sio.Println(i.Session, "Setting up ffmpeg binaries")

	return nil
}

func (i *Installer) Test(ctx context.Context) error {
	sio.Println(i.Session, "Testing ffmpeg")
	//ffmpegBin := filepath.Join(i.PREFIX, "ffmpeg", "")
	//cmd := exec.Command(i.PREFIX, "-version")

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
