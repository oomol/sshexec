package exec

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

const MyJSONData = `
{
  "mountPoints": [
    {
      "hostPath": "/Users/localuser/Downloads",
      "containerPath": "/oomol-driver/downloads",
      "containerDirName": "downloads"
    },
    {
      "hostPath": "/Users/localuser/Desktop",
      "containerPath": "/oomol-driver/desktop",
      "containerDirName": "desktop"
    }
  ],
  "currentMountPoint": null
}
`

func TestExecPathCover(t *testing.T) {
	p := filepath.Join("/tmp", "mount-point.json")
	jsonFile, err := os.Create(p)
	if err != nil {
		logrus.Fatalf("Failed to create json file: %v", err)
	}
	_, _ = jsonFile.WriteString(MyJSONData)
	MyJSONFile = p

	homeDir, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf("Failed to get user home dir: %v", err)
	}
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

func TestExecPathCover2(t *testing.T) {
	p := filepath.Join("/tmp", "mount-point.json")
	jsonFile, err := os.Create(p)
	if err != nil {
		logrus.Fatalf("Failed to create json file: %v", err)
	}
	_, _ = jsonFile.WriteString(MyJSONData)
	MyJSONFile = p

	testArgsArray := []string{
		"test",
		"/oomol-driver/oomol-storage",
		"/oomol-driver/sessions",
		"/oomol-driver/desktop",
		"/oomol-driver/downloads",
	}

	sanitizers, err := DoArgsSanitizers(testArgsArray)
	if err != nil {
		return
	}
	t.Log(sanitizers)
}
