package core

import (
	"context"
	"embed"
	"fmt"
	"github.com/alpha-omega-corp/cloud/core/database"
	"github.com/alpha-omega-corp/cloud/core/server"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
	"github.com/uptrace/bun/migrate"
	"github.com/uptrace/bunrouter"
	"github.com/uptrace/bunrouter/extra/bunrouterotel"
	"github.com/uptrace/bunrouter/extra/reqlog"
	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type App struct {
	name          string
	dbHandler     *database.Handler
	configHandler *ConfigHandler

	fs embed.FS
}

func NewApp(efs embed.FS, name string) *App {
	return &App{
		name:      name,
		fs:        efs,
		dbHandler: nil,
	}
}

func (app *App) CreateApi(init func(router *bunrouter.Router, configHandler *ConfigHandler)) os.Signal {
	appCli := &cli.Command{
		Usage: "cloud application cli",
		Commands: []*cli.Command{
			app.newHttpCommand(init),
		},
	}

	if err := appCli.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("app start error: %v\n", err)
	}

	// Create keyboard listener
	ch := make(chan os.Signal, 3)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
	)

	return <-ch
}

func (app *App) CreateApp(init func(config *Config, db *bun.DB, grpc *grpc.Server), models ...any) {
	appCli := &cli.Command{
		Usage: "cloud application cli",
		Commands: []*cli.Command{
			app.newGrpcCommand(init),
			app.migrateCommand(),
		},
	}

	if err := appCli.Run(context.Background(), os.Args); err != nil {
		log.Fatalf("app start error: %v\n", err)
	}
}

func (app *App) newGrpcCommand(init func(config *Config, db *bun.DB, grpc *grpc.Server)) *cli.Command {
	return app.createCommand("app", "server", func(ctx context.Context, cmd *cli.Command) {
		if err := server.NewGRPC(*app.configHandler.config.Url, app.dbHandler, func(db *bun.DB, grpc *grpc.Server) {
			init(app.configHandler.config, db, grpc)
			fmt.Printf("server start success\n")
		}); err != nil {
			panic(err)
		}
	})
}

func (app *App) newHttpCommand(init func(router *bunrouter.Router, configHandler *ConfigHandler)) *cli.Command {
	return app.createCommand("app", "server", func(ctx context.Context, cmd *cli.Command) {
		r := bunrouter.New(
			bunrouter.WithMiddleware(reqlog.NewMiddleware(
				reqlog.WithEnabled(true),
				reqlog.WithVerbose(true),
			)),
			bunrouter.Use(bunrouterotel.NewMiddleware(
				bunrouterotel.WithClientIP(),
			)))

		init(r, app.configHandler)

		// Listen and serve
		handler := otelhttp.NewHandler(r, "")
		httpSrv := &http.Server{
			Addr:         *app.configHandler.config.Url,
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
	})
}

func (app *App) migrateCommand() *cli.Command {
	return app.createCommand("db", "migration", func(ctx context.Context, cmd *cli.Command) {
		db := app.dbHandler.Database()

		migrator := migrate.NewMigrator(db, migrate.NewMigrations())
		if err := migrator.Init(ctx); err != nil {
			panic(err)
		}

		if err := db.ResetModel(ctx, app.dbHandler.Models); err != nil {
			panic(err)
		}

		fixture := dbfixture.New(db)
		if err := fixture.Load(ctx, os.DirFS("cmd/fixtures"), "fixture.yml"); err != nil {
			fmt.Printf("load fixture error: %v\n", err)
			panic(err)
		}
	})
}

func (app *App) loadConfig(env string, name string) {
	configFile, err := app.fs.ReadFile(GetConfigPath(env))
	if err != nil {
		log.Fatalf("read config file error: %v\n", err)
	}

	app.configHandler = NewConfigHandler(context.Background(), name, configFile)
}

func (app *App) createCommand(category string, name string, action func(ctx context.Context, cmd *cli.Command), models ...any) *cli.Command {
	return &cli.Command{
		Name:     name,
		Category: category,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "env",
				Aliases: []string{"e"},
				Value:   "local",
				Usage:   "environment to select configuration file",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			env := cmd.String("env")
			app.loadConfig(env, app.name)

			if app.configHandler.config.Dsn != nil {
				app.dbHandler = database.NewHandler(*app.configHandler.config.Dsn)
				app.dbHandler.WithModels(models...)
			}

			action(ctx, cmd)

			return nil
		},
	}
}
