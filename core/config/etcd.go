package config

import (
	"context"
	"github.com/alpha-omega-corp/cloud/core/types"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Handler interface {
	LoadAs(ctx context.Context, name string) (config types.Config, err error)
	Read(key string, format string) (err error)
	GetConfig(name string) (config types.Config, err error)
}

type handler struct {
	Handler

	viper         *viper.Viper
	etcd          *clientv3.Client
	initialConfig []byte
	host          string
}

func NewHandler(initialConfig []byte, err error) Handler {
	if err != nil {
		panic(err)
	}

	host := GetEnv() + ":2380"
	config := clientv3.Config{
		Endpoints:   []string{host},
		DialTimeout: 5 * time.Second,
	}

	etcd, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}

	return &handler{
		etcd:          etcd,
		viper:         viper.New(),
		initialConfig: initialConfig,
		host:          host,
	}
}

func (m *handler) LoadAs(ctx context.Context, name string) (config types.Config, err error) {
	_, err = m.etcd.Put(ctx, "config_"+name, string(m.initialConfig))
	environment, err := m.GetConfig(name)
	if err != nil {
		return
	}

	return environment, nil
}

func (m *handler) Read(key string, format string) (err error) {
	err = m.viper.AddRemoteProvider("etcd3", "http://"+m.host, key)
	if err != nil {
		return
	}

	m.viper.SetConfigType(format)
	err = m.viper.ReadRemoteConfig()

	return
}

func (m *handler) GetConfig(name string) (config types.Config, err error) {
	var cfg types.Config

	err = m.Read("config_"+name, "yaml")
	if err != nil {
		return
	}

	err = m.viper.Unmarshal(&cfg)
	if err != nil {
		return
	}

	cfg.Env = m.viper
	config = cfg

	return
}
