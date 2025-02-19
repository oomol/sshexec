package hash

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
)

func CmpFileChecksum(path string, wantSum string) (string, error) {
	sum, err := sha256sum(path)
	if err != nil {
		return "", fmt.Errorf("sha256sum(%q) failed: %v", path, err)
	}

	if sum != wantSum {
		return sum, errors.New("checksum mismatch")
	}

	return sum, nil
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
