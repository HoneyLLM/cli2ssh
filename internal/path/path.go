package path

import (
	"os"
	"path"
)

func GetDefaultHostKeyPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return ""
	}

	dataDirPath := path.Join(homeDir, ".cli2ssh")
	err = os.MkdirAll(dataDirPath, 0700)
	if err != nil {
		return ""
	}

	return path.Join(dataDirPath, "id_ed25519")
}
