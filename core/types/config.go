package types

import "github.com/spf13/viper"

type Config struct {
	Url *string `mapstruct:"url"`
	Dsn *string `mapstruct:"dsn"`
	Env *viper.Viper
}
