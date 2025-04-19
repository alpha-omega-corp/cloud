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

type Handler struct {
	viper         *viper.Viper
	etcd          *clientv3.Client
	initialConfig []byte
	host          string
}

func NewHandler(file []byte) *Handler {
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

	return &Handler{
		host:          host,
		etcd:          etcd,
		initialConfig: file,
		viper:         viper.New(),
	}
}

func (h *Handler) LoadAs(ctx context.Context, name string) (config *types.Config) {
	_, err := h.etcd.Put(ctx, "config_"+name, string(h.initialConfig))
	if err != nil {
		panic(err)
	}

	cfg, err := h.GetConfig(name)
	if err != nil {
		panic(err)
	}

	return cfg
}

func (h *Handler) GetConfig(name string) (config *types.Config, err error) {
	var cfg *types.Config

	err = h.Read("config_"+name, "yaml")
	if err != nil {
		return nil, err
	}

	err = h.viper.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	cfg.Env = h.viper
	config = cfg

	return config, nil
}

func (h *Handler) Read(key string, format string) (err error) {
	h.viper = viper.New()

	err = h.viper.AddRemoteProvider("etcd3", "http://"+h.host, key)
	if err != nil {
		return
	}

	h.viper.SetConfigType(format)
	err = h.viper.ReadRemoteConfig()

	return
}
