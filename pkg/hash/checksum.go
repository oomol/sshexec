package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func CmpFileChecksum(path string, wantSum string) error {
	sum, err := sha256sum(path)
	if err != nil {
		return fmt.Errorf("sha256sum(%s) failed: %v", path, err)
	}

	if sum != wantSum {
		logrus.Errorf("checksum mismatch: got %s, want %s", sum, wantSum)
		return errors.New("checksum mismatch")
	}
	logrus.Infoln("checksum matched")

	return nil
}

func sha256sum(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, f); err != nil {
		return "", err
	}

	// get the 32 bytes hash
	bytes := hash.Sum(nil)[:32]

	return hex.EncodeToString(bytes), nil
}
