package ffmpeg

import (
	"context"
	"github.com/gliderlabs/ssh"
	"os/exec"
)

type Runner struct {
	File    string
	Args    []string
	Envs    []string
	Session ssh.Session
}

func (r *Runner) Run(ctx context.Context) error {
	cmd := exec.CommandContext(ctx, r.File, r.Args...)
	cmd.Env = append(cmd.Env, r.Envs...)
	cmd.Stdout = r.Session
	cmd.Stderr = r.Session.Stderr()
	if err := cmd.Start(); err != nil {
		return err
	}

	if err := cmd.Wait(); err != nil {
		return err
	}

	return nil
}
