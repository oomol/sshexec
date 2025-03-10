package sshd

import (
	"sshd/pkg/define"
	"sshd/pkg/handler"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
)

func SSHExec() error {
	addr := define.Addr + ":" + define.Port
	logrus.Infof("Starting SSH server at %s", addr)

	return ssh.ListenAndServe(addr, nil, handler.WithMiddleware(
		handler.RunFFMPEG,
		handler.InstallFFMPEG,
		handler.Sanitizers,
	))
}
