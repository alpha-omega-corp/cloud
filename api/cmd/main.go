package main

import (
	"embed"
	"github.com/alpha-omega-corp/cloud/api/clients/docker"
	"github.com/alpha-omega-corp/cloud/api/clients/github"
	"github.com/alpha-omega-corp/cloud/api/clients/user"
	"github.com/alpha-omega-corp/cloud/api/middlewares"
	"github.com/alpha-omega-corp/cloud/core"
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
	core.NewApp(embedFS, "api").
		CreateApi(func(configHandler *core.ConfigHandler, router *bunrouter.Router) {

			// Router middlewares
			router.Use(bunrouterotel.NewMiddleware())
			router.Use(middlewares.NewCorsMiddleware())
			router.Use(middlewares.NewErrorHandler)

			// Register clients
			configUser, err := configHandler.WithConfig("user")
			if err != nil {
				log.Fatal(err.Error())
			}
			userClient := user.RegisterClient(user.NewClient(configUser), router)
			router.Use(middlewares.NewAuthMiddleware(userClient).Auth)

			configDocker, err := configHandler.WithConfig("docker")
			if err != nil {
				log.Fatal(err.Error())
			}
			docker.RegisterClient(docker.NewClient(configDocker), router)

			configGithub, err := configHandler.WithConfig("github")
			if err != nil {
				log.Fatal(err.Error())
			}
			github.RegisterClient(github.NewClient(configGithub), router)

			env := *configHandler.GetConfig().Env

			uptrace.ConfigureOpentelemetry(
				uptrace.WithDSN(env.GetString("uptrace_dsn")),
				uptrace.WithServiceName(env.GetString("uptrace_name")),
				uptrace.WithServiceVersion(env.GetString("uptrace_version")),
				uptrace.WithDeploymentEnvironment(env.GetString("uptrace_env")),
			)
		})
}
