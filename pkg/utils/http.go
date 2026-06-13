package utils

func GetConfigName(env string) string {
	if env == "docker" {
		return "config-docker"
	}
	return "config-local"
}
