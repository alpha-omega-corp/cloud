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
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"log"
	"os"
)

type App struct {
	name      string
	config    *types.Config
	dbHandler *database.Handler
	dbModels  []any
	configFS  embed.FS
}

func NewApp(efs embed.FS, name string) *App {
	return &App{
		name:      name,
		configFS:  efs,
		config:    nil,
		dbHandler: nil,
	}
}

func (app *App) Start(init func(config *types.Config, db *bun.DB, grpc *grpc.Server), models ...any) {
	app.dbModels = append(app.dbModels, models...)

	appCli := &cli.Command{
		Name:  "app",
		Usage: "cloud application cli",
		Commands: []*cli.Command{
			app.startCommand(init),
			app.migrateCommand(),
		},
	}

	if err := appCli.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("app start error: %v\n", err)
	}
}

func (app *App) startCommand(init func(config *types.Config, db *bun.DB, grpc *grpc.Server)) *cli.Command {
	return app.createCommand("start", func(ctx context.Context, cmd *cli.Command) {
		if err := srv.NewGRPC(app.config.Url, app.dbHandler, func(db *bun.DB, grpc *grpc.Server) {
			init(app.config, db, grpc)
			fmt.Printf("server start success\n")
		}); err != nil {
			panic(err)
		}
	})
}

func (app *App) migrateCommand() *cli.Command {
	return app.createCommand("migration", func(ctx context.Context, cmd *cli.Command) {
		db := app.dbHandler.Database()

		migrator := migrate.NewMigrator(db, migrate.NewMigrations())
		if err := migrator.Init(ctx); err != nil {
			panic(err)
		}

		if err := db.ResetModel(ctx, app.dbModels...); err != nil {
			panic(err)
		}

		fixture := dbfixture.New(db)
		if err := fixture.Load(ctx, os.DirFS("cmd/fixtures"), "fixture.yml"); err != nil {
			fmt.Printf("load fixture error: %v\n", err)
			panic(err)
		}
	})
}

func (app *App) loadConfig(env string) {
	configFile, err := app.configFS.ReadFile(config.GetConfigPath(env))
	if err != nil {
		log.Fatalf("read config file error: %v\n", err)
	}

	app.config = config.NewHandler(configFile).LoadAs(context.Background(), app.name)
	app.dbHandler = database.NewHandler(app.config.Dsn)
}

func (app *App) createCommand(name string, action func(ctx context.Context, cmd *cli.Command)) *cli.Command {
	return &cli.Command{
		Name: name,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env",
				Aliases: []string{"e"},
				Value:   "local",
				Usage:   "choose environment name",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			env := cmd.String("env")
			app.loadConfig(env)

			app.dbHandler.Database().RegisterModel(app.dbModels...)

			action(ctx, cmd)

			return nil
		},
	}
}
