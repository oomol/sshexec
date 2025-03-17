package define

const (
	FFMPEG                 = "ffmpeg"
	FFPROBE                = "ffprobe"
	InstallFFMPEG          = "install_ffmpeg"
	StudioDir              = ".oomol-studio"
	FFReleaseURLForVentura = "https://github.com/oomol/builded/releases/download/v1.7/ffmpeg_macos_arm64_ventura.tar.xz"
	FFReleaseURLForSequoia = "https://github.com/oomol/builded/releases/download/v1.7/ffmpeg_macos_arm64_sequoia.tar.xz"
	FFMSha256ForVentura    = "79ea28ac973aae590761cff3dde0922f89ff153a8320e4e0a85befe02034cb7c"
	FFMSha256ForSequoia    = "ea0bde7f9959fe3dd21fc1bd48e41a6ab88a8074dcff36d1f541d9dbe561fde2"
)

const (
	DYLD_LIBRARY_PATH = "DYLD_LIBRARY_PATH" //nolint:stylecheck
)

var (
	Addr = "127.0.0.1"
	Port = "5322"
)
