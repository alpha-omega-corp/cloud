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
	"github.com/alpha-omega-corp/cloud/core/types"
	_ "github.com/spf13/viper/remote"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"os"
)

var (
	//go:embed config
	embedFS embed.FS
)

type App struct {
	name      string
	config    types.Config
	dbHandler *database.Handler
	models    interface{}
}

func NewApp(efs embed.FS, name string) *App {
	configFile, err := efs.ReadFile(config.GetConfigPath())
	if err != nil {
		panic(err)
	}

	cfg := config.NewHandler(configFile).LoadAs(context.Background(), name)

	return &App{
		name:      name,
		config:    cfg,
		dbHandler: database.NewHandler(cfg.Dsn),
	}
}

func (app *App) WithModels(models ...interface{}) *App {
	app.dbHandler.Database().RegisterModel(models)
	app.models = models

	return app
}

func (app *App) Start(init func(config types.Config, db *bun.DB, grpc *grpc.Server)) {
	if err := srv.NewGRPC(app.config.Url, app.dbHandler, func(db *bun.DB, grpc *grpc.Server) {
		init(app.config, db, grpc)
	}); err != nil {
		panic(err)
	}

	appCli := &cli.App{
		Name:  app.name,
		Usage: "cloud application cli",
		Commands: []*cli.Command{
			app.startCommand(),
			app.migrateCommand(),
		},
	}

	if err := appCli.Run(os.Args); err != nil {
		panic(err)
	}
}

func (app *App) startCommand() *cli.Command {
	return &cli.Command{
		Name: "server",
		Action: func(c *cli.Context) error {
			return nil
		},
	}
}

func (app *App) migrateCommand() *cli.Command {
	db := app.dbHandler.Database()
	db.RegisterModel(app.models)

	return &cli.Command{
		Name: "db",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					migrator := migrate.NewMigrator(db, migrate.NewMigrations())
					return migrator.Init(c.Context)
				},
			},
			{
				Name: "reset",
				Action: func(c *cli.Context) error {
					if err := db.ResetModel(c.Context, app.models); err != nil {
						return err
					}

					fixture := dbfixture.New(db)
					if err := fixture.Load(c.Context, os.DirFS("cmd/migrations/fixtures"), "fixture.yml"); err != nil {
						panic(err)
					}

					return nil
				},
			},
		},
	}
}

func main() {
	NewApp(embedFS, "user").WithModels(
		(*models.UserToRole)(nil),
	).Start(func(config types.Config, db *bun.DB, grpc *grpc.Server) {
		auth := utils.NewAuthWrapper(config.Env.GetString("secret"))
		proto.RegisterUserServiceServer(grpc, pkg.NewServer(db, auth))
	})
}
