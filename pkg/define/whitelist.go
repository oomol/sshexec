package define

var Whitelist = []string{
	FFPROBEBin,
	FFMPEGBin,
	InstallFFMPEGVersion6,
}

func IsWhitelisted(command string) bool {
	for _, item := range Whitelist {
		if item == command {
			return true
		}
	}
	return false
}
