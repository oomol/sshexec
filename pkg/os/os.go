package os

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"

	"github.com/sirupsen/logrus"
)

func getMacOSVersion() (string, error) {
	var out bytes.Buffer
	cmd := exec.Command("sw_vers", "-productVersion")
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run sw_vers: %v", err)
	}
	return strings.TrimSpace(out.String()), nil
}

func IsSequoia() bool {
	version, err := getMacOSVersion()
	if err != nil {
		logrus.Errorf("failed to get macOS version: %v", err)
		return false
	}
	return strings.HasPrefix(version, "15")
}

func IsSonoma() bool {
	version, err := getMacOSVersion()
	if err != nil {
		logrus.Errorf("failed to get macOS version: %v", err)
		return false
	}
	return strings.HasPrefix(version, "14")
}

func IsVentura() bool {
	version, err := getMacOSVersion()
	if err != nil {
		logrus.Errorf("failed to get macOS version: %v", err)
		return false
	}
	return strings.HasPrefix(version, "13")
}
