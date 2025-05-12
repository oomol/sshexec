package ffmpeg

import (
	"context"
	"fmt"

	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/Code-Hex/pget"
	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"

	"sshd/pkg/archiver"
	"sshd/pkg/define"
	"sshd/pkg/hash"
	"sshd/pkg/utils"
)

type Stubber struct {
	Session   ssh.Session
	Installer define.InstallOpts
	Runner    define.RunOpts
	Version   define.Version
}

const VERSION = "7"
const NAME = "ffmpeg_macos_arm64"

func New(s ssh.Session) *Stubber {
	stdioHome, err := utils.GetStudioHomeDir()
	if err != nil {
		logrus.Errorf("GetStudioHomeDir error: %v", err)
		return nil
	}

	logrus.Infof("GetStudioHomeDir: %q", stdioHome)

	return &Stubber{
		Session: s,
		Version: define.Version{
			PkgName: NAME,
			PkgVer:  VERSION,
		},
		Installer: define.InstallOpts{
			URL:       define.FFReleaseURLForVentura,
			Sha256Sum: define.FFMSha256ForVentura,
			Prefix:    filepath.Join(stdioHome, define.HostShared, NAME, VERSION),
		},
		Runner: define.RunOpts{
			FFMPEGPath:  filepath.Join(stdioHome, define.HostShared, NAME, VERSION, define.FFMPEG),
			FFPROBEPath: filepath.Join(stdioHome, define.HostShared, NAME, VERSION, define.FFPROBE),
		},
	}
}

func (l *Stubber) Run(ctx context.Context, target string, args, envs []string) error {
	l.Runner.Args = args
	l.Runner.Envs = envs

	cmd := exec.CommandContext(ctx, l.Runner.FFMPEGPath, l.Runner.Args...)

	if target == define.FFPROBE {
		cmd = exec.CommandContext(ctx, l.Runner.FFPROBEPath, l.Runner.Args...)
	}

	if l.Runner.Envs != nil {
		cmd.Env = append(cmd.Env, l.Runner.Envs...)
	}

	logrus.Infof("full cmdline: %q", cmd.Args)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		// logger.Fatalf(l.Session, "cmd.StdoutPipe() error: %s", err)
		return fmt.Errorf("cmd.StdoutPipe() error: %w", err)
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		// logger.Fatalf(l.Session, "cmd.StderrPipe() error: %s", err)
		return fmt.Errorf("cmd.StderrPipe() error: %w", err)
	}

	stdIn, err := cmd.StdinPipe()
	if err != nil {
		// logger.Fatalf(l.Session, "cmd.StdinPipe() error: %s", err)
		return fmt.Errorf("cmd.StdinPipe() error: %w", err)
	}

	if err = cmd.Start(); err != nil {
		// logger.Fatalf(l.Session, "cmd.Start() error: %v", err.Error())
		return fmt.Errorf("cmd.Start() error: %w", err)
	}

	// Copy cmd stdout to ssh session
	go func() {
		_, _ = io.Copy(l.Session, stdOut)
		logrus.Infof("Copy cmd stdout to ssh session's stdout finished")
		_ = stdOut.Close()
	}()

	// Copy cmd stderr to ssh session's stderr
	go func() {
		_, _ = io.Copy(l.Session.Stderr(), stdErr)
		logrus.Infof("Copy cmd stderr to ssh session's stderr finished")
		_ = stdErr.Close()
	}()

	// Copy stdin from session to cmd stdin
	go func() {
		_, err := io.Copy(stdIn, l.Session)
		if err != nil {
			logrus.Errorf("io.Copy(stdIn, l.Session) error: %v", err)
		}
		logrus.Infof("Copy stdin from session to cmd stdin finished")
		_ = stdIn.Close()
	}()

	if err = cmd.Wait(); err != nil {
		return fmt.Errorf("cmd.Wait() error: %w", err)
	}

	logrus.Infoln("cmdline execute finished")

	return nil
}

func (l *Stubber) CleanUp(ctx context.Context) error {
	return nil
}

func (l *Stubber) Setup(ctx context.Context) error {
	if err := os.Chmod(l.Runner.FFMPEGPath, 0755); err != nil {
		return fmt.Errorf("chmod ffmpeg failed: %w", err)
	}
	if err := os.Chmod(l.Runner.FFPROBEPath, 0755); err != nil {
		return fmt.Errorf("chmod ffprobe failed: %w", err)
	}

	return nil
}

func (l *Stubber) Test(ctx context.Context) error {
	if err := exec.CommandContext(ctx, l.Runner.FFMPEGPath, "-version").Start(); err != nil {
		// TODO: we should report the error into ssh session in caller
		return fmt.Errorf("test ffmpeg failed: %w", err)
	}

	if err := exec.CommandContext(ctx, l.Runner.FFPROBEPath, "-version").Start(); err != nil {
		return fmt.Errorf("test ffprobe failed: %w", err)
	}

	return nil
}

func (l *Stubber) Download(ctx context.Context) error {
	logrus.Infof("Download ffmpeg binaries from %q", l.Installer.URL)

	tempDir, err := os.MkdirTemp("", "dirs")
	if err != nil {
		return fmt.Errorf("create temp dir failed: %w", err)
	}

	logrus.Infof("temp dir: %q", tempDir)

	outFile := filepath.Join(tempDir, "ffmpeg.tar.xz")

	pGet := pget.New()
	if err = pGet.Run(ctx, "1.0", []string{
		"-p", "4",
		"-o", outFile,
		l.Installer.URL,
	}); err != nil {
		return fmt.Errorf("download ffmpeg failed: %w", err)
	}

	sha256sum, err := hash.Sha256sumFile(outFile)
	if err != nil {
		return fmt.Errorf("sha256sum file failed: %w", err)
	}

	if sha256sum != l.Installer.Sha256Sum {
		return fmt.Errorf("sha256sum mismatch: %q != %q", sha256sum, l.Installer.Sha256Sum)
	} else {
		logrus.Infof("sha256sum match: %q", sha256sum)
	}

	l.Installer.TarBar = outFile

	return nil
}

func (l *Stubber) Unpack(ctx context.Context) error {
	logrus.Infof("Unpack %q to %q", l.Installer.TarBar, l.Installer.Prefix)
	return archiver.Unarchive(l.Installer.TarBar, l.Installer.Prefix)
}
