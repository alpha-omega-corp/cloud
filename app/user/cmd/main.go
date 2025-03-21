package main

import (
	"context"
	"embed"
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

var (
	//go:embed config
	embedFS embed.FS
)

func main() {
	configHandler := config.NewHandler(embedFS.ReadFile("config/config.yaml"))
	cfg, err := configHandler.LoadAs(context.Background(), "user")
	if err != nil {
		panic(err)
	}

	dbHandler := database.NewHandler(cfg.Dsn)
	dbHandler.Database().RegisterModel(
		(*models.UserToRole)(nil),
	)

	if err := srv.NewGRPC(cfg.Url, dbHandler, func(db *bun.DB, grpc *grpc.Server) {
		auth := utils.NewAuthWrapper(cfg.Env.GetString("secret"))
		proto.RegisterUserServiceServer(grpc, pkg.NewServer(db, auth))
	}); err != nil {
		panic(err)
	}
}
