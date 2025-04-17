package main

import (
	"embed"
	"github.com/alpha-omega-corp/cloud/api/pkg/user"
	"github.com/alpha-omega-corp/cloud/core"
	"github.com/alpha-omega-corp/cloud/core/config"
	"github.com/uptrace/bunrouter"
	"log"
)

var (
	//go:embed config
	embedFS embed.FS
)

func main() {
	core.NewApp(embedFS, "gateway").
		CreateApi(func(router *bunrouter.Router, configHandler config.Handler) {
			configUser, err := configHandler.GetConfig("user")
			if err != nil {
				log.Fatal(err.Error())
			}

			svcUser := user.NewClient(configUser)
			user.RegisterClient(svcUser, router)
		})
}
