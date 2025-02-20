# install.sh 说明

`install.sh` 会使用标准的 ssh client 发送 `install_ffmpeg` 给 `sshexec`, `sshexec` 会安装 oomol 提供的 ffmpeg 二进制文件
`install.sh` 会将 `caller.sh` 链接到 `/usr/bin/ffmpeg` 和 `/usr/bin/ffprobe`，在 Container 中调用 `ffmpeg` 和 `ffprobe` 时，会通过 sshexec 被转发到 Host 端的 ffmpeg 和 ffprobe。

# caller.sh 说明

`caller.sh` 运行在 Container 中，用于呼叫 Host 端的二进制文件，目前支持的二进制文件有 `ffmpeg` 和 `whisper`。

Caller.sh 需要被链接为 Host 端所支持的二进制文件名，例如 `ffmpeg` 或者 `whisper`。

```shell
ln -s /path/to/caller.sh /usr/local/bin/ffmpeg
ln -s /path/to/caller.sh /usr/local/bin/whisper
```

Container 调用方直接调用 `ffmpeg` 或者 `whisper`，Host 端对应的二进制文件会被直接调起，sshexec 会保证 `stdin/stdout` 传输到与 Container 完全同步。