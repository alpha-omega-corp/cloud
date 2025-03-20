package config

import (
	"github.com/alpha-omega-corp/cloud/core/types"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	"strings"
)

type Handler interface {
	Read(key string, format string) (err error)
	Environment(name string) (env *types.Environment, err error)
}

type handler struct {
	Handler

	viper *viper.Viper
	host  string
}

func NewHandler(host string) Handler {
	return &handler{
		viper: nil,
		host:  host,
	}
}

func (m *handler) Read(key string, format string) (err error) {
	m.handle()
	err = m.viper.AddRemoteProvider("etcd3", m.host, key)
	if err != nil {
		return
	}

	m.viper.SetConfigType(format)
	err = m.viper.ReadRemoteConfig()

	return
}

func (m *handler) Environment(name string) (env *types.Environment, err error) {
	var envVars types.EnvVars
	err = m.Read("env_"+name, "yaml")
	err = m.viper.Unmarshal(&envVars)
	if err != nil {
		return
	}

	var hostConfig types.ConfigHost
	err = m.Read(strings.ToLower(envVars.Host), "yaml")
	err = m.viper.Unmarshal(&hostConfig)
	if err != nil {
		return
	}

	if envVars.Config != "" {
		err = m.Read(strings.ToLower(envVars.Config), "env")
		if err != nil {
			return
		}

		env = &types.Environment{
			Vars: envVars,
			Host: hostConfig,
			Config: types.Config{
				Viper: m.viper,
			},
		}
	} else {
		env = &types.Environment{
			Vars: envVars,
			Host: hostConfig,
		}
	}

	return
}

func (m *handler) handle() {
	m.viper = viper.New()
	return
}
