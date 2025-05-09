package define

type Version struct {
	PkgName string
	PkgVer  string
}

type RunOpts struct {
	FFMPEGPath  string
	FFPROBEPath string
	Args        []string
	Envs        []string
}

type InstallOpts struct {
	URL       string
	Prefix    string
	Sha256Sum string
	TarBar    string
}
