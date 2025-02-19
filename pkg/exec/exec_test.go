package exec

import (
	"os"
	"path/filepath"
	"testing"
)

func TestExecPathCover(t *testing.T) {
	MyJsonFile = "/tmp/mount-point.json"
	homeDir, err := os.UserHomeDir()

	path, err := ContainerPath2HostPath("test")
	if err != nil {
		t.Fatalf("ContainerPath2HostPath failed: %v", err)
	}
	if path != "test" {
		t.Error("ContainerPath2HostPath failed")
	}

	path, err = ContainerPath2HostPath("/oomol-driver/oomol-storage")
	if err != nil {
		t.Errorf("ContainerPath2HostPath failed: %v", err)
	}
	if path != filepath.Join(homeDir, ooStorage) {
		t.Error("ContainerPath2HostPath failed")
	}

	path, err = ContainerPath2HostPath("/oomol-driver/sessions")
	if err != nil {
		t.Errorf("ContainerPath2HostPath failed: %v", err)
	}
	if path != filepath.Join(homeDir, ooHomePrefix, ooSessions) {
		t.Error("ContainerPath2HostPath failed")
	}

	path, err = ContainerPath2HostPath("/oomol-driver/desktop")
	if err != nil {
		t.Errorf("ContainerPath2HostPath failed: %v", err)
	}
	if path != "/Users/localuser/Desktop" {
		t.Error("ContainerPath2HostPath failed")
	}

	path, err = ContainerPath2HostPath("/oomol-driver/downloads")
	if err != nil {
		t.Errorf("ContainerPath2HostPath failed: %v", err)
	}
	if path != "/Users/localuser/Downloads" {
		t.Error("ContainerPath2HostPath failed")
	}

}
