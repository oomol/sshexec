package logger

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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
	_, _ = fmt.Fprintf(s.Stderr(), f, v...)
}

func SetupLogger() error {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	logFile := filepath.Join(homeDir, ".oomol-studio", "ovm-krun", "logs", "sshexec.log")

	if runtime.GOARCH == "amd64" {
		logFile = filepath.Join(homeDir, ".oomol-studio", "ovm", "logs", "sshexec.log")
	}

	logDir := filepath.Dir(logFile)

	logrus.Infof("Try to make logDir dir: %q", logDir)
	if err = os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("failed to create log directory %q: %v", logDir, err)
	}

	logrus.Infof("Try to open log file: %q", logFile)
	fd, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return fmt.Errorf("failed to open log file %q: %v", logFile, err)
	}

	logrus.SetOutput(fd)
	logrus.SetLevel(logrus.InfoLevel)
	return nil
}
