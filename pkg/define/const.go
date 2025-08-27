package define

const (
	// CurrentVersion Major version
	CurrentVersion = "v1.0.12"

	FFMPEG6ForARM64PkgName    = "ffmpeg_macos_arm64"
	FFMPEG6ForARM64Version    = "v6.0.1"
	FFMPEG6ForARM64ReleaseURL = "https://static.oomol.com/sshexec/v1.0.11/jellyfin-ffmpeg_6.0.1-8_portable_macarm64-gpl.tar.xz"
	FFMPEG6ForARM64Sha256     = "efce8779d5f35127ec7dcd500669e00a2e13f2b099c7bf11ab5acee467dc5d57"

	FFMPEG6ForAMD64PkgName    = "ffmpeg_macos_amd64"
	FFMPEG6ForAMD64Version    = "v6.0.1"
	FFMPEG6ForAMD64ReleaseURL = "https://static.oomol.com/sshexec/v1.0.11/jellyfin-ffmpeg_6.0.1-8_portable_mac64-gpl.tar.xz"
	FFMPEG6ForAMD64Sha256     = "aa0351fff7a682d2da85ccfa5ab7d9bc27121ac778e7845c2d95fe1b1d752655"
)

const (
	InstallFFMPEGVersion6 = "install_ffmpeg_6"
	ShowCurrentVersion    = "show_version"
	StudioDir             = ".oomol-studio"
	HostShared            = "host-shared"

	FFMPEGBin  = "ffmpeg"
	FFPROBEBin = "ffprobe"
)

var (
	Addr = "127.0.0.1"
	Port = "5322"
)
