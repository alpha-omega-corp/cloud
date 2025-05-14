package main

import (
	"embed"
	"github.com/alpha-omega-corp/cloud/api/clients/docker"
	"github.com/alpha-omega-corp/cloud/api/clients/github"
	"github.com/alpha-omega-corp/cloud/api/clients/user"
	"github.com/alpha-omega-corp/cloud/api/middlewares"
	"github.com/alpha-omega-corp/cloud/core"
	"github.com/alpha-omega-corp/cloud/core/config"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/bunrouterotel"
	"log"
)

var (
	//go:embed config
	embedFS embed.FS
)

func main() {
	core.NewApp(embedFS, "gateway").
		CreateApi(func(router *bunrouter.Router, configHandler *config.Handler) {

			// Router middlewares
			router.Use(middlewares.NewCorsMiddleware())
			router.Use(bunrouterotel.NewMiddleware())
			router.Use(middlewares.NewErrorHandler)

			// Create user service
			configUser, err := configHandler.GetConfig("user")
			if err != nil {
				log.Fatal(err.Error())
			}
			user.RegisterClient(user.NewClient(configUser), router)

			configDocker, err := configHandler.GetConfig("docker")
			if err != nil {
				log.Fatal(err.Error())
			}
			docker.RegisterClient(docker.NewClient(configDocker), router)

			configGithub, err := configHandler.GetConfig("github")
			if err != nil {
				log.Fatal(err.Error())
			}
			github.RegisterClient(github.NewClient(configGithub), router)
		})
}
