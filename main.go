package main

import (
	"github.com/sirupsen/logrus"
	"sshd/pkg/sshd"
)

func main() {
	if err := sshd.SSHD(); err != nil {
		logrus.Errorf("Error: %s", err)
		return
	}
}
