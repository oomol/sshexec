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
	return RunningELF(ctx, r.File, r.Args, r.Envs, r.Session)
}

func RunningELF(ctx context.Context, elf string, args, envs []string, session ssh.Session) error {
	cmd := exec.CommandContext(ctx, elf, args...)
	cmd.Env = append(cmd.Env, envs...)
	cmd.Stdout = session
	cmd.Stderr = session.Stderr()
	return cmd.Run()
}
