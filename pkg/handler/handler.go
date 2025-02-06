package handler

type Installer interface {
	Download() error // download binaries assets
	Unpack() error   // unpack the tar file
	Setup() error    // setting up the symbolic links, generate the helper scripts, etc...
}

type Server interface {
	Serve() error // A function that responsible for ssh client command line
}
