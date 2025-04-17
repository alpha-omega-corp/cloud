package main

import (
	"embed"
	"github.com/alpha-omega-corp/cloud/app/user/pkg"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/models"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/proto"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/utils"
	"github.com/alpha-omega-corp/cloud/core"
	"github.com/alpha-omega-corp/cloud/core/types"
	_ "github.com/spf13/viper/remote"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
)

var (
	//go:embed config
	embedFS embed.FS
)

func main() {
	core.NewApp(embedFS, "user").
		CreateServer(func(config *types.Config, db *bun.DB, grpc *grpc.Server) {
			auth := utils.NewAuthWrapper(config.Env.GetString("secret"))
			proto.RegisterUserServiceServer(grpc, pkg.NewServer(db, auth))
		}, []interface{}{
			(*models.UserToRole)(nil),
			(*models.User)(nil),
			(*models.Role)(nil),
			(*models.Service)(nil),
			(*models.Permission)(nil),
		}...)
}
