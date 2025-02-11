package sshd

import (
	"context"
	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"path/filepath"
	"sshd/pkg/define"
	"sshd/pkg/ffmpeg"
	"sshd/pkg/sio"
)

type Middleware func(next ssh.Handler) ssh.Handler

// WithMiddleware composes the provided Middleware and returns an ssh.Option.
// This is useful if you manually create an ssh.Server and want to set the
// Server.Handler.
//
// Notice that middlewares are composed from first to last, which means the last one is executed first.
func WithMiddleware(mw ...Middleware) ssh.Option {
	return func(s *ssh.Server) error {
		h := func(ssh.Session) {}
		for _, m := range mw {
			h = m(h)
		}
		s.Handler = h
		return nil
	}
}

func RunFFMPEG(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		if s.Command()[0] == define.FFMPEG || s.Command()[0] == define.FFPROBE {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				// The context is canceled when the client's connection closes or I/ O operation fails.
				<-s.Context().Done()
				cancel()
			}()

			studioHome, err := ffmpeg.GetStudioHomeDir()
			if err != nil {
				sio.Fatalf(s, "GetStudioHomeDir error: %s\n", err.Error())
			}
			// ffmpeg or ffprobe
			ffBin := filepath.Join(studioHome, "host-shared", "ffmpeg", "bin", s.Command()[0])

			env := define.DYLD_LIBRARY_PATH + "=" + filepath.Join(studioHome, "host-shared", "ffmpeg", "lib")

			ffmpegELF := ffmpeg.Runner{
				File:    ffBin,           // ffmpeg
				Args:    s.Command()[1:], // Pass the rest of the command to ffmpeg
				Envs:    []string{env},
				Session: s,
			}

			logrus.Infof("Run ffmpeg: %q with args %q with env %q", ffmpegELF.File, ffmpegELF.Args, ffmpegELF.Envs)
			if err = ffmpegELF.Run(ctx); err != nil {
				sio.Fatalf(s, "Run ffmpeg error: %s\n", err.Error())
			}
		} else {
			next(s)
		}
	}
}

func InstallFFMPEG(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		if s.Command()[0] == define.InstallFFMPEG {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()

			go func() {
				// The context is canceled when the client's connection closes or I/ O operation fails.
				<-s.Context().Done()
				cancel()
			}()

			// If given command is define.InstallFFMPEG, try to install ffmpeg

			sio.Println(s, "Try to install ffmpeg")

			stduioHome, err := ffmpeg.GetStudioHomeDir()
			if err != nil {
				sio.Fatalf(s, "GetStudioHomeDir error: %s \n", err.Error())
			}

			sio.Printf(s, "GetStudioHomeDir: %q \n", stduioHome)
			ffmpegInstaller := ffmpeg.Installer{
				PREFIX:  filepath.Join(stduioHome, "host-shared"),
				URL:     define.FFReleaseURL,
				Session: s,
			}

			if err = ffmpegInstaller.Download(ctx); err != nil {
				sio.Fatalf(s, "Download ffmpeg error: %s\n", err.Error())
			}

			if err = ffmpegInstaller.Unpack(ctx); err != nil {
				sio.Fatalf(s, "Unpack ffmpeg error: %s\n", err.Error())
			}

			if err = ffmpegInstaller.Setup(ctx); err != nil {
				sio.Fatalf(s, "Setup ffmpeg error: %s\n", err.Error())
			}
			if err = ffmpegInstaller.Test(ctx); err != nil {
				sio.Fatalf(s, "Test ffmpeg error: %s\n", err.Error())
			}
		} else {
			next(s)
		}
	}
}

func ExecCmd(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		logrus.Info("ExecCmd Middleware")
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		go func() {
			// The context is canceled when the client's connection closes or I/ O operation fails.
			<-s.Context().Done()
			cancel()
		}()

		str := s.Command()
		if len(str) == 0 {
			logrus.Warn("Empty command")
			_ = s.Exit(128)
			return
		}

		logrus.Infof("Command: %q", str)
		cmd := exec.CommandContext(ctx, str[0], str[1:]...)
		stdOut, err := cmd.StdoutPipe()
		if err != nil {
			logrus.Errorf("cmd.StdoutPipe() error: %s", err)
			_ = s.Exit(100)
			return
		}

		stdErr, err := cmd.StderrPipe()
		if err != nil {
			logrus.Errorf("cmd.StderrPipe() error: %s", err)
			_ = s.Exit(100)
			return
		}

		err = cmd.Start()
		if err != nil {
			logrus.Errorf("cmd.Start() error: %s", err)
			_ = s.Exit(127) //nolint:mnd
		}

		// Copy cmd stdout to ssh session
		go func() {
			_, _ = io.Copy(s, stdOut)
		}()

		// Copy cmd stderr to ssh session's stderr
		go func() {
			_, _ = io.Copy(s.Stderr(), stdErr)
		}()

		if err = cmd.Wait(); err != nil {
			_ = s.Exit(cmd.ProcessState.ExitCode())
			logrus.Errorf("cmd.Wait() error: %s", err)
		} else {
			logrus.Infof("Command %q finished", str)
		}
		next(s)
	}
}

func Sanitizers(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		logrus.Info("Commandline sanitizing")
		str := s.Command()
		if len(str) == 0 {
			sio.Fatalf(s, "Empty command, Support commands: %q", define.Whitelist)
		}

		// Sanitizing the command with whitelist
		if define.IsWhitelisted(str[0]) {
			next(s)
		} else {
			sio.Fatalf(s, "Command %q not allowed, Support commands: %q\n", str[0], define.Whitelist)
		}
	}
}

func SSHExec() error {
	addr := define.Addr + ":" + define.Port
	logrus.Infof("Starting SSH server at %s", addr)
	return ssh.ListenAndServe(addr, nil, WithMiddleware(
		//ExecCmd,
		RunFFMPEG,
		InstallFFMPEG,
		Sanitizers,
	))
}
