package sshd

import (
	"sshd/pkg/define"
	"sshd/pkg/handler"
	"sshd/pkg/handler/ffmpeg"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
)

func Server() error {
	addr := define.Addr + ":" + define.Port
	logrus.Infof("Starting SSH server at %s", addr)

	return ssh.ListenAndServe(addr, nil, handler.WithMiddleware(
		ffmpeg.Run,
		ffmpeg.Install,
		handler.ValidateCmdline,
	))
}
