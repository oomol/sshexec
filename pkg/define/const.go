package define

const (
	FFMPEG                 = "ffmpeg"
	FFPROBE                = "ffprobe"
	InstallFFMPEG          = "install_ffmpeg"
	StudioDir              = ".oomol-studio"
	HostShared             = "host-shared"
	FFReleaseURLForVentura = "https://github.com/jellyfin/jellyfin-ffmpeg/releases/download/v7.1.1-2/jellyfin-ffmpeg_7.1.1-2_portable_macarm64-gpl.tar.xz"
	FFMSha256ForVentura    = "8ff4ed4eba151346d6d1ee894bbfcbd9f507dc82b34260f4567954136449001d"
)

const (
	DYLD_LIBRARY_PATH = "DYLD_LIBRARY_PATH" //nolint:stylecheck
)

var (
	Addr = "127.0.0.1"
	Port = "5322"
)
