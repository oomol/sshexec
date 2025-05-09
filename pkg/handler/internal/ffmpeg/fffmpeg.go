package ffmpeg

import (
	"context"
	"errors"
	"fmt"

	"io"
	"os"
	"os/exec"
	"path/filepath"

	"sshd/pkg/decompress"
	myexec "sshd/pkg/exec"
	"sshd/pkg/hash"
	"sshd/pkg/sio"

	"github.com/Code-Hex/pget"
	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
)

type Runner struct {
	File    string
	Args    []string
	Envs    []string
	Session ssh.Session
}

type Installer struct {
	URL         string
	PREFIX      string
	FFMPEGTarXZ string
	Sha256Sum   string
	Session     ssh.Session
}

func (r *Runner) Run(ctx context.Context) error {
	return myexec.Exec(r.Session, ctx, r.File, r.Envs, r.Args)
}

// CleanUp is a function that responsible for cleaning up the resources generated by the Runner
func (r *Runner) CleanUp(ctx context.Context) error {
	return nil
}

func (i *Installer) Download(ctx context.Context) error {
	sio.Printf(i.Session, "Downloading ffmpeg binaries from %q\n", i.URL)
	ffmpegTarXZFile := filepath.Join(os.TempDir(), "ffmpeg.tar.xz")
	if err := os.RemoveAll(ffmpegTarXZFile); err != nil {
		return fmt.Errorf("remove %q failed: %v", ffmpegTarXZFile, err)
	}

	logrus.Infof("Download ffmpeg binaries from %q using pGet", i.URL)
	pGet := pget.New()
	if err := pGet.Run(ctx, "1.0", []string{
		"-p", "4",
		"-o", ffmpegTarXZFile,
		i.URL,
	}); err != nil {
		return fmt.Errorf("download ffmpeg failed: %w", err)
	}

	logrus.Infof("Do sum check with %q", i.Sha256Sum)
	if sum, err := hash.CmpFileChecksum(ffmpegTarXZFile, i.Sha256Sum); err != nil {
		return fmt.Errorf("checksum failed: %w, want %q ,got %q", err, i.Sha256Sum, sum)
	}

	sio.Println(i.Session, "Download successful")
	i.FFMPEGTarXZ = ffmpegTarXZFile
	logrus.Infof("Download successful")
	return nil
}

func (i *Installer) Unpack(ctx context.Context) error {
	if i.FFMPEGTarXZ == "" {
		return errors.New("ffmpeg tar.xz file not found")
	}
	sio.Println(i.Session, "Unpacking ffmpeg binaries")
	if err := os.RemoveAll(filepath.Join(i.PREFIX, "ffmpeg")); err != nil {
		return err
	}
	err := decompress.Extract(ctx, i.FFMPEGTarXZ, i.PREFIX)
	if err != nil {
		return fmt.Errorf("unpack ffmpeg failed: %w", err)
	}

	return nil
}

func (i *Installer) Setup(ctx context.Context) error {
	sio.Println(i.Session, "Setting up ffmpeg binaries")
	return nil
}

func (i *Installer) Test(ctx context.Context) error {
	ffBin := filepath.Join(i.PREFIX, "ffmpeg", "bin", "ffmpeg")
	sio.Printf(i.Session, "Testing %q\n", ffBin)
	if err := os.Chmod(ffBin, 0755); err != nil {
		return fmt.Errorf("chmod ffmpeg failed: %w", err)
	}

	cmd := exec.CommandContext(ctx, ffBin, "-version")
	cmd.Env = append(cmd.Env, fmt.Sprintf("DYLD_LIBRARY_PATH=%s", filepath.Join(i.PREFIX, "ffmpeg", "libs")))
	cmd.Stdout = io.MultiWriter(i.Session, os.Stdout)
	cmd.Stderr = io.MultiWriter(i.Session.Stderr(), os.Stderr)

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("test ffmpeg failed: %w", err)
	}

	sio.Println(i.Session, "Test successful")
	return nil
}

func (i *Installer) CleanUp(ctx context.Context) error {
	if i.FFMPEGTarXZ != "" {
		if err := os.Remove(i.FFMPEGTarXZ); err != nil {
			return fmt.Errorf("remove install package failed: %w", err)
		}
	}
	return nil
}
