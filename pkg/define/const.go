package define

const (
	FFMPEG                 = "ffmpeg"
	FFPROBE                = "ffprobe"
	SHELL                  = "bash"
	InstallFFMPEG          = "install_ffmpeg"
	StudioDir              = ".oomol-studio"
	FFReleaseURLForVentura = "https://github.com/oomol/builded/releases/download/v1.4/ffmpeg_macos_arm64_ventura.tar.xz"
	FFReleaseURLForSequoia = "https://github.com/oomol/builded/releases/download/v1.7/ffmpeg_macos_arm64_sequoia.tar.xz"
	FFMSha256ForVentura    = "4fb4effeb4e4e19ec1ca9f971d01c2c7fb682edabe238df38960caf9e3c421bd"
	FFMSha256ForSequoia    = "ea0bde7f9959fe3dd21fc1bd48e41a6ab88a8074dcff36d1f541d9dbe561fde2"
)

const (
	DYLD_LIBRARY_PATH = "DYLD_LIBRARY_PATH"
)

var (
	Addr = "127.0.0.1"
	Port = "5321"
)
