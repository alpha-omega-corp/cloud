package config

import (
	"flag"
	"fmt"
)

var Environment = map[string]string{
	"local":  "local",
	"docker": "docker",
}

func GetConfigPath() string {
	cmd := flag.String("env", "local", "")
	flag.Parse()
	fmt.Printf("enivonrment : \"%v\"\n", *cmd)

	path := "config/config." + Environment[*cmd] + ".yaml"

	return path
}
