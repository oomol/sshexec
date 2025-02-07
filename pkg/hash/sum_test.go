package hash

import "testing"

func TestCheckSum(t *testing.T) {
	sum, err := sha256sum("/tmp/ffmpeg_macos_arm64_ventura.tar.xz")
	if err != nil {
		t.Errorf("sha256sum() failed: %v", err)
	}
	t.Log(sum)
}
