package main

import (
	"github.com/alpha-omega-corp/cloud/app/user/pkg"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/models"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/proto"
	"github.com/alpha-omega-corp/cloud/app/user/pkg/utils"
	"github.com/alpha-omega-corp/cloud/core/config"
	"github.com/alpha-omega-corp/cloud/core/database"
	srv "github.com/alpha-omega-corp/cloud/core/server"
	_ "github.com/spf13/viper/remote"
	"github.com/uptrace/bun"
	"google.golang.org/grpc"
)

func main() {
	cHandler := config.NewHandler()
	env, err := cHandler.Environment("user")
	if err != nil {
		panic(err)
	}

	dbHandler := database.NewHandler(env.Host.Dsn)
	dbHandler.Database().RegisterModel(
		(*models.UserToRole)(nil),
	)

	if err := srv.NewGRPC(env.Host.Url, dbHandler, func(db *bun.DB, grpc *grpc.Server) {
		auth := utils.NewAuthWrapper(env.Config.Viper.GetString("secret"))
		proto.RegisterUserServiceServer(grpc, pkg.NewServer(db, auth))
	}); err != nil {
		panic(err)
	}
}
