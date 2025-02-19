//  SPDX-FileCopyrightText: 2024-2025 OOMOL, Inc. <https://www.oomol.com>
//  SPDX-License-Identifier: MPL-2.0

package exec

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
)

type MountPoint struct {
	HostPath         string `json:"hostPath"`
	ContainerPath    string `json:"containerPath"`
	ContainerDirName string `json:"containerDirName"`
	Undeletable      bool   `json:"undeletable,omitempty"`
}

var MyJSONFile string

type DataStruct struct {
	MountPoints       []MountPoint `json:"mountPoints"`
	CurrentMountPoint MountPoint   `json:"currentMountPoint"`
}

const (
	oomolStorage  = "/oomol-driver/oomol-storage"
	oomolSessions = "/oomol-driver/sessions"
)

// loadJSON loads the mount-point.json file from the given path using json.Unmarshal
func loadJSON(path string) (*DataStruct, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var data DataStruct
	decoder := json.NewDecoder(file)
	if err = decoder.Decode(&data); err != nil {
		return nil, err
	}
	return &data, nil
}

const (
	ooHomePrefix    = ".oomol-studio"
	ooAppConfig     = "app-config"
	ooStorage       = "oomol-storage"
	ooMountJSONFile = "mount-point.json"

	ooSessions = "sessions"
)

func ContainerPath2HostPath(arg string) (string, error) {
	// MacOS Host home directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}

	if MyJSONFile == "" {
		MyJSONFile = filepath.Join(homeDir, ooHomePrefix, ooAppConfig, ooStorage, ooMountJSONFile)
	}

	logrus.Infof("Load MountPoint json file: %q", MyJSONFile)
	jsonData, err := loadJSON(MyJSONFile)
	if err != nil {
		return "", fmt.Errorf("failed to load json file: %v", err)
	}

	logrus.Infof("Process string: %q", arg)

	// /oomol-driver/sessions --> $HOME/.oomol-studio/sessions
	if strings.Contains(arg, oomolSessions) {
		path := filepath.Join(homeDir, ooHomePrefix, ooSessions)
		newArg := strings.Replace(arg, oomolSessions, path, 1)
		logrus.Warnf("%q --> %q", arg, newArg)
		return newArg, nil
	}

	// /oomol-driver/ oomol-storage --> $HOME/oomol-storage
	if strings.Contains(arg, oomolStorage) {
		path := filepath.Join(homeDir, ooStorage)
		newArg := strings.Replace(arg, oomolStorage, path, 1)
		logrus.Warnf("%q --> %q", arg, newArg)
		return newArg, nil
	}

	// Other binding directories
	if jsonData != nil {
		for _, mountPoint := range jsonData.MountPoints {
			// Check if the argument starts with the mountPoint.ContainerPath the replace
			if strings.Contains(arg, mountPoint.ContainerPath) {
				newArg := strings.Replace(arg, mountPoint.ContainerPath, mountPoint.HostPath, 1)
				logrus.Warnf("%q --> %q", arg, newArg)
				return newArg, nil
			}
		}
	}
	return arg, nil
}

func DoArgsSanitizers(args []string) ([]string, error) {
	newArgs := make([]string, 0)
	for _, arg := range args {
		singleArg, err := ContainerPath2HostPath(arg)
		if err != nil {
			return nil, fmt.Errorf("failed to convert container path to host path: %v", err)
		}
		newArgs = append(newArgs, singleArg)
	}
	return newArgs, nil
}
