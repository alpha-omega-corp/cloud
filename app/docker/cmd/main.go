package main

import (
	"context"
	"embed"
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

var (
	//go:embed config
	embedFS embed.FS
)

func main() {
	configHandler := config.NewHandler(embedFS.ReadFile("config/config.yaml"))
	cfg, err := configHandler.LoadAs(context.Background(), "docker")
	if err != nil {
		panic(err)
	}

	dbHandler := database.NewHandler(cfg.Dsn)
	dbHandler.Database().RegisterModel(
		(*models.Dockerfile)(nil),
	)

	if err := srv.NewGRPC(cfg.Url, dbHandler, func(db *bun.DB, grpc *grpc.Server) {
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

		proto.RegisterDockerServiceServer(grpc, pkg.NewServer(cfg, dockerClient, db))
	}); err != nil {
		panic(err)
	}
}
