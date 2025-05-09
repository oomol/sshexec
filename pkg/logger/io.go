package logger

import (
	"fmt"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
)

// Fatalf formats according to the given format, prints to the session's STDERR
// followed by an exit 1.
//
// Notice that this might cause formatting issues if you don't add a \r\n in the end of your string.
func Fatalf(s ssh.Session, f string, v ...interface{}) {
	Errorf(s, f, v...)
	_ = s.Exit(1)
}

func Errorf(s ssh.Session, f string, v ...interface{}) {
	logrus.Errorf(f, v...)
	_, _ = fmt.Fprintf(s.Stderr(), f, v...)
}

func Infof(s ssh.Session, f string, v ...interface{}) {
	logrus.Infof(f, v...)
	_, _ = fmt.Fprintf(s, f, v...)
}
