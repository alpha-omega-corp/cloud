package main

import (
	"github.com/alpha-omega-corp/cloud/app/docker/pkg"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/models"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/proto"
	"github.com/alpha-omega-corp/cloud/core/config"
	"github.com/alpha-omega-corp/cloud/core/database"
	srv "github.com/alpha-omega-corp/cloud/core/server"

	"github.com/docker/docker/client"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
)

func main() {
	env, err := config.NewHandler().Environment("docker")
	if err != nil {
		panic(err)
	}

	dbHandler := database.NewHandler(env.Host.Dsn)
	dbHandler.Database().RegisterModel(
		(*models.Dockerfile)(nil),
	)

	if err := srv.NewGRPC(env.Host.Url, dbHandler, func(db *bun.DB, grpc *grpc.Server) {
		dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
		if err != nil {
			panic(err)
		}

		defer func(cli *client.Client) {
			err := cli.Close()
			if err != nil {
				panic(err)
			}
		}(dockerClient)

		proto.RegisterDockerServiceServer(grpc, pkg.NewServer(env.Config, dockerClient, db))
	}); err != nil {
		panic(err)
	}
}
