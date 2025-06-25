package main

import (
	"embed"
	"github.com/alpha-omega-corp/cloud/app/github/pkg"
	"github.com/alpha-omega-corp/cloud/app/github/pkg/proto"
	"github.com/alpha-omega-corp/cloud/app/github/pkg/utils"
	"github.com/alpha-omega-corp/cloud/core"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
)

var (
	//go:embed config
	embedFS embed.FS
)

func main() {
	core.NewApp(embedFS, "github").
		CreateApp(func(config *core.Config, db *bun.DB, grpc *grpc.Server) {
			client := utils.NewGithubApiClient(config)

			proto.RegisterGithubServiceServer(grpc, pkg.NewServer(client))
		}, []interface{}{
			// Empty
		}...)
}
