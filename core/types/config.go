package types

import "github.com/spf13/viper"

type ConfigHost struct {
	Name string `mapstruct:"name"`
	Url  string `mapstruct:"url"`
	Repo string `mapstruct:"repo"`
	Dsn  string `mapstruct:"dsn"`
}

type Config struct {
	Viper *viper.Viper
}
