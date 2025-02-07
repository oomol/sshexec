package decompress

import (
	"context"
	"testing"
)

func TestDecompress(t *testing.T) {
	err := Decompress(context.Background(), "/tmp/ffmpeg_macos_arm64_ventura.tar.xz", "/tmp/ffmpeg_macos_arm64_ventura.tar")
	if err != nil {
		t.Errorf("Decompress() failed: %v", err)
	}
}

func TestExtract(t *testing.T) {
	err := Extract(context.Background(), "/tmp/ffmpeg_macos_arm64_ventura.tar", "/tmp/ffmpeg_macos_arm64_ventura")
	if err != nil {
		t.Errorf("Extract() failed: %v", err)
	}
}
