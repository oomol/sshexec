package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sshd/pkg/flags"
	os2 "sshd/pkg/os"
	"sshd/pkg/sshd"
	"time"

	"github.com/sirupsen/logrus"
)

func init() {
	_ = os.Unsetenv("http_proxy")
	_ = os.Unsetenv("https_proxy")
	_ = os.Unsetenv("HTTP_PROXY")
	_ = os.Unsetenv("HTTPS_PROXY")
	_ = os.Unsetenv("ftp_proxy")
	_ = os.Unsetenv("FTP_PROXY")
}

func setupLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		DisableColors:   false,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = "/tmp"
	}

	logFile := filepath.Join(homeDir, ".oomol-studio", "ovm-krun", "logs", "sshexec.log")
	if runtime.GOARCH == "amd64" {
		logFile = filepath.Join(homeDir, ".oomol-studio", "ovm", "logs", "sshexec.log")
	}

	if err = os.MkdirAll(filepath.Dir(logFile), 0755); err != nil {
		logFile = filepath.Join(homeDir, "sshexec.log")
	}

	fd, _ := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, os.ModePerm)
	logrus.SetOutput(fd)

	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	setupLogger()
	// assume always has parent pid
	ppid := os.Getppid()
	logrus.Infof("PPID: %d", ppid)
	flags.SetOomolStudioPID(ppid)

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
