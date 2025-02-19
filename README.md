# Host utils Installer

> 注意这并不是一个标准的 SSHD 服务，只是使用 OpenSSH 协议调用一组 `Handler` 将 ffmpeg 或者 whisper 等
> 需要调用MacOS硬件资源能力的程序安装到本机，并提供对应的调用接口。

# 什么是 Handler
- Handler 是一组脚本或者一串代码，用于快速安装并设置 ffmpeg 或者 whisper
- Handler 也提供对外调用的接口（通过标准 SSH 协议）

# 如何对外提供调用接口
使用 ssh 直接执行命令

```shell
ssh -p <PORT> localhost ffmpeg [args...]
```

```shell
ssh -p <PORT> localhost whisper [args...]
```

> 这个程序不提供 Handler 之外的访问，也就是不信任容器内的 SSH 访问请求，如
> `ssh -p <PORT> localhost cat .ssh/id_ed25519` (获取 macos host 私钥)，则会被拒绝执行
> 因为此程序只提供 Handler 的安装逻辑 ，也只对外提供 Handler 的调用接口。