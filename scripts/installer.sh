#! /usr/bin/env bash
set -o pipefail

RED="\033[31m"
YELLOW="\033[33m"
GREEN="\033[32m"
RESET="\033[0m"

log_err() {
	echo -e "${RED}[ERROR]${RESET} $*" >&2
	exit 100
}

log_warn() {
	echo -e "${YELLOW}[WARN]${RESET} $*" >&2
}

log_std() {
	echo -e "${GREEN}[INFO]${RESET} $*"
}

sshexec_version="unknown"

get_platform() {
	arch=$(uname -m)
	platform=unknown
	log_std "get sshexec version..."
	sshexec_version="$(ssh -o ConnectTimeout=5 \
	                      -o ServerAliveInterval=2 \
	                      -o ServerAliveCountMax=2 \
	                      -o StrictHostKeyChecking=no \
	                      -q root@192.168.127.254 -p5322 show_version | xargs | tr -d '\r' | tr -d '\n')"
	log_std "sshexec_version: $sshexec_version"



	if [[ -z "$arch" ]]; then
		log_err "uname -m return empty"
	fi

	case "$arch" in
      aarch64|arm64)
          platform="linux_arm64"
          [[ "$OO_HOST_PLATFORM" == "darwin" ]] && platform="macos_arm64"
          ;;
      x86_64|amd64)
          platform="linux_amd64"
          [[ "$OO_HOST_PLATFORM" == "win32" ]] && platform="wsl2_amd64"
          [[ "$OO_HOST_PLATFORM" == "darwin" ]] && platform="macos_amd64"
          ;;
      *)
          log_err "unsupported architecture: $arch"
          ;;
  esac
}

# Fallback to install native ffmpeg
install_native_ffmpeg_linux() {
	if ffmpeg -version; then
	  log_std "ffmpeg installed before"
	  return
	fi

	if [[ "$platform" == "linux_arm64" ]] || [[ "$platform" == "macos_arm64" ]] || [[ "$platform" == "wsl2_arm64" ]]; then
		log_std "Install ffmpeg for linux-arm64"
		local url="https://static.oomol.com/sshexec/v1.0.11/jellyfin-ffmpeg_6.0.1-8_portable_linuxarm64-gpl.tar.xz"
		local ffmpeg_tar="/tmp/$(basename "$url")"

		log_std "download ffmpeg"
		wget "$url" --output-document="$ffmpeg_tar"

		log_std "extract ffmpeg bins to /usr/bin"
		tar -xvf "$ffmpeg_tar" -C /usr/bin/

		log_std "chmod +x /usr/bin/{ffmpeg,ffprobe}"
		chmod +x /usr/bin/ffmpeg
		chmod +x /usr/bin/ffprobe

		log_std "Install ffmpeg for linux-arm64 done"
	elif [[ "$platform" == "linux_amd64" ]] || [[ "$platform" == "macos_amd64" ]] || [[ "$platform" == "wsl2_amd64" ]]; then
		log_std "Install ffmpeg for linux-amd64"

		local url="https://static.oomol.com/sshexec/v1.0.11/jellyfin-ffmpeg_6.0.1-8_portable_linux64-gpl.tar.xz"
		local ffmpeg_tar="/tmp/$(basename "$url")"

		log_std "download ffmpeg"
		wget "$url" --output-document="$ffmpeg_tar"

		log_std "extract ffmpeg bins to /usr/bin"
		tar -xvf "$ffmpeg_tar" -C /usr/bin/

		log_std "chmod +x /usr/bin/{ffmpeg,ffprobe}"
		chmod +x /usr/bin/ffmpeg
		chmod +x /usr/bin/ffprobe

		/usr/bin/ffmpeg -version
		
		log_std "Install ffmpeg for linux-amd64 done"
	else
		log_err "platform not support"
	fi
}

# Install the Linux version of ffmpeg in compatibility mode
setup_macos_compat() {
	log_std "Install the Linux version of ffmpeg in compatibility mode"
	install_native_ffmpeg_linux
}

# Use caller to install ffmpeg on the macOS host, and use caller to call ffmpeg on the macOS host from the container
setup_macos_host_v1dot0() {
	if [[ "$platform" == "macos_arm64" ]]; then
		caller_name=caller-arm64
	elif [[ "$platform" == "macos_amd64" ]]; then
		caller_name=caller-amd64
	fi

	log_std "Download caller version: $sshexec_version"
	wget "https://static.oomol.com/sshexec/$sshexec_version/$caller_name" --output-document "/usr/bin/caller"
	chmod +x /usr/bin/caller

	log_std "create soft-link to /usr/bin/{ffmpeg,ffprobe,install_ffmpeg_6}"
	ln -sf /usr/bin/caller /usr/bin/ffmpeg
	ln -sf /usr/bin/caller /usr/bin/ffprobe
	ln -sf /usr/bin/caller /usr/bin/install_ffmpeg_6

	log_std "call /usr/bin/install_ffmpeg_6 to install ffmpeg for host"
	/usr/bin/install_ffmpeg_6
}

setup_wsl() {
	install_native_ffmpeg_linux
}

setup_macos() {
	if [[ $sshexec_version == "v1.0"* ]]; then
		log_std "setup_macos_host_v1dot0"
		setup_macos_host_v1dot0
	else
		setup_macos_compat
	fi
}

setup_linux() {
	install_native_ffmpeg_linux
}

setup() {
	if [[ "$platform" == "macos_"* ]]; then
		setup_macos
	elif [[ "$platform" == "wsl2_"* ]]; then
		setup_wsl
	elif [[ "$platform" == "linux_"* ]]; then
		setup_linux
	else
		log_err "unsupport platform: $platform"
	fi
}

main() {
	get_platform
	set -e
  set -o pipefail
	setup
	set +e
  set +o pipefail
}

main
