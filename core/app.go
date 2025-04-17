package core

import (
	"context"
	"embed"
	"fmt"
	"github.com/alpha-omega-corp/cloud/core/config"
	"github.com/alpha-omega-corp/cloud/core/database"
	"github.com/alpha-omega-corp/cloud/core/httputils"
	srv "github.com/alpha-omega-corp/cloud/core/server"
	"github.com/alpha-omega-corp/cloud/core/types"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/migrate"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"github.com/urfave/cli/v3"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	name string

	dbHandler *database.Handler
	dbModels  []any

	config        *types.Config
	configFS      embed.FS
	ConfigHandler *config.Handler
}

func NewApp(efs embed.FS, name string) *App {
	return &App{
		name:      name,
		configFS:  efs,
		config:    nil,
		dbHandler: nil,
	}
}

func (app *App) CreateApi(init func(router *bunrouter.Router, configHandler config.Handler)) {
	appCli := &cli.Command{
		Usage: "cloud application cli",
		Commands: []*cli.Command{
			app.createCommand("app", "server", func(ctx context.Context, cmd *cli.Command) {
				router := bunrouter.New(
					bunrouter.WithMiddleware(reqlog.NewMiddleware(
						reqlog.WithEnabled(true),
						reqlog.WithVerbose(true),
					)))

				// Create clients
				init(router, *app.ConfigHandler)

				// Listen and serve
				var handler http.Handler
				handler = httputils.ExitOnPanicHandler{Next: router}

				httpSrv := &http.Server{
					Addr:         "0.0.0.0:3000",
					ReadTimeout:  60 * time.Second,
					WriteTimeout: 60 * time.Second,
					IdleTimeout:  60 * time.Second,
					Handler:      handler,
				}

				go func() {
					if err := httpSrv.ListenAndServe(); err != nil && err.Error() != "http: Server closed" {
						log.Printf("ListenAndServe failed: %s", err)
					}
				}()

				fmt.Printf("listening on http://%s\n", httpSrv.Addr)

				// Create keyboard listener
				ch := make(chan os.Signal, 3)
				signal.Notify(
					ch,
					syscall.SIGINT,
					syscall.SIGQUIT,
					syscall.SIGTERM,
				)
			}),
		},
	}

	if err := appCli.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("app start error: %v\n", err)
	}
}

func (app *App) CreateServer(init func(config *types.Config, db *bun.DB, grpc *grpc.Server), models ...any) {
	app.dbModels = append(app.dbModels, models...)

	appCli := &cli.Command{
		Usage: "cloud application cli",
		Commands: []*cli.Command{
			app.serverCommand(init),
			app.migrateCommand(),
		},
	}

	if err := appCli.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("app start error: %v\n", err)
	}
}

func (app *App) serverCommand(init func(config *types.Config, db *bun.DB, grpc *grpc.Server)) *cli.Command {
	return app.createCommand("app", "server", func(ctx context.Context, cmd *cli.Command) {
		if err := srv.NewGRPC(app.config.Url, app.dbHandler, func(db *bun.DB, grpc *grpc.Server) {
			init(app.config, db, grpc)
			fmt.Printf("server start success\n")
		}); err != nil {
			panic(err)
		}
	})
}

func (app *App) migrateCommand() *cli.Command {
	return app.createCommand("db", "migration", func(ctx context.Context, cmd *cli.Command) {
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
}

func (app *App) createCommand(category string, name string, action func(ctx context.Context, cmd *cli.Command)) *cli.Command {
	return &cli.Command{
		Name:     name,
		Category: category,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env",
				Aliases: []string{"e"},
				Value:   "local",
				Usage:   "environment for configuration file",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			env := cmd.String("env")
			app.loadConfig(env)

			app.dbHandler = database.NewHandler(app.config.Dsn)
			app.dbHandler.Database().RegisterModel(app.dbModels...)

			action(ctx, cmd)

			return nil
		},
	}
}
