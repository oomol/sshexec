#! /usr/bin/env bash
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

get_ver() {
	sshexec_version="$(ssh -o StrictHostKeyChecking=no -q root@192.168.127.254 -p5322 show_version | xargs)"
}

############ START OF MACOS ARM64 ############
install_ffmpeg_arm64_v1.0.11() {
	wget https://static.oomol.com/sshexec/v1.0.11/caller-arm64 --output-document=/usr/bin/caller
	chmod +x /usr/bin/caller
	ln -sf /usr/bin/caller /usr/bin/ffmpeg
	ln -sf /usr/bin/caller /usr/bin/ffprobe
	ln -sf /usr/bin/caller /usr/bin/install_ffmpeg_6
	/usr/bin/install_ffmpeg_6
}

# macos arm64 sshexec(from stable version) support install ffmpeg into host and calling ffmpeg from host
install_ffmpeg_arm64_old() {
	install_ffmpeg_arm64_v1.0.11
}

setup-macos-arm64() {
	get_ver
	if [[ "$sshexec_version" == "v1.0.11"* ]]; then
		install_ffmpeg_arm64_v1.0.11
	else
		# macos arm64 sshexec(from stable version) support install ffmpeg into host and calling ffmpeg from host
		install_ffmpeg_arm64_old
	fi
}

############ END OF MACOS ARM64 ############

############ START OF MACOS AMD64 ############
install_ffmpeg_amd64_v1.0.11() {
	wget https://static.oomol.com/sshexec/v1.0.11/caller-amd64 --output-document=/usr/bin/caller
	chmod +x /usr/bin/caller
	ln -sf /usr/bin/caller /usr/bin/ffmpeg
	ln -sf /usr/bin/caller /usr/bin/ffprobe
	ln -sf /usr/bin/caller /usr/bin/install_ffmpeg_6
	/usr/bin/install_ffmpeg_6
}

install_ffmpeg_amd64_old() {
	sudo apt update
	sudo apt install -y ffmpeg
}

setup-macos-amd64() {
	get_ver
	if [[ "$sshexec_version" == "v1.0.11"* ]]; then
		install_ffmpeg_amd64_v1.0.11
	else
		# from sshexec v1.0.11 support install ffmpeg into host and calling ffmpeg from host
		# For compatibility reasons, use apt to install ffmpeg as a fallback mechanism
		install_ffmpeg_amd64_old
	fi
}

############ END OF MACOS AMD64 ############

############ START OF WSL2 AMD64 ############
setup-wsl-amd64() {
	echo "Install ffmpeg"
	_ver=v1.0.11
	ffmpeg_tar="/tmp/jellyfin-ffmpeg_6.0.1-8_portable_linux64-gpl.tar.xz"
	wget https://static.oomol.com/sshexec/$_ver/jellyfin-ffmpeg_6.0.1-8_portable_linux64-gpl.tar.xz --output-document="$ffmpeg_tar"
	tar -xvf "$ffmpeg_tar" -C /usr/bin/
	chmod +x /usr/bin/ffmpeg
	chmod +x /usr/bin/ffprobe
	echo "Install ffmpeg done"
}

############ END OF WSL2 AMD64 ############

setup_ffmpeg() {
	if [[ "$platform" == macos-aarch64 ]]; then
		setup-macos-arm64
	elif [[ "$platform" == wsl2-x86_64 ]]; then
		setup-wsl-amd64
	elif [[ "$platform" == macos-x86_64 ]]; then
		setup-macos-amd64
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
