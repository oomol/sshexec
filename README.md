[![Build and test](https://github.com/oomol/sshexec/actions/workflows/build.yaml/badge.svg)](https://github.com/oomol/sshexec/actions/workflows/build.yaml)

> Note that this is not a standard SSHD service, but a set of `Handlers` that use the OpenSSH protocol to call ffmpeg or whisper, etc.
> Programs that need to call MacOS hardware resource capabilities are installed on the local machine and provide corresponding calling interfaces.

# What is a Handler
- A Handler is a set of scripts or a string of code that is used to quickly install and set up ffmpeg or whisper
- A Handler also provides an interface for external calls (through the standard SSH protocol)

# How to provide external call interface
Use ssh to directly execute commands

```shell
ssh -p <PORT> localhost ffmpeg [args...]
```

```shell
ssh -p <PORT> localhost whisper [args...]
```

> This program does not provide access outside of the Handler, that is, it does not trust SSH access requests within the container, such as
> `ssh -p <PORT> localhost cat .ssh/id_ed25519` (obtaining the private key of the macos host), which will be rejected
> Because this program only provides the installation logic of the Handler and only provides the Handler's calling interface to the outside world.

# TODO
- [ ] support whisper installer
- [ ] save sshexec's stderr into file, keep stderr strictly the same from host process