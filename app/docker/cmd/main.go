package main

import (
	"embed"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/models"
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/proto"
	"github.com/alpha-omega-corp/cloud/core"
	"github.com/alpha-omega-corp/cloud/core/types"
	"github.com/docker/docker/client"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
)

var (
	//go:embed config
	embedFS embed.FS
)

func main() {
	core.NewApp(embedFS, "docker").
		CreateApp(func(config *types.Config, db *bun.DB, grpc *grpc.Server) {
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

			proto.RegisterDockerServiceServer(grpc, pkg.NewServer(config, dockerClient, db))
		}, []interface{}{
			(*models.Dockerfile)(nil),
		}...)
}
