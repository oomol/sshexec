package main

import (
	"context"
	"os"

	"golang.org/x/sync/errgroup"

	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v3"

	"sshd/pkg/define"
	"sshd/pkg/handler"
	"sshd/pkg/handler/ffmpeg"
	"sshd/pkg/logger"
	os2 "sshd/pkg/os"
)

const (
	ppid   = "ppid"
	listen = "listen"
)

func main() {
	app := cli.Command{
		Flags: []cli.Flag{
			&cli.IntFlag{
				Name:  ppid,
				Usage: "parent process id",
			},
			&cli.StringFlag{
				Name:  listen,
				Usage: "listen address",
				Value: define.Addr + ":" + define.Port,
			},
		},
		Before: func(ctx context.Context, command *cli.Command) (context.Context, error) {
			return ctx, logger.SetupLogger()
		},
		Action: start,
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func start(ctx context.Context, command *cli.Command) error {
	g, ctx := errgroup.WithContext(ctx)

	parentPid := os.Getppid()
	if command.IsSet("ppid") {
		parentPid = command.Int(ppid)
	}

	os2.WatchPPID(g, ctx, parentPid)
	os2.ListenSignal(g, ctx)

	g.Go(func() error {
		errChan := make(chan error, 1)
		go func() {
			logrus.Infof("Start ssh server on %q", command.String(listen))
			errChan <- ssh.ListenAndServe(command.String(listen), nil, handler.WithMiddleware(
				ffmpeg.Run,
				ffmpeg.Install,
				handler.ValidateCmdline,
			))
		}()

		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		}
	})

	return g.Wait()
}
