package sshd

import (
	"context"
	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
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

func InstallFFMPEG(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		str := s.Command()
		if len(str) == 0 {
			logrus.Warn("SSH Client give empty command")
			sio.Fatalln(s, "Empty command")
		}

		// If given command is define.InstallFFMPEG, try to install ffmpeg
		if str[0] == define.InstallFFMPEG {
			sio.Println(s, "Try to install ffmpeg")

			d, err := ffmpeg.GetStudioHomeDir()
			if err != nil {
				sio.Fatalf(s, "GetStudioHomeDir error: %s \n", err.Error())
			}

			sio.Printf(s, "GetStudioHomeDir: %q \n", d)
			ffmpegInstaller := ffmpeg.Installer{
				PREFIX:  d,
				URL:     define.FFReleaseURL,
				Session: s,
			}

			if err = ffmpegInstaller.Download(); err != nil {
				sio.Fatalln(s, "Download ffmpeg error: %s", err.Error())
			}

			if err = ffmpegInstaller.Unpack(); err != nil {
				sio.Fatalln(s, "Unpack ffmpeg error: %s", err.Error())
			}

			if err = ffmpegInstaller.Setup(); err != nil {
				sio.Fatalln(s, "Setup ffmpeg error: %s", err.Error())
			}
		} else {
			next(s)
		}
	}
}

func MyMiddleware(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		next(s)
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

func SSHD() error {
	addr := define.Addr + ":" + define.Port
	logrus.Infof("Starting SSH server at %s", addr)
	return ssh.ListenAndServe(addr, nil, WithMiddleware(
		ExecCmd,
		InstallFFMPEG,
		MyMiddleware,
	))
}
