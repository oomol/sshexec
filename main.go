package main

import (
	"github.com/sirupsen/logrus"
	"sshd/pkg/sshd"
)

func main() {
	if err := sshd.SSHExec(); err != nil {
		logrus.Fatalf("Error: %s", err)
	}
}
