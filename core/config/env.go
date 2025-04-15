package config

var Environment = map[string]string{
	"local":  "local",
	"docker": "docker",
}

func GetConfigPath(env string) string {
	path := "config/config." + Environment[env] + ".yaml"

	return path
}
