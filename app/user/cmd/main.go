package main

import (
	"context"
	"embed"
	"fmt"
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
	"log"
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
	dbModels  interface{}
}

func NewApp(efs embed.FS, name string) *App {
	configFile, err := efs.ReadFile(config.GetConfigPath())
	if err != nil {
		fmt.Printf("read config file error: %v\n", err)
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

	return app
}

func (app *App) Start(init func(config types.Config, db *bun.DB, grpc *grpc.Server)) {
	appCli := &cli.App{
		Name:  app.name,
		Usage: "cloud application cli",
		Commands: []*cli.Command{
			app.startCommand(init),
			app.migrateCommand(),
		},
	}

	if err := appCli.Run(os.Args); err != nil {
		fmt.Printf("app start error: %v\n", err)
		panic(err)
	}
}

func (app *App) bootstrap() {

}

func (app *App) startCommand(init func(config types.Config, db *bun.DB, grpc *grpc.Server)) *cli.Command {
	return &cli.Command{
		Name: "server",
		Action: func(c *cli.Context) error {
			if err := srv.NewGRPC(app.config.Url, app.dbHandler, func(db *bun.DB, grpc *grpc.Server) {
				init(app.config, db, grpc)
			}); err != nil {
				panic(err)
			}

			return nil
		},
	}
}

func (app *App) migrateCommand() *cli.Command {
	// TODO : pass models (err = passing struct to interface)
	app.dbModels = ((*models.Service)(nil))

	if app.dbModels == nil {
		log.Fatalf("db models is nil, crashing...")
	}

	db := app.dbHandler.Database()
	db.RegisterModel(app.dbModels)

	return &cli.Command{
		Name: "db",
		Subcommands: []*cli.Command{
			{
				Name:  "migration",
				Usage: "run migration",
				Action: func(c *cli.Context) error {
					migrator := migrate.NewMigrator(db, migrate.NewMigrations())

					if err := migrator.Init(c.Context); err != nil {
						return err
					}

					if err := db.ResetModel(c.Context, app.dbModels); err != nil {
						return err
					}

					fixture := dbfixture.New(db)
					if err := fixture.Load(c.Context, os.DirFS("cmd/fixtures"), "fixture.yml"); err != nil {
						fmt.Printf("load fixture error: %v\n", err)
						panic(err)
					}

					return nil
				},
			},
		},
	}
}

func main() {
	/*
		WithModels(
					(*models.User)(nil),
					(*models.Role)(nil),
					(*models.UserToRole)(nil),
					(*models.Service)(nil),
					(*models.Permission)(nil),
				).
	*/
	NewApp(embedFS, "user").Start(func(config types.Config, db *bun.DB, grpc *grpc.Server) {
		auth := utils.NewAuthWrapper(config.Env.GetString("secret"))
		proto.RegisterUserServiceServer(grpc, pkg.NewServer(db, auth))
	})
}
