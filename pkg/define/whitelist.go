package define

var Whitelist = []string{
	FFMPEG,
	FFPROBE,
	InstallFFMPEG,
}

func IsWhitelisted(command string) bool {
	for _, item := range Whitelist {
		if item == command {
			return true
		}
	}
	return false
}
