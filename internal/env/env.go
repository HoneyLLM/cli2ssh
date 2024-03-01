package env

import "os"

func getEnvWithDefault(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

var (
	Host        = getEnvWithDefault("CLI2SSH_HOST", "localhost")
	Port        = getEnvWithDefault("CLI2SSH_HOST", "2222")
	HostKeyPath = getEnvWithDefault("CLI2SSH_HOST_KEY_PATH", "")
)
