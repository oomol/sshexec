package define

const (
	ARCH = "arm64"

	FFMPEGPkgName = "ffmpeg_macos_arm64"

	// FFMPEG VERSION v7.1.1 info
	FFMPEG7Version = "v7.1.1"
	FF7ReleaseURL  = "https://github.com/jellyfin/jellyfin-ffmpeg/releases/download/v7.1.1-2/jellyfin-ffmpeg_7.1.1-2_portable_macarm64-gpl.tar.xz"
	FF7Sha256      = "8ff4ed4eba151346d6d1ee894bbfcbd9f507dc82b34260f4567954136449001d"

	// 	// FFMPEG VERSION v6.0.1 info
	FFMPEG6Version = "v6.0.1"
	FF6ReleaseURL  = "https://github.com/jellyfin/jellyfin-ffmpeg/releases/download/v6.0.1-8/jellyfin-ffmpeg_6.0.1-8_portable_macarm64-gpl.tar.xz"
	FF6Sha256      = "efce8779d5f35127ec7dcd500669e00a2e13f2b099c7bf11ab5acee467dc5d57"
)

const (
	InstallFFMPEGVersion7 = "install_ffmpeg_7"
	InstallFFMPEGVersion6 = "install_ffmpeg_6"
	StudioDir             = ".oomol-studio"
	HostShared            = "host-shared"

	FFMPEGBin  = "ffmpeg"
	FFPROBEBin = "ffprobe"
)

var (
	Addr = "127.0.0.1"
	Port = "5322"
)
