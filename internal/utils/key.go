package utils

import (
	"encoding/base64"

	"github.com/charmbracelet/ssh"
)

func StringifyPublicKey(key ssh.PublicKey) string {
	return key.Type() + " " + base64.StdEncoding.EncodeToString(key.Marshal())
}
