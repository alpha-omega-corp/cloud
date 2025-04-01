package core

import (
	"context"
	"embed"
	"fmt"
	"github.com/alpha-omega-corp/cloud/core/config"
	"github.com/alpha-omega-corp/cloud/core/database"
	srv "github.com/alpha-omega-corp/cloud/core/server"
	"github.com/alpha-omega-corp/cloud/core/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"log"
	"os"
)

type App struct {
	name      string
	config    types.Config
	dbHandler *database.Handler
	dbModels  []any
}

func NewApp(efs embed.FS, name string) *App {
	configFile, err := efs.ReadFile(config.GetConfigPath())
	if err != nil {
		fmt.Printf("read config file error: %v\n", err)
		panic(err)
	}

	fmt.Printf("read config file success: %v\n", string(configFile))

	cfg := config.NewHandler(configFile).LoadAs(context.Background(), name)

	fmt.Println(cfg.Url)

	return &App{
		name:      name,
		config:    cfg,
		dbHandler: database.NewHandler(cfg.Dsn),
	}
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
		log.Fatalf("app start error: %v\n", err)
	}
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

func (app *App) WithModels(models ...any) *App {
	// Register pivot models first
	app.dbModels = append(app.dbModels, models...)
	app.dbHandler.Database().RegisterModel(app.dbModels...)

	return app
}

func (app *App) migrateCommand() *cli.Command {
	db := app.dbHandler.Database()

	if app.dbModels == nil {
		log.Println("db models is nil")
	}

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

					if err := db.ResetModel(c.Context, app.dbModels...); err != nil {
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
