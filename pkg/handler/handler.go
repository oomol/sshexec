package handler

import "context"

type Installer interface {
	Download(ctx context.Context) error // download binaries assets
	Unpack(ctx context.Context) error   // unpack the tar file
	Setup(ctx context.Context) error    // setting up the symbolic links, generate the helper scripts, etc...
	Test(ctx context.Context) error
}

type Runner interface {
	Runner(ctx context.Context) error // A function that responsible for ssh client command line
}
