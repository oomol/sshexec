package sshd

import (
	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	"sshd/pkg/define"
	"sshd/pkg/handler"
)

func SSHExec() error {
	addr := define.Addr + ":" + define.Port
	logrus.Infof("Starting SSH server at %s", addr)
	return ssh.ListenAndServe(addr, nil, handler.WithMiddleware(
		//ExecCmd,
		handler.RunFFMPEG,
		handler.InstallFFMPEG,
		handler.Sanitizers,
	))
}
