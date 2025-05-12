package ffmpeg

import (
	"sshd/pkg/define"
	"sshd/pkg/exec"
	"sshd/pkg/logger"
	"sshd/pkg/provider/ffmpeg"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
)

const RunFFMPEGStage = "ffmpeg run handler"

func Run(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		targetBin := s.Command()[0]
		if targetBin == define.FFPROBE || targetBin == define.FFMPEG {
			logrus.Infof("run middleware: %q", RunFFMPEGStage)
			stubber := ffmpeg.New(s)
			args, err := exec.DoArgsSanitizers(s.Command()[1:])
			if err != nil {
				logger.Fatalf(s, "DoArgsSanitizers error: %v", err)
				return
			}

			if err = stubber.Run(s.Context(), targetBin, args, nil); err != nil {
				logger.Fatalf(s, "Run with error: %v", err)
				return
			}
		}

		// next handler if no error
		next(s)
	}
}

const InstallStage = "ffmpeg install handler"

func Install(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		// if the command is not define.InstallFFMPEG, do nothing and run next handler
		if s.Command()[0] == define.InstallFFMPEG {
			logrus.Infof("run middleware: %q", InstallStage)
			stubber := ffmpeg.New(s)
			if err := stubber.Download(s.Context()); err != nil {
				logger.Fatalf(s, "Download ffmpeg error: %v", err)
				return
			}
			logger.Infof(s, "Download ffmpeg success\n")

			if err := stubber.Unpack(s.Context()); err != nil {
				logger.Fatalf(s, "Unpack ffmpeg error: %v", err)
				return
			}
			logger.Infof(s, "Unpack ffmpeg success\n")

			if err := stubber.Setup(s.Context()); err != nil {
				logger.Fatalf(s, "Setup ffmpeg error: %v", err)
				return
			}
			logger.Infof(s, "Setup ffmpeg success\n")

			if err := stubber.Test(s.Context()); err != nil {
				logger.Fatalf(s, "Test ffmpeg package error: %v", err)
				return
			}
			logger.Infof(s, "Test ffmpeg package success\n")

			if err := stubber.CleanUp(s.Context()); err != nil {
				logger.Fatalf(s, "Clean up ffmpeg error: %v", err)
				return
			}
			logger.Infof(s, "Clean up ffmpeg success\n")
		}

		// next handler if no error
		next(s)
	}
}
