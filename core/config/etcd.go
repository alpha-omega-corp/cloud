package config

import (
	"bytes"
	"context"
	"fmt"
	"github.com/alpha-omega-corp/cloud/core/types"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
	clientv3 "go.etcd.io/etcd/client/v3"
	"time"
)

type Handler interface {
	Init(initialConfig []byte, err error) Handler
	LoadAs(ctx context.Context, name string) (config *types.Config)
	GetConfig(name string) (config types.Config, err error)
	Read(key string, format string) (err error)
}

type handler struct {
	Handler

	viper         *viper.Viper
	etcd          *clientv3.Client
	initialConfig []byte
	host          string
}

func NewHandler(file []byte) Handler {
	v := viper.New()

	v.SetConfigType("yaml")
	err := v.ReadConfig(bytes.NewBuffer(file))
	if err != nil {
		panic(err)
	}

	host := v.GetString("kvs")

	fmt.Printf("config host > %s\n", host)

	config := clientv3.Config{
		Endpoints:   []string{host},
		DialTimeout: 5 * time.Second,
	}

	etcd, err := clientv3.New(config)
	if err != nil {
		panic(err)
	}

	return &handler{
		host:          host,
		etcd:          etcd,
		initialConfig: file,
		viper:         viper.New(),
	}
}

func (m *handler) LoadAs(ctx context.Context, name string) (config *types.Config) {
	_, err := m.etcd.Put(ctx, "config_"+name, string(m.initialConfig))
	if err != nil {
		panic(err)
	}

	cfg, err := m.GetConfig(name)
	if err != nil {
		panic(err)
	}

	return &cfg
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

func (m *handler) Read(key string, format string) (err error) {
	err = m.viper.AddRemoteProvider("etcd3", "http://"+m.host, key)
	if err != nil {
		return
	}

	m.viper.SetConfigType(format)
	err = m.viper.ReadRemoteConfig()

	return
}
