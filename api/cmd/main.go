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
	"github.com/uptrace/uptrace-go/uptrace"
	"log"
)

var (
	//go:embed config
	embedFS embed.FS
)

func main() {
	core.NewApp(embedFS, "gateway").
		CreateApi(func(router *bunrouter.Router, configHandler *config.Handler) {

			configApi, err := configHandler.GetConfig("api")
			if err != nil {
				log.Fatal(err)
			}

			env := *configApi.Env

			uptrace.ConfigureOpentelemetry(
				uptrace.WithDSN(env.GetString("uptrace_dsn")),
				uptrace.WithServiceName(env.GetString("uptrace_name")),
				uptrace.WithServiceVersion(env.GetString("uptrace_version")),
				uptrace.WithDeploymentEnvironment(env.GetString("uptrace_env")),
			)

			// Router middlewares
			router.Use(bunrouterotel.NewMiddleware())
			router.Use(middlewares.NewCorsMiddleware())
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
