package sio

import (
	"fmt"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
)

// Fatal prints to the given session's STDERR and exits 1.
func Fatal(s ssh.Session, v ...interface{}) {
	Error(s, v...)
	_ = s.Exit(1)
	_ = s.Close()
}

// Fatalf formats according to the given format, prints to the session's STDERR
// followed by an exit 1.
//
// Notice that this might cause formatting issues if you don't add a \r\n in the end of your string.
func Fatalf(s ssh.Session, f string, v ...interface{}) {
	Errorf(s, f, v...)
	_ = s.Exit(1)
	_ = s.Close()
}

// Fatalln formats according to the default format, prints to the session's
// STDERR, followed by a new line and an exit 1.
func Fatalln(s ssh.Session, v ...interface{}) {
	Errorln(s, v...)
	Errorf(s, "\r")
	_ = s.Exit(1)
	_ = s.Close()
}

// Error prints the given error the the session's STDERR.
func Error(s ssh.Session, v ...interface{}) {
	logrus.Error(v...)
	_, _ = fmt.Fprint(s.Stderr(), v...)
}

// Errorf formats according to the given format and prints to the session's STDERR.
func Errorf(s ssh.Session, f string, v ...interface{}) {
	logrus.Errorf(f, v...)
	_, _ = fmt.Fprintf(s.Stderr(), f, v...)
}

// Errorf formats according to the default format and prints to the session's STDERR.
func Errorln(s ssh.Session, v ...interface{}) {
	logrus.Errorln(v...)
	_, _ = fmt.Fprintln(s.Stderr(), v...)
}

// Print writes to the session's STDOUT followed.
func Print(s ssh.Session, v ...interface{}) {
	logrus.Error(v...)
	_, _ = fmt.Fprint(s, v...)
}

// Printf formats according to the given format and writes to the session's STDOUT.
func Printf(s ssh.Session, f string, v ...interface{}) {
	logrus.Infof(f, v...)
	_, _ = fmt.Fprintf(s, f, v...)
}

// Println formats according to the default format and writes to the session's STDOUT.
func Println(s ssh.Session, v ...interface{}) {
	logrus.Infoln(v...)
	_, _ = fmt.Fprintln(s, v...)
}
