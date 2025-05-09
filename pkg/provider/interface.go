package provider

import "context"

type Provider interface {
	Download(ctx context.Context) error                                // download binaries assets
	Unpack(ctx context.Context) error                                  // unpack the tar file
	Setup(ctx context.Context) error                                   // setting up the symbolic links, generate the helper scripts, etc...
	Test(ctx context.Context) error                                    // test the binaries
	Run(ctx context.Context, target string, args, envs []string) error // A function that responsible for ssh client command line
	CleanUp(ctx context.Context) error
}
