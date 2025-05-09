package handler

import (
	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	"sshd/pkg/define"
	"sshd/pkg/logger"
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

func ValidateCmdline(next ssh.Handler) ssh.Handler {
	return func(s ssh.Session) {
		logrus.Infof("run middleware Sanitizers")
		// Parameter parsing follows the openssh standard implementation
		// https://stackoverflow.com/questions/53465980/how-to-keep-parameter-with-spaces-when-running-remote-script-file-with-ssh
		logrus.Infof("Validate string: %q", s.Command())
		str := s.Command()
		if len(str) == 0 {
			logger.Fatalf(s, "Empty command, Support commands: %q", define.Whitelist)
			return
		}

		// if the command is in the whitelist, then we can execute it, otherwise we will report an error
		if define.IsWhitelisted(str[0]) {
			next(s)
		} else {
			logger.Fatalf(s, "Command %q not allowed, Support commands: %q\n", str[0], define.Whitelist)
			return
		}
	}
}
