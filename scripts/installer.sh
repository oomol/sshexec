#! /usr/bin/env bash
set -e
set -u

get_platform() {
	arch=$(uname -m)
	platform=unknown

	if [[ -z "$arch" ]]; then
		echo "uname -m return empty"
		return
	fi

	# For wsl2
	if [[ "$arch" == x86_64 ]] && [[ -d "/usr/lib/wsl" ]]; then
		platform="wsl2-$arch"
		return
	fi

	# For MacOS-x86_64
	if [[ "$arch" == x86_64 ]]; then
		platform="macos-$arch"
		return
	fi

	# For MacOS-aarch64
	if [[ "$arch" == aarch64 ]] || [[ $arch == arm64 ]]; then
		platform="macos-$arch"
		return
	fi
}

setup-macos-arm64() {
	wget https://static.oomol.com/sshexec/1.0.7/caller-arm64 --output-document=/usr/bin/caller
	chmod +x /usr/bin/caller
	ln -sf /usr/bin/caller /usr/bin/ffmpeg
	ln -sf /usr/bin/caller /usr/bin/ffprobe
	ln -sf /usr/bin/caller /usr/bin/install_ffmpeg_6
	/usr/bin/install_ffmpeg_6
}

setup-macos-x86_64() {
	wget https://static.oomol.com/sshexec/1.0.7/caller-amd64 --output-document=/usr/bin/caller
	chmod +x /usr/bin/caller
	ln -sf /usr/bin/caller /usr/bin/ffmpeg
	ln -sf /usr/bin/caller /usr/bin/ffprobe
	ln -sf /usr/bin/caller /usr/bin/install_ffmpeg_6
	/usr/bin/install_ffmpeg_6
}

setup-wsl2-x86_64() {
	ffmpeg_tar="/tmp/jellyfin-ffmpeg_6.0.1-8_portable_linux64-gpl.tar.xz"
	echo "Install ffmpeg"
	wget https://static.oomol.com/sshexec/v1.0.11/jellyfin-ffmpeg_6.0.1-8_portable_linux64-gpl.tar.xz --output-document="$ffmpeg_tar"
	tar -xvf "$ffmpeg_tar" -C /usr/bin/
	chmod +x /usr/bin/ffmpeg
	chmod +x /usr/bin/ffprobe
	echo "Install ffmpeg done"
}

setup_ffmpeg() {
	if [[ "$platform" == macos-aarch64 ]]; then
		setup-macos-arm64
	elif [[ "$platform" == wsl2-x86_64 ]]; then
		setup-wsl2-x86_64
	elif [[ "$platform" == macos-x86_64 ]]; then
		setup-macos-x86_64
	else
		echo "unsupport platform: $platform"
		exit 100
	fi
}

main() {
	get_platform
	if [[ "$platform" == "unknown" ]]; then
		echo "unknown platform"
		exit 100
	fi
	setup_ffmpeg
}

main
