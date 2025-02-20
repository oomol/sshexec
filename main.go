package main

import (
	"fmt"
	"os"
	"sshd/pkg/flags"
	os2 "sshd/pkg/os"
	"sshd/pkg/sshd"
	"time"

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

	errChan := make(chan error, 1)
	go func() {
		errChan <- waitForStudioExit()
	}()

	go func() {
		errChan <- sshd.SSHExec()
	}()

	err := <-errChan
	logrus.Fatalf("main func error: %v", err)
}

func waitForStudioExit() error {
	for {
		isRunning, err := os2.IsProcessAliveV4(flags.OomolStudioPID)
		if err != nil {
			return fmt.Errorf("error within os2.IsProcessAliveV4: %v", err)
		}
		if !isRunning {
			return fmt.Errorf("OomolStudioPID [%d] not found exit", flags.OomolStudioPID)
		}
		time.Sleep(300 * time.Millisecond)
	}
}
