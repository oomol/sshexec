package main

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func init() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     false,
		DisableColors:   true,
		TimestampFormat: "2006-01-02 15:04:05.000",
	})
	logrus.SetOutput(os.Stderr)
}

const (
	addr     = "127.0.0.1"
	user     = "oomol"
	port     = "5322"
	endpoint = user + "@" + addr
)

func main() {
	targetName := filepath.Base(os.Args[0])

	fullFFMPEGArgsList := append([]string{targetName}, os.Args[1:]...)

	var argsBuilder strings.Builder
	// Wrap the arguments with single quotes
	for _, arg := range fullFFMPEGArgsList {
		str := fmt.Sprintf("%s%s%s", "'", arg, "'")
		argsBuilder.WriteString(str)
		argsBuilder.WriteString(" ")
	}
	fullFFMPEGArgString := argsBuilder.String()
	_, _ = fmt.Fprintf(os.Stderr, "ffmpeg cmdline: %q\n", fullFFMPEGArgString)

	var finalArgs strings.Builder
	finalArgs.WriteString(fullFFMPEGArgString)
	fullArgs := finalArgs.String()

	cmd := exec.Command("ssh", "-q", "-o", "StrictHostKeyChecking=no",
		"-p", port, endpoint, fullArgs)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	_, _ = fmt.Fprintf(os.Stderr, "full cmdline: %q\n", cmd.Args)
	if err := cmd.Run(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
	}
}
