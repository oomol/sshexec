package main

import (
	"os"
	"sshd/pkg/flags"
	"sshd/pkg/sshd"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	logrus.SetLevel(logrus.InfoLevel)
}
func main() {
	// assume always has parent pid
	flags.SetOomolStudioPID(os.Getppid())

	if err := sshd.SSHExec(); err != nil {
		logrus.Fatalf("Error: %s", err)
	}
}
