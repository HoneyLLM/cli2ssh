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
	Host        = getEnvWithDefault("HOST", "localhost")
	Port        = getEnvWithDefault("PORT", "2222")
	HostKeyPath = getEnvWithDefault("HOST_KEY_PATH", "")
)
