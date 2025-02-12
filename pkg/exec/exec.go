package exec

import (
	"context"
	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	"io"
	"os/exec"
	"sshd/pkg/sio"
)

// Exec is a function that responsible for executing the target binary with the provided args and envs over the ssh session
func Exec(s ssh.Session, ctx context.Context, targetBin string, envs, args []string) error {
	cmd := exec.CommandContext(ctx, targetBin, args...)
	cmd.Env = append(cmd.Env, envs...)
	logrus.Infof("Run %q with args %q with env %q", targetBin, args, envs)

	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		sio.Fatalf(s, "cmd.StdoutPipe() error: %s", err)
		return err
	}

	stdErr, err := cmd.StderrPipe()
	if err != nil {
		sio.Fatalf(s, "cmd.StderrPipe() error: %s", err)
		return err
	}

	stdIn, err := cmd.StdinPipe()
	if err != nil {
		sio.Fatalf(s, "cmd.StdinPipe() error: %s", err)
		return err
	}

	// Copy cmd stdout to ssh session
	go func() {
		_, _ = io.Copy(s, stdOut)
	}()

	// Copy cmd stderr to ssh session's stderr
	go func() {
		_, _ = io.Copy(s.Stderr(), stdErr)
	}()

	// Copy stdin from session to cmd stdin
	go func() {
		_, _ = io.Copy(stdIn, s)
	}()

	if err = cmd.Start(); err != nil {
		sio.Fatalf(s, "cmd.Start() error: %v", err.Error())
		return err
	}

	if err = cmd.Wait(); err != nil {
		sio.Fatalf(s, "cmd.Wait() error: %v", err.Error())
		return err
	}

	return nil

}
