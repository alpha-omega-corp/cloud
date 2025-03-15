package main

import (
	"github.com/alpha-omega-corp/cloud/app/docker/pkg/models"
	"github.com/alpha-omega-corp/cloud/core/config"
	"github.com/alpha-omega-corp/cloud/core/database"
	_ "github.com/spf13/viper/remote"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	cHandler := config.NewHandler()

	env, err := cHandler.Environment("docker")
	if err != nil {
		panic(err)
	}

	dbHandler := database.NewHandler(env.Host.Dsn)

	defer func(db *bun.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal(err)
		}
	}(dbHandler.Database())

	appCli := &cli.App{
		Name:  "user-svc",
		Usage: "bootstrap the service",
		Commands: []*cli.Command{
			migrateCommand(dbHandler.Database()),
		},
	}

	if err := appCli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func migrateCommand(db *bun.DB) *cli.Command {
	db.RegisterModel(
		(*models.Dockerfile)(nil),
	)

	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
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
				Name:  "reset",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					if err := db.ResetModel(c.Context,
						(*models.Dockerfile)(nil),
					); err != nil {
						return err
					}

					return nil
				},
			},
		},
	}
}
