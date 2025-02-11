package sshd

import (
	"context"
	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
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

			runner := ffmpeg.Runner{
				File:    ffBin,           // ffmpeg/ffprobe is the target binary
				Args:    s.Command()[1:], // Pass the rest of the command to ffmpeg
				Envs:    []string{env},
				Session: s,
			}

			logrus.Infof("Run ffmpeg: %q with args %q with env %q", runner.File, runner.Args, runner.Envs)
			if err = runner.Run(ctx); err != nil {
				logrus.Errorf("Run cmd error: %s", err.Error())
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

//func ExecCmd(next ssh.Handler) ssh.Handler {
//	return handler.ExecHandler()
//}

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
