package ffmpeg

import (
	"context"
	"github.com/gliderlabs/ssh"
	"sshd/pkg/handler"
)

type Runner struct {
	File    string
	Args    []string
	Envs    []string
	Session ssh.Session
}

func (r *Runner) Run(ctx context.Context) error {
	return handler.ExecHandler(r.Session, ctx, r.File, r.Envs, r.Args)
}
