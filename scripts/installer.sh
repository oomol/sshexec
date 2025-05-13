#! /usr/bin/env bash
set -e
set -u

log() {
  msg="LOG: $S_NAME> $*"
  echo "$(date +'%Y-%m-%dT%H:%M:%S.%3N%z'):" "$msg" >&2
}

err() {
  msg="ERROR: $S_NAME> $*"
  echo "$(date +'%Y-%m-%dT%H:%M:%S.%3N%z'):" "$msg" >&2
  exit 100
}

warn() {
  msg="WARN: $S_NAME> $*"
  echo "$(date +'%Y-%m-%dT%H:%M:%S.%3N%z'):" "$msg" >&2
}

get_platform() {
  arch=$(uname -m)
  platform=unknown

  if [[ -z "$arch" ]]; then
    warn "uname -m return empty"
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

# Setup ffmpeg binaries logic
setup_ffmpeg_for_macos_aarch64() {
  wget https://github.com/oomol/sshexec/releases/download/v1.0.6/caller-arm64 --output-document=/usr/bin/caller
  ln -sf /usr/bin/caller /usr/bin/ffmpeg
  ln -sf /usr/bin/caller /usr/bin/ffprobe
}

setup_ffmpeg_for_wsl2_x86_64() {
  wget https://github.com/oomol/builded/releases/download/v1.7/ffmpeg-wsl2_x86_64.tar.xz --output-document=/tmp/ffmpeg-wsl2_x86_64.tar.xz
  tar -xvf /tmp/ffmpeg-wsl2_x86_64.tar.xz -C /tmp/
  echo "Install ffmpeg"
  cp /tmp/ffmpeg/ffmpeg /usr/bin/
  cp /tmp/ffmpeg/ffprobe /usr/bin/
  echo "Install ffmpeg done"
}

setup_ffmpeg() {
  if [[ "$platform" == macos-aarch64 ]]; then
    setup_ffmpeg_for_macos_aarch64
  elif [[ "$platform" == wsl2-x86_64 ]]; then
    setup_ffmpeg_for_wsl2_x86_64
  else
    err "unsupport platform: $platform"
  fi
}

main() {
  get_platform
  if [[ "$platform" == "unknown" ]]; then
    err "unknown platform"
  fi
  setup_ffmpeg
}

main
