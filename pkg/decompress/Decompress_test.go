package decompress

import (
	"context"
	"testing"
)

func TestUnpack(t *testing.T) {
	err := Extract(context.Background(), "/tmp/ffmpeg_macos_arm64_ventura.tar.xz", "/tmp/ffmpeg_macos_arm64_ventura")
	if err != nil {
		t.Errorf("Decompress() failed: %v", err)
	}
}
