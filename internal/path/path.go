package path

import (
	"os"
	"path"
)

func GetDefaultHostKeyPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	dataDirPath := path.Join(homeDir, ".cli2ssh")
	err = os.MkdirAll(dataDirPath, 0700)
	if err != nil {
		return "", err
	}

	return path.Join(dataDirPath, "id_ed25519"), nil
}
