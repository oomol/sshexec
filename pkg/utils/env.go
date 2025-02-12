package utils

import (
	"errors"
	"os"
	"path/filepath"
	"sshd/pkg/define"
)

func GetStudioHomeDir() (string, error) {
	d, err := os.UserHomeDir()
	if err != nil {
		return "", errors.New("failed to get user home directory")
	}
	d = filepath.Join(d, define.StudioDir)
	return d, nil
}
