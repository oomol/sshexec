package ffmpeg

import (
	"sshd/pkg/define"
	"sshd/pkg/exec"
	slog "sshd/pkg/logger"
	"sshd/pkg/provider/ffmpeg"
	"sync"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
)

const RunFFMPEGStage = "ffmpeg run handler"

func Run(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		targetBin := s.Command()[0]
		if targetBin == define.FFPROBEBin || targetBin == define.FFMPEGBin {
			logrus.Infof("run middleware: %q\n", RunFFMPEGStage)

			stubber := ffmpeg.NewVersion6(s)

			args, err := exec.DoArgsSanitizers(s.Command()[1:])
			if err != nil {
				slog.Fatalf(s, "DoArgsSanitizers error: %v\r\n", err)
				return
			}

			if err = stubber.Run(s.Context(), targetBin, args, nil); err != nil {
				slog.Fatalf(s, "Run with error: %v\r\n", err)
				return
			}
		}

		// next handler if no error
		next(s)
	}
}

const InstallStage = "ffmpeg install handler"

var lock sync.Mutex

func Install(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		// if the command is not InstallFFMPEGVersion6, do nothing and run the next handler
		if s.Command()[0] == define.InstallFFMPEGVersion6 {
			lock.Lock()
			defer lock.Unlock()

			slog.Infof(s, "run middleware: %q\r\n", InstallStage)
			stubber := ffmpeg.NewVersion6(s)

			// we first test the ffmpeg call be called, if ffmpeg can be called without error, just return
			err := stubber.Test(s.Context())
			if err == nil {
				slog.Infof(s, "ffmpeg package installed before\r\n")
				return
			}

			if err := stubber.Download(s.Context()); err != nil {
				slog.Fatalf(s, "Download ffmpeg error: %v\r\n", err)
				return
			}

			slog.Infof(s, "Download ffmpeg success\r\n")

			if err := stubber.Unpack(s.Context()); err != nil {
				slog.Fatalf(s, "Unpack ffmpeg error: %v\r\n", err)
				return
			}

			slog.Infof(s, "Unpack ffmpeg success\r\n")

			if err := stubber.Setup(s.Context()); err != nil {
				slog.Fatalf(s, "Setup ffmpeg error: %v", err)
				return
			}

			slog.Infof(s, "Setup ffmpeg success\r\n")

			if err := stubber.Test(s.Context()); err != nil {
				slog.Fatalf(s, "Test ffmpeg package error: %v\r\n", err)
				return
			}

			slog.Infof(s, "Test ffmpeg package success\r\n")

			if err := stubber.CleanUp(s.Context()); err != nil {
				slog.Fatalf(s, "Clean up ffmpeg error: %v\r\n", err)
				return
			}

			slog.Infof(s, "Clean up ffmpeg success\r\n")

			return
		}

		// next handler if no error
		next(s)
	}
}
